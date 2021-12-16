package cmd

import (
	"fmt"
	"path"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ahwhy/myGolang/historys/week08/cloudstation/store"
	"github.com/ahwhy/myGolang/historys/week08/cloudstation/store/provider/aliyun"
	"github.com/spf13/cobra"
)

const (
	// BuckName todo
	defaultBuckName = "cloud-station"
	defaultEndpoint = "https://oss-cn-chengdu.aliyuncs.com"
	defaultALIAK    = "LTAI5tPdJRm5driecxQFt6tK"
	defaultALISK    = "******"
)

var (
	buckName       string
	uploadFilePath string
	bucketEndpoint string
)

func init() {
	uploadCmd.PersistentFlags().StringVarP(&uploadFilePath, "file_path", "f", "", "upload file path")
	uploadCmd.PersistentFlags().StringVarP(&buckName, "bucket_name", "b", defaultBuckName, "upload oss bucket name")
	uploadCmd.PersistentFlags().StringVarP(&bucketEndpoint, "bucket_endpoint", "e", defaultEndpoint, "upload oss endpoint")
	RootCmd.AddCommand(uploadCmd)
}

func getAccessKeyFromInput() {
	prompt := &survey.Password{
		Message: "请输入access key: ",
	}
	survey.AskOne(prompt, &aliAccessKey)
}

func getProvider() (p store.Uploader, err error) {
	switch ossProvider {
	case "aliyun":
		fmt.Printf("上传云厂商: 阿里云[%s]\n", defaultEndpoint)
		if aliAccessID == "" {
			aliAccessID = defaultALIAK
		}
		if aliAccessKey == "" {
			aliAccessKey = defaultALISK
		}
		fmt.Printf("上传用户: %s\n", aliAccessID)
		getAccessKeyFromInput()
		p, err = aliyun.NewUploader(bucketEndpoint, aliAccessID, aliAccessKey)
		return
	case "qcloud":
		return nil, fmt.Errorf("not impl")
	default:
		return nil, fmt.Errorf("unknown oss privier options [aliyun/qcloud]")
	}
}

// uploadCmd represents the start command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "上传文件到中转站",
	Long:  `上传文件到中转站`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := getProvider()
		if err != nil {
			return err
		}
		if uploadFilePath == "" {
			return fmt.Errorf("upload file path is missing")
		}
		day := time.Now().Format("20060102")
		fn := path.Base(uploadFilePath)
		ok := fmt.Sprintf("%s/%s", day, fn)
		err = p.UploadFile(buckName, ok, uploadFilePath)
		if err != nil {
			return err
		}
		return nil
	},
}
