package k8sJob

import "sync"

type JobManager struct {
	JobMutex sync.RWMutex
	// 增量更新
	ActiveJobs map[string]DeployJob
}

func NewJobManager() *JobManager {
	return &JobManager{
		JobMutex:   sync.RWMutex{},
		ActiveJobs: map[string]DeployJob{},
	}
}

// Sync 增量更新
func (jm *JobManager) Sync(jobs []DeployJob) {
	// 注册需要运行的任务至jb
	jb := make(map[string]DeployJob)
	for _, v := range jobs {
		jb[v.GetName()] = v
	}

	jm.JobMutex.Lock()
	defer jm.JobMutex.Unlock()

	// 如果jm中的任务，jb中也有，则代表任务不变，从jb中删除
	// 反之，如果jm中的任务，jb中没有，则代表任务需要被停止
	for k, jmv := range jm.ActiveJobs {
		if jbv, ok := jb[k]; ok && jmv.GetKind() == jbv.GetKind() {
			delete(jb, k)
		} else {
			go jm.ActiveJobs[k].stop()
		}
	}

	// 启动新增任务，即jb中剩余任务
	for k, v := range jb {
		jm.ActiveJobs[k] = v
		go jm.ActiveJobs[k].start()
	}
}
