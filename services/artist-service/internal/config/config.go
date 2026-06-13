package config

import (
	"os"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	Server   ServerConfig    `yaml:"server"`
	Logger   LoggerConfig    `yaml:"logger"`
	Database DatabaseConfig  `yaml:"database"`
	Kafka    KafkaConfig     `yaml:"kafka"`
	Otel     TelemetryConfig `yaml:"otel"`
}

type ServerConfig struct {
	Address         string `yaml:"address"`
	Port            string `yaml:"port"`
	OpenAPIPath     string `yaml:"openapiPath"`
	SwaggerPath     string `yaml:"swaggerPath"`
	ErrorStack      bool   `yaml:"errorStack"`
	ErrorLogEnabled bool   `yaml:"errorLogEnabled"`
	ErrorLogPattern string `yaml:"errorLogPattern"`
	LogLevel        string
}

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

type DatabaseConfig struct {
	Host            string               `yaml:"host"`
	Port            string               `yaml:"port"`
	User            string               `yaml:"user"`
	Password        string               `yaml:"password"`
	DBName          string               `yaml:"dbname"`
	SSLMode         string               `yaml:"sslmode"`
	MaxOpenConns    int                  `yaml:"max_open_conns"`
	MaxIdleConns    int                  `yaml:"max_idle_conns"`
	ConnMaxLifetime int                  `yaml:"conn_max_lifetime"`
	Logger          DatabaseLoggerConfig `yaml:"logger"`
	Debug           bool                 `yaml:"debug"`
	Timezone        string               `yaml:"timezone"`
}

type DatabaseLoggerConfig struct {
	Level  string `yaml:"level"`
	Stdout bool   `yaml:"stdout"`
}

type KafkaConfig struct {
	Brokers      []string          `yaml:"brokers"`
	RequiredAcks int               `yaml:"required_acks"`
	MaxRetry     int               `yaml:"max_retry"`
	Topics       KafkaTopicsConfig `yaml:"topics"`
}

type KafkaTopicsConfig struct {
	ActivityLogs string `yaml:"activity_logs"`
}

type TelemetryConfig struct {
	ServiceName    string  `mapstructure:"service_name" yaml:"service_name"`
	ServiceVersion string  `mapstructure:"service_version" yaml:"service_version"`
	Environment    string  `mapstructure:"environment" yaml:"environment"`
	ExporterType   string  `mapstructure:"exporter_type" yaml:"exporter_type"`
	Endpoint       string  `mapstructure:"endpoint" yaml:"endpoint"`
	SampleRate     float64 `mapstructure:"sample_rate" yaml:"sample_rate"`
	Insecure       bool    `mapstructure:"insecure" yaml:"insecure"`
}

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
	config.Database.Host = getEnv("DB_HOST", config.Database.Host)
	config.Database.Port = getEnv("DB_PORT", config.Database.Port)
	config.Database.User = getEnv("DB_USER", config.Database.User)
	config.Database.Password = getEnv("DB_PASSWORD", config.Database.Password)
	config.Database.DBName = getEnv("DB_NAME", config.Database.DBName)
	config.Database.SSLMode = getEnv("DB_SSL_MODE", config.Database.SSLMode)

	return &config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
