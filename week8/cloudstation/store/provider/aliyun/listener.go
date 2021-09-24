package aliyun

import (
	"fmt"
	"time"

	"github.com/ahwhy/myGolang/week8/cloudstation/tool"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

// OssProgressListener is the progress listener
type OssProgressListener struct {
	bar     *progressbar.ProgressBar
	startAt time.Time
}

// NewOssProgressListener todo
func NewOssProgressListener() *OssProgressListener {
	return &OssProgressListener{}
}

// ProgressChanged todo
// func (p *OssProgressListener) ProgressChanged(event *oss.ProgressEvent) {
// 	switch event.EventType {
// 	case oss.TransferStartedEvent:
// 		p.bar = progressbar.DefaultBytes(
// 			event.TotalBytes,
// 			"文件上传中",
// 		)
// 	case oss.TransferDataEvent:
// 		p.bar.Add64(event.RwBytes)
// 	case oss.TransferCompletedEvent:
// 		fmt.Printf("\n上传完成\n")
// 	case oss.TransferFailedEvent:
// 		fmt.Printf("\n上传失败\n")
// 	default:
// 	}
// }

// ProgressChanged todo
func (p *OssProgressListener) ProgressChanged(event *oss.ProgressEvent) {
	switch event.EventType {
	case oss.TransferStartedEvent:
		p.bar = progressbar.NewOptions64(event.TotalBytes,
			progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionShowBytes(true),
			progressbar.OptionSetWidth(30),
			progressbar.OptionSetDescription("开始上传:"),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "=",
				SaucerHead:    ">",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}),
		)
		p.startAt = time.Now()
		fmt.Println()
		fmt.Printf("文件大小: %s\n", tool.HumanBytesLoaded(event.TotalBytes))
	// case oss.TransferDataEvent:
	// 	p.bar.Add64(event.RwBytes)
	case oss.TransferCompletedEvent:
		fmt.Printf("\n上传完成: 耗时%d秒\n", int(time.Since(p.startAt).Seconds()))
	case oss.TransferFailedEvent:
		fmt.Printf("\n上传失败: \n")
	default:
	}
}
