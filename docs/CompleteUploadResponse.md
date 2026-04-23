# CompleteUploadResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**Status** | **string** |  | 
**Sha256** | **string** | SHA-256 hex digest of the assembled file. | 
**Bytes** | **int64** |  | 

## Methods

### NewCompleteUploadResponse

`func NewCompleteUploadResponse(id string, status string, sha256 string, bytes int64, ) *CompleteUploadResponse`

NewCompleteUploadResponse instantiates a new CompleteUploadResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCompleteUploadResponseWithDefaults

`func NewCompleteUploadResponseWithDefaults() *CompleteUploadResponse`

NewCompleteUploadResponseWithDefaults instantiates a new CompleteUploadResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *CompleteUploadResponse) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *CompleteUploadResponse) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *CompleteUploadResponse) SetId(v string)`

SetId sets Id field to given value.


### GetStatus

`func (o *CompleteUploadResponse) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *CompleteUploadResponse) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *CompleteUploadResponse) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetSha256

`func (o *CompleteUploadResponse) GetSha256() string`

GetSha256 returns the Sha256 field if non-nil, zero value otherwise.

### GetSha256Ok

`func (o *CompleteUploadResponse) GetSha256Ok() (*string, bool)`

GetSha256Ok returns a tuple with the Sha256 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSha256

`func (o *CompleteUploadResponse) SetSha256(v string)`

SetSha256 sets Sha256 field to given value.


### GetBytes

`func (o *CompleteUploadResponse) GetBytes() int64`

GetBytes returns the Bytes field if non-nil, zero value otherwise.

### GetBytesOk

`func (o *CompleteUploadResponse) GetBytesOk() (*int64, bool)`

GetBytesOk returns a tuple with the Bytes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBytes

`func (o *CompleteUploadResponse) SetBytes(v int64)`

SetBytes sets Bytes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


