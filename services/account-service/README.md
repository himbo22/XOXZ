# package
go install github.com/google/wire/cmd/wire@latest
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/joho/godotenv/cmd/godotenv@latest

go get github.com/spf13/viper
go get gorm.io/gorm
go get go.uber.org/zap
go get github.com/labstack/echo/v5
go get github.com/swaggo/swag
go get github.com/redis/go-redis/v9
go get github.com/google/uuid
go get github.com/golang-jwt/jwt/v5
go get gorm.io/plugin/opentelemetry/tracing




/internal
  - /adapter -> third parties