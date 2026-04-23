# JobSourceInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Original file name of the source. | [optional] 
**Bytes** | Pointer to **int64** |  | [optional] 
**Sha256** | **string** | SHA-256 hex digest of the source file. | 
**DurationSeconds** | Pointer to **float64** |  | [optional] 
**Resolution** | Pointer to **string** |  | [optional] 
**FrameRate** | Pointer to **string** |  | [optional] 

## Methods

### NewJobSourceInfo

`func NewJobSourceInfo(sha256 string, ) *JobSourceInfo`

NewJobSourceInfo instantiates a new JobSourceInfo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewJobSourceInfoWithDefaults

`func NewJobSourceInfoWithDefaults() *JobSourceInfo`

NewJobSourceInfoWithDefaults instantiates a new JobSourceInfo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *JobSourceInfo) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *JobSourceInfo) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *JobSourceInfo) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *JobSourceInfo) HasName() bool`

HasName returns a boolean if a field has been set.

### GetBytes

`func (o *JobSourceInfo) GetBytes() int64`

GetBytes returns the Bytes field if non-nil, zero value otherwise.

### GetBytesOk

`func (o *JobSourceInfo) GetBytesOk() (*int64, bool)`

GetBytesOk returns a tuple with the Bytes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBytes

`func (o *JobSourceInfo) SetBytes(v int64)`

SetBytes sets Bytes field to given value.

### HasBytes

`func (o *JobSourceInfo) HasBytes() bool`

HasBytes returns a boolean if a field has been set.

### GetSha256

`func (o *JobSourceInfo) GetSha256() string`

GetSha256 returns the Sha256 field if non-nil, zero value otherwise.

### GetSha256Ok

`func (o *JobSourceInfo) GetSha256Ok() (*string, bool)`

GetSha256Ok returns a tuple with the Sha256 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSha256

`func (o *JobSourceInfo) SetSha256(v string)`

SetSha256 sets Sha256 field to given value.


### GetDurationSeconds

`func (o *JobSourceInfo) GetDurationSeconds() float64`

GetDurationSeconds returns the DurationSeconds field if non-nil, zero value otherwise.

### GetDurationSecondsOk

`func (o *JobSourceInfo) GetDurationSecondsOk() (*float64, bool)`

GetDurationSecondsOk returns a tuple with the DurationSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDurationSeconds

`func (o *JobSourceInfo) SetDurationSeconds(v float64)`

SetDurationSeconds sets DurationSeconds field to given value.

### HasDurationSeconds

`func (o *JobSourceInfo) HasDurationSeconds() bool`

HasDurationSeconds returns a boolean if a field has been set.

### GetResolution

`func (o *JobSourceInfo) GetResolution() string`

GetResolution returns the Resolution field if non-nil, zero value otherwise.

### GetResolutionOk

`func (o *JobSourceInfo) GetResolutionOk() (*string, bool)`

GetResolutionOk returns a tuple with the Resolution field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResolution

`func (o *JobSourceInfo) SetResolution(v string)`

SetResolution sets Resolution field to given value.

### HasResolution

`func (o *JobSourceInfo) HasResolution() bool`

HasResolution returns a boolean if a field has been set.

### GetFrameRate

`func (o *JobSourceInfo) GetFrameRate() string`

GetFrameRate returns the FrameRate field if non-nil, zero value otherwise.

### GetFrameRateOk

`func (o *JobSourceInfo) GetFrameRateOk() (*string, bool)`

GetFrameRateOk returns a tuple with the FrameRate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFrameRate

`func (o *JobSourceInfo) SetFrameRate(v string)`

SetFrameRate sets FrameRate field to given value.

### HasFrameRate

`func (o *JobSourceInfo) HasFrameRate() bool`

HasFrameRate returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


