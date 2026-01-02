package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

const (
	envStand       = "STAND"
	envDbMainPass  = "DB_MAIN_PASS"
	envDbRedisPass = "DB_REDIS_PASS"
)

const (
	devFileName  = "dev.yaml"
	prodFileName = "prod.yaml"

	dev  = "DEV"
	prod = "PROD"
)

type MainDatabase struct {
	Host                      string
	Port                      string
	Username                  string
	Password                  string
	Database                  string
	SslMode                   string `yaml:"ssl_mode"`
	MaxTransactionRetries     int    `yaml:"max_transaction_retries"`
	MaxOpenConnections        int    `yaml:"max_open_connections"`
	MaxIdleConnections        int    `yaml:"max_idle_connections"`
	MaxLifetimeConnectionsMin int    `yaml:"max_lifetime_connections"`
	MigrationFilePath         string `yaml:"migration_file_path"`
}

type RedisDatabase struct {
	Address  string `yaml:"address"`
	Database int    `yaml:"db"`
	Password string
}

type Config struct {
	Server struct {
		Host string
		Port string
	} `yaml:"server"`
	MainDatabase  *MainDatabase  `yaml:"psql"`
	RedisDatabase *RedisDatabase `yaml:"rediska"`
	DefaultClient struct {
		Timeout            int `yaml:"timeout"`
		MaxIdleConn        int `yaml:"max_idle_conn"`
		MaxIdleConnPerHost int `yaml:"max_idle_conn_per_host"`
		IdleConnTimeout    int `yaml:"idle_conn_timeout"`
	} `yaml:"default_client"`
	Services struct {
		Youkassa struct {
			URL    string `yaml:"url"`
			ShopID string
			Token  string
		}
		UserService struct {
			URL string `yaml:"url"`
		} `yaml:"user_service"`
		SubscriptionService struct {
			URL string `yaml:"url"`
		} `yaml:"subscription_service"`
		Front struct {
			PaymentLink string `yaml:"payment_link"`
		} `yaml:"front"`
	} `yaml:"services"`
}

func IsDevStand() bool {
	stand := os.Getenv(envStand)
	return stand == "" || stand == dev
}

var AppConfig *Config

func Load() *Config {
	var configFilePath string

	if IsDevStand() {
		configFilePath = devFileName
	} else {
		configFilePath = prodFileName
	}

	data, err := os.ReadFile("resources/config/" + configFilePath)
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("failed to parse config file: %v", err)
	}

	cfg.MainDatabase.Password = getEnv(envDbMainPass)
	cfg.RedisDatabase.Password = getEnv(envDbRedisPass)

	AppConfig = &cfg
	return AppConfig
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("environment variable '%s' is missing", key)
	}
	return value
}
