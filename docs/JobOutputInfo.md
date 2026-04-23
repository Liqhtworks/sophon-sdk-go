# JobOutputInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**State** | **string** |  | 
**Container** | **string** | Output container format (\&quot;mp4\&quot; or \&quot;mkv\&quot;). | 
**Audio** | **bool** | Whether the output file actually contains audio. Reflects the muxed result, not the request flag — a video-only source with audio requested will report false.  | 
**TargetHeight** | Pointer to **int32** | Customer-requested output height, echoed back. Null when the job ran at source dimensions (passthrough).  | [optional] 
**Width** | Pointer to **int32** | Actual encoded output width in pixels (post-ffprobe). Null until the job completes or if the probe failed.  | [optional] 
**Height** | Pointer to **int32** | Actual encoded output height in pixels. See &#x60;width&#x60;. | [optional] 
**Bytes** | Pointer to **int64** |  | [optional] 
**Sha256** | Pointer to **string** |  | [optional] 
**RetentionExpiresAt** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewJobOutputInfo

`func NewJobOutputInfo(state string, container string, audio bool, ) *JobOutputInfo`

NewJobOutputInfo instantiates a new JobOutputInfo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewJobOutputInfoWithDefaults

`func NewJobOutputInfoWithDefaults() *JobOutputInfo`

NewJobOutputInfoWithDefaults instantiates a new JobOutputInfo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetState

`func (o *JobOutputInfo) GetState() string`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *JobOutputInfo) GetStateOk() (*string, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *JobOutputInfo) SetState(v string)`

SetState sets State field to given value.


### GetContainer

`func (o *JobOutputInfo) GetContainer() string`

GetContainer returns the Container field if non-nil, zero value otherwise.

### GetContainerOk

`func (o *JobOutputInfo) GetContainerOk() (*string, bool)`

GetContainerOk returns a tuple with the Container field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContainer

`func (o *JobOutputInfo) SetContainer(v string)`

SetContainer sets Container field to given value.


### GetAudio

`func (o *JobOutputInfo) GetAudio() bool`

GetAudio returns the Audio field if non-nil, zero value otherwise.

### GetAudioOk

`func (o *JobOutputInfo) GetAudioOk() (*bool, bool)`

GetAudioOk returns a tuple with the Audio field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAudio

`func (o *JobOutputInfo) SetAudio(v bool)`

SetAudio sets Audio field to given value.


### GetTargetHeight

`func (o *JobOutputInfo) GetTargetHeight() int32`

GetTargetHeight returns the TargetHeight field if non-nil, zero value otherwise.

### GetTargetHeightOk

`func (o *JobOutputInfo) GetTargetHeightOk() (*int32, bool)`

GetTargetHeightOk returns a tuple with the TargetHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTargetHeight

`func (o *JobOutputInfo) SetTargetHeight(v int32)`

SetTargetHeight sets TargetHeight field to given value.

### HasTargetHeight

`func (o *JobOutputInfo) HasTargetHeight() bool`

HasTargetHeight returns a boolean if a field has been set.

### GetWidth

`func (o *JobOutputInfo) GetWidth() int32`

GetWidth returns the Width field if non-nil, zero value otherwise.

### GetWidthOk

`func (o *JobOutputInfo) GetWidthOk() (*int32, bool)`

GetWidthOk returns a tuple with the Width field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWidth

`func (o *JobOutputInfo) SetWidth(v int32)`

SetWidth sets Width field to given value.

### HasWidth

`func (o *JobOutputInfo) HasWidth() bool`

HasWidth returns a boolean if a field has been set.

### GetHeight

`func (o *JobOutputInfo) GetHeight() int32`

GetHeight returns the Height field if non-nil, zero value otherwise.

### GetHeightOk

`func (o *JobOutputInfo) GetHeightOk() (*int32, bool)`

GetHeightOk returns a tuple with the Height field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHeight

`func (o *JobOutputInfo) SetHeight(v int32)`

SetHeight sets Height field to given value.

### HasHeight

`func (o *JobOutputInfo) HasHeight() bool`

HasHeight returns a boolean if a field has been set.

### GetBytes

`func (o *JobOutputInfo) GetBytes() int64`

GetBytes returns the Bytes field if non-nil, zero value otherwise.

### GetBytesOk

`func (o *JobOutputInfo) GetBytesOk() (*int64, bool)`

GetBytesOk returns a tuple with the Bytes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBytes

`func (o *JobOutputInfo) SetBytes(v int64)`

SetBytes sets Bytes field to given value.

### HasBytes

`func (o *JobOutputInfo) HasBytes() bool`

HasBytes returns a boolean if a field has been set.

### GetSha256

`func (o *JobOutputInfo) GetSha256() string`

GetSha256 returns the Sha256 field if non-nil, zero value otherwise.

### GetSha256Ok

`func (o *JobOutputInfo) GetSha256Ok() (*string, bool)`

GetSha256Ok returns a tuple with the Sha256 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSha256

`func (o *JobOutputInfo) SetSha256(v string)`

SetSha256 sets Sha256 field to given value.

### HasSha256

`func (o *JobOutputInfo) HasSha256() bool`

HasSha256 returns a boolean if a field has been set.

### GetRetentionExpiresAt

`func (o *JobOutputInfo) GetRetentionExpiresAt() time.Time`

GetRetentionExpiresAt returns the RetentionExpiresAt field if non-nil, zero value otherwise.

### GetRetentionExpiresAtOk

`func (o *JobOutputInfo) GetRetentionExpiresAtOk() (*time.Time, bool)`

GetRetentionExpiresAtOk returns a tuple with the RetentionExpiresAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRetentionExpiresAt

`func (o *JobOutputInfo) SetRetentionExpiresAt(v time.Time)`

SetRetentionExpiresAt sets RetentionExpiresAt field to given value.

### HasRetentionExpiresAt

`func (o *JobOutputInfo) HasRetentionExpiresAt() bool`

HasRetentionExpiresAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


