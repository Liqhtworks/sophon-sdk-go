/*
SOPHON Encoding API

REST API for submitting, monitoring, and retrieving SOPHON encoding jobs.  Authentication is via Bearer API key or session cookie. All POST endpoints require an Idempotency-Key header. List endpoints use opaque cursor-based pagination.  ---  ## Integration example  A real-world walkthrough of how [Daisy](https://daisy.so) wires SOPHON into two production flows — user-uploaded video compression and automatic post-generation encoding after video rendering. Both converge on the same adapter and state machine; only the source differs.  The patterns below are the ones that transfer cleanly to any integration.  ### 1. One thin adapter, one method per endpoint  Keep the HTTP surface boring. Axios (or your stack's equivalent), a per-endpoint idempotency key, and no enum for profile names:  ```ts @Injectable() export class SophonService {   private client() {     return axios.create({       baseURL: process.env.SOPHON_BASE_URL,       headers: { Authorization: `Bearer ${process.env.SOPHON_API_KEY}` },       timeout: 60_000,     });   }    async createUploadSession(req, idempotencyKey) { /_* POST /v1/uploads *_/ }   async uploadChunk(uploadId, partNumber, bytes) { /_* PUT /v1/uploads/{id}/parts/{n} *_/ }   async completeUpload(uploadId, idempotencyKey) { /_* POST /v1/uploads/{id}/complete *_/ }   async createJob(req, idempotencyKey) { /_* POST /v1/jobs *_/ }   async getJob(id) { /_* GET /v1/jobs/{id} *_/ }   async downloadOutputStream(jobId) { /_* GET /v1/jobs/{id}/output *_/ } } ```  **Suffix idempotency keys per endpoint.** SOPHON scopes dedupe per route but a shared key collides across retries that hit different endpoints. Do this:  ```ts const base = `video:${video.id}:v1`; await sophon.createUploadSession(req,  `${base}:create-upload`); await sophon.completeUpload(uploadId,  `${base}:complete-upload`); await sophon.createJob(req,            `${base}:create-job`); ```  **Profile names are strings, not an enum.** We add and rename profiles (`sophon-espresso` → `sophon-auto` → future variants). A TypeScript union will drift; let the server validate.  ### 2. Model your pipeline as a state machine  Persist a single `sophonState` JSON column per row. `jobId === null` routes to dispatch; anything else polls that job:  ```ts interface SophonState {   jobId: string | null;          // null = not dispatched; string = poll it   uploadId?: string;             // persist between upload + createJob   profile?: string;              // sophon-auto | sophon-espresso | ...   dispatchRetries: number;       // 3 strikes → fallback   downloadRetries: number;   lastError?: { stage, code, message, at }; }  // In your cron (5-second tick is plenty): if (state.jobId === null) {   await dispatch(video, state);  // upload + createJob } else {   await poll(video, state);      // getJob + (if completed) downloadAndComplete } ```  Persisting `uploadId` between the upload completion and the `createJob` call matters — a crash in that window otherwise re-uploads the file.  ### 3. Stream for large sources; buffer for small  User-uploaded sources can be 1 GB+. Stream S3 → SOPHON in chunks equal to `session.chunk_size` from the createUploadSession response:  ```ts async uploadStream(stream, fileName, mimeType, fileSize) {   const session = await this.createUploadSession({     file_name: fileName, file_size: fileSize, mime_type: mimeType,   });   let partIndex = 0, buffer = Buffer.alloc(0);   for await (const chunk of stream) {     buffer = Buffer.concat([buffer, chunk]);     while (buffer.length >= session.chunk_size) {       await this.uploadChunk(session.id, partIndex++,         buffer.subarray(0, session.chunk_size));       buffer = buffer.subarray(session.chunk_size);     }   }   if (buffer.length > 0) {     await this.uploadChunk(session.id, partIndex, buffer);   }   return this.completeUpload(session.id); } ```  Generated outputs from a model run are typically <30 MB — for those, a buffered upload path is simpler and avoids managing a stream lifetime.  ### 4. Always keep a fallback URL  Before a row enters your encoding state, make sure the source is already playable from your CDN. Every SOPHON failure then degrades to \"use the original\" — the user's video never disappears because SOPHON is slow or down. This is the single most important invariant:  ```ts await videoRepository.update({ id: video.id }, {   videoUrl: sourceCloudfrontUrl,   // fallback URL, stays intact   status: VideoStatus.EncodingPending,   sophonState: { jobId: null, profile, dispatchRetries: 0, downloadRetries: 0 },   sourceFileSize: sourceBytes, }); ```  On any terminal failure (structured `retryable: false`, retry budget exhausted, 404 on getJob, 23h stuck-row guard), flip status back to `Done` with `videoUrl` unchanged. SOPHON is enhancement, not a delivery dependency.  ### 5. Handle the \"no-gain\" success path  `sophon-auto` runs a pre-probe and, when it decides the output wouldn't be smaller than the source, returns `final_artifact: \"original\"` and `saved_percent: 0`. Skip the output download — the source already lives in your bucket:  ```ts if (job.status === 'completed') {   if (job.final_artifact === 'original') {     // Persist outputFileSize = sourceFileSize so your UI shows     // \"no reduction\" instead of a missing value.     await completeWithFallbackOutput(video, job.output?.bytes ?? null);     return;   }   await downloadAndComplete(video, state, job.output?.bytes ?? null); } ```  ### 6. Finalize by streaming into your own storage  `GET /v1/jobs/{id}/output` returns a 302 to a presigned URL with a 24h TTL. Stream that directly into your bucket — no temp file, no buffering:  ```ts const { stream } = await sophon.downloadOutputStream(state.jobId); const outputKey = `encoded/${video.userId}/${video.id}.mp4`; await fileService.uploadStream(outputKey, stream, 'video/mp4'); await videoRepository.update({ id: video.id }, {   videoUrl: fileService.cloudfrontUrl(outputKey),   outputFileSize: sophonOutputBytes,   status: VideoStatus.Done, }); ```  ### 7. Failure taxonomy  | Error | Handling | |---|---| | Structured `retryable: false` from SOPHON | Terminal. Fall back to `Done` with source URL. | | Retryable upload / createJob failure | Increment `dispatchRetries`; after 3, fall back. | | Retryable download failure | Increment `downloadRetries`; after 3, fall back. | | `getJob` → HTTP 404 | Terminal. Job expired or never created. Fall back. | | Transient poll network error | Do nothing; next tick retries. Don't burn retry budget. | | Row stuck in encode state > 23h | Fall back (safety net against orphans). |  ### Minimal config  ```bash SOPHON_API_KEY=sk_live_... SOPHON_BASE_URL=https://api.liqhtworks.xyz ``` 

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package sophon

import (
	"encoding/json"
	"fmt"
)

// JobProfile Encoding profile ID. Coffee-themed naming: prep time maps to encode speed (espresso = fast, cortado = medium, americano = slow). The naming suffixes compose:  - bare name → 8-bit HEVC Main (universal decoder compatibility,   default) - `-10bit` suffix → HEVC Main10 output. Requires a decoder that   supports Main10 (modern phones, modern TVs, Safari, Chrome   with hardware decode). Older / low-end devices may not play   Main10 output. Pick only when you know the downstream   pipeline supports it. - `-hq` suffix → quality-biased 8-bit variant for heavy source   formats (ProRes, DNxHD, high-bitrate camera originals,   mastering-grade intermediates). Files will be larger than   the standard tier; pick when preserving detail matters more   than bitrate. Broad device compatibility (8-bit Main). - `-hq-10bit` suffix → combines HQ with Main10 to preserve   10-bit depth end-to-end. Same Main10 compatibility caveat   as above; pick for ProRes 422/4444, DNxHD, BRAW, or camera   masters where detail AND bit depth matter and you control   the downstream pipeline.  **For broad audience playback, pick `sophon-auto` or an explicit 8-bit coffee profile.** `sophon-auto` produces 8-bit Main output regardless of source bit depth.  If you're not sure which to pick, use `sophon-auto` — the API picks per-source settings tuned for consistent output regardless of what you submit, and automatically re-encodes at stricter settings if the first pass doesn't hold up.  **8-bit (standard, default):** - `sophon-espresso` — fastest, lowest compression - `sophon-cortado` — balanced speed and quality - `sophon-americano` — slowest, highest compression  **8-bit HQ** (max quality preservation for heavy formats): - `sophon-espresso-hq` - `sophon-cortado-hq` - `sophon-americano-hq`  **10-bit (HEVC Main10):** - `sophon-espresso-10bit` - `sophon-cortado-10bit` - `sophon-americano-10bit`  **10-bit HQ** (max quality preservation AND preserves 10-bit depth): - `sophon-espresso-hq-10bit` - `sophon-cortado-hq-10bit` - `sophon-americano-hq-10bit`  **Adaptive dispatcher:** - `sophon-auto` — content-adaptive. The API probes each source,   picks tuned settings, and re-encodes at stricter settings if   the first pass doesn't hold up. `profile` on the job response   stays `sophon-auto`; `effective_profile_id` records the   concrete variant the API actually ran. 
type JobProfile string

// List of JobProfile
const (
	SOPHON_ESPRESSO JobProfile = "sophon-espresso"
	SOPHON_CORTADO JobProfile = "sophon-cortado"
	SOPHON_AMERICANO JobProfile = "sophon-americano"
	SOPHON_ESPRESSO_HQ JobProfile = "sophon-espresso-hq"
	SOPHON_CORTADO_HQ JobProfile = "sophon-cortado-hq"
	SOPHON_AMERICANO_HQ JobProfile = "sophon-americano-hq"
	SOPHON_ESPRESSO_10BIT JobProfile = "sophon-espresso-10bit"
	SOPHON_CORTADO_10BIT JobProfile = "sophon-cortado-10bit"
	SOPHON_AMERICANO_10BIT JobProfile = "sophon-americano-10bit"
	SOPHON_ESPRESSO_HQ_10BIT JobProfile = "sophon-espresso-hq-10bit"
	SOPHON_CORTADO_HQ_10BIT JobProfile = "sophon-cortado-hq-10bit"
	SOPHON_AMERICANO_HQ_10BIT JobProfile = "sophon-americano-hq-10bit"
	SOPHON_AUTO JobProfile = "sophon-auto"
)

// All allowed values of JobProfile enum
var AllowedJobProfileEnumValues = []JobProfile{
	"sophon-espresso",
	"sophon-cortado",
	"sophon-americano",
	"sophon-espresso-hq",
	"sophon-cortado-hq",
	"sophon-americano-hq",
	"sophon-espresso-10bit",
	"sophon-cortado-10bit",
	"sophon-americano-10bit",
	"sophon-espresso-hq-10bit",
	"sophon-cortado-hq-10bit",
	"sophon-americano-hq-10bit",
	"sophon-auto",
}

func (v *JobProfile) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := JobProfile(value)
	for _, existing := range AllowedJobProfileEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid JobProfile", value)
}

// NewJobProfileFromValue returns a pointer to a valid JobProfile
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewJobProfileFromValue(v string) (*JobProfile, error) {
	ev := JobProfile(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for JobProfile: valid values are %v", v, AllowedJobProfileEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v JobProfile) IsValid() bool {
	for _, existing := range AllowedJobProfileEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to JobProfile value
func (v JobProfile) Ptr() *JobProfile {
	return &v
}

type NullableJobProfile struct {
	value *JobProfile
	isSet bool
}

func (v NullableJobProfile) Get() *JobProfile {
	return v.value
}

func (v *NullableJobProfile) Set(val *JobProfile) {
	v.value = val
	v.isSet = true
}

func (v NullableJobProfile) IsSet() bool {
	return v.isSet
}

func (v *NullableJobProfile) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableJobProfile(val *JobProfile) *NullableJobProfile {
	return &NullableJobProfile{value: val, isSet: true}
}

func (v NullableJobProfile) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableJobProfile) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

