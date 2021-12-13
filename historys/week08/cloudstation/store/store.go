package store

type Uploader interface {
	UploadFile(bucketName, objectKey, localFilePath string) error
}
