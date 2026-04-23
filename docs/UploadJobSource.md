# UploadJobSource

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | [**JobSourceType**](JobSourceType.md) |  | 
**UploadId** | **string** | ID of a completed upload session. | 

## Methods

### NewUploadJobSource

`func NewUploadJobSource(type_ JobSourceType, uploadId string, ) *UploadJobSource`

NewUploadJobSource instantiates a new UploadJobSource object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUploadJobSourceWithDefaults

`func NewUploadJobSourceWithDefaults() *UploadJobSource`

NewUploadJobSourceWithDefaults instantiates a new UploadJobSource object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetType

`func (o *UploadJobSource) GetType() JobSourceType`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *UploadJobSource) GetTypeOk() (*JobSourceType, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *UploadJobSource) SetType(v JobSourceType)`

SetType sets Type field to given value.


### GetUploadId

`func (o *UploadJobSource) GetUploadId() string`

GetUploadId returns the UploadId field if non-nil, zero value otherwise.

### GetUploadIdOk

`func (o *UploadJobSource) GetUploadIdOk() (*string, bool)`

GetUploadIdOk returns a tuple with the UploadId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUploadId

`func (o *UploadJobSource) SetUploadId(v string)`

SetUploadId sets UploadId field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


