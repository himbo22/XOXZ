package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConfig    `yaml:"server"`
	Logger  LoggerConfig    `yaml:"logger"`
	Mongo   MongoDBConfig   `yaml:"mongo"`
	Redis   RedisConfig     `yaml:"redis"`
	Kafka   KafkaConfig     `yaml:"kafka"`
	Otel    TelemetryConfig `yaml:"otel"`
	LiveKit LiveKitConfig   `yaml:"livekit"`
}

type MongoDBConfig struct {
	URI              string        `yaml:"uri"`
	DBName           string        `yaml:"dbname"`
	MaxPoolSize      uint64        `yaml:"max_pool_size"`
	MinPoolSize      uint64        `yaml:"min_pool_size"`
	MaxConnsIdleTime time.Duration `yaml:"max_conns_idle_time"`
}

type LiveKitConfig struct {
	APIKey    string        `yaml:"api_key"`
	APISecret string        `yaml:"api_secret"`
	Endpoint  string        `yaml:"endpoint"`
	Timeout   time.Duration `yaml:"timeout"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Address         string `yaml:"address"`
	Port            string `yaml:"port"`
	OpenAPIPath     string `yaml:"openapiPath"`
	SwaggerPath     string `yaml:"swaggerPath"`
	ErrorStack      bool   `yaml:"errorStack"`
	ErrorLogEnabled bool   `yaml:"errorLogEnabled"`
	ErrorLogPattern string `yaml:"errorLogPattern"`
	LogLevel        string // Derived from logger config
}

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Path                 string   `yaml:"path"`
	File                 string   `yaml:"file"`
	Prefix               string   `yaml:"prefix"`
	Level                string   `yaml:"level"`
	TimeFormat           string   `yaml:"timeFormat"`
	CtxKeys              []string `yaml:"ctxKeys"`
	Header               bool     `yaml:"header"`
	StSkip               int      `yaml:"stSkip"`
	Stdout               bool     `yaml:"stdout"`
	RotateSize           int      `yaml:"rotateSize"`
	RotateExpire         int      `yaml:"rotateExpire"`
	RotateBackupLimit    int      `yaml:"rotateBackupLimit"`
	RotateBackupExpire   int      `yaml:"rotateBackupExpire"`
	RotateBackupCompress int      `yaml:"rotateBackupCompress"`
	RotateCheckInterval  string   `yaml:"rotateCheckInterval"`
	StdoutColorDisabled  bool     `yaml:"stdoutColorDisabled"`
	WriterColorEnable    bool     `yaml:"writerColorEnable"`
	Flags                int      `yaml:"flags"`
}

type DatabaseLoggerConfig struct {
	Level  string `yaml:"level"`
	Stdout bool   `yaml:"stdout"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Default RedisDefaultConfig `yaml:"default"`
}

type RedisDefaultConfig struct {
	Address         string        `yaml:"address"`
	Password        string        `yaml:"password"`
	DB              int           `yaml:"db"`
	IdleTimeout     time.Duration `yaml:"idleTimeout"`
	MaxConnLifetime time.Duration `yaml:"maxConnLifetime"`
	WaitTimeout     time.Duration `yaml:"waitTimeout"`
	DialTimeout     time.Duration `yaml:"dialTimeout"`
	ReadTimeout     time.Duration `yaml:"readTimeout"`
	WriteTimeout    time.Duration `yaml:"writeTimeout"`
	MaxActive       int           `yaml:"maxActive"`
}

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	Brokers      []string          `yaml:"brokers"`
	RequiredAcks int               `yaml:"required_acks"`
	MaxRetry     int               `yaml:"max_retry"`
	Topics       KafkaTopicsConfig `yaml:"topics"`
}

// KafkaTopicsConfig holds Kafka topic names
type KafkaTopicsConfig struct {
	ActivityLogs string `yaml:"activity_logs"`
}

type TelemetryConfig struct {
	ServiceName    string  `mapstructure:"service_name" yaml:"service_name"`
	ServiceVersion string  `mapstructure:"service_version" yaml:"service_version"`
	Environment    string  `mapstructure:"environment" yaml:"environment"` // production | staging | development
	ExporterType   string  `mapstructure:"exporter_type" yaml:"exporter_type"`
	Endpoint       string  `mapstructure:"endpoint" yaml:"endpoint"`       // localhost:4317 or localhost:4318
	SampleRate     float64 `mapstructure:"sample_rate" yaml:"sample_rate"` // 1.0 = 100%, 0.1 = 10%
	Insecure       bool    `mapstructure:"insecure" yaml:"insecure"`       // disable TLS
}

var ApiKey string
var ApiSecret string

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	config.Server.LogLevel = config.Logger.Level

	// Load sensitive data from environment variables
	config.LiveKit.APIKey = getEnv("API_KEY", config.LiveKit.APIKey)
	config.LiveKit.APISecret = getEnv("API_SECRET", config.LiveKit.APISecret)

	ApiKey = getEnv("API_KEY", config.LiveKit.APIKey)
	ApiSecret = getEnv("API_SECRET", config.LiveKit.APISecret)

	return &config, nil
}

// getEnv gets environment variable with fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
