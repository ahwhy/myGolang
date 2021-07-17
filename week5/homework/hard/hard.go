package main

import (
	"log"
	"sync"
	"time"
)

/*
- interface
- map
- 发布系统，能发布k8s node(jtype字段) 不同的类型的任务
- 增量更新：开启新的，停掉旧的
- 原有的是a,b,c ，b,c,d --> 开启d, 停掉a
*/

type jobManager struct {
	jobMutex sync.RWMutex
	// 增量更新
	activeJobs map[string]deployJob
}

func NewJobManager() *jobManager {
	return &jobManager{
		jobMutex:   sync.RWMutex{},
		activeJobs: map[string]deployJob{},
	}
}

func (jm *jobManager) sync(jobs []deployJob) {
	// ***增量更新
	jb := make(map[string]deployJob)

	for _, v := range jobs {
		jb[v.GetName()] = v
	}

	jm.jobMutex.Lock()
	defer jm.jobMutex.Unlock()
	for k, jmv := range jm.activeJobs {
		if jbv, ok := jb[k]; ok && jmv.GetKind() == jbv.GetKind() {
			delete(jb, k)
		} else {
			go jm.activeJobs[k].stop()
		}
	}
	for k, v := range jb {
		jm.activeJobs[k] = v
		go jm.activeJobs[k].start()
	}
}

type deployJob interface {
	start()
	stop()
	SetFlag(b bool)
	GetName() string
	GetFlag() bool
	GetKind() string
}

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

func main() {
	jm := NewJobManager()
	go jm.sync([]deployJob{NewK8sD("a"), NewK8sD("b"), NewK8sD("c")})
	time.Sleep(10 * time.Second)
	// go jm.sync([]deployJob{NewHost("A"), NewHost("B"), NewHost("C")})
	// time.Sleep(10 * time.Second)
	go jm.sync([]deployJob{NewK8sD("b"), NewK8sD("c"), NewK8sD("d")})
	time.Sleep(10 * time.Second)
	// go jm.sync([]deployJob{NewK8sD("c"), NewK8sD("d"), NewK8sD("e")})
	// time.Sleep(10 * time.Second)
	// go jm.sync([]deployJob{NewHost("A"), NewHost("B"), NewHost("C"), NewHost("D")})
	// time.Sleep(10 * time.Second)
	// go jm.sync([]deployJob{NewK8sD("c"), NewK8sD("d"), NewK8sD("e"), NewK8sD("f")})
	// time.Sleep(10 * time.Second)
}
