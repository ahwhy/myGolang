package k8sJob

import (
	"log"
	"time"
)

type k8sD struct {
	Name string
	Flag bool
	Kind string
}

func NewK8sD(name string) *k8sD {
	return &k8sD{
		Name: name,
		Kind: "K8sD",
	}
}

func (k *k8sD) start() {
	k.SetFlag(true)
	log.Printf("Pod %s 开始运行 \n", k.Name)
	count := 0

	for {
		if k.Flag {
			count++
			log.Printf("Pod %s 正在运行 %d 秒\n", k.Name, count)
			time.Sleep(1 * time.Second)
		} else {
			log.Printf("Pod %s 退出运行 \n", k.Name)
			break
		}
	}
}

func (k *k8sD) stop() {
	k.SetFlag(false)
	log.Printf("Pod %s 退出运行 \n", k.Name)
}

func (k *k8sD) SetFlag(b bool) {
	k.Flag = b
}

func (k *k8sD) GetName() string {
	return k.Name
}

func (k *k8sD) GetFlag() bool {
	return k.Flag
}

func (k *k8sD) GetKind() string {
	return k.Kind
}
