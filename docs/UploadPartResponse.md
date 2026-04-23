# UploadPartResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PartNumber** | **int64** |  | 
**Received** | **bool** |  | 

## Methods

### NewUploadPartResponse

`func NewUploadPartResponse(partNumber int64, received bool, ) *UploadPartResponse`

NewUploadPartResponse instantiates a new UploadPartResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUploadPartResponseWithDefaults

`func NewUploadPartResponseWithDefaults() *UploadPartResponse`

NewUploadPartResponseWithDefaults instantiates a new UploadPartResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPartNumber

`func (o *UploadPartResponse) GetPartNumber() int64`

GetPartNumber returns the PartNumber field if non-nil, zero value otherwise.

### GetPartNumberOk

`func (o *UploadPartResponse) GetPartNumberOk() (*int64, bool)`

GetPartNumberOk returns a tuple with the PartNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartNumber

`func (o *UploadPartResponse) SetPartNumber(v int64)`

SetPartNumber sets PartNumber field to given value.


### GetReceived

`func (o *UploadPartResponse) GetReceived() bool`

GetReceived returns the Received field if non-nil, zero value otherwise.

### GetReceivedOk

`func (o *UploadPartResponse) GetReceivedOk() (*bool, bool)`

GetReceivedOk returns a tuple with the Received field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReceived

`func (o *UploadPartResponse) SetReceived(v bool)`

SetReceived sets Received field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


