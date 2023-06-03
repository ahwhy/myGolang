# Golang-Prometheus  Golang的Prometheus

## 一、Exporter

### 1. 数据格式
- Prometheus是拉取数据的监控模型, 它对客户端暴露的数据格式要求如下
```
	# HELP go_goroutines Number of goroutines that currently exist.
	# TYPE go_goroutines gauge
	go_goroutines 16
```

### 2. 简单实现
- simple implementation
```golang
	func HelloHandler(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "lexporter_request_count{user=\"admin\"} 1000" )
	}
	http.HandleFunc("/metrics", HelloHandler)
	http.ListenAndServe(":8050", nil)
```

### 3. 使用SDK
- promhttp工具包
	- `github.com/prometheus/client_golang/prometheus/promhttp`
	- 通过访问 http://127.0.0.1:8050/metrics，获得默认的监控指标数据
	- Go 客户端库默认在暴露的全局默认指标注册表中注册了一些关于 promhttp 处理器和运行时间相关的默认指标
		- go_*: 以 go_ 为前缀的指标是关于 Go 运行时相关的指标，比如垃圾回收时间、goroutine 数量等，这些都是 Go 客户端库特有的，其他语言的客户端库可能会暴露各自语言的其他运行时指标
		- promhttp_*: 来自 promhttp 工具包的相关指标，用于跟踪对指标请求的处理
```golang
	// Serve the default Prometheus metrics registry over HTTP on /metrics.
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8050", nil)
```

- 自定义指标
	- Prometheus的Server端，需要上述数据格式
	- Prometheus的Client端本身也提供一些简单数据二次加工的能力，这种能力被描述为4种指标类型
		- Gauges(仪表盘): Gauge类型代表一种样本数据可以任意变化的指标，即可增可减
		- Counters(计数器): counter类型代表一种样本数据单调递增的指标，即只增不减，除非监控系统发生了重置
		- Histograms(直方图): 创建直方图指标比 counter 和 gauge 都要复杂，因为需要配置把观测值归入的 bucket 的数量，以及每个 bucket 的上边界
			- Prometheus 中的直方图是累积的，所以每一个后续的 bucket 都包含前一个 bucket 的观察计数，所有 bucket 的下限都从 0 开始的，所以我们不需要明确配置每个 bucket 的下限，只需要配置上限即可
		- Summaries(摘要): 与Histogram类似类型，用于表示一段时间内的数据采样结果(通常是请求持续时间或响应大小等)，但它直接存储了分位数(通过客户端计算，然后展示出来)，而不是通过区间计算
	- [ Metrics_Demo ](../prometheus/metrics/metrics_test.go)

- 指标标签
	- Prometheus将指标的标签分为2类
		- 静态标签: constLabels，在指标创建时，就提前声明好，采集过程中永不变动
		- 动态标签: variableLabels，用于在指标的收集过程中动态补充标签，比如kafka集群的exporter 需要动态补充 instance_id
			- 支持动态标签的构造函数 `NewGaugeVec()`、 `NewCounterVec()`、`NewSummaryVec()`、`NewHistogramVec()`
			- `func NewGaugeVec(opts GaugeOpts, labelNames []string) *GaugeVec`

- 指标注册
	- 指标采集完成后，需要注册给Prometheus的Http Handler才能将指标暴露出去
	- Prometheus 定义了一个注册表的接口
```golang
	// 指标注册接口
	type Registerer interface {
		// 注册采集器, 有异常会报错
		Register(Collector) error
		// 注册采集器，有异常会panic
		MustRegister(...Collector)
		// 注销该采集器
		Unregister(Collector) bool
	}
```

- 默认注册表
	- Prometheus 实现了一个默认的Registerer对象，也就是默认注册表
	- 通过提供的MustRegister可以将自定义指标注册进去
```golang
	var (
		defaultRegistry              = NewRegistry()
		DefaultRegisterer Registerer = defaultRegistry
		DefaultGatherer   Gatherer   = defaultRegistry
	)

	// 在默认的注册表中注册temp指标
	prometheus.MustRegister(temp)
	prometheus.Register()
	prometheus.Unregister()
```

- 自定义注册表
	- Prometheus 默认的Registerer对象，会添加一些默认指标的采集，比如go运行时和当前process相关信息
		- 通过 `NewProcessCollector()`、`NewGoCollector()` 控制是否添加
	- 使用自定义的注册表的方式，控制暴露的指标
		- 使用`NewRegistry()`创建一个全新的注册表
		- 通过注册表对象的`MustRegister()`把指标注册到自定义的注册表中
	- 暴露指标的时候必须通过调用 `promhttp.HandleFor()` 函数来创建一个专门针对自定义注册表的 HTTP 处理器
		- 并且需要在 `promhttp.HandlerOpts` 配置对象的 `Registry` 字段中传递注册表对象
```golang
	...
	// 添加 process 和 Go 运行时指标到自定义的注册表中
	registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	registry.MustRegister(prometheus.NewGoCollector())
	...
	// 暴露指标
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}))
	http.ListenAndServe(":8050", nil)
```

- 采集器
	-  为了能更好的模块化，很多第三方Collecotor的标准写法是，需要把指标采集封装为一个Collector对象
	- [ Metrics_Collector ](../prometheus/metrics/metrics_collector.go)
```golang
	type Collector interface {
		// 指标的一些描述信息, 就是# 标识的那部分
		// 注意这里使用的是指针, 因为描述信息 全局存储一份就可以了
		Describe(chan<- *Desc)
		// 指标的数据, 比如 promhttp_metric_handler_errors_total{cause="gathering"} 0
		// 这里没有使用指针, 因为每次采集的值都是独立的
		Collect(chan<- Metric)
	}
```


```golang
```