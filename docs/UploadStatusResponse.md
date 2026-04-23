# UploadStatusResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**Status** | **string** |  | 
**FileName** | **string** |  | 
**TotalChunks** | **int32** |  | 
**ReceivedChunks** | **[]int32** | Array of 0-indexed part numbers that have been received. | 
**ExpiresAt** | **time.Time** |  | 
**SourceWidth** | Pointer to **int32** | Source media width in pixels, populated from ffprobe after upload assembly. Null for uploads in &#x60;initiated&#x60;/&#x60;uploading&#x60; state or when probe failed.  | [optional] 
**SourceHeight** | Pointer to **int32** | Source media height in pixels. See &#x60;source_width&#x60;. | [optional] 
**SourceDurationSeconds** | Pointer to **float32** | Source media duration in seconds, from ffprobe after upload assembly. Used by the webapp free-tier budget check to compute realistic billable_seconds (5-second ceiling rounding).  | [optional] 

## Methods

### NewUploadStatusResponse

`func NewUploadStatusResponse(id string, status string, fileName string, totalChunks int32, receivedChunks []int32, expiresAt time.Time, ) *UploadStatusResponse`

NewUploadStatusResponse instantiates a new UploadStatusResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUploadStatusResponseWithDefaults

`func NewUploadStatusResponseWithDefaults() *UploadStatusResponse`

NewUploadStatusResponseWithDefaults instantiates a new UploadStatusResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *UploadStatusResponse) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UploadStatusResponse) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UploadStatusResponse) SetId(v string)`

SetId sets Id field to given value.


### GetStatus

`func (o *UploadStatusResponse) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *UploadStatusResponse) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *UploadStatusResponse) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetFileName

`func (o *UploadStatusResponse) GetFileName() string`

GetFileName returns the FileName field if non-nil, zero value otherwise.

### GetFileNameOk

`func (o *UploadStatusResponse) GetFileNameOk() (*string, bool)`

GetFileNameOk returns a tuple with the FileName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFileName

`func (o *UploadStatusResponse) SetFileName(v string)`

SetFileName sets FileName field to given value.


### GetTotalChunks

`func (o *UploadStatusResponse) GetTotalChunks() int32`

GetTotalChunks returns the TotalChunks field if non-nil, zero value otherwise.

### GetTotalChunksOk

`func (o *UploadStatusResponse) GetTotalChunksOk() (*int32, bool)`

GetTotalChunksOk returns a tuple with the TotalChunks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalChunks

`func (o *UploadStatusResponse) SetTotalChunks(v int32)`

SetTotalChunks sets TotalChunks field to given value.


### GetReceivedChunks

`func (o *UploadStatusResponse) GetReceivedChunks() []int32`

GetReceivedChunks returns the ReceivedChunks field if non-nil, zero value otherwise.

### GetReceivedChunksOk

`func (o *UploadStatusResponse) GetReceivedChunksOk() (*[]int32, bool)`

GetReceivedChunksOk returns a tuple with the ReceivedChunks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReceivedChunks

`func (o *UploadStatusResponse) SetReceivedChunks(v []int32)`

SetReceivedChunks sets ReceivedChunks field to given value.


### GetExpiresAt

`func (o *UploadStatusResponse) GetExpiresAt() time.Time`

GetExpiresAt returns the ExpiresAt field if non-nil, zero value otherwise.

### GetExpiresAtOk

`func (o *UploadStatusResponse) GetExpiresAtOk() (*time.Time, bool)`

GetExpiresAtOk returns a tuple with the ExpiresAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiresAt

`func (o *UploadStatusResponse) SetExpiresAt(v time.Time)`

SetExpiresAt sets ExpiresAt field to given value.


### GetSourceWidth

`func (o *UploadStatusResponse) GetSourceWidth() int32`

GetSourceWidth returns the SourceWidth field if non-nil, zero value otherwise.

### GetSourceWidthOk

`func (o *UploadStatusResponse) GetSourceWidthOk() (*int32, bool)`

GetSourceWidthOk returns a tuple with the SourceWidth field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSourceWidth

`func (o *UploadStatusResponse) SetSourceWidth(v int32)`

SetSourceWidth sets SourceWidth field to given value.

### HasSourceWidth

`func (o *UploadStatusResponse) HasSourceWidth() bool`

HasSourceWidth returns a boolean if a field has been set.

### GetSourceHeight

`func (o *UploadStatusResponse) GetSourceHeight() int32`

GetSourceHeight returns the SourceHeight field if non-nil, zero value otherwise.

### GetSourceHeightOk

`func (o *UploadStatusResponse) GetSourceHeightOk() (*int32, bool)`

GetSourceHeightOk returns a tuple with the SourceHeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSourceHeight

`func (o *UploadStatusResponse) SetSourceHeight(v int32)`

SetSourceHeight sets SourceHeight field to given value.

### HasSourceHeight

`func (o *UploadStatusResponse) HasSourceHeight() bool`

HasSourceHeight returns a boolean if a field has been set.

### GetSourceDurationSeconds

`func (o *UploadStatusResponse) GetSourceDurationSeconds() float32`

GetSourceDurationSeconds returns the SourceDurationSeconds field if non-nil, zero value otherwise.

### GetSourceDurationSecondsOk

`func (o *UploadStatusResponse) GetSourceDurationSecondsOk() (*float32, bool)`

GetSourceDurationSecondsOk returns a tuple with the SourceDurationSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSourceDurationSeconds

`func (o *UploadStatusResponse) SetSourceDurationSeconds(v float32)`

SetSourceDurationSeconds sets SourceDurationSeconds field to given value.

### HasSourceDurationSeconds

`func (o *UploadStatusResponse) HasSourceDurationSeconds() bool`

HasSourceDurationSeconds returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


