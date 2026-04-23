# CreateUploadRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**FileName** | **string** |  | 
**FileSize** | **int64** | Total file size in bytes. | 
**MimeType** | **string** |  | 

## Methods

### NewCreateUploadRequest

`func NewCreateUploadRequest(fileName string, fileSize int64, mimeType string, ) *CreateUploadRequest`

NewCreateUploadRequest instantiates a new CreateUploadRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateUploadRequestWithDefaults

`func NewCreateUploadRequestWithDefaults() *CreateUploadRequest`

NewCreateUploadRequestWithDefaults instantiates a new CreateUploadRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFileName

`func (o *CreateUploadRequest) GetFileName() string`

GetFileName returns the FileName field if non-nil, zero value otherwise.

### GetFileNameOk

`func (o *CreateUploadRequest) GetFileNameOk() (*string, bool)`

GetFileNameOk returns a tuple with the FileName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFileName

`func (o *CreateUploadRequest) SetFileName(v string)`

SetFileName sets FileName field to given value.


### GetFileSize

`func (o *CreateUploadRequest) GetFileSize() int64`

GetFileSize returns the FileSize field if non-nil, zero value otherwise.

### GetFileSizeOk

`func (o *CreateUploadRequest) GetFileSizeOk() (*int64, bool)`

GetFileSizeOk returns a tuple with the FileSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFileSize

`func (o *CreateUploadRequest) SetFileSize(v int64)`

SetFileSize sets FileSize field to given value.


### GetMimeType

`func (o *CreateUploadRequest) GetMimeType() string`

GetMimeType returns the MimeType field if non-nil, zero value otherwise.

### GetMimeTypeOk

`func (o *CreateUploadRequest) GetMimeTypeOk() (*string, bool)`

GetMimeTypeOk returns a tuple with the MimeType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMimeType

`func (o *CreateUploadRequest) SetMimeType(v string)`

SetMimeType sets MimeType field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


