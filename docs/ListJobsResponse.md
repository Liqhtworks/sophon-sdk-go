# ListJobsResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Jobs** | [**[]JobResponse**](JobResponse.md) |  | 
**NextCursor** | Pointer to **string** | Opaque cursor for the next page. Null when no more results. | [optional] 
**HasMore** | **bool** |  | 

## Methods

### NewListJobsResponse

`func NewListJobsResponse(jobs []JobResponse, hasMore bool, ) *ListJobsResponse`

NewListJobsResponse instantiates a new ListJobsResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewListJobsResponseWithDefaults

`func NewListJobsResponseWithDefaults() *ListJobsResponse`

NewListJobsResponseWithDefaults instantiates a new ListJobsResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetJobs

`func (o *ListJobsResponse) GetJobs() []JobResponse`

GetJobs returns the Jobs field if non-nil, zero value otherwise.

### GetJobsOk

`func (o *ListJobsResponse) GetJobsOk() (*[]JobResponse, bool)`

GetJobsOk returns a tuple with the Jobs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJobs

`func (o *ListJobsResponse) SetJobs(v []JobResponse)`

SetJobs sets Jobs field to given value.


### GetNextCursor

`func (o *ListJobsResponse) GetNextCursor() string`

GetNextCursor returns the NextCursor field if non-nil, zero value otherwise.

### GetNextCursorOk

`func (o *ListJobsResponse) GetNextCursorOk() (*string, bool)`

GetNextCursorOk returns a tuple with the NextCursor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextCursor

`func (o *ListJobsResponse) SetNextCursor(v string)`

SetNextCursor sets NextCursor field to given value.

### HasNextCursor

`func (o *ListJobsResponse) HasNextCursor() bool`

HasNextCursor returns a boolean if a field has been set.

### GetHasMore

`func (o *ListJobsResponse) GetHasMore() bool`

GetHasMore returns the HasMore field if non-nil, zero value otherwise.

### GetHasMoreOk

`func (o *ListJobsResponse) GetHasMoreOk() (*bool, bool)`

GetHasMoreOk returns a tuple with the HasMore field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHasMore

`func (o *ListJobsResponse) SetHasMore(v bool)`

SetHasMore sets HasMore field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


