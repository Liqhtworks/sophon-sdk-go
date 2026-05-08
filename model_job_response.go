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

// checks if the JobResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &JobResponse{}

// JobResponse struct for JobResponse
type JobResponse struct {
	Id string `json:"id"`
	Status JobStatus `json:"status"`
	StatusReason *string `json:"status_reason,omitempty"`
	Attempt int32 `json:"attempt"`
	// Whether the job can still be retried (attempt < max_attempts and not terminal).
	Retryable bool `json:"retryable"`
	// Public profile ID submitted by the customer. For adaptive jobs this stays `sophon-auto`; see `effective_profile_id` for the worker's resolved concrete profile. 
	Profile JobProfile `json:"profile"`
	// Concrete profile resolved by the worker. Omitted until dispatch resolves. On explicit-profile jobs this equals `profile`. On `sophon-auto` jobs this is a variant identifier recording which path the API routed the source through; exact encoder settings for a given variant may be updated between releases as the adaptive logic is tuned. 
	EffectiveProfileId *string `json:"effective_profile_id,omitempty"`
	Source JobSourceInfo `json:"source"`
	Progress JobProgress `json:"progress"`
	Output JobOutputInfo `json:"output"`
	// Arbitrary JSON object attached to a job. Keys and values are passed through unchanged to webhook deliveries and echoed on job reads. The serialized representation must not exceed 16 KiB. Free-form; SDKs surface this as a `Record<string, unknown>` / `dict[str, Any]` / `map[string]interface{}` depending on language. 
	Metadata map[string]interface{} `json:"metadata"`
	CreatedAt time.Time `json:"created_at"`
	StartedAt *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Error *string `json:"error,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _JobResponse JobResponse

// NewJobResponse instantiates a new JobResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewJobResponse(id string, status JobStatus, attempt int32, retryable bool, profile JobProfile, source JobSourceInfo, progress JobProgress, output JobOutputInfo, metadata map[string]interface{}, createdAt time.Time) *JobResponse {
	this := JobResponse{}
	this.Id = id
	this.Status = status
	this.Attempt = attempt
	this.Retryable = retryable
	this.Profile = profile
	this.Source = source
	this.Progress = progress
	this.Output = output
	this.Metadata = metadata
	this.CreatedAt = createdAt
	return &this
}

// NewJobResponseWithDefaults instantiates a new JobResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewJobResponseWithDefaults() *JobResponse {
	this := JobResponse{}
	return &this
}

// GetId returns the Id field value
func (o *JobResponse) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *JobResponse) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *JobResponse) SetId(v string) {
	o.Id = v
}

// GetStatus returns the Status field value
func (o *JobResponse) GetStatus() JobStatus {
	if o == nil {
		var ret JobStatus
		return ret
	}

	return o.Status
}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
func (o *JobResponse) GetStatusOk() (*JobStatus, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Status, true
}

// SetStatus sets field value
func (o *JobResponse) SetStatus(v JobStatus) {
	o.Status = v
}

// GetStatusReason returns the StatusReason field value if set, zero value otherwise.
func (o *JobResponse) GetStatusReason() string {
	if o == nil || IsNil(o.StatusReason) {
		var ret string
		return ret
	}
	return *o.StatusReason
}

// GetStatusReasonOk returns a tuple with the StatusReason field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JobResponse) GetStatusReasonOk() (*string, bool) {
	if o == nil || IsNil(o.StatusReason) {
		return nil, false
	}
	return o.StatusReason, true
}

// HasStatusReason returns a boolean if a field has been set.
func (o *JobResponse) HasStatusReason() bool {
	if o != nil && !IsNil(o.StatusReason) {
		return true
	}

	return false
}

// SetStatusReason gets a reference to the given string and assigns it to the StatusReason field.
func (o *JobResponse) SetStatusReason(v string) {
	o.StatusReason = &v
}

// GetAttempt returns the Attempt field value
func (o *JobResponse) GetAttempt() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Attempt
}

// GetAttemptOk returns a tuple with the Attempt field value
// and a boolean to check if the value has been set.
func (o *JobResponse) GetAttemptOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Attempt, true
}

// SetAttempt sets field value
func (o *JobResponse) SetAttempt(v int32) {
	o.Attempt = v
}

// GetRetryable returns the Retryable field value
func (o *JobResponse) GetRetryable() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.Retryable
}

// GetRetryableOk returns a tuple with the Retryable field value
// and a boolean to check if the value has been set.
func (o *JobResponse) GetRetryableOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Retryable, true
}

// SetRetryable sets field value
func (o *JobResponse) SetRetryable(v bool) {
	o.Retryable = v
}

// GetProfile returns the Profile field value
func (o *JobResponse) GetProfile() JobProfile {
	if o == nil {
		var ret JobProfile
		return ret
	}

	return o.Profile
}

// GetProfileOk returns a tuple with the Profile field value
// and a boolean to check if the value has been set.
func (o *JobResponse) GetProfileOk() (*JobProfile, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Profile, true
}

// SetProfile sets field value
func (o *JobResponse) SetProfile(v JobProfile) {
	o.Profile = v
}

// GetEffectiveProfileId returns the EffectiveProfileId field value if set, zero value otherwise.
func (o *JobResponse) GetEffectiveProfileId() string {
	if o == nil || IsNil(o.EffectiveProfileId) {
		var ret string
		return ret
	}
	return *o.EffectiveProfileId
}

// GetEffectiveProfileIdOk returns a tuple with the EffectiveProfileId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JobResponse) GetEffectiveProfileIdOk() (*string, bool) {
	if o == nil || IsNil(o.EffectiveProfileId) {
		return nil, false
	}
	return o.EffectiveProfileId, true
}

// HasEffectiveProfileId returns a boolean if a field has been set.
func (o *JobResponse) HasEffectiveProfileId() bool {
	if o != nil && !IsNil(o.EffectiveProfileId) {
		return true
	}

	return false
}

// SetEffectiveProfileId gets a reference to the given string and assigns it to the EffectiveProfileId field.
func (o *JobResponse) SetEffectiveProfileId(v string) {
	o.EffectiveProfileId = &v
}

// GetSource returns the Source field value
func (o *JobResponse) GetSource() JobSourceInfo {
	if o == nil {
		var ret JobSourceInfo
		return ret
	}

	return o.Source
}

// GetSourceOk returns a tuple with the Source field value
// and a boolean to check if the value has been set.
func (o *JobResponse) GetSourceOk() (*JobSourceInfo, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Source, true
}

// SetSource sets field value
func (o *JobResponse) SetSource(v JobSourceInfo) {
	o.Source = v
}

// GetProgress returns the Progress field value
func (o *JobResponse) GetProgress() JobProgress {
	if o == nil {
		var ret JobProgress
		return ret
	}

	return o.Progress
}

// GetProgressOk returns a tuple with the Progress field value
// and a boolean to check if the value has been set.
func (o *JobResponse) GetProgressOk() (*JobProgress, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Progress, true
}

// SetProgress sets field value
func (o *JobResponse) SetProgress(v JobProgress) {
	o.Progress = v
}

// GetOutput returns the Output field value
func (o *JobResponse) GetOutput() JobOutputInfo {
	if o == nil {
		var ret JobOutputInfo
		return ret
	}

	return o.Output
}

// GetOutputOk returns a tuple with the Output field value
// and a boolean to check if the value has been set.
func (o *JobResponse) GetOutputOk() (*JobOutputInfo, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Output, true
}

// SetOutput sets field value
func (o *JobResponse) SetOutput(v JobOutputInfo) {
	o.Output = v
}

// GetMetadata returns the Metadata field value
func (o *JobResponse) GetMetadata() map[string]interface{} {
	if o == nil {
		var ret map[string]interface{}
		return ret
	}

	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
func (o *JobResponse) GetMetadataOk() (map[string]interface{}, bool) {
	if o == nil {
		return map[string]interface{}{}, false
	}
	return o.Metadata, true
}

// SetMetadata sets field value
func (o *JobResponse) SetMetadata(v map[string]interface{}) {
	o.Metadata = v
}

// GetCreatedAt returns the CreatedAt field value
func (o *JobResponse) GetCreatedAt() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value
// and a boolean to check if the value has been set.
func (o *JobResponse) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CreatedAt, true
}

// SetCreatedAt sets field value
func (o *JobResponse) SetCreatedAt(v time.Time) {
	o.CreatedAt = v
}

// GetStartedAt returns the StartedAt field value if set, zero value otherwise.
func (o *JobResponse) GetStartedAt() time.Time {
	if o == nil || IsNil(o.StartedAt) {
		var ret time.Time
		return ret
	}
	return *o.StartedAt
}

// GetStartedAtOk returns a tuple with the StartedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JobResponse) GetStartedAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.StartedAt) {
		return nil, false
	}
	return o.StartedAt, true
}

// HasStartedAt returns a boolean if a field has been set.
func (o *JobResponse) HasStartedAt() bool {
	if o != nil && !IsNil(o.StartedAt) {
		return true
	}

	return false
}

// SetStartedAt gets a reference to the given time.Time and assigns it to the StartedAt field.
func (o *JobResponse) SetStartedAt(v time.Time) {
	o.StartedAt = &v
}

// GetCompletedAt returns the CompletedAt field value if set, zero value otherwise.
func (o *JobResponse) GetCompletedAt() time.Time {
	if o == nil || IsNil(o.CompletedAt) {
		var ret time.Time
		return ret
	}
	return *o.CompletedAt
}

// GetCompletedAtOk returns a tuple with the CompletedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JobResponse) GetCompletedAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CompletedAt) {
		return nil, false
	}
	return o.CompletedAt, true
}

// HasCompletedAt returns a boolean if a field has been set.
func (o *JobResponse) HasCompletedAt() bool {
	if o != nil && !IsNil(o.CompletedAt) {
		return true
	}

	return false
}

// SetCompletedAt gets a reference to the given time.Time and assigns it to the CompletedAt field.
func (o *JobResponse) SetCompletedAt(v time.Time) {
	o.CompletedAt = &v
}

// GetError returns the Error field value if set, zero value otherwise.
func (o *JobResponse) GetError() string {
	if o == nil || IsNil(o.Error) {
		var ret string
		return ret
	}
	return *o.Error
}

// GetErrorOk returns a tuple with the Error field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JobResponse) GetErrorOk() (*string, bool) {
	if o == nil || IsNil(o.Error) {
		return nil, false
	}
	return o.Error, true
}

// HasError returns a boolean if a field has been set.
func (o *JobResponse) HasError() bool {
	if o != nil && !IsNil(o.Error) {
		return true
	}

	return false
}

// SetError gets a reference to the given string and assigns it to the Error field.
func (o *JobResponse) SetError(v string) {
	o.Error = &v
}

func (o JobResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o JobResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["status"] = o.Status
	if !IsNil(o.StatusReason) {
		toSerialize["status_reason"] = o.StatusReason
	}
	toSerialize["attempt"] = o.Attempt
	toSerialize["retryable"] = o.Retryable
	toSerialize["profile"] = o.Profile
	if !IsNil(o.EffectiveProfileId) {
		toSerialize["effective_profile_id"] = o.EffectiveProfileId
	}
	toSerialize["source"] = o.Source
	toSerialize["progress"] = o.Progress
	toSerialize["output"] = o.Output
	toSerialize["metadata"] = o.Metadata
	toSerialize["created_at"] = o.CreatedAt
	if !IsNil(o.StartedAt) {
		toSerialize["started_at"] = o.StartedAt
	}
	if !IsNil(o.CompletedAt) {
		toSerialize["completed_at"] = o.CompletedAt
	}
	if !IsNil(o.Error) {
		toSerialize["error"] = o.Error
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *JobResponse) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"id",
		"status",
		"attempt",
		"retryable",
		"profile",
		"source",
		"progress",
		"output",
		"metadata",
		"created_at",
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

	varJobResponse := _JobResponse{}

	err = json.Unmarshal(data, &varJobResponse)

	if err != nil {
		return err
	}

	*o = JobResponse(varJobResponse)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "id")
		delete(additionalProperties, "status")
		delete(additionalProperties, "status_reason")
		delete(additionalProperties, "attempt")
		delete(additionalProperties, "retryable")
		delete(additionalProperties, "profile")
		delete(additionalProperties, "effective_profile_id")
		delete(additionalProperties, "source")
		delete(additionalProperties, "progress")
		delete(additionalProperties, "output")
		delete(additionalProperties, "metadata")
		delete(additionalProperties, "created_at")
		delete(additionalProperties, "started_at")
		delete(additionalProperties, "completed_at")
		delete(additionalProperties, "error")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableJobResponse struct {
	value *JobResponse
	isSet bool
}

func (v NullableJobResponse) Get() *JobResponse {
	return v.value
}

func (v *NullableJobResponse) Set(val *JobResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableJobResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableJobResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableJobResponse(val *JobResponse) *NullableJobResponse {
	return &NullableJobResponse{value: val, isSet: true}
}

func (v NullableJobResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableJobResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


