package configs

// import "time"

// type PostgresConfig struct {
// 	Connection    string        `yaml:"connection"`
// 	MaxConns      int32         `yaml:"max_conns"`
// 	LogLevel      string        `yaml:"log_level"` // must follow pgx.LogLevel format
// 	RetryCount    int           `yaml:"retry_count"`
// 	RetryInterval time.Duration `yaml:"retry_interval"`
// 	ShardID       int           `yaml:"shard_id"`

// 	MaxConnIdleTime time.Duration `yaml:"max_conn_idle_time"`
// }

import (
	"io/ioutil"
	"time"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

// Config structure
type (
	Config struct {
		AppEnv   string   `yaml:"app_env" envconfig:"APP_ENV"`
		Service  string   `yaml:"service" envconfig:"SERVICE"`
		LogLevel string   `yaml:"log_level" envconfig:"LOG_LEVEL"`
		Postgres Postgres `yaml:"postgres" envconfig:"POSTGRES"`
		Address  string   `yaml:"address" envconfig:"ADDRESS"`
	}
	Postgres struct {
		Username        string        `yaml:"username" envconfig:"PDB_USERNAME"`
		Password        string        `yaml:"password" envconfig:"PDB_PASSWORD"`
		Host            string        `yaml:"host" envconfig:"PDB_HOST"`
		Port            string        `yaml:"port" envconfig:"PDB_PORT"`
		DBName          string        `yaml:"db_name" envconfig:"PDB_DBNAME"`
		Connection      string        `yaml:"connection"`
		MaxConns        int32         `yaml:"max_conns"`
		LogLevel        string        `yaml:"log_level"` // must follow pgx.LogLevel format
		RetryCount      int           `yaml:"retry_count"`
		RetryInterval   time.Duration `yaml:"retry_interval"`
		MaxConnIdleTime time.Duration `yaml:"max_conn_idle_time"`
	}
)

// LoadConfigFile load default config from file
func LoadConfigFile(path string) Config {
	c := Config{}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return c
	}
	if err = yaml.Unmarshal(content, &c); err != nil {
		return c
	}
	return c
}

// LoadConfigEnv load config from environment variables
func LoadConfigEnv(c *Config) {
	if err := envconfig.Process("", c); err != nil {
		panic(err)
	}
}
