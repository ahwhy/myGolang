# Golang-Sync  Golang的同步基元

## 一、Golang的标准库 strconv包

### 1. strconv
- sync
	- sync包提供了基本的同步基元，如互斥锁
	- 除了Once和WaitGroup类型，大部分都是适用于低水平程序线程，高水平的同步使用channel通信更好一些
	- 本包的类型的值不应被拷贝

### 2. Once
```go
	// Once 是只执行一次动作的对象
	type Once struct { ... }
	// Do 方法当且仅当第一次被调用时才执行函数f
	// 给定变量 `var once Once`
	// 如果once.Do(f)被多次调用，只有第一次调用会执行f，即使f每次调用Do 提供的f值不同；需要给每个要执行仅一次的函数都建立一个Once类型的实例
	// Do用于必须刚好运行一次的初始化；因为f是没有参数的，因此可能需要使用闭包来提供给Do方法调用 `config.once.Do(func() { config.init(filename) })`
	// 因为只有f返回后Do方法才会返回，f若引起了Do的调用，会导致死锁
	func (o *Once) Do(f func())

	// For Example
	var once sync.Once
	onceBody := func() {
		fmt.Println("Only once")
	}
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			once.Do(onceBody)
			done <- true
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
	// Output: Only once
```

### 3. Mutex
```go
	// Locker 接口代表一个可以加锁和解锁的对象
	type Locker interface {
		Lock()
		Unlock()
	}

	// Mutex 是一个互斥锁，可以创建为其他结构体的字段；零值为解锁状态
	// Mutex类型的锁和线程无关，可以由不同的线程加锁和解锁
	type Mutex struct { ... }
	// Lock 方法锁住m，如果m已经加锁，则阻塞直到m解锁
	func (m *Mutex) Lock()
	// Unlock 方法解锁m，如果m未加锁会导致运行时错误
	func (m *Mutex) Unlock()

	// RWMutex 是读写互斥锁；该锁可以被同时多个读取者持有或唯一个写入者持有
	// RWMutex可以创建为其他结构体的字段；零值为解锁状态
	// RWMutex类型的锁也和线程无关，可以由不同的线程加读取锁/写入和解读取锁/写入锁。
	type RWMutex struct { ... }
	// Lock 方法将rw锁定为写入状态，禁止其他线程读取或者写入
	func (rw *RWMutex) Lock()
	// Unlock 方法解除rw的写入锁状态，如果m未加写入锁会导致运行时错误
	func (rw *RWMutex) Unlock()
	// RLock 方法将rw锁定为读取状态，禁止其他线程写入，但不禁止读取
	func (rw *RWMutex) RLock()
	// Runlock 方法解除rw的读取锁状态，如果m未加读取锁会导致运行时错误
	func (rw *RWMutex) RUnlock()
	// Rlocker 方法返回一个互斥锁，通过调用rw.Rlock和rw.Runlock实现了Locker接口
	func (rw *RWMutex) RLocker() Locker
```

### 4. Cond
```go
	// Cond实现了一个条件变量，一个线程集合地，供线程等待或者宣布某事件的发生
	// 每个Cond实例都有一个相关的锁(一般是*Mutex或*RWMutex类型的值)，它必须在改变条件时或者调用Wait方法时保持锁定
	// Cond可以创建为其他结构体的字段，Cond在开始使用后不能被拷贝
	type Cond struct {
		// 在观测或更改条件时L会冻结
		L Locker
		...
	}
	// NewCond 使用锁l创建一个*Cond
	func NewCond(l Locker) *Cond
	// Broadcast 唤醒所有等待c的线程；调用者在调用本方法时，建议(但并非必须)保持c.L的锁定
	func (c *Cond) Broadcast()
	// Signal唤醒等待c的一个线程(如果存在)
	func (c *Cond) Signal()
	// Wait自行解锁c.L并阻塞当前线程，在之后线程恢复执行时，Wait方法会在返回前锁定c.L；和其他系统不同，Wait除非被Broadcast或者Signal唤醒，不会主动返回
	func (c *Cond) Wait()
```

### 5. WaitGroup
```go
	// WaitGroup 用于等待一组线程的结束
	// 父线程调用Add方法来设定应等待的线程的数量，每个被等待的线程在结束时应调用Done方法
	// 同时，主线程里可以调用Wait方法阻塞至所有线程结束
	type WaitGroup struct { ... }
	// Add方法向内部计数加上delta，delta可以是负数；如果内部计数器变为0，Wait方法阻塞等待的所有线程都会释放，如果计数器小于0，方法panic
	// 注意Add加上正数的调用应在Wait之前，否则Wait可能只会等待很少的线程
	// 一般来说本方法应在创建新的线程或者其他应等待的事件之前调用
	func (wg *WaitGroup) Add(delta int)
	// Done 方法减少WaitGroup计数器的值，应在线程的最后执行
	func (wg *WaitGroup) Done()
	// Wait 方法阻塞直到WaitGroup计数器减为0
	func (wg *WaitGroup) Wait()
```

### 6. Pool
```go
	// Pool是一个可以分别存取的临时对象的集合
	// Pool中保存的任何item都可能随时不做通告的释放掉；如果Pool持有该对象的唯一引用，这个item就可能被回收
	// Pool可以安全的被多个线程同时使用
	// Pool的目的是缓存申请但未使用的item用于之后的重用，以减轻GC的压力；也就是说，让创建高效而线程安全的空闲列表更容易；但Pool并不适用于所有空闲列表
	type Pool struct {
		// 可选参数New指定一个函数在Get方法可能返回nil时来生成一个值
		// 该参数不能在调用Get方法时被修改
		New func() interface{}
		...
	}
	// Get 方法从池中选择任意一个item，删除其在池中的引用计数，并提供给调用者
	// Get方法也可能选择无视内存池，将其当作空的；调用者不应认为Get的返回这和传递给Put的值之间有任何关系
	// 假使Get方法没有取得item：如p.New非nil，Get返回调用p.New的结果；否则返回nil
	func (p *Pool) Get() interface{}
	// Put 方法将x放入池中
	func (p *Pool) Put(x interface{})
```

## 二、Golang的标准库 atomic包

### 1. sync/atomic
- sync/atomic
	- atomic包提供了底层的原子级内存操作，对于同步算法的实现很有用
	- 应通过通信来共享内存，而不通过共享内存实现通信
