package conf

func newConfig() *Config {
	return &Config{
		App:    newDefaultAPP(),
		Script: newDefaultScript(),
	}
}

type Config struct {
	App    *app    `toml:"app"`
	Script *script `toml:"script"`
}

func newDefaultAPP() *app {
	return &app{
		Name: "yCmdb",
		HTTP: newDefaultHTTP(),
	}
}

type app struct {
	Name string `toml:"name" env:"APP_NAME"`
	HTTP *http  `toml:"http"`
}

func newDefaultHTTP() *http {
	return &http{
		Host: "127.0.0.1",
		Port: "8010",
	}
}

type http struct {
	Host      string `toml:"host" env:"HTTP_HOST"`
	Port      string `toml:"port" env:"HTTP_PORT"`
	EnableSSL bool   `toml:"enable_ssl" env:"HTTP_ENABLE_SSL"`
	CertFile  string `toml:"cert_file" env:"HTTP_CERT_FILE"`
	KeyFile   string `toml:"key_file" env:"HTTP_KEY_FILE"`
}

func (a *http) Addr() string {
	return a.Host + ":" + a.Port
}

func newDefaultScript() *script {
	return &script{
		HomeDir: "modules",
	}
}

type script struct {
	HomeDir string `toml:"home_dir" env:"SCRIPT_HOME_DIR"`
}
