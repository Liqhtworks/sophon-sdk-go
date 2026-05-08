/*
SOPHON Encoding API

REST API for submitting, monitoring, and retrieving SOPHON encoding jobs.  Authentication is via Bearer API key or session cookie. All POST endpoints require an Idempotency-Key header. List endpoints use opaque cursor-based pagination.  ---  ## Integration example  A real-world walkthrough of how [Daisy](https://daisy.so) wires SOPHON into two production flows — user-uploaded video compression and automatic post-generation encoding after video rendering. Both converge on the same adapter and state machine; only the source differs.  The patterns below are the ones that transfer cleanly to any integration.  ### 1. One thin adapter, one method per endpoint  Keep the HTTP surface boring. Axios (or your stack's equivalent), a per-endpoint idempotency key, and no enum for profile names:  ```ts @Injectable() export class SophonService {   private client() {     return axios.create({       baseURL: process.env.SOPHON_BASE_URL,       headers: { Authorization: `Bearer ${process.env.SOPHON_API_KEY}` },       timeout: 60_000,     });   }    async createUploadSession(req, idempotencyKey) { /_* POST /v1/uploads *_/ }   async uploadChunk(uploadId, partNumber, bytes) { /_* PUT /v1/uploads/{id}/parts/{n} *_/ }   async completeUpload(uploadId, idempotencyKey) { /_* POST /v1/uploads/{id}/complete *_/ }   async createJob(req, idempotencyKey) { /_* POST /v1/jobs *_/ }   async getJob(id) { /_* GET /v1/jobs/{id} *_/ }   async downloadOutputStream(jobId) { /_* GET /v1/jobs/{id}/output *_/ } } ```  **Suffix idempotency keys per endpoint.** SOPHON scopes dedupe per route but a shared key collides across retries that hit different endpoints. Do this:  ```ts const base = `video:${video.id}:v1`; await sophon.createUploadSession(req,  `${base}:create-upload`); await sophon.completeUpload(uploadId,  `${base}:complete-upload`); await sophon.createJob(req,            `${base}:create-job`); ```  **Profile names are strings, not an enum.** We add and rename profiles (`sophon-espresso` → `sophon-auto` → future variants). A TypeScript union will drift; let the server validate.  ### 2. Model your pipeline as a state machine  Persist a single `sophonState` JSON column per row. `jobId === null` routes to dispatch; anything else polls that job:  ```ts interface SophonState {   jobId: string | null;          // null = not dispatched; string = poll it   uploadId?: string;             // persist between upload + createJob   profile?: string;              // sophon-auto | sophon-espresso | ...   dispatchRetries: number;       // 3 strikes → fallback   downloadRetries: number;   lastError?: { stage, code, message, at }; }  // In your cron (5-second tick is plenty): if (state.jobId === null) {   await dispatch(video, state);  // upload + createJob } else {   await poll(video, state);      // getJob + (if completed) downloadAndComplete } ```  Persisting `uploadId` between the upload completion and the `createJob` call matters — a crash in that window otherwise re-uploads the file.  ### 3. Stream for large sources; buffer for small  User-uploaded sources can be 1 GB+. Stream S3 → SOPHON in chunks equal to `session.chunk_size` from the createUploadSession response:  ```ts async uploadStream(stream, fileName, mimeType, fileSize) {   const session = await this.createUploadSession({     file_name: fileName, file_size: fileSize, mime_type: mimeType,   });   let partIndex = 0, buffer = Buffer.alloc(0);   for await (const chunk of stream) {     buffer = Buffer.concat([buffer, chunk]);     while (buffer.length >= session.chunk_size) {       await this.uploadChunk(session.id, partIndex++,         buffer.subarray(0, session.chunk_size));       buffer = buffer.subarray(session.chunk_size);     }   }   if (buffer.length > 0) {     await this.uploadChunk(session.id, partIndex, buffer);   }   return this.completeUpload(session.id); } ```  Generated outputs from a model run are typically <30 MB — for those, a buffered upload path is simpler and avoids managing a stream lifetime.  ### 4. Always keep a fallback URL  Before a row enters your encoding state, make sure the source is already playable from your CDN. Every SOPHON failure then degrades to \"use the original\" — the user's video never disappears because SOPHON is slow or down. This is the single most important invariant:  ```ts await videoRepository.update({ id: video.id }, {   videoUrl: sourceCloudfrontUrl,   // fallback URL, stays intact   status: VideoStatus.EncodingPending,   sophonState: { jobId: null, profile, dispatchRetries: 0, downloadRetries: 0 },   sourceFileSize: sourceBytes, }); ```  On any terminal failure (structured `retryable: false`, retry budget exhausted, 404 on getJob, 23h stuck-row guard), flip status back to `Done` with `videoUrl` unchanged. SOPHON is enhancement, not a delivery dependency.  ### 5. Handle the \"no-gain\" success path  `sophon-auto` runs a pre-probe and, when it decides the output wouldn't be smaller than the source, returns `final_artifact: \"original\"` and `saved_percent: 0`. Skip the output download — the source already lives in your bucket:  ```ts if (job.status === 'completed') {   if (job.final_artifact === 'original') {     // Persist outputFileSize = sourceFileSize so your UI shows     // \"no reduction\" instead of a missing value.     await completeWithFallbackOutput(video, job.output?.bytes ?? null);     return;   }   await downloadAndComplete(video, state, job.output?.bytes ?? null); } ```  ### 6. Finalize by streaming into your own storage  `GET /v1/jobs/{id}/output` returns a 302 to a presigned URL with a 24h TTL. Stream that directly into your bucket — no temp file, no buffering:  ```ts const { stream } = await sophon.downloadOutputStream(state.jobId); const outputKey = `encoded/${video.userId}/${video.id}.mp4`; await fileService.uploadStream(outputKey, stream, 'video/mp4'); await videoRepository.update({ id: video.id }, {   videoUrl: fileService.cloudfrontUrl(outputKey),   outputFileSize: sophonOutputBytes,   status: VideoStatus.Done, }); ```  ### 7. Failure taxonomy  | Error | Handling | |---|---| | Structured `retryable: false` from SOPHON | Terminal. Fall back to `Done` with source URL. | | Retryable upload / createJob failure | Increment `dispatchRetries`; after 3, fall back. | | Retryable download failure | Increment `downloadRetries`; after 3, fall back. | | `getJob` → HTTP 404 | Terminal. Job expired or never created. Fall back. | | Transient poll network error | Do nothing; next tick retries. Don't burn retry budget. | | Row stuck in encode state > 23h | Fall back (safety net against orphans). |  ### Minimal config  ```bash SOPHON_API_KEY=sk_live_... SOPHON_BASE_URL=https://api.liqhtworks.xyz ``` 

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package sophon

import (
	"encoding/json"
)

// checks if the CreateJobOutputOptions type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CreateJobOutputOptions{}

// CreateJobOutputOptions Optional output shaping knobs for a new job.
type CreateJobOutputOptions struct {
	Container *OutputContainer `json:"container,omitempty"`
	// When true, audio is included in the output. MKV preserves source audio streams unchanged. MP4 preserves broadly compatible source audio codecs when possible, and may normalize incompatible codecs to AAC for playback compatibility. When false, the output is video only. 
	Audio *bool `json:"audio,omitempty"`
	// Target output height in pixels. When set, output is scaled down (aspect ratio preserved, width derived from source, both dims rounded to even). If absent or larger than source height, output uses source dimensions. Billing tier is determined by the actual encoded output, not by this requested value. 
	TargetHeight *int32 `json:"target_height,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _CreateJobOutputOptions CreateJobOutputOptions

// NewCreateJobOutputOptions instantiates a new CreateJobOutputOptions object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateJobOutputOptions() *CreateJobOutputOptions {
	this := CreateJobOutputOptions{}
	var container OutputContainer = MP4
	this.Container = &container
	var audio bool = false
	this.Audio = &audio
	return &this
}

// NewCreateJobOutputOptionsWithDefaults instantiates a new CreateJobOutputOptions object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateJobOutputOptionsWithDefaults() *CreateJobOutputOptions {
	this := CreateJobOutputOptions{}
	var container OutputContainer = MP4
	this.Container = &container
	var audio bool = false
	this.Audio = &audio
	return &this
}

// GetContainer returns the Container field value if set, zero value otherwise.
func (o *CreateJobOutputOptions) GetContainer() OutputContainer {
	if o == nil || IsNil(o.Container) {
		var ret OutputContainer
		return ret
	}
	return *o.Container
}

// GetContainerOk returns a tuple with the Container field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateJobOutputOptions) GetContainerOk() (*OutputContainer, bool) {
	if o == nil || IsNil(o.Container) {
		return nil, false
	}
	return o.Container, true
}

// HasContainer returns a boolean if a field has been set.
func (o *CreateJobOutputOptions) HasContainer() bool {
	if o != nil && !IsNil(o.Container) {
		return true
	}

	return false
}

// SetContainer gets a reference to the given OutputContainer and assigns it to the Container field.
func (o *CreateJobOutputOptions) SetContainer(v OutputContainer) {
	o.Container = &v
}

// GetAudio returns the Audio field value if set, zero value otherwise.
func (o *CreateJobOutputOptions) GetAudio() bool {
	if o == nil || IsNil(o.Audio) {
		var ret bool
		return ret
	}
	return *o.Audio
}

// GetAudioOk returns a tuple with the Audio field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateJobOutputOptions) GetAudioOk() (*bool, bool) {
	if o == nil || IsNil(o.Audio) {
		return nil, false
	}
	return o.Audio, true
}

// HasAudio returns a boolean if a field has been set.
func (o *CreateJobOutputOptions) HasAudio() bool {
	if o != nil && !IsNil(o.Audio) {
		return true
	}

	return false
}

// SetAudio gets a reference to the given bool and assigns it to the Audio field.
func (o *CreateJobOutputOptions) SetAudio(v bool) {
	o.Audio = &v
}

// GetTargetHeight returns the TargetHeight field value if set, zero value otherwise.
func (o *CreateJobOutputOptions) GetTargetHeight() int32 {
	if o == nil || IsNil(o.TargetHeight) {
		var ret int32
		return ret
	}
	return *o.TargetHeight
}

// GetTargetHeightOk returns a tuple with the TargetHeight field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateJobOutputOptions) GetTargetHeightOk() (*int32, bool) {
	if o == nil || IsNil(o.TargetHeight) {
		return nil, false
	}
	return o.TargetHeight, true
}

// HasTargetHeight returns a boolean if a field has been set.
func (o *CreateJobOutputOptions) HasTargetHeight() bool {
	if o != nil && !IsNil(o.TargetHeight) {
		return true
	}

	return false
}

// SetTargetHeight gets a reference to the given int32 and assigns it to the TargetHeight field.
func (o *CreateJobOutputOptions) SetTargetHeight(v int32) {
	o.TargetHeight = &v
}

func (o CreateJobOutputOptions) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CreateJobOutputOptions) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Container) {
		toSerialize["container"] = o.Container
	}
	if !IsNil(o.Audio) {
		toSerialize["audio"] = o.Audio
	}
	if !IsNil(o.TargetHeight) {
		toSerialize["target_height"] = o.TargetHeight
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *CreateJobOutputOptions) UnmarshalJSON(data []byte) (err error) {
	varCreateJobOutputOptions := _CreateJobOutputOptions{}

	err = json.Unmarshal(data, &varCreateJobOutputOptions)

	if err != nil {
		return err
	}

	*o = CreateJobOutputOptions(varCreateJobOutputOptions)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "container")
		delete(additionalProperties, "audio")
		delete(additionalProperties, "target_height")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableCreateJobOutputOptions struct {
	value *CreateJobOutputOptions
	isSet bool
}

func (v NullableCreateJobOutputOptions) Get() *CreateJobOutputOptions {
	return v.value
}

func (v *NullableCreateJobOutputOptions) Set(val *CreateJobOutputOptions) {
	v.value = val
	v.isSet = true
}

func (v NullableCreateJobOutputOptions) IsSet() bool {
	return v.isSet
}

func (v *NullableCreateJobOutputOptions) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCreateJobOutputOptions(val *CreateJobOutputOptions) *NullableCreateJobOutputOptions {
	return &NullableCreateJobOutputOptions{value: val, isSet: true}
}

func (v NullableCreateJobOutputOptions) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCreateJobOutputOptions) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


