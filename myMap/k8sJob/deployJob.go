package k8sJob

/*
- interface
- map
- 发布系统，能发布k8s node(jtype字段) 不同的类型的任务
- 增量更新：开启新的，停掉旧的
- 原有的是a,b,c ，b,c,d --> 开启d, 停掉a
*/

type DeployJob interface {
	start()
	stop()
	SetFlag(b bool)
	GetName() string
	GetFlag() bool
	GetKind() string
}
