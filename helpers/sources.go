package helpers

import sophon "github.com/Liqhtworks/sophon-sdk-go"

// jobSourceConstructors namespaces typed constructors for the oneOf
// CreateJobRequest.Source field. Use the package-level JobSource value.
type jobSourceConstructors struct{}

// JobSource is the recommended way to build a CreateJobRequest.Source
// value. Each method returns a properly-typed oneOf variant so callers
// do not have to remember the discriminator value.
//
//	req := sophon.CreateJobRequest{
//	    Source:  helpers.JobSource.Upload(uploadID),
//	    Profile: sophon.SOPHON_AUTO,
//	}
//
// Today only the upload source is supported; additional source kinds
// will be added here as the API surface grows.
var JobSource jobSourceConstructors

// Upload returns an UploadJobSource referencing a completed upload session.
func (jobSourceConstructors) Upload(uploadID string) sophon.UploadJobSource {
	return sophon.UploadJobSource{Type: sophon.UPLOAD, UploadId: uploadID}
}

// UploadJobSource is a free-function shorthand equivalent to
// JobSource.Upload(uploadID). Prefer the JobSource.Upload form in new code.
func UploadJobSource(uploadID string) sophon.UploadJobSource {
	return JobSource.Upload(uploadID)
}
