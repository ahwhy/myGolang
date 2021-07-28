package aliyun_test

import (
	"testing"

	"github.com/ahwhy/myGolang/week8/cloudstation/store/provider/aliyun"
	"github.com/stretchr/testify/assert"
)

var (
	bucketName    = "cloud-station"
	objectKey     = "test.txt"
	localFilePath = "test.txt"

	endpoint = "https://oss-cn-chengdu.aliyuncs.com"
	ak       = "LTAI5tPdJRm5driecxQFt6tK"
	sk       = "edsFTXwFh4bkodhiTdBF4mPVGSJaPb"
)

// DDD的开发流程  DDD -> 测试驱动开发
func TestUploadFile(t *testing.T) {
	should := assert.New(t)

	uploader, err := aliyun.NewUploader(endpoint, ak, sk)
	if should.NoError(err) {
		err = uploader.UploadFile(bucketName, objectKey, localFilePath)
		should.NoError(err)
	}
}
