package aliyun

import (
	"fmt"

	"github.com/ahwhy/myGolang/cloudstation/store"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

func NewUploader(endpt, accessID, accessKey string) (store.Uploader, error) {
	p := &aliyun{
		Endpoint:  endpt,
		AccessID:  accessID,
		AccessKey: accessKey,

		listener: NewOssProgressListener(),
	}

	if err := p.validate(); err != nil {
		return nil, err
	}

	return p, nil
}

type aliyun struct {
	Endpoint  string `validate:"required,url"`
	AccessID  string `validate:"required"`
	AccessKey string `validate:"required"`

	listener oss.ProgressListener
}

func (p *aliyun) validate() error {
	return validate.Struct(p)
}

func (p *aliyun) UploadFile(bucketName, objectKey, localFilePath string) error {
	return fmt.Errorf("not impl")
}
