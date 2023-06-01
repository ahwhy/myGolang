package conf

import (
	"github.com/caarlos0/env/v6"
)

var (
	global *Config
)

// C 全局配置对象
func C() *Config {
	if global == nil {
		panic("Load Config first")
	}

	return global
}

// LoadConfigFromEnv 从环境变量中加载配置
func LoadConfigFromEnv() error {
	global = NewDefaultConfig()

	return env.Parse(global)
}