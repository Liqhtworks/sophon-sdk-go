# CreateUploadResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**ChunkSize** | **int64** | Size of each chunk in bytes. Tiered by file size: &lt;64 MB &#x3D; whole file, &lt;&#x3D;1 GB &#x3D; 8 MB, &lt;&#x3D;10 GB &#x3D; 16 MB, &gt;10 GB &#x3D; 32 MB.  | 
**TotalChunks** | **int64** |  | 
**ExpiresAt** | **time.Time** | Upload session expiry (24 hours from creation). | 

## Methods

### NewCreateUploadResponse

`func NewCreateUploadResponse(id string, chunkSize int64, totalChunks int64, expiresAt time.Time, ) *CreateUploadResponse`

NewCreateUploadResponse instantiates a new CreateUploadResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateUploadResponseWithDefaults

`func NewCreateUploadResponseWithDefaults() *CreateUploadResponse`

NewCreateUploadResponseWithDefaults instantiates a new CreateUploadResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *CreateUploadResponse) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *CreateUploadResponse) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *CreateUploadResponse) SetId(v string)`

SetId sets Id field to given value.


### GetChunkSize

`func (o *CreateUploadResponse) GetChunkSize() int64`

GetChunkSize returns the ChunkSize field if non-nil, zero value otherwise.

### GetChunkSizeOk

`func (o *CreateUploadResponse) GetChunkSizeOk() (*int64, bool)`

GetChunkSizeOk returns a tuple with the ChunkSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChunkSize

`func (o *CreateUploadResponse) SetChunkSize(v int64)`

SetChunkSize sets ChunkSize field to given value.


### GetTotalChunks

`func (o *CreateUploadResponse) GetTotalChunks() int64`

GetTotalChunks returns the TotalChunks field if non-nil, zero value otherwise.

### GetTotalChunksOk

`func (o *CreateUploadResponse) GetTotalChunksOk() (*int64, bool)`

GetTotalChunksOk returns a tuple with the TotalChunks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalChunks

`func (o *CreateUploadResponse) SetTotalChunks(v int64)`

SetTotalChunks sets TotalChunks field to given value.


### GetExpiresAt

`func (o *CreateUploadResponse) GetExpiresAt() time.Time`

GetExpiresAt returns the ExpiresAt field if non-nil, zero value otherwise.

### GetExpiresAtOk

`func (o *CreateUploadResponse) GetExpiresAtOk() (*time.Time, bool)`

GetExpiresAtOk returns a tuple with the ExpiresAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiresAt

`func (o *CreateUploadResponse) SetExpiresAt(v time.Time)`

SetExpiresAt sets ExpiresAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


