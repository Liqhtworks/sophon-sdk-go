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

// checks if the UploadPartResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &UploadPartResponse{}

// UploadPartResponse struct for UploadPartResponse
type UploadPartResponse struct {
	PartNumber int64 `json:"part_number"`
	Received bool `json:"received"`
	AdditionalProperties map[string]interface{}
}

type _UploadPartResponse UploadPartResponse

// NewUploadPartResponse instantiates a new UploadPartResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUploadPartResponse(partNumber int64, received bool) *UploadPartResponse {
	this := UploadPartResponse{}
	this.PartNumber = partNumber
	this.Received = received
	return &this
}

// NewUploadPartResponseWithDefaults instantiates a new UploadPartResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUploadPartResponseWithDefaults() *UploadPartResponse {
	this := UploadPartResponse{}
	return &this
}

// GetPartNumber returns the PartNumber field value
func (o *UploadPartResponse) GetPartNumber() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.PartNumber
}

// GetPartNumberOk returns a tuple with the PartNumber field value
// and a boolean to check if the value has been set.
func (o *UploadPartResponse) GetPartNumberOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PartNumber, true
}

// SetPartNumber sets field value
func (o *UploadPartResponse) SetPartNumber(v int64) {
	o.PartNumber = v
}

// GetReceived returns the Received field value
func (o *UploadPartResponse) GetReceived() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.Received
}

// GetReceivedOk returns a tuple with the Received field value
// and a boolean to check if the value has been set.
func (o *UploadPartResponse) GetReceivedOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Received, true
}

// SetReceived sets field value
func (o *UploadPartResponse) SetReceived(v bool) {
	o.Received = v
}

func (o UploadPartResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o UploadPartResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["part_number"] = o.PartNumber
	toSerialize["received"] = o.Received

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *UploadPartResponse) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"part_number",
		"received",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err;
	}

	for _, requiredProperty := range(requiredProperties) {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varUploadPartResponse := _UploadPartResponse{}

	err = json.Unmarshal(data, &varUploadPartResponse)

	if err != nil {
		return err
	}

	*o = UploadPartResponse(varUploadPartResponse)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "part_number")
		delete(additionalProperties, "received")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableUploadPartResponse struct {
	value *UploadPartResponse
	isSet bool
}

func (v NullableUploadPartResponse) Get() *UploadPartResponse {
	return v.value
}

func (v *NullableUploadPartResponse) Set(val *UploadPartResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableUploadPartResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableUploadPartResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUploadPartResponse(val *UploadPartResponse) *NullableUploadPartResponse {
	return &NullableUploadPartResponse{value: val, isSet: true}
}

func (v NullableUploadPartResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUploadPartResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


