package cspmodel_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ahwhy/myGolang/concurrent/cspmodel"
)

func TestBarrierMode(t *testing.T) {
	cspmodel.BarrierMode()
}

func TestCallBackMode(t *testing.T) {
	cspmodel.CallBackMode()
}

func TestPipelineMode(t *testing.T) {
	cspmodel.PipelineMode()
}

func TestProducerConsumerMode(t *testing.T) {
	cspmodel.ProducerConsumerMode()
}

func TestPubSubMode(t *testing.T) {
	p := cspmodel.NewPublisher(100*time.Millisecond, 10)
	defer p.Close()

	// 订阅所有
	all := p.Subscribe()

	// 通过过滤订阅一部分信息
	golang := p.SubscribeTopic(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "golang")
		}
		return false
	})

	python := p.SubscribeTopic(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "python")
		}
		return false
	})

	// 发布者 发布信息
	p.Publish("hello,   python!")
	p.Publish("godbybe, python!")
	p.Publish("hello,   golang!")

	// 订阅者查看消息
	go func() {
		for msg := range all {
			fmt.Println("all:", msg)
		}
	}()

	// 订阅者查看消息
	go func() {
		for msg := range golang {
			fmt.Println("golang:", msg)
		}
	}()

	// 订阅者查看消息
	go func() {
		for msg := range python {
			fmt.Println("python:", msg)
		}
	}()

	// 运行一定时间后退出
	time.Sleep(3 * time.Second)
}
