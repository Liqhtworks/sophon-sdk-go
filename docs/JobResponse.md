# JobResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**Status** | [**JobStatus**](JobStatus.md) |  | 
**StatusReason** | Pointer to **string** |  | [optional] 
**Attempt** | **int32** |  | 
**Retryable** | **bool** | Whether the job can still be retried (attempt &lt; max_attempts and not terminal). | 
**Profile** | [**JobProfile**](JobProfile.md) | Public profile ID submitted by the customer. For adaptive jobs this stays &#x60;sophon-auto&#x60;; see &#x60;effective_profile_id&#x60; for the worker&#39;s resolved concrete profile.  | 
**EffectiveProfileId** | Pointer to **string** | Concrete profile resolved by the worker. Omitted until dispatch resolves. On explicit-profile jobs this equals &#x60;profile&#x60;. On &#x60;sophon-auto&#x60; jobs this is a variant identifier recording which path the API routed the source through; exact encoder settings for a given variant may be updated between releases as the adaptive logic is tuned.  | [optional] 
**Source** | [**JobSourceInfo**](JobSourceInfo.md) |  | 
**Progress** | [**JobProgress**](JobProgress.md) |  | 
**Output** | [**JobOutputInfo**](JobOutputInfo.md) |  | 
**Metadata** | **map[string]interface{}** | Arbitrary JSON object attached to a job. Keys and values are passed through unchanged to webhook deliveries and echoed on job reads. The serialized representation must not exceed 16 KiB. Free-form; SDKs surface this as a &#x60;Record&lt;string, unknown&gt;&#x60; / &#x60;dict[str, Any]&#x60; / &#x60;map[string]interface{}&#x60; depending on language.  | 
**CreatedAt** | **time.Time** |  | 
**StartedAt** | Pointer to **time.Time** |  | [optional] 
**CompletedAt** | Pointer to **time.Time** |  | [optional] 
**Error** | Pointer to **string** |  | [optional] 

## Methods

### NewJobResponse

`func NewJobResponse(id string, status JobStatus, attempt int32, retryable bool, profile JobProfile, source JobSourceInfo, progress JobProgress, output JobOutputInfo, metadata map[string]interface{}, createdAt time.Time, ) *JobResponse`

NewJobResponse instantiates a new JobResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewJobResponseWithDefaults

`func NewJobResponseWithDefaults() *JobResponse`

