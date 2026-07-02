package config

import (
	"crypto/rsa"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.yaml.in/yaml/v3"
)

type Config struct {
	Server   ServerConfig    `yaml:"server"`
	Logger   LoggerConfig    `yaml:"logger"`
	Database DatabaseConfig  `yaml:"database"`
	Redis    RedisConfig     `yaml:"redis"`
	Kafka    KafkaConfig     `yaml:"kafka"`
	Auth     AuthConfig      `yaml:"auth"`
	Email    EmailConfig     `yaml:"email"`
	Otel     TelemetryConfig `yaml:"otel"`
	OAuth    OAuthConfig     `yaml:"oauth"`
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

// DatabaseConfig holds database configuration
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

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Default RedisDefaultConfig `yaml:"default"`
}

type RedisDefaultConfig struct {
	Address         string `yaml:"address"`
	Password        string `yaml:"password"`
	DB              int    `yaml:"db"`
	IdleTimeout     string `yaml:"idleTimeout"`
	MaxConnLifetime string `yaml:"maxConnLifetime"`
	WaitTimeout     string `yaml:"waitTimeout"`
	DialTimeout     string `yaml:"dialTimeout"`
	ReadTimeout     string `yaml:"readTimeout"`
	WriteTimeout    string `yaml:"writeTimeout"`
	MaxActive       int    `yaml:"maxActive"`
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

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JwtKeyId                 string `yaml:"jwtKeyId"` // Kong credential key (iss claim), e.g. bitzap-key
	AccessTokenExpireMinute  int    `yaml:"accessTokenExpireMinute"`
	RefreshTokenExpireMinute int    `yaml:"refreshTokenExpireMinute"`
	GracePeriodExpireSecond  int    `yaml:"gracePeriodExpireSecond"`
	GoogleClientId           string `yaml:"googleClientId"`
	RsaPublicKey             string `yaml:"rsaPublicKey"`
	RsaPrivateKey            string `yaml:"rsaPrivateKey"`
}

type EmailConfig struct {
	MailjetAPIKey    string `yaml:"mailjet_api_key" env:"MAILJET_API_KEY"`
	MailjetSecretKey string `yaml:"mailjet_secret_key" env:"MAILJET_SECRET_KEY"`
	FromEmail        string `yaml:"from_email" env:"FROM_EMAIL"`
	FromName         string `yaml:"from_name" env:"FROM_NAME"`
	AppURL           string `yaml:"app_url" env:"APP_URL"`
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

type OAuthConfig struct {
	Google OAuthProviderConfig `yaml:"google"`
	Github OAuthProviderConfig `yaml:"github"`
}

type OAuthProviderConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RedirectURL  string `yaml:"redirect_url"`
}

var (
	AppPrivateKey          *rsa.PrivateKey
	AppPublicKey           *rsa.PublicKey
	RefreshTokenExpiryTime time.Duration
	AccessTokenExpiryTime  time.Duration
)

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
	config.Email.MailjetAPIKey = getEnv("MAILJET_API_KEY", config.Email.MailjetAPIKey)
	config.Email.MailjetSecretKey = getEnv("MAILJET_SECRET_KEY", config.Email.MailjetSecretKey)
	config.Auth.GoogleClientId = getEnv("GOOGLE_CLIENT_ID", config.Auth.GoogleClientId)
	pubKeyStr := getEnv("RSA_PUBLIC_KEY", config.Auth.RsaPublicKey)
	privKeyStr := getEnv("RSA_PRIVATE_KEY", config.Auth.RsaPrivateKey)

	privKeyStr = strings.ReplaceAll(privKeyStr, "\\n", "\n")
	pubKeyStr = strings.ReplaceAll(pubKeyStr, "\\n", "\n")

	AppPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(privKeyStr))
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	AppPublicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(pubKeyStr))
	if err != nil {
		log.Fatalf("Failed to parse public key: %v", err)
	}

	if config.Auth.RefreshTokenExpireMinute == 0 {
		config.Auth.RefreshTokenExpireMinute = 60
	}
	if config.Auth.AccessTokenExpireMinute == 0 {
		config.Auth.AccessTokenExpireMinute = 15
	}
	// production
	RefreshTokenExpiryTime = time.Duration(config.Auth.RefreshTokenExpireMinute) * time.Minute
	AccessTokenExpiryTime = time.Duration(config.Auth.AccessTokenExpireMinute) * time.Minute

	// test: 1s
	//RefreshTokenExpiryTime = time.Duration(1) * time.Minute
	//AccessTokenExpiryTime = time.Duration(1) * time.Second

	return &config, nil
}

// getEnv gets environment variable with fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
