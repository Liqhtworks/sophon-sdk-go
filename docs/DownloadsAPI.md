# \DownloadsAPI

All URIs are relative to *https://api.liqhtworks.xyz*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Download**](DownloadsAPI.md#Download) | **Get** /v1/downloads/{token} | Download an output file via signed token



## Download

> *os.File Download(ctx, token).Execute()

Download an output file via signed token



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
	token := "token_example" // string | HMAC-signed download token encoding the object key and expiry.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DownloadsAPI.Download(context.Background(), token).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DownloadsAPI.Download``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `Download`: *os.File
	fmt.Fprintf(os.Stdout, "Response from `DownloadsAPI.Download`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**token** | **string** | HMAC-signed download token encoding the object key and expiry. | 

### Other Parameters

Other parameters are passed through a pointer to a apiDownloadRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[***os.File**](*os.File.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: video/mp4, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

