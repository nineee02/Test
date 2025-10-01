package config

import (
	"errors"
	"time"

	"github.com/spf13/viper"
)

type Configuration struct {
	Server ServerConfig `mapstructure:"server"`
	Logger Logger       `mapstructure:"logger"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
	API    APIConfig    `mapstructure:"api"`
}

type ServerConfig struct {
	Address      string   `mspstructure:"address"`
	Port         string   `mapstructure:"port"`
	Host         string   `mapstructure:"host"`
	ServerURL    string   `mapstructure:"server_url"`
	AllowHeaders []string `mapstructure:"allow_headers"`
	AllowMethods []string `mapstructure:"allow_methods"`
}

type MySQLConfig struct {
	Host            string        `mapstructure:"host"`
	Port            string        `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DBname          string        `mapstructure:"db_name"`
	TablePrefix     string        `mapstructure:"table_prefix"`
	Debug           bool          `mapstructure:"debug"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type APIConfig struct {
	NotificationUrl string `mapstructure:"notification_url"`
	AesKey          string `mapstructure:"aes_key"`
	SoftKey         string `mapstructure:"soft_key"`
	JwtSecret       string `mapstructure:"jwt_secret"`
}

type Logger struct {
	Development       bool   `mapstructure:"development"`
	DisableCaller     bool   `mapstructure:"disable_caller"`
	DisableStacktrace bool   `mapstructure:"disable_stacktrace"`
	Encoding          string `mapstructure:"encoding"`
	Level             string `mapstructure:"level"`
	Path              string `mapstructure:"path"`
}

func New() (*Configuration, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("config/")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	var cfg Configuration
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
