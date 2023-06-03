package metrics_test

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ahwhy/myGolang/prometheus/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/expfmt"
)

var (
	// Create global registry
	Registry = prometheus.NewRegistry()
)

// Go Process
func TestAll(t *testing.T) {
	Registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		metrics.NewDemoCollector(),
	)

	// Expose /metrics HTTP endpoint using the created custom registry.
	http.Handle("/metrics", promhttp.HandlerFor(Registry, promhttp.HandlerOpts{Registry: Registry}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Gauges 仪表盘
func TestGauges(t *testing.T) {
	queueLength := prometheus.NewGauge(prometheus.GaugeOpts{
		// Namespace, Subsystem, Name 会拼接成指标的名称: myGolang_prome_metrics_queue_length
		Namespace: "myGolang",
		Subsystem: "prome_metrics",
		Name:      "queue_length", // 必填参数
		// 指标的描信息
		Help: "The number of items in the queue.",
		// 指标的静态标签 
		ConstLabels: map[string]string{
			"module": "http-server",
		},
	})
	// 使用 Set() 设置指定的值
	queueLength.Set(1000)
	// 增加或减少
	queueLength.Inc()   // +1 gauge增加1
	queueLength.Dec()   // -1 gauge减少1
	queueLength.Add(66) // 增加66个增量
	queueLength.Sub(88) // 减少88

	queueLength2 := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		// Namespace, Subsystem, Name 会拼接成指标的名称: myGolang_prome_metrics_queue_length
		Namespace: "myGolang",
		Subsystem: "prome_metrics",
		Name:      "queue_length2", // 必填参数
		// 指标的描信息
		Help: "The number of items in the queue.",
		// 指标的静态标签 
		ConstLabels: map[string]string{
			"module": "http-server",
		},
	}, []string{"instance_id", "instance_name"})
	// 指标的动态标签 
	queueLength2.WithLabelValues("rm_001", "kafka01").Set(100)

	// 注册
	Registry.MustRegister(queueLength, queueLength2)

	// 获取注册所有数据
	data, err := Registry.Gather()
	if err != nil {
		panic(err)
	}

	// 编码输出
	enc := expfmt.NewEncoder(os.Stdout, expfmt.FmtText)
	for _, v := range data {
		fmt.Println(enc.Encode(v))
	}
}

// Counters 计数器
func TestCounters(t *testing.T) {
	totalRequests := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "The total number of handled HTTP requests.",
	})

	// Counters 只增不减
	for i := 0; i < 10; i++ {
		totalRequests.Inc()
	}
	totalRequests.Add(66)

	// 注册
	Registry.MustRegister(totalRequests)

	// 获取注册所有数据
	data, err := Registry.Gather()
	if err != nil {
		panic(err)
	}

	// 编码输出
	enc := expfmt.NewEncoder(os.Stdout, expfmt.FmtText)
	fmt.Println(enc.Encode(data[0]))
}

// Histograms 直方图
func TestHistograms(t *testing.T) {
	requestDurations := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "A histogram of the HTTP request durations in seconds.",
		// Bucket: 第一个 bucket 包括所有在 0.1s 内完成的请求，最后一个包括所有在 1.6s 内完成的请求。
		// 同 Buckets: []float64{0.1, 0.2, 0.4, 0.8, 1.6}
		Buckets: prometheus.ExponentialBuckets(0.0000001, 2, 10),
	})

	// Add go runtime metrics and process collectors.
	Registry.MustRegister(
		requestDurations,
	)

	go func() {
		for {
			// Record fictional latency.
			now := time.Now()
			requestDurations.(prometheus.ExemplarObserver).ObserveWithExemplar(
				time.Since(now).Seconds(), prometheus.Labels{"dummyID": fmt.Sprint(rand.Intn(1000))},
			)
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	// 获取注册所有数据
	data, err := Registry.Gather()
	if err != nil {
		panic(err)
	}

	// 编码输出
	enc := expfmt.NewEncoder(os.Stdout, expfmt.FmtText)
	fmt.Println(enc.Encode(data[0]))

	// Expose /metrics HTTP endpoint using the created custom registry.
	http.Handle(
		"/metrics", promhttp.HandlerFor(
			Registry,
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
			}),
	)
	// To test: curl -H 'Accept: application/openmetrics-text' localhost:8080/metrics
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

// Summaries 摘要
func TestSummaries(t *testing.T) {
	requestDurations := prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "http_request_duration_seconds",
		Help: "A summary of the HTTP request durations in seconds.",
		Objectives: map[float64]float64{
			0.5:  0.05,  // 第50个百分位数，最大绝对误差为0.05。
			0.9:  0.01,  // 第90个百分位数，最大绝对误差为0.01。
			0.99: 0.001, // 第90个百分位数，最大绝对误差为0.001。
		},
	})

	// 添加值
	for _, v := range []float64{0.01, 0.02, 0.3, 0.4, 0.6, 0.7, 5.5, 11} {
		requestDurations.Observe(v)
	}

	// 注册
	Registry.MustRegister(requestDurations)

	// 获取注册所有数据
	data, err := Registry.Gather()
	if err != nil {
		panic(err)
	}

	// 编码输出
	enc := expfmt.NewEncoder(os.Stdout, expfmt.FmtText)
	fmt.Println(enc.Encode(data[0]))
}
