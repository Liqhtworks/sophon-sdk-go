# CreateJobRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Source** | [**UploadJobSource**](UploadJobSource.md) |  | 
**Profile** | [**JobProfile**](JobProfile.md) |  | 
**Output** | Pointer to [**CreateJobOutputOptions**](CreateJobOutputOptions.md) |  | [optional] 
**WebhookIds** | Pointer to **[]string** | IDs of registered webhook endpoints to notify on job state changes. | [optional] [default to {}]
**Metadata** | Pointer to **map[string]interface{}** | Arbitrary key-value metadata attached to the job. Max 16 KiB serialized. | [optional] [default to {}]

## Methods

### NewCreateJobRequest

`func NewCreateJobRequest(source UploadJobSource, profile JobProfile, ) *CreateJobRequest`

NewCreateJobRequest instantiates a new CreateJobRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateJobRequestWithDefaults

`func NewCreateJobRequestWithDefaults() *CreateJobRequest`

NewCreateJobRequestWithDefaults instantiates a new CreateJobRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSource

`func (o *CreateJobRequest) GetSource() UploadJobSource`

GetSource returns the Source field if non-nil, zero value otherwise.

### GetSourceOk

`func (o *CreateJobRequest) GetSourceOk() (*UploadJobSource, bool)`

GetSourceOk returns a tuple with the Source field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSource

`func (o *CreateJobRequest) SetSource(v UploadJobSource)`

SetSource sets Source field to given value.


### GetProfile

`func (o *CreateJobRequest) GetProfile() JobProfile`

GetProfile returns the Profile field if non-nil, zero value otherwise.

### GetProfileOk

`func (o *CreateJobRequest) GetProfileOk() (*JobProfile, bool)`

GetProfileOk returns a tuple with the Profile field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProfile

`func (o *CreateJobRequest) SetProfile(v JobProfile)`

SetProfile sets Profile field to given value.


### GetOutput

`func (o *CreateJobRequest) GetOutput() CreateJobOutputOptions`

GetOutput returns the Output field if non-nil, zero value otherwise.

### GetOutputOk

`func (o *CreateJobRequest) GetOutputOk() (*CreateJobOutputOptions, bool)`

GetOutputOk returns a tuple with the Output field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutput

`func (o *CreateJobRequest) SetOutput(v CreateJobOutputOptions)`

SetOutput sets Output field to given value.

### HasOutput

`func (o *CreateJobRequest) HasOutput() bool`

HasOutput returns a boolean if a field has been set.

### GetWebhookIds

`func (o *CreateJobRequest) GetWebhookIds() []string`

GetWebhookIds returns the WebhookIds field if non-nil, zero value otherwise.

### GetWebhookIdsOk

`func (o *CreateJobRequest) GetWebhookIdsOk() (*[]string, bool)`

GetWebhookIdsOk returns a tuple with the WebhookIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWebhookIds

`func (o *CreateJobRequest) SetWebhookIds(v []string)`

SetWebhookIds sets WebhookIds field to given value.

### HasWebhookIds

`func (o *CreateJobRequest) HasWebhookIds() bool`

HasWebhookIds returns a boolean if a field has been set.

### GetMetadata

`func (o *CreateJobRequest) GetMetadata() map[string]interface{}`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *CreateJobRequest) GetMetadataOk() (*map[string]interface{}, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *CreateJobRequest) SetMetadata(v map[string]interface{})`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *CreateJobRequest) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


