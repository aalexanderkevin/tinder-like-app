package config

import (
	"sync"

	"github.com/jinzhu/configor"
)

type DBConfig struct {
	Client     string `default:"postgres" env:"POSTGRES_CLIENT"`
	Host       string `default:"127.0.0.1" env:"POSTGRES_HOST"`
	Username   string `default:"root" env:"POSTGRES_USER"`
	Password   string `required:"true" env:"POSTGRES_PASSWORD"`
	Port       uint   `default:"5432" env:"POSTGRES_PORT"`
	Database   string `default:"gits" env:"POSTGRES_DATABASE"`
	Migrations struct {
		Path string `default:"database/migrations" env:"POSTGRES_MIGRATION_PATH"`
	}
	MaxIdleConnections int  `default:"25" env:"POSTGRES_MAX_IDLE_CONN"`
	MaxOpenConnections int  `default:"0" env:"POSTGRES_MAX_OPEN_CONN"`
	MaxConnLifeTime    int  `default:"90" env:"POSTGRES_MAX_CONN_LIFETIME"`
	Debug              bool `default:"false" env:"POSTGRES_DEBUG"`
}

type RedisConfig struct {
	Host        string `default:"localhost" env:"REDIS_HOST"`
	Port        string `default:"6379" env:"REDIS_PORT"`
	MaxIdle     int    `default:"50" env:"REDIS_MAX_IDLE"`
	MaxActive   int    `default:"10000" env:"REDIS_MAX_ACTIVE"`
	IdleTimeout int    `default:"240" env:"REDIS_IDLE_TIMEOUT"`
}

type Config struct {
	Service struct {
		Host     string `default:"0.0.0.0" env:"SERVICE_HOST"`
		Port     string `default:"8080" env:"SERVICE_PORT"`
		LogLevel string `default:"DEBUG" env:"LOG_LEVEL"`
	}
	DB    DBConfig
	Redis RedisConfig
}

var config *Config
var configLock = &sync.Mutex{}

func Instance() Config {
	if config == nil {
		err := Load()
		if err != nil {
			panic(err)
		}
	}
	return *config
}

func Load() error {
	tmpConfig := Config{}
	err := configor.Load(&tmpConfig)
	if err != nil {
		return err
	}

	configLock.Lock()
	defer configLock.Unlock()
	config = &tmpConfig

	return nil
}
