package config

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Cors     CorsConfig
	Password PasswordConfig
	Logger   LoggerConfig
	Otp      OtpConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port    string
	RunMode string
	Domain  string
}

type LoggerConfig struct {
	FilePath string
	Encoding string
	Level    string
	Logger   string
}

type PostgresConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DbName          string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}
type RedisConfig struct {
	Host               string
	Port               string
	Password           string
	Db                 string
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	IdleCheckFrequency time.Duration
	PoolSize           int
	PoolTimeout        time.Duration
}
type CorsConfig struct {
	AllowOrigins string
}

type PasswordConfig struct {
	IncludeChars     bool
	IncludeDigits    bool
	MinLength        int
	MaxLength        int
	IncludeUppercase bool
	IncludeLowercase bool
}

type OtpConfig struct {
	ExpireTime time.Duration
	Digits     int
	Limiter    time.Duration
}

type JWTConfig struct {
	AccessTokenExpireDuration  time.Duration
	RefreshTokenExpireDuration time.Duration
	Secret                     string
	RefreshSecret              string
}

func GetConfig() *Config {
	cfgPath := getConfigPath(os.Getenv("APP_ENV"))
	v, err := loadConfig(cfgPath, "yml")
	if err != nil {
		log.Fatalf("Error in load config %v", err)
	}

	cfg, err := parseConfig(v)
	if err != nil {
		log.Fatalf("Error in parse config %v", err)
	}

	return cfg
}

func parseConfig(v *viper.Viper) (*Config, error) {
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Printf("Unable parse config %v", err)
		return nil, err
	}

	return &cfg, nil
}

func loadConfig(fileName, fileType string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType(fileType)
	v.SetConfigName(fileName)
	v.AddConfigPath("../config")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if err == err.(viper.ConfigFileNotFoundError) {
			return nil, errors.New("Config file not found")
		}
		return nil, err
	}

	return v, nil
}

func getConfigPath(env string) string {
	if env == "docker" {
		return "config-docker.yml"
	} else if env == "production" {
		return "config-production.yml"
	} else {
		return "config-development.yml"
	}
}
