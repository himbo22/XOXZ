# ==================== Build Stage ====================
# Alpine is a good builder stage because the binary is built statically (CGO_ENABLED=0).
FROM golang:1.25-alpine AS builder

WORKDIR /build

# 1. Optimize CI/CD pipeline caching.
# If dependencies do not change, Docker reuses this layer.
COPY auth-svc/go.mod auth-svc/go.sum ./auth-svc/
COPY common-svc/kafka/go.mod common-svc/kafka/go.sum ./common-svc/kafka/

WORKDIR /build/auth-svc
RUN go mod download

# 2. Copy source code
WORKDIR /build
COPY . .

# 3. Build a static binary that does not depend on C libraries.
WORKDIR /build/auth-svc
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static'" \
    -tags netgo \
    -o /app-bin ./cmd/main.go

# ==================== Runtime Stage ====================
# Use the static image: very small (~2MB), zero attack surface, with TLS and timezone support.
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /

# 4. Copy the binary and configs with the correct ownership.
COPY --from=builder --chown=nonroot:nonroot /app-bin /auth-svc
COPY --from=builder --chown=nonroot:nonroot /build/auth-svc/configs /configs

# 5. Run as the preconfigured nonroot user.
USER nonroot:nonroot

EXPOSE 8080

CMD ["/auth-svc"]
