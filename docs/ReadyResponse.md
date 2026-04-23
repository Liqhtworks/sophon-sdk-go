# ReadyResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Ready** | **bool** |  | 
**ChecksFailed** | Pointer to **[]string** | Names of failed readiness checks (database, disk_critical, draining, workers_dead). | [optional] 

## Methods

### NewReadyResponse

`func NewReadyResponse(ready bool, ) *ReadyResponse`

NewReadyResponse instantiates a new ReadyResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReadyResponseWithDefaults

`func NewReadyResponseWithDefaults() *ReadyResponse`

NewReadyResponseWithDefaults instantiates a new ReadyResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetReady

`func (o *ReadyResponse) GetReady() bool`

GetReady returns the Ready field if non-nil, zero value otherwise.

### GetReadyOk

`func (o *ReadyResponse) GetReadyOk() (*bool, bool)`

GetReadyOk returns a tuple with the Ready field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReady

`func (o *ReadyResponse) SetReady(v bool)`

SetReady sets Ready field to given value.


### GetChecksFailed

`func (o *ReadyResponse) GetChecksFailed() []string`

GetChecksFailed returns the ChecksFailed field if non-nil, zero value otherwise.

### GetChecksFailedOk

`func (o *ReadyResponse) GetChecksFailedOk() (*[]string, bool)`

GetChecksFailedOk returns a tuple with the ChecksFailed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChecksFailed

`func (o *ReadyResponse) SetChecksFailed(v []string)`

SetChecksFailed sets ChecksFailed field to given value.

### HasChecksFailed

`func (o *ReadyResponse) HasChecksFailed() bool`

HasChecksFailed returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


