package conf

func newConfig() *Config {
	return &Config{}
}

type Config struct {
	FEISHU_BOT_URL   string `env:"FEISHU_BOT_URL"`
	DINGDING_BOT_URL string `env:"DINGDING_BOT_URL"`
	WECHAT_BOT_URL   string `env:"WECHAT_BOT_URL"`

	DEPLOY_ID        string `env:"DEPLOY_ID"`
	BUILD_ID         string `env:"BUILD_ID"`
	MCENTER_BUILD_ID string `env:"MCENTER_BUILD_ID"`
	SERVICE_ID       string `env:"SERVICE_ID"`
	PIPELINE_TASK_ID string `env:"PIPELINE_TASK_ID"`

	DEPLOY_JOB_ID     string `env:"DEPLOY_JOB_ID"`
	BUILD_JOB_ID      string `env:"BUILD_JOB_ID"`
	CICD_PIPELINE_ID  string `env:"CICD_PIPELINE_ID"`
	MPAAS_PIPELINE_ID string `env:"MPAAS_PIPELINE_ID"`

	DEVCLOUD_DEPLOY_APPROVAL_ID string `env:"DEVCLOUD_DEPLOY_APPROVAL_ID"`
}
