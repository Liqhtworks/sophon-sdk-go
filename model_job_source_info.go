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

// checks if the JobSourceInfo type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &JobSourceInfo{}

// JobSourceInfo struct for JobSourceInfo
type JobSourceInfo struct {
	// Original file name of the source.
	Name *string `json:"name,omitempty"`
	Bytes *int64 `json:"bytes,omitempty"`
	// SHA-256 hex digest of the source file.
	Sha256 string `json:"sha256"`
	DurationSeconds *float64 `json:"duration_seconds,omitempty"`
	Resolution *string `json:"resolution,omitempty"`
	FrameRate *string `json:"frame_rate,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _JobSourceInfo JobSourceInfo

// NewJobSourceInfo instantiates a new JobSourceInfo object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewJobSourceInfo(sha256 string) *JobSourceInfo {
	this := JobSourceInfo{}
	this.Sha256 = sha256
	return &this
}

// NewJobSourceInfoWithDefaults instantiates a new JobSourceInfo object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewJobSourceInfoWithDefaults() *JobSourceInfo {
	this := JobSourceInfo{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *JobSourceInfo) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JobSourceInfo) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *JobSourceInfo) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *JobSourceInfo) SetName(v string) {
	o.Name = &v
}

// GetBytes returns the Bytes field value if set, zero value otherwise.
func (o *JobSourceInfo) GetBytes() int64 {
	if o == nil || IsNil(o.Bytes) {
		var ret int64
		return ret
	}
	return *o.Bytes
}

// GetBytesOk returns a tuple with the Bytes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JobSourceInfo) GetBytesOk() (*int64, bool) {
	if o == nil || IsNil(o.Bytes) {
		return nil, false
	}
	return o.Bytes, true
}

// HasBytes returns a boolean if a field has been set.
func (o *JobSourceInfo) HasBytes() bool {
	if o != nil && !IsNil(o.Bytes) {
		return true
	}

	return false
}

// SetBytes gets a reference to the given int64 and assigns it to the Bytes field.
func (o *JobSourceInfo) SetBytes(v int64) {
	o.Bytes = &v
}

// GetSha256 returns the Sha256 field value
func (o *JobSourceInfo) GetSha256() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Sha256
}

// GetSha256Ok returns a tuple with the Sha256 field value
// and a boolean to check if the value has been set.
func (o *JobSourceInfo) GetSha256Ok() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Sha256, true
}

// SetSha256 sets field value
func (o *JobSourceInfo) SetSha256(v string) {
	o.Sha256 = v
}

// GetDurationSeconds returns the DurationSeconds field value if set, zero value otherwise.
func (o *JobSourceInfo) GetDurationSeconds() float64 {
	if o == nil || IsNil(o.DurationSeconds) {
		var ret float64
		return ret
	}
	return *o.DurationSeconds
}

// GetDurationSecondsOk returns a tuple with the DurationSeconds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JobSourceInfo) GetDurationSecondsOk() (*float64, bool) {
	if o == nil || IsNil(o.DurationSeconds) {
		return nil, false
	}
	return o.DurationSeconds, true
}

// HasDurationSeconds returns a boolean if a field has been set.
func (o *JobSourceInfo) HasDurationSeconds() bool {
	if o != nil && !IsNil(o.DurationSeconds) {
		return true
	}

	return false
}

// SetDurationSeconds gets a reference to the given float64 and assigns it to the DurationSeconds field.
func (o *JobSourceInfo) SetDurationSeconds(v float64) {
	o.DurationSeconds = &v
}

// GetResolution returns the Resolution field value if set, zero value otherwise.
func (o *JobSourceInfo) GetResolution() string {
	if o == nil || IsNil(o.Resolution) {
		var ret string
		return ret
	}
	return *o.Resolution
}

// GetResolutionOk returns a tuple with the Resolution field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JobSourceInfo) GetResolutionOk() (*string, bool) {
	if o == nil || IsNil(o.Resolution) {
		return nil, false
	}
	return o.Resolution, true
}

// HasResolution returns a boolean if a field has been set.
func (o *JobSourceInfo) HasResolution() bool {
	if o != nil && !IsNil(o.Resolution) {
		return true
	}

	return false
}

// SetResolution gets a reference to the given string and assigns it to the Resolution field.
func (o *JobSourceInfo) SetResolution(v string) {
	o.Resolution = &v
}

// GetFrameRate returns the FrameRate field value if set, zero value otherwise.
func (o *JobSourceInfo) GetFrameRate() string {
	if o == nil || IsNil(o.FrameRate) {
		var ret string
		return ret
	}
	return *o.FrameRate
}

// GetFrameRateOk returns a tuple with the FrameRate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JobSourceInfo) GetFrameRateOk() (*string, bool) {
	if o == nil || IsNil(o.FrameRate) {
		return nil, false
	}
	return o.FrameRate, true
}

// HasFrameRate returns a boolean if a field has been set.
func (o *JobSourceInfo) HasFrameRate() bool {
	if o != nil && !IsNil(o.FrameRate) {
		return true
	}

	return false
}

// SetFrameRate gets a reference to the given string and assigns it to the FrameRate field.
func (o *JobSourceInfo) SetFrameRate(v string) {
	o.FrameRate = &v
}

func (o JobSourceInfo) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o JobSourceInfo) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Bytes) {
		toSerialize["bytes"] = o.Bytes
	}
	toSerialize["sha256"] = o.Sha256
	if !IsNil(o.DurationSeconds) {
		toSerialize["duration_seconds"] = o.DurationSeconds
	}
	if !IsNil(o.Resolution) {
		toSerialize["resolution"] = o.Resolution
	}
	if !IsNil(o.FrameRate) {
		toSerialize["frame_rate"] = o.FrameRate
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *JobSourceInfo) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"sha256",
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

	varJobSourceInfo := _JobSourceInfo{}

	err = json.Unmarshal(data, &varJobSourceInfo)

	if err != nil {
		return err
	}

	*o = JobSourceInfo(varJobSourceInfo)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "bytes")
		delete(additionalProperties, "sha256")
		delete(additionalProperties, "duration_seconds")
		delete(additionalProperties, "resolution")
		delete(additionalProperties, "frame_rate")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableJobSourceInfo struct {
	value *JobSourceInfo
	isSet bool
}

func (v NullableJobSourceInfo) Get() *JobSourceInfo {
	return v.value
}

func (v *NullableJobSourceInfo) Set(val *JobSourceInfo) {
	v.value = val
	v.isSet = true
}

func (v NullableJobSourceInfo) IsSet() bool {
	return v.isSet
}

func (v *NullableJobSourceInfo) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableJobSourceInfo(val *JobSourceInfo) *NullableJobSourceInfo {
	return &NullableJobSourceInfo{value: val, isSet: true}
}

func (v NullableJobSourceInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableJobSourceInfo) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


