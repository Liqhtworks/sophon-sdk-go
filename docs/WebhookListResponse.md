# WebhookListResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Webhooks** | [**[]WebhookListItem**](WebhookListItem.md) |  | 

## Methods

### NewWebhookListResponse

`func NewWebhookListResponse(webhooks []WebhookListItem, ) *WebhookListResponse`

NewWebhookListResponse instantiates a new WebhookListResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWebhookListResponseWithDefaults

`func NewWebhookListResponseWithDefaults() *WebhookListResponse`

NewWebhookListResponseWithDefaults instantiates a new WebhookListResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWebhooks

`func (o *WebhookListResponse) GetWebhooks() []WebhookListItem`

GetWebhooks returns the Webhooks field if non-nil, zero value otherwise.

### GetWebhooksOk

`func (o *WebhookListResponse) GetWebhooksOk() (*[]WebhookListItem, bool)`

GetWebhooksOk returns a tuple with the Webhooks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWebhooks

`func (o *WebhookListResponse) SetWebhooks(v []WebhookListItem)`

SetWebhooks sets Webhooks field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


