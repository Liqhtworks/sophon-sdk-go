# WebhookDeliveryPayload

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**EventId** | **string** | Unique delivery event ID for deduplication. | 
**Type** | **string** | Event type. | 
**Timestamp** | **time.Time** | ISO 8601 timestamp of the event. | 
**JobId** | **string** | The job that reached a terminal state. | 
**Status** | **string** | Terminal job status. | 
**Metadata** | **map[string]interface{}** | Opaque metadata from the original job submission. | 

## Methods

### NewWebhookDeliveryPayload

`func NewWebhookDeliveryPayload(eventId string, type_ string, timestamp time.Time, jobId string, status string, metadata map[string]interface{}, ) *WebhookDeliveryPayload`

NewWebhookDeliveryPayload instantiates a new WebhookDeliveryPayload object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWebhookDeliveryPayloadWithDefaults

`func NewWebhookDeliveryPayloadWithDefaults() *WebhookDeliveryPayload`

NewWebhookDeliveryPayloadWithDefaults instantiates a new WebhookDeliveryPayload object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEventId

`func (o *WebhookDeliveryPayload) GetEventId() string`

GetEventId returns the EventId field if non-nil, zero value otherwise.

### GetEventIdOk

`func (o *WebhookDeliveryPayload) GetEventIdOk() (*string, bool)`

GetEventIdOk returns a tuple with the EventId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEventId

`func (o *WebhookDeliveryPayload) SetEventId(v string)`

SetEventId sets EventId field to given value.


### GetType

`func (o *WebhookDeliveryPayload) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *WebhookDeliveryPayload) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *WebhookDeliveryPayload) SetType(v string)`

SetType sets Type field to given value.


### GetTimestamp

`func (o *WebhookDeliveryPayload) GetTimestamp() time.Time`

GetTimestamp returns the Timestamp field if non-nil, zero value otherwise.

### GetTimestampOk

`func (o *WebhookDeliveryPayload) GetTimestampOk() (*time.Time, bool)`

GetTimestampOk returns a tuple with the Timestamp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimestamp

`func (o *WebhookDeliveryPayload) SetTimestamp(v time.Time)`

SetTimestamp sets Timestamp field to given value.


### GetJobId

`func (o *WebhookDeliveryPayload) GetJobId() string`

GetJobId returns the JobId field if non-nil, zero value otherwise.

### GetJobIdOk

`func (o *WebhookDeliveryPayload) GetJobIdOk() (*string, bool)`

GetJobIdOk returns a tuple with the JobId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJobId

`func (o *WebhookDeliveryPayload) SetJobId(v string)`

SetJobId sets JobId field to given value.


### GetStatus

`func (o *WebhookDeliveryPayload) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *WebhookDeliveryPayload) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *WebhookDeliveryPayload) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetMetadata

`func (o *WebhookDeliveryPayload) GetMetadata() map[string]interface{}`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *WebhookDeliveryPayload) GetMetadataOk() (*map[string]interface{}, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *WebhookDeliveryPayload) SetMetadata(v map[string]interface{})`

SetMetadata sets Metadata field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


