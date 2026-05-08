package helpers

import sophon "github.com/Liqhtworks/sophon-sdk-go"

type jobSourceConstructors struct{}

var JobSource jobSourceConstructors

func (jobSourceConstructors) Upload(uploadID string) sophon.UploadJobSource {
	return sophon.UploadJobSource{Type: sophon.UPLOAD, UploadId: uploadID}
}

func UploadJobSource(uploadID string) sophon.UploadJobSource {
	return JobSource.Upload(uploadID)
}
