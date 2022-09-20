# Golang-Context  Golang的上下文

## 一、Golang的标准库 Context包

- Context 包定义了上下文类型，该上下文类型跨越 API 边界和进程之间传递截止期限，取消信号和其他请求范围值
	- 对服务器的传入请求应创建一个 Context，对服务器的传出调用应接受 Context
		- 它们之间的函数调用链必须传播 Context，可以用使用 `WithCancel`，`WithDeadline`，`WithTimeout` 或 `WithValue` 创建的派生上下文替换
		- 当 Context 被取消时，从它派生的所有 Context 也被取消
	- 即使函数允许，也不要传递 nil Context
		- 如果不确定要使用哪个Context，请传递 `context.TODO`
	- 使用上下文值仅适用于传输进程和 API 的请求范围数据，而不用于将可选参数传递给函数
	- 相同的上下文可以传递给在不同 goroutine 中运行的函数; 
		- 上下文对于多个 goroutine 同时使用是安全的。
```go
	var Canceled = errors.New("context canceled")

	var DeadlineExceeded error = deadlineExceededError{}

	// CancelFunc 通知操作放弃其工作，且不会等待工作停止，在第一次调用之后，对CancelFunc的后续调用不起作用
	type CancelFunc func()

	type Context interface {
		Deadline() (deadline time.Time, ok bool)
		Done() <-chan struct{}
		Err() error
		Value(key interface{}) interface{}
	}

	// Background 返回non-nil(非零)，空的 Context; 其从未被取消，没有值，也没有最后期限
	// 通常由主函数，初始化和测试使用，并作为传入请求的top-level Context (顶级上下文)
	func Background() Context

	// TODO 返回非零空的上下文，一般在不清楚使用哪个Context 或 它尚不可用时(因为周围的函数尚未扩展为接受Context参数)
	func TODO() Context

	// WithCancel 返回父上下文的副本，同时返回一个 CancelFunc
	// 取消这个上下文可以释放与它相关的资源，因此只要在这个Context 中运行的操作完成，代码就应该立即调用 cancel，下同
	func WithCancel(parent Context) (ctx Context, cancel CancelFunc)

	// WithDeadline 返回父上下文的副本，并将截止日期调整为不晚于deadline，同时返回一个 CancelFunc
	func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)

	// WithTimeout 返回 WithDeadline(parent, time.Now().Add(timeout))
	func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)

	// WithValue 返回父键的副本，其中与键关联的值是val
	// 提供的密钥必须具有可比性，不应该是字符串类型或任何其他内置类型，以避免使用上下文的包之间发生冲突
	func WithValue(parent Context, key, val interface{}) Context
```

- 示例
```go
	// context.WithCancel
	gen := func(ctx context.Context) <-chan int {
		dst := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case <-ctx.Done():
					return // returning not to leak the goroutine
				case dst <- n:
					n++
				}
			}
		}()
		return dst
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers
	for n := range gen(ctx) {
		fmt.Print(n)
		if n == 5 {
			break
		}
	}
	// Output: 
	// 1 2 3 4 5 

	// context.WithDeadline
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(50 * time.Millisecond)) // context.WithTimeout(context.Background(), 50 * time.Millisecond)
	defer cancel()
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("overslept")
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}
	// Output: 
	// context deadline exceeded

	// context.WithTimeout
	ctx, cancel := context.WithTimeout(context.Background(), 50 * time.Millisecond)
	defer cancel()
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("overslept")
	case <-ctx.Done():
		fmt.Println(ctx.Err()) // prints "context deadline exceeded"
	}
	// Output: 
	// context deadline exceeded

	// context.WithValue
	type favContextKey string
	f := func(ctx context.Context, k favContextKey) {
		if v := ctx.Value(k); v != nil {
			fmt.Println("found value:", v)
			return
		}
		fmt.Println("key not found:", k)
	}

	k := favContextKey("language")
	ctx := context.WithValue(context.Background(), k, "Go")
	f(ctx, k)
	f(ctx, favContextKey("color"))
	// Output:
	// found value: Go
	// key not found: color
```
