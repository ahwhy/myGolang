package conf

func NewDefaultConfig() *Config {
	return &Config{
		Mode: CMD_MODE,
		CmdConfig: &cmdConfig{
			CmdPath:  "/bin/mqadmin",
			Template: CMD_TEMP,
		},
		FileConfig: &fileConfig{
			Path: "modules/rocketmq/sample/data.txt",
		},
	}
}

type Config struct {
	Mode MODE `toml:"mode"`

	CmdConfig  *cmdConfig  `toml:"cmd"`
	FileConfig *fileConfig `toml:"file"`
}

type cmdConfig struct {
	CmdPath  string `toml:"path" env:"CMD_PATH"`
	Target   string `toml:"target" env:"CMD_TARGET"`
	Template string `toml:"template" env:"CMD_TEMPATE"`
}

type fileConfig struct {
	Path string `toml:"path" env:"APP_NAME"`
}

type MODE string

const (
	// 通过执行命令获取数据
	CMD_MODE MODE = "cmd"
	// 通过文件获取数据
	FILE_MODE MODE = "file"
)

var (
	CMD_TEMP = "%s  consumerProgress -n %s |grep -v OFFLINE|egrep -v group_con_gps_sync_test|egrep -v group_con_gps_storage_trace"
)