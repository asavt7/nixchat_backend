package config

import (
	"github.com/spf13/viper"
	"strings"
	"time"
)

type (
	Config struct {
		Logger   LoggerConfig   `mapstructure:"logger"`
		Auth     AuthConfig     `mapstructure:"auth"`
		Postgres PostgresConfig `mapstructure:"pg"`
		HTTP     HTTPConfig     `mapstructure:"http"`
		Redis    RedisConfig    `mapstructure:"redis"`
	}
	LoggerConfig struct {
		Level string `mapstructure:"level"`
	}
	AuthConfig struct {
		AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
		SigningKey      string        `mapstructure:"jwtSigningKey"`
	}
	HTTPConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}
	PostgresConfig struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"db"`
		SSLMode  string `mapstructure:"sslmode"`
	}
	RedisConfig struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	}
)

// Init populates Config struct with values from config file and env variables
// Env vars takes precedence over configs from files
func Init(pathToConfig string) (*Config, error) {

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(pathToConfig)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	if err = bindingEnvs(); err != nil {
		return nil, err
	}

	cfg, err := unmarshalConfigs()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func bindingEnvs() error {
	err := viper.BindEnv("auth.jwtSigningKey", "AUTH_JWR_SIGHING_KEY", "AUTH_JWRSIGHINGKEY")
	if err != nil {
		return err
	}
	return nil
}

func unmarshalConfigs() (*Config, error) {
	var cfg Config

	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, err
	/*
		err := viper.UnmarshalKey("pg", &cfg.Postgres)
		if err != nil {
			return &cfg, err
		}
		err = viper.UnmarshalKey("logger", &cfg.Logger)
		if err != nil {
			return &cfg, err
		}
		err = viper.UnmarshalKey("http", &cfg.HTTP)
		if err != nil {
			return &cfg, err
		}
		err = viper.UnmarshalKey("auth", &cfg.Auth)
		if err != nil {
			return &cfg, err
		}
		err = viper.UnmarshalKey("redis", &cfg.Redis)
		if err != nil {
			return &cfg, err
		}

		return &cfg, nil */
}
