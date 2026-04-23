# \UploadsAPI

All URIs are relative to *https://api.liqhtworks.xyz*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CancelUpload**](UploadsAPI.md#CancelUpload) | **Delete** /v1/uploads/{id} | Cancel an upload session
[**CompleteUpload**](UploadsAPI.md#CompleteUpload) | **Post** /v1/uploads/{id}/complete | Finalize a chunked upload
[**CreateUpload**](UploadsAPI.md#CreateUpload) | **Post** /v1/uploads | Initialize a chunked upload session
[**GetUpload**](UploadsAPI.md#GetUpload) | **Get** /v1/uploads/{id} | Get upload session status
[**UploadPart**](UploadsAPI.md#UploadPart) | **Put** /v1/uploads/{id}/parts/{part_number} | Upload a single chunk



## CancelUpload

> CancelUpload(ctx, id).Execute()

Cancel an upload session



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/Liqhtworks/sophon-sdk-go"
)

func main() {
	id := "id_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.UploadsAPI.CancelUpload(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UploadsAPI.CancelUpload``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiCancelUploadRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[sessionCookie](../README.md#sessionCookie), [bearerApiKey](../README.md#bearerApiKey)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CompleteUpload

> CompleteUploadResponse CompleteUpload(ctx, id).IdempotencyKey(idempotencyKey).Execute()

Finalize a chunked upload



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/Liqhtworks/sophon-sdk-go"
)

func main() {
	id := "id_example" // string | 
	idempotencyKey := "idempotencyKey_example" // string | Client-generated UUID or string for exactly-once semantics. Required on all POST endpoints. Replaying the same key with the same request body returns the original response without side effects. 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.UploadsAPI.CompleteUpload(context.Background(), id).IdempotencyKey(idempotencyKey).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UploadsAPI.CompleteUpload``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CompleteUpload`: CompleteUploadResponse
	fmt.Fprintf(os.Stdout, "Response from `UploadsAPI.CompleteUpload`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiCompleteUploadRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **idempotencyKey** | **string** | Client-generated UUID or string for exactly-once semantics. Required on all POST endpoints. Replaying the same key with the same request body returns the original response without side effects.  | 

### Return type

[**CompleteUploadResponse**](CompleteUploadResponse.md)

### Authorization

[sessionCookie](../README.md#sessionCookie), [bearerApiKey](../README.md#bearerApiKey)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateUpload

> CreateUploadResponse CreateUpload(ctx).IdempotencyKey(idempotencyKey).CreateUploadRequest(createUploadRequest).Execute()

Initialize a chunked upload session



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/Liqhtworks/sophon-sdk-go"
)

func main() {
	idempotencyKey := "idempotencyKey_example" // string | Client-generated UUID or string for exactly-once semantics. Required on all POST endpoints. Replaying the same key with the same request body returns the original response without side effects. 
	createUploadRequest := *openapiclient.NewCreateUploadRequest("FileName_example", int64(123), "MimeType_example") // CreateUploadRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.UploadsAPI.CreateUpload(context.Background()).IdempotencyKey(idempotencyKey).CreateUploadRequest(createUploadRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UploadsAPI.CreateUpload``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateUpload`: CreateUploadResponse
	fmt.Fprintf(os.Stdout, "Response from `UploadsAPI.CreateUpload`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateUploadRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **idempotencyKey** | **string** | Client-generated UUID or string for exactly-once semantics. Required on all POST endpoints. Replaying the same key with the same request body returns the original response without side effects.  | 
 **createUploadRequest** | [**CreateUploadRequest**](CreateUploadRequest.md) |  | 

### Return type

[**CreateUploadResponse**](CreateUploadResponse.md)

### Authorization

[sessionCookie](../README.md#sessionCookie), [bearerApiKey](../README.md#bearerApiKey)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetUpload

> UploadStatusResponse GetUpload(ctx, id).Execute()

Get upload session status



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/Liqhtworks/sophon-sdk-go"
)

func main() {
	id := "id_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.UploadsAPI.GetUpload(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UploadsAPI.GetUpload``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetUpload`: UploadStatusResponse
	fmt.Fprintf(os.Stdout, "Response from `UploadsAPI.GetUpload`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetUploadRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**UploadStatusResponse**](UploadStatusResponse.md)

### Authorization

[sessionCookie](../README.md#sessionCookie), [bearerApiKey](../README.md#bearerApiKey)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UploadPart

> UploadPartResponse UploadPart(ctx, id, partNumber).Body(body).Execute()

Upload a single chunk



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/Liqhtworks/sophon-sdk-go"
)

func main() {
	id := "id_example" // string | 
	partNumber := int32(56) // int32 | 
	body := os.NewFile(1234, "some_file") // *os.File | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.UploadsAPI.UploadPart(context.Background(), id, partNumber).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UploadsAPI.UploadPart``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UploadPart`: UploadPartResponse
	fmt.Fprintf(os.Stdout, "Response from `UploadsAPI.UploadPart`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 
**partNumber** | **int32** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiUploadPartRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **body** | ***os.File** |  | 

### Return type

[**UploadPartResponse**](UploadPartResponse.md)

### Authorization

[sessionCookie](../README.md#sessionCookie), [bearerApiKey](../README.md#bearerApiKey)

### HTTP request headers

- **Content-Type**: application/octet-stream
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

