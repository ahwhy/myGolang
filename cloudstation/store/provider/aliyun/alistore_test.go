package aliyun_test

import (
	"testing"
	"time"

	"github.com/ahwhy/myGolang/cloudstation/store/provider/aliyun"
	"github.com/schollz/progressbar/v3"
	"github.com/stretchr/testify/assert"
)

var (
	bucketName    = "cloud-station"
	objectKey     = "test.txt"
	localFilePath = "test.txt"

	endpoint = "https://oss-cn-chengdu.aliyuncs.com"
	ak       = "LTAI5tPdJRm5driecxQFt6tK"
	sk       = "******"
)

func TestUploadFile(t *testing.T) {
	should := assert.New(t)

	uploader, err := aliyun.NewUploader(endpoint, ak, sk)
	if should.NoError(err) {
		err = uploader.UploadFile(bucketName, objectKey, localFilePath)
		should.NoError(err)
	}
}

func TestProgressbar(t *testing.T) {
	bar := progressbar.Default(100)
	for i := 0; i < 100; i++ {
		bar.Add(1)
		time.Sleep(40 * time.Millisecond)
	}
}
