FROM golang:1.25-alpine AS builder

WORKDIR /build
COPY artist-service/go.mod artist-service/go.sum ./artist-service/
COPY common-service ./common-service/

WORKDIR /build/artist-service
RUN go mod download

WORKDIR /build
COPY artist-service ./artist-service/

WORKDIR /build/artist-service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static'" \
    -tags netgo \
    -o /app-bin ./cmd/main.go

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /
COPY --from=builder --chown=nonroot:nonroot /app-bin /artist-svc
COPY --from=builder --chown=nonroot:nonroot /build/artist-service/configs /configs

USER nonroot:nonroot
EXPOSE 8080

CMD ["/artist-svc"]
