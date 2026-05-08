/*
SOPHON Encoding API

REST API for submitting, monitoring, and retrieving SOPHON encoding jobs.  Authentication is via Bearer API key or session cookie. All POST endpoints require an Idempotency-Key header. List endpoints use opaque cursor-based pagination.  ---  ## Integration example  A real-world walkthrough of how [Daisy](https://daisy.so) wires SOPHON into two production flows — user-uploaded video compression and automatic post-generation encoding after video rendering. Both converge on the same adapter and state machine; only the source differs.  The patterns below are the ones that transfer cleanly to any integration.  ### 1. One thin adapter, one method per endpoint  Keep the HTTP surface boring. Axios (or your stack's equivalent), a per-endpoint idempotency key, and no enum for profile names:  ```ts @Injectable() export class SophonService {   private client() {     return axios.create({       baseURL: process.env.SOPHON_BASE_URL,       headers: { Authorization: `Bearer ${process.env.SOPHON_API_KEY}` },       timeout: 60_000,     });   }    async createUploadSession(req, idempotencyKey) { /_* POST /v1/uploads *_/ }   async uploadChunk(uploadId, partNumber, bytes) { /_* PUT /v1/uploads/{id}/parts/{n} *_/ }   async completeUpload(uploadId, idempotencyKey) { /_* POST /v1/uploads/{id}/complete *_/ }   async createJob(req, idempotencyKey) { /_* POST /v1/jobs *_/ }   async getJob(id) { /_* GET /v1/jobs/{id} *_/ }   async downloadOutputStream(jobId) { /_* GET /v1/jobs/{id}/output *_/ } } ```  **Suffix idempotency keys per endpoint.** SOPHON scopes dedupe per route but a shared key collides across retries that hit different endpoints. Do this:  ```ts const base = `video:${video.id}:v1`; await sophon.createUploadSession(req,  `${base}:create-upload`); await sophon.completeUpload(uploadId,  `${base}:complete-upload`); await sophon.createJob(req,            `${base}:create-job`); ```  **Profile names are strings, not an enum.** We add and rename profiles (`sophon-espresso` → `sophon-auto` → future variants). A TypeScript union will drift; let the server validate.  ### 2. Model your pipeline as a state machine  Persist a single `sophonState` JSON column per row. `jobId === null` routes to dispatch; anything else polls that job:  ```ts interface SophonState {   jobId: string | null;          // null = not dispatched; string = poll it   uploadId?: string;             // persist between upload + createJob   profile?: string;              // sophon-auto | sophon-espresso | ...   dispatchRetries: number;       // 3 strikes → fallback   downloadRetries: number;   lastError?: { stage, code, message, at }; }  // In your cron (5-second tick is plenty): if (state.jobId === null) {   await dispatch(video, state);  // upload + createJob } else {   await poll(video, state);      // getJob + (if completed) downloadAndComplete } ```  Persisting `uploadId` between the upload completion and the `createJob` call matters — a crash in that window otherwise re-uploads the file.  ### 3. Stream for large sources; buffer for small  User-uploaded sources can be 1 GB+. Stream S3 → SOPHON in chunks equal to `session.chunk_size` from the createUploadSession response:  ```ts async uploadStream(stream, fileName, mimeType, fileSize) {   const session = await this.createUploadSession({     file_name: fileName, file_size: fileSize, mime_type: mimeType,   });   let partIndex = 0, buffer = Buffer.alloc(0);   for await (const chunk of stream) {     buffer = Buffer.concat([buffer, chunk]);     while (buffer.length >= session.chunk_size) {       await this.uploadChunk(session.id, partIndex++,         buffer.subarray(0, session.chunk_size));       buffer = buffer.subarray(session.chunk_size);     }   }   if (buffer.length > 0) {     await this.uploadChunk(session.id, partIndex, buffer);   }   return this.completeUpload(session.id); } ```  Generated outputs from a model run are typically <30 MB — for those, a buffered upload path is simpler and avoids managing a stream lifetime.  ### 4. Always keep a fallback URL  Before a row enters your encoding state, make sure the source is already playable from your CDN. Every SOPHON failure then degrades to \"use the original\" — the user's video never disappears because SOPHON is slow or down. This is the single most important invariant:  ```ts await videoRepository.update({ id: video.id }, {   videoUrl: sourceCloudfrontUrl,   // fallback URL, stays intact   status: VideoStatus.EncodingPending,   sophonState: { jobId: null, profile, dispatchRetries: 0, downloadRetries: 0 },   sourceFileSize: sourceBytes, }); ```  On any terminal failure (structured `retryable: false`, retry budget exhausted, 404 on getJob, 23h stuck-row guard), flip status back to `Done` with `videoUrl` unchanged. SOPHON is enhancement, not a delivery dependency.  ### 5. Handle the \"no-gain\" success path  `sophon-auto` runs a pre-probe and, when it decides the output wouldn't be smaller than the source, returns `final_artifact: \"original\"` and `saved_percent: 0`. Skip the output download — the source already lives in your bucket:  ```ts if (job.status === 'completed') {   if (job.final_artifact === 'original') {     // Persist outputFileSize = sourceFileSize so your UI shows     // \"no reduction\" instead of a missing value.     await completeWithFallbackOutput(video, job.output?.bytes ?? null);     return;   }   await downloadAndComplete(video, state, job.output?.bytes ?? null); } ```  ### 6. Finalize by streaming into your own storage  `GET /v1/jobs/{id}/output` returns a 302 to a presigned URL with a 24h TTL. Stream that directly into your bucket — no temp file, no buffering:  ```ts const { stream } = await sophon.downloadOutputStream(state.jobId); const outputKey = `encoded/${video.userId}/${video.id}.mp4`; await fileService.uploadStream(outputKey, stream, 'video/mp4'); await videoRepository.update({ id: video.id }, {   videoUrl: fileService.cloudfrontUrl(outputKey),   outputFileSize: sophonOutputBytes,   status: VideoStatus.Done, }); ```  ### 7. Failure taxonomy  | Error | Handling | |---|---| | Structured `retryable: false` from SOPHON | Terminal. Fall back to `Done` with source URL. | | Retryable upload / createJob failure | Increment `dispatchRetries`; after 3, fall back. | | Retryable download failure | Increment `downloadRetries`; after 3, fall back. | | `getJob` → HTTP 404 | Terminal. Job expired or never created. Fall back. | | Transient poll network error | Do nothing; next tick retries. Don't burn retry budget. | | Row stuck in encode state > 23h | Fall back (safety net against orphans). |  ### Minimal config  ```bash SOPHON_API_KEY=sk_live_... SOPHON_BASE_URL=https://api.liqhtworks.xyz ``` 

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package sophon

import (
	"encoding/json"
	"time"
	"fmt"
)

// checks if the JobOutputInfo type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &JobOutputInfo{}

// JobOutputInfo struct for JobOutputInfo
type JobOutputInfo struct {
	State string `json:"state"`
	// Output container format (\"mp4\" or \"mkv\").
	Container string `json:"container"`
	// Whether the output file actually contains audio. Reflects the muxed result, not the request flag — a video-only source with audio requested will report false. 
	Audio bool `json:"audio"`
	// Customer-requested output height, echoed back. Null when the job ran at source dimensions (passthrough). 
	TargetHeight *int32 `json:"target_height,omitempty"`
	// Actual encoded output width in pixels (post-ffprobe). Null until the job completes or if the probe failed. 
	Width *int32 `json:"width,omitempty"`
	// Actual encoded output height in pixels. See `width`.
	Height *int32 `json:"height,omitempty"`
	Bytes *int64 `json:"bytes,omitempty"`
	Sha256 *string `json:"sha256,omitempty"`
	RetentionExpiresAt *time.Time `json:"retention_expires_at,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _JobOutputInfo JobOutputInfo

// NewJobOutputInfo instantiates a new JobOutputInfo object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewJobOutputInfo(state string, container string, audio bool) *JobOutputInfo {
	this := JobOutputInfo{}
	this.State = state
	this.Container = container
	this.Audio = audio
	return &this
}

// NewJobOutputInfoWithDefaults instantiates a new JobOutputInfo object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewJobOutputInfoWithDefaults() *JobOutputInfo {
	this := JobOutputInfo{}
	return &this
}

// GetState returns the State field value
func (o *JobOutputInfo) GetState() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.State
}

// GetStateOk returns a tuple with the State field value
// and a boolean to check if the value has been set.
func (o *JobOutputInfo) GetStateOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.State, true
}

// SetState sets field value
func (o *JobOutputInfo) SetState(v string) {
	o.State = v
}

// GetContainer returns the Container field value
func (o *JobOutputInfo) GetContainer() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Container
}

// GetContainerOk returns a tuple with the Container field value
// and a boolean to check if the value has been set.
func (o *JobOutputInfo) GetContainerOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Container, true
}

// SetContainer sets field value
func (o *JobOutputInfo) SetContainer(v string) {
	o.Container = v
}

// GetAudio returns the Audio field value
func (o *JobOutputInfo) GetAudio() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.Audio
}

// GetAudioOk returns a tuple with the Audio field value
// and a boolean to check if the value has been set.
func (o *JobOutputInfo) GetAudioOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Audio, true
}

// SetAudio sets field value
func (o *JobOutputInfo) SetAudio(v bool) {
	o.Audio = v
}

// GetTargetHeight returns the TargetHeight field value if set, zero value otherwise.
func (o *JobOutputInfo) GetTargetHeight() int32 {
	if o == nil || IsNil(o.TargetHeight) {
		var ret int32
		return ret
	}
	return *o.TargetHeight
}

// GetTargetHeightOk returns a tuple with the TargetHeight field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JobOutputInfo) GetTargetHeightOk() (*int32, bool) {
	if o == nil || IsNil(o.TargetHeight) {
		return nil, false
	}
	return o.TargetHeight, true
}

// HasTargetHeight returns a boolean if a field has been set.
func (o *JobOutputInfo) HasTargetHeight() bool {
	if o != nil && !IsNil(o.TargetHeight) {
		return true
	}

	return false
}

// SetTargetHeight gets a reference to the given int32 and assigns it to the TargetHeight field.
func (o *JobOutputInfo) SetTargetHeight(v int32) {
	o.TargetHeight = &v
}

// GetWidth returns the Width field value if set, zero value otherwise.
func (o *JobOutputInfo) GetWidth() int32 {
	if o == nil || IsNil(o.Width) {
		var ret int32
		return ret
	}
	return *o.Width
}

// GetWidthOk returns a tuple with the Width field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JobOutputInfo) GetWidthOk() (*int32, bool) {
	if o == nil || IsNil(o.Width) {
		return nil, false
	}
	return o.Width, true
}

// HasWidth returns a boolean if a field has been set.
func (o *JobOutputInfo) HasWidth() bool {
	if o != nil && !IsNil(o.Width) {
		return true
	}

	return false
}

// SetWidth gets a reference to the given int32 and assigns it to the Width field.
func (o *JobOutputInfo) SetWidth(v int32) {
	o.Width = &v
}

// GetHeight returns the Height field value if set, zero value otherwise.
func (o *JobOutputInfo) GetHeight() int32 {
	if o == nil || IsNil(o.Height) {
		var ret int32
		return ret
	}
	return *o.Height
}

// GetHeightOk returns a tuple with the Height field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JobOutputInfo) GetHeightOk() (*int32, bool) {
	if o == nil || IsNil(o.Height) {
		return nil, false
	}
	return o.Height, true
}

// HasHeight returns a boolean if a field has been set.
func (o *JobOutputInfo) HasHeight() bool {
	if o != nil && !IsNil(o.Height) {
		return true
	}

	return false
}

// SetHeight gets a reference to the given int32 and assigns it to the Height field.
func (o *JobOutputInfo) SetHeight(v int32) {
	o.Height = &v
}

// GetBytes returns the Bytes field value if set, zero value otherwise.
func (o *JobOutputInfo) GetBytes() int64 {
	if o == nil || IsNil(o.Bytes) {
		var ret int64
		return ret
	}
	return *o.Bytes
}

// GetBytesOk returns a tuple with the Bytes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JobOutputInfo) GetBytesOk() (*int64, bool) {
	if o == nil || IsNil(o.Bytes) {
		return nil, false
	}
	return o.Bytes, true
}

// HasBytes returns a boolean if a field has been set.
func (o *JobOutputInfo) HasBytes() bool {
	if o != nil && !IsNil(o.Bytes) {
		return true
	}

	return false
}

// SetBytes gets a reference to the given int64 and assigns it to the Bytes field.
func (o *JobOutputInfo) SetBytes(v int64) {
	o.Bytes = &v
}

// GetSha256 returns the Sha256 field value if set, zero value otherwise.
func (o *JobOutputInfo) GetSha256() string {
	if o == nil || IsNil(o.Sha256) {
		var ret string
		return ret
	}
	return *o.Sha256
}

// GetSha256Ok returns a tuple with the Sha256 field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JobOutputInfo) GetSha256Ok() (*string, bool) {
	if o == nil || IsNil(o.Sha256) {
		return nil, false
	}
	return o.Sha256, true
}

// HasSha256 returns a boolean if a field has been set.
func (o *JobOutputInfo) HasSha256() bool {
	if o != nil && !IsNil(o.Sha256) {
		return true
	}

	return false
}

// SetSha256 gets a reference to the given string and assigns it to the Sha256 field.
func (o *JobOutputInfo) SetSha256(v string) {
	o.Sha256 = &v
}

// GetRetentionExpiresAt returns the RetentionExpiresAt field value if set, zero value otherwise.
func (o *JobOutputInfo) GetRetentionExpiresAt() time.Time {
	if o == nil || IsNil(o.RetentionExpiresAt) {
		var ret time.Time
		return ret
	}
	return *o.RetentionExpiresAt
}

// GetRetentionExpiresAtOk returns a tuple with the RetentionExpiresAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JobOutputInfo) GetRetentionExpiresAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.RetentionExpiresAt) {
		return nil, false
	}
	return o.RetentionExpiresAt, true
}

// HasRetentionExpiresAt returns a boolean if a field has been set.
func (o *JobOutputInfo) HasRetentionExpiresAt() bool {
	if o != nil && !IsNil(o.RetentionExpiresAt) {
		return true
	}

	return false
}

// SetRetentionExpiresAt gets a reference to the given time.Time and assigns it to the RetentionExpiresAt field.
func (o *JobOutputInfo) SetRetentionExpiresAt(v time.Time) {
	o.RetentionExpiresAt = &v
}

func (o JobOutputInfo) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o JobOutputInfo) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["state"] = o.State
	toSerialize["container"] = o.Container
	toSerialize["audio"] = o.Audio
	if !IsNil(o.TargetHeight) {
		toSerialize["target_height"] = o.TargetHeight
	}
	if !IsNil(o.Width) {
		toSerialize["width"] = o.Width
	}
	if !IsNil(o.Height) {
		toSerialize["height"] = o.Height
	}
	if !IsNil(o.Bytes) {
		toSerialize["bytes"] = o.Bytes
	}
	if !IsNil(o.Sha256) {
		toSerialize["sha256"] = o.Sha256
	}
	if !IsNil(o.RetentionExpiresAt) {
		toSerialize["retention_expires_at"] = o.RetentionExpiresAt
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *JobOutputInfo) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"state",
		"container",
		"audio",
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

	varJobOutputInfo := _JobOutputInfo{}

	err = json.Unmarshal(data, &varJobOutputInfo)

	if err != nil {
		return err
	}

	*o = JobOutputInfo(varJobOutputInfo)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "state")
		delete(additionalProperties, "container")
		delete(additionalProperties, "audio")
		delete(additionalProperties, "target_height")
		delete(additionalProperties, "width")
		delete(additionalProperties, "height")
		delete(additionalProperties, "bytes")
		delete(additionalProperties, "sha256")
		delete(additionalProperties, "retention_expires_at")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableJobOutputInfo struct {
	value *JobOutputInfo
	isSet bool
}

func (v NullableJobOutputInfo) Get() *JobOutputInfo {
	return v.value
}

func (v *NullableJobOutputInfo) Set(val *JobOutputInfo) {
	v.value = val
	v.isSet = true
}

func (v NullableJobOutputInfo) IsSet() bool {
	return v.isSet
}

func (v *NullableJobOutputInfo) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableJobOutputInfo(val *JobOutputInfo) *NullableJobOutputInfo {
	return &NullableJobOutputInfo{value: val, isSet: true}
}

func (v NullableJobOutputInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableJobOutputInfo) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


