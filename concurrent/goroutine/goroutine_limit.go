package goroutine

func NewGlimit(limit int) *Glimit {
	return &Glimit{
		limit: limit,
		ch:    make(chan struct{}, limit), // 缓冲长度为limit，运行的协程数量不会超过这个值
	}
}

type Glimit struct {
	limit int
	ch    chan struct{}
}

func (g *Glimit) Run(f func()) {
	g.ch <- struct{}{} //创建子协程前往管道里send一个数据
	go func() {
		f()
		<-g.ch //子协程退出时从管理里取出一个数据
	}()
}
