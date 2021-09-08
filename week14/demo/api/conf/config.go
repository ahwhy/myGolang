package conf

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// app 配置项
type app struct {
	Name string `toml:"name" env:"APP_NAME"`
	Host string `toml:"host" env:"APP_HOST"`
	Port string `toml:"port" env:"APP_PORT"`
}

func newDefaultAPP() *app {
	return &app{
		Name: "demo",
		Host: "127.0.0.1",
		Port: "8050",
	}
}

func (a *app) Addr() string {
	return fmt.Sprintf("%s:%s", a.Host, a.Port)
}

// mysql 配置项
// MySQL todo
type mysql struct {
	Host        string `toml:"host" env:"D_MYSQL_HOST"`
	Port        string `toml:"port" env:"D_MYSQL_PORT"`
	UserName    string `toml:"username" env:"D_MYSQL_USERNAME"`
	Password    string `toml:"password" env:"D_MYSQL_PASSWORD"`
	Database    string `toml:"database" env:"D_MYSQL_DATABASE"`
	MaxOpenConn int    `toml:"max_open_conn" env:"D_MYSQL_MAX_OPEN_CONN"`
	MaxIdleConn int    `toml:"max_idle_conn" env:"D_MYSQL_MAX_IDLE_CONN"`
	MaxLifeTime int    `toml:"max_life_time" env:"D_MYSQL_MAX_LIFE_TIME"`
	MaxIdleTime int    `toml:"max_idle_time" env:"D_MYSQL_MAX_idle_TIME"`
	lock        sync.Mutex
}

// newDefaultMySQL todo
func newDefaultMySQL() *mysql {
	return &mysql{
		Database:    "go_course",
		Host:        "127.0.0.1",
		Port:        "3306",
		MaxOpenConn: 200,
		MaxIdleConn: 50,
		MaxLifeTime: 1800,
		MaxIdleTime: 600,
	}
}

var (
	db *sql.DB
)

// getDBConn use to get db connection pool
func (m *mysql) getDBConn() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&multiStatements=true", m.UserName, m.Password, m.Host, m.Port, m.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to mysql<%s> error, %s", dsn, err.Error())
	}

	db.SetMaxOpenConns(m.MaxOpenConn)
	db.SetMaxIdleConns(m.MaxIdleConn)
	db.SetConnMaxLifetime(time.Second * time.Duration(m.MaxLifeTime))
	db.SetConnMaxIdleTime(time.Second * time.Duration(m.MaxIdleTime))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping mysql<%s> error, %s", dsn, err.Error())
	}

	return db, nil
}

// GetDB todo
func (m *mysql) GetDB() (*sql.DB, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// 加载全局数据库单例
	if db == nil {
		conn, err := m.getDBConn()
		if err != nil {
			return nil, err
		}
		db = conn
	}

	return db, nil
}

// log 配置项
type log struct {
	Level   string    `toml:"level" env:"LOG_LEVEL"`
	PathDir string    `toml:"path_dir" env:"LOG_PATH_DIR"`
	Format  LogFormat `toml:"format" env:"LOG_FORMAT"`
	To      LogTo     `toml:"to" env:"LOG_TO"`
}

// newDefaultLog todo
func newDefaultLog() *log {
	return &log{
		Level:  "debug",
		Format: "text",
		To:     "stdout",
	}
}

// Config 应用配置
type Config struct {
	// app 配置项
	App *app `toml:"app"`

	// mysql 配置项
	MySQL *mysql `toml:"mysql"`

	// log 配置项
	Log *log `toml:"log"`
}

func newConfig() *Config {
	return &Config{
		App:   newDefaultAPP(),
		Log:   newDefaultLog(),
		MySQL: newDefaultMySQL(),
	}
}
