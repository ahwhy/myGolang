package k8sJob

import (
	"log"
	"time"
)

type host struct {
	Name string
	Flag bool
	Kind string
	// Pod  []k8sD
}

func NewHost(name string) *host {
	return &host{
		Name: name,
		Kind: "Host",
	}
}

func (h *host) start() {
	h.SetFlag(true)
	log.Printf("Host %s 开始运行 \n", h.Name)
	count := 0
	
	for {
		count++
		if h.Flag {
			log.Printf("Host %s 正在运行 %d 秒\n", h.Name, count)
			time.Sleep(1 * time.Second)
		} else {
			log.Printf("Host %s 退出运行 \n", h.Name)
			break
		}
	}
}

func (h *host) stop() {
	h.SetFlag(false)
	log.Printf("Host %s 退出运行 \n", h.Name)
}

func (h *host) SetFlag(b bool) {
	h.Flag = b
}

func (h *host) GetName() string {
	return h.Name
}

func (h *host) GetFlag() bool {
	return h.Flag
}

func (h *host) GetKind() string {
	return h.Kind
}