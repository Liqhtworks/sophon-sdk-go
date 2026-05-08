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

// checks if the UploadStatusResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &UploadStatusResponse{}

// UploadStatusResponse struct for UploadStatusResponse
type UploadStatusResponse struct {
	Id string `json:"id"`
	Status string `json:"status"`
	FileName string `json:"file_name"`
	TotalChunks int32 `json:"total_chunks"`
	// Array of 0-indexed part numbers that have been received.
	ReceivedChunks []int32 `json:"received_chunks"`
	ExpiresAt time.Time `json:"expires_at"`
	// Source media width in pixels, populated from ffprobe after upload assembly. Null for uploads in `initiated`/`uploading` state or when probe failed. 
	SourceWidth *int32 `json:"source_width,omitempty"`
	// Source media height in pixels. See `source_width`.
	SourceHeight *int32 `json:"source_height,omitempty"`
	// Source media duration in seconds, from ffprobe after upload assembly. Used by the webapp free-tier budget check to compute realistic billable_seconds (5-second ceiling rounding). 
	SourceDurationSeconds *float32 `json:"source_duration_seconds,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _UploadStatusResponse UploadStatusResponse

// NewUploadStatusResponse instantiates a new UploadStatusResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUploadStatusResponse(id string, status string, fileName string, totalChunks int32, receivedChunks []int32, expiresAt time.Time) *UploadStatusResponse {
	this := UploadStatusResponse{}
	this.Id = id
	this.Status = status
	this.FileName = fileName
	this.TotalChunks = totalChunks
	this.ReceivedChunks = receivedChunks
	this.ExpiresAt = expiresAt
	return &this
}

// NewUploadStatusResponseWithDefaults instantiates a new UploadStatusResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUploadStatusResponseWithDefaults() *UploadStatusResponse {
	this := UploadStatusResponse{}
	return &this
}

// GetId returns the Id field value
func (o *UploadStatusResponse) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *UploadStatusResponse) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *UploadStatusResponse) SetId(v string) {
	o.Id = v
}

// GetStatus returns the Status field value
func (o *UploadStatusResponse) GetStatus() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Status
}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
func (o *UploadStatusResponse) GetStatusOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Status, true
}

// SetStatus sets field value
func (o *UploadStatusResponse) SetStatus(v string) {
	o.Status = v
}

// GetFileName returns the FileName field value
func (o *UploadStatusResponse) GetFileName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.FileName
}

// GetFileNameOk returns a tuple with the FileName field value
// and a boolean to check if the value has been set.
func (o *UploadStatusResponse) GetFileNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.FileName, true
}

// SetFileName sets field value
func (o *UploadStatusResponse) SetFileName(v string) {
	o.FileName = v
}

// GetTotalChunks returns the TotalChunks field value
func (o *UploadStatusResponse) GetTotalChunks() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.TotalChunks
}

// GetTotalChunksOk returns a tuple with the TotalChunks field value
// and a boolean to check if the value has been set.
func (o *UploadStatusResponse) GetTotalChunksOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TotalChunks, true
}

// SetTotalChunks sets field value
func (o *UploadStatusResponse) SetTotalChunks(v int32) {
	o.TotalChunks = v
}

// GetReceivedChunks returns the ReceivedChunks field value
func (o *UploadStatusResponse) GetReceivedChunks() []int32 {
	if o == nil {
		var ret []int32
		return ret
	}

	return o.ReceivedChunks
}

// GetReceivedChunksOk returns a tuple with the ReceivedChunks field value
// and a boolean to check if the value has been set.
func (o *UploadStatusResponse) GetReceivedChunksOk() ([]int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.ReceivedChunks, true
}

// SetReceivedChunks sets field value
func (o *UploadStatusResponse) SetReceivedChunks(v []int32) {
	o.ReceivedChunks = v
}

// GetExpiresAt returns the ExpiresAt field value
func (o *UploadStatusResponse) GetExpiresAt() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.ExpiresAt
}

// GetExpiresAtOk returns a tuple with the ExpiresAt field value
// and a boolean to check if the value has been set.
func (o *UploadStatusResponse) GetExpiresAtOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ExpiresAt, true
}

// SetExpiresAt sets field value
func (o *UploadStatusResponse) SetExpiresAt(v time.Time) {
	o.ExpiresAt = v
}

// GetSourceWidth returns the SourceWidth field value if set, zero value otherwise.
func (o *UploadStatusResponse) GetSourceWidth() int32 {
	if o == nil || IsNil(o.SourceWidth) {
		var ret int32
		return ret
	}
	return *o.SourceWidth
}

// GetSourceWidthOk returns a tuple with the SourceWidth field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UploadStatusResponse) GetSourceWidthOk() (*int32, bool) {
	if o == nil || IsNil(o.SourceWidth) {
		return nil, false
	}
	return o.SourceWidth, true
}

// HasSourceWidth returns a boolean if a field has been set.
func (o *UploadStatusResponse) HasSourceWidth() bool {
	if o != nil && !IsNil(o.SourceWidth) {
		return true
	}

	return false
}

// SetSourceWidth gets a reference to the given int32 and assigns it to the SourceWidth field.
func (o *UploadStatusResponse) SetSourceWidth(v int32) {
	o.SourceWidth = &v
}

// GetSourceHeight returns the SourceHeight field value if set, zero value otherwise.
func (o *UploadStatusResponse) GetSourceHeight() int32 {
	if o == nil || IsNil(o.SourceHeight) {
		var ret int32
		return ret
	}
	return *o.SourceHeight
}

// GetSourceHeightOk returns a tuple with the SourceHeight field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UploadStatusResponse) GetSourceHeightOk() (*int32, bool) {
	if o == nil || IsNil(o.SourceHeight) {
		return nil, false
	}
	return o.SourceHeight, true
}

// HasSourceHeight returns a boolean if a field has been set.
func (o *UploadStatusResponse) HasSourceHeight() bool {
	if o != nil && !IsNil(o.SourceHeight) {
		return true
	}

	return false
}

// SetSourceHeight gets a reference to the given int32 and assigns it to the SourceHeight field.
func (o *UploadStatusResponse) SetSourceHeight(v int32) {
	o.SourceHeight = &v
}

// GetSourceDurationSeconds returns the SourceDurationSeconds field value if set, zero value otherwise.
func (o *UploadStatusResponse) GetSourceDurationSeconds() float32 {
	if o == nil || IsNil(o.SourceDurationSeconds) {
		var ret float32
		return ret
	}
	return *o.SourceDurationSeconds
}

// GetSourceDurationSecondsOk returns a tuple with the SourceDurationSeconds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UploadStatusResponse) GetSourceDurationSecondsOk() (*float32, bool) {
	if o == nil || IsNil(o.SourceDurationSeconds) {
		return nil, false
	}
	return o.SourceDurationSeconds, true
}

// HasSourceDurationSeconds returns a boolean if a field has been set.
func (o *UploadStatusResponse) HasSourceDurationSeconds() bool {
	if o != nil && !IsNil(o.SourceDurationSeconds) {
		return true
	}

	return false
}

// SetSourceDurationSeconds gets a reference to the given float32 and assigns it to the SourceDurationSeconds field.
func (o *UploadStatusResponse) SetSourceDurationSeconds(v float32) {
	o.SourceDurationSeconds = &v
}

func (o UploadStatusResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o UploadStatusResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["status"] = o.Status
	toSerialize["file_name"] = o.FileName
	toSerialize["total_chunks"] = o.TotalChunks
	toSerialize["received_chunks"] = o.ReceivedChunks
	toSerialize["expires_at"] = o.ExpiresAt
	if !IsNil(o.SourceWidth) {
		toSerialize["source_width"] = o.SourceWidth
	}
	if !IsNil(o.SourceHeight) {
		toSerialize["source_height"] = o.SourceHeight
	}
	if !IsNil(o.SourceDurationSeconds) {
		toSerialize["source_duration_seconds"] = o.SourceDurationSeconds
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *UploadStatusResponse) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"id",
		"status",
		"file_name",
		"total_chunks",
		"received_chunks",
		"expires_at",
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

	varUploadStatusResponse := _UploadStatusResponse{}

	err = json.Unmarshal(data, &varUploadStatusResponse)

	if err != nil {
		return err
	}

	*o = UploadStatusResponse(varUploadStatusResponse)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "id")
		delete(additionalProperties, "status")
		delete(additionalProperties, "file_name")
		delete(additionalProperties, "total_chunks")
		delete(additionalProperties, "received_chunks")
		delete(additionalProperties, "expires_at")
		delete(additionalProperties, "source_width")
		delete(additionalProperties, "source_height")
		delete(additionalProperties, "source_duration_seconds")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableUploadStatusResponse struct {
	value *UploadStatusResponse
	isSet bool
}

func (v NullableUploadStatusResponse) Get() *UploadStatusResponse {
	return v.value
}

func (v *NullableUploadStatusResponse) Set(val *UploadStatusResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableUploadStatusResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableUploadStatusResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUploadStatusResponse(val *UploadStatusResponse) *NullableUploadStatusResponse {
	return &NullableUploadStatusResponse{value: val, isSet: true}
}

func (v NullableUploadStatusResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUploadStatusResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


