# CreateJobOutputOptions

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Container** | Pointer to [**OutputContainer**](OutputContainer.md) |  | [optional] [default to MP4]
**Audio** | Pointer to **bool** | When true, audio is included in the output. MKV preserves source audio streams unchanged. MP4 preserves broadly compatible source audio codecs when possible, and may normalize incompatible codecs to AAC for playback compatibility. When false, the output is video only.  | [optional] [default to false]
**TargetHeight** | Pointer to **int32** | Target output height in pixels. When set, output is scaled down (aspect ratio preserved, width derived from source, both dims rounded to even). If absent or larger than source height, output uses source dimensions. Billing tier is determined by the actual encoded output, not by this requested value.  | [optional] 

## Methods

### NewCreateJobOutputOptions

`func NewCreateJobOutputOptions() *CreateJobOutputOptions`

NewCreateJobOutputOptions instantiates a new CreateJobOutputOptions object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateJobOutputOptionsWithDefaults

`func NewCreateJobOutputOptionsWithDefaults() *CreateJobOutputOptions`

NewCreateJobOutputOptionsWithDefaults instantiates a new CreateJobOutputOptions object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetContainer

`func (o *CreateJobOutputOptions) GetContainer() OutputContainer`

GetContainer returns the Container field if non-nil, zero value otherwise.

### GetContainerOk

`func (o *CreateJobOutputOptions) GetContainerOk() (*OutputContainer, bool)`

GetContainerOk returns a tuple with the Container field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContainer

`func (o *CreateJobOutputOptions) SetContainer(v OutputContainer)`

SetContainer sets Container field to given value.

### HasContainer

`func (o *CreateJobOutputOptions) HasContainer() bool`

HasContainer returns a boolean if a field has been set.

### GetAudio

`func (o *CreateJobOutputOptions) GetAudio() bool`

GetAudio returns the Audio field if non-nil, zero value otherwise.

### GetAudioOk

`func (o *CreateJobOutputOptions) GetAudioOk() (*bool, bool)`

GetAudioOk returns a tuple with the Audio field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAudio

`func (o *CreateJobOutputOptions) SetAudio(v bool)`

SetAudio sets Audio field to given value.

### HasAudio

`func (o *CreateJobOutputOptions) HasAudio() bool`

HasAudio returns a boolean if a field has been set.

### GetTargetHeight

`func (o *CreateJobOutputOptions) GetTargetHeight() int32`

GetTargetHeight returns the TargetHeight field if non-nil, zero value otherwise.

### GetTargetHeightOk

`func (o *CreateJobOutputOptions) GetTargetHeightOk() (*int32, bool)`

GetTargetHeightOk returns a tuple with the TargetHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTargetHeight

`func (o *CreateJobOutputOptions) SetTargetHeight(v int32)`

SetTargetHeight sets TargetHeight field to given value.

### HasTargetHeight

`func (o *CreateJobOutputOptions) HasTargetHeight() bool`

HasTargetHeight returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