NewJobResponseWithDefaults instantiates a new JobResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *JobResponse) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *JobResponse) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *JobResponse) SetId(v string)`

SetId sets Id field to given value.


### GetStatus

`func (o *JobResponse) GetStatus() JobStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *JobResponse) GetStatusOk() (*JobStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *JobResponse) SetStatus(v JobStatus)`

SetStatus sets Status field to given value.


### GetStatusReason

`func (o *JobResponse) GetStatusReason() string`

GetStatusReason returns the StatusReason field if non-nil, zero value otherwise.

### GetStatusReasonOk

`func (o *JobResponse) GetStatusReasonOk() (*string, bool)`

GetStatusReasonOk returns a tuple with the StatusReason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusReason

`func (o *JobResponse) SetStatusReason(v string)`

SetStatusReason sets StatusReason field to given value.

### HasStatusReason

`func (o *JobResponse) HasStatusReason() bool`

HasStatusReason returns a boolean if a field has been set.

### GetAttempt

`func (o *JobResponse) GetAttempt() int32`

GetAttempt returns the Attempt field if non-nil, zero value otherwise.

### GetAttemptOk

`func (o *JobResponse) GetAttemptOk() (*int32, bool)`

GetAttemptOk returns a tuple with the Attempt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttempt

`func (o *JobResponse) SetAttempt(v int32)`

SetAttempt sets Attempt field to given value.


### GetRetryable

`func (o *JobResponse) GetRetryable() bool`

GetRetryable returns the Retryable field if non-nil, zero value otherwise.

### GetRetryableOk

`func (o *JobResponse) GetRetryableOk() (*bool, bool)`

GetRetryableOk returns a tuple with the Retryable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRetryable

`func (o *JobResponse) SetRetryable(v bool)`

SetRetryable sets Retryable field to given value.


### GetProfile

`func (o *JobResponse) GetProfile() JobProfile`

GetProfile returns the Profile field if non-nil, zero value otherwise.

### GetProfileOk

`func (o *JobResponse) GetProfileOk() (*JobProfile, bool)`

GetProfileOk returns a tuple with the Profile field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProfile

`func (o *JobResponse) SetProfile(v JobProfile)`

SetProfile sets Profile field to given value.


### GetEffectiveProfileId

`func (o *JobResponse) GetEffectiveProfileId() string`

GetEffectiveProfileId returns the EffectiveProfileId field if non-nil, zero value otherwise.

### GetEffectiveProfileIdOk

`func (o *JobResponse) GetEffectiveProfileIdOk() (*string, bool)`

GetEffectiveProfileIdOk returns a tuple with the EffectiveProfileId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEffectiveProfileId

`func (o *JobResponse) SetEffectiveProfileId(v string)`

SetEffectiveProfileId sets EffectiveProfileId field to given value.

### HasEffectiveProfileId

`func (o *JobResponse) HasEffectiveProfileId() bool`

HasEffectiveProfileId returns a boolean if a field has been set.

### GetSource

`func (o *JobResponse) GetSource() JobSourceInfo`

GetSource returns the Source field if non-nil, zero value otherwise.

### GetSourceOk

`func (o *JobResponse) GetSourceOk() (*JobSourceInfo, bool)`

GetSourceOk returns a tuple with the Source field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSource

`func (o *JobResponse) SetSource(v JobSourceInfo)`

SetSource sets Source field to given value.


### GetProgress

`func (o *JobResponse) GetProgress() JobProgress`

GetProgress returns the Progress field if non-nil, zero value otherwise.

### GetProgressOk

`func (o *JobResponse) GetProgressOk() (*JobProgress, bool)`

GetProgressOk returns a tuple with the Progress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProgress

`func (o *JobResponse) SetProgress(v JobProgress)`

SetProgress sets Progress field to given value.


### GetOutput

`func (o *JobResponse) GetOutput() JobOutputInfo`

GetOutput returns the Output field if non-nil, zero value otherwise.

### GetOutputOk

`func (o *JobResponse) GetOutputOk() (*JobOutputInfo, bool)`

GetOutputOk returns a tuple with the Output field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutput

`func (o *JobResponse) SetOutput(v JobOutputInfo)`

SetOutput sets Output field to given value.


### GetMetadata

`func (o *JobResponse) GetMetadata() map[string]interface{}`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *JobResponse) GetMetadataOk() (*map[string]interface{}, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *JobResponse) SetMetadata(v map[string]interface{})`

SetMetadata sets Metadata field to given value.


### GetCreatedAt

`func (o *JobResponse) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *JobResponse) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *JobResponse) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetStartedAt

`func (o *JobResponse) GetStartedAt() time.Time`

GetStartedAt returns the StartedAt field if non-nil, zero value otherwise.

### GetStartedAtOk

`func (o *JobResponse) GetStartedAtOk() (*time.Time, bool)`

GetStartedAtOk returns a tuple with the StartedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartedAt

`func (o *JobResponse) SetStartedAt(v time.Time)`

SetStartedAt sets StartedAt field to given value.

### HasStartedAt

`func (o *JobResponse) HasStartedAt() bool`

HasStartedAt returns a boolean if a field has been set.

### GetCompletedAt

`func (o *JobResponse) GetCompletedAt() time.Time`

GetCompletedAt returns the CompletedAt field if non-nil, zero value otherwise.

### GetCompletedAtOk

`func (o *JobResponse) GetCompletedAtOk() (*time.Time, bool)`

GetCompletedAtOk returns a tuple with the CompletedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompletedAt

`func (o *JobResponse) SetCompletedAt(v time.Time)`

SetCompletedAt sets CompletedAt field to given value.

### HasCompletedAt

`func (o *JobResponse) HasCompletedAt() bool`

HasCompletedAt returns a boolean if a field has been set.

### GetError

`func (o *JobResponse) GetError() string`

GetError returns the Error field if non-nil, zero value otherwise.

### GetErrorOk

`func (o *JobResponse) GetErrorOk() (*string, bool)`

GetErrorOk returns a tuple with the Error field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetError

`func (o *JobResponse) SetError(v string)`

SetError sets Error field to given value.

### HasError

`func (o *JobResponse) HasError() bool`

HasError returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


