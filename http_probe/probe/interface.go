package probe

import (
	"fmt"
	"time"
)

type Service interface {
	DoHttpProbe(string) string
}

// 传递 time.Duration类型，返回ms 单位的字符串
func MsdurationStr(d time.Duration) string {
	return fmt.Sprintf("%dms", int(d/time.Millisecond)) // s/1000 = ms
}
