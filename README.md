# XOXZ

A **Weverse.io clone** built with Go microservices — a fan community platform where artists connect with their fans through profiles, media, livestreaming, and community features.

## Architecture

```
                         ┌─────────────────┐
                         │     Kong API     │
                         │    Gateway       │
                         └────────┬────────┘
                                  │
              ┌───────────────────┼───────────────────┐
              │                   │                   │
     ┌────────▼────────┐ ┌───────▼────────┐ ┌────────▼────────┐
     │  Account Service │ │  Artist Service │ │  Media Service  │
     │  :8080           │ │  :8090          │ │  :8180 / :50051 │
     │  Auth, Profiles, │ │  Artist/Group   │ │  MinIO Storage, │
     │  Roles, Perms    │ │  Profiles       │ │  Presigned URLs │
     └────────┬─────────┘ └────────────────┘ └────────┬────────┘
              │                                       │
              │                              ┌────────▼────────┐
              │                              │ Livestream      │
              │                              │ Service         │
              │                              │ :8280           │
              │                              │ LiveKit Rooms,  │
              │                              │ Stream Lifecycle│
              │                              └─────────────────┘
              │
     ┌────────▼───────────────────────────────────────┐
     │         Common Service (shared modules)         │
     │  Telemetry  │  Protobuf  │  Logger  │  Helpers  │
     └──────────────────────────────────────────────────┘
```

## Services

| Service | Module | Port | Dependencies | Weverse Feature |
|---|---|---|---|---|
| **Account** | `services/account-service` | `8080` | PostgreSQL, Redis, Google OAuth, RSA JWT | User sign-up/sign-in, profiles, roles & permissions |
| **Artist** | `services/artist-service` | `8090` | PostgreSQL | Artist & group profiles, discography |
| **Media** | `services/media-service` | `8180`, gRPC `50051` | MinIO (S3) | Photo/video uploads, CDN delivery, presigned URLs |
| **Livestream** | `services/livestream-service` | `8280` | MongoDB, Redis, LiveKit | Live broadcasting, WebRTC streaming, viewer count |
| **Community** (planned) | — | — | PostgreSQL, Kafka | Posts, comments, likes, fan community feed |
| **Notifications** (planned) | — | — | Kafka, FCM/APNs | Push notifications, in-app alerts |
| **Search** (planned) | — | — | Elasticsearch | Artist, post, and user search |
| **Shop** (planned) | — | — | PostgreSQL | Merchandise, albums, ticketing |
| **Common** | `services/common-service/*` | — | — | Shared telemetry, protobuf, logger, helpers |

## Repository Layout

```
app/                           Frontend playground
docs/                          Architecture notes
infra/                         Kong, Nginx, OpenTelemetry, Grafana, LiveKit config
services/account-service       User authentication, sessions, profiles, roles, permissions
services/artist-service        Artist and group profile APIs
services/media-service         Presigned uploads, MinIO storage, media gRPC API
services/livestream-service    LiveKit rooms, stream lifecycle, webhook handling
services/common-service        Shared Go modules (telemetry, protobuf, logger, helpers)
docker-compose.yaml            Infrastructure stack (PostgreSQL, Redis, Kafka, MinIO, etc.)
```

## Local Infrastructure

`docker-compose.yaml` runs the infrastructure stack. Application services are started from their own service directories.

### Start Infrastructure

```sh
docker compose up -d account-postgres livestream-mongo redis redisinsight kafka kafka-ui minio minio-setup cdn-nginx kong otel-collector prometheus tempo loki grafana livekit ingress
```

### Useful Local URLs

| Tool | URL | Purpose |
|---|---|---|
| Kong API Gateway | `http://localhost:8000` | Entry point for all API requests |
| Kong Admin API | `http://localhost:8001` | Gateway configuration |
| RedisInsight | `http://localhost:5540` | Redis GUI |
| Kafka UI | `http://localhost:8010` | Kafka management |
| MinIO S3 API | `http://localhost:9000` | Object storage API |
| MinIO Console | `http://localhost:9001` | MinIO management UI |
| CDN/Nginx | `http://localhost:8030` | Media delivery (API) |
| CDN/Nginx UI | `http://localhost:8031` | Media delivery (UI) |
| Grafana | `http://localhost:3333` | Observability dashboard |
| Prometheus | `http://localhost:9090` | Metrics |
| Tempo | `http://localhost:3200` | Distributed tracing |
| Loki | `http://localhost:3100` | Log aggregation |
| LiveKit | `ws://localhost:7880` | WebRTC media server |

## Run a Service

Each service is an independent Go module:

```sh
cd services/account-service
go mod download
go run ./cmd/main.go
```

### Common Commands

```sh
make run          # Start the service
make test         # Run tests
make wire         # Generate dependency injection
make swagger      # Generate Swagger docs
make migrate-up   # Run database migrations
make migrate-down # Rollback migrations
```

### Seed Data

```sh
cd services/account-service
go run ./cmd/seeder/main.go
```

## API Routes

### Internal (all services)

```
GET  /api/v1/internal/ping
GET  /api/v1/internal/health
GET  /api/v1/internal/swagger/*
```

### Account — Auth & Profile

```
POST   /api/v1/public/auth/google           # Google OAuth login
POST   /api/v1/public/auth/session/refresh  # Refresh access token
POST   /api/v1/public/auth/session/logout   # Logout
POST   /api/v1/public/auth/session/revoke-all  # Revoke all sessions
GET    /api/v1/public/profile/:username      # Get user profile
GET    /api/v1/public/profile/me             # Get my profile
PUT    /api/v1/public/profile/me             # Update my profile
PUT    /api/v1/public/profile/me/avatar      # Update avatar
POST   /api/v1/public/artist/account         # Register as artist
DELETE /api/v1/public/artist/account         # Remove artist role
```

### Artist — Profiles & Groups

```
POST   /api/v1/public/artists               # Create artist
GET    /api/v1/public/artists                # List artists
GET    /api/v1/public/artists/:id            # Get artist details
PUT    /api/v1/public/artists/:id            # Update artist
DELETE /api/v1/public/artists/:id            # Delete artist
```

### Media — Uploads

```
POST   /api/v1/public/create-presigned-url   # Generate presigned upload URL
```

### Livestream — Broadcasting

```
POST   /api/v1/public/streams               # Create stream room
PUT    /api/v1/public/streams/:id            # Update stream
POST   /api/v1/public/webhook/livekit        # LiveKit webhook receiver
```

## gRPC

Media service also runs a gRPC server on `localhost:50051`.

```
MediaService.GeneratePresignedURL   → Generate upload URL
MediaService.CommitFile             → Confirm file & move to permanent storage
MediaService.DeleteFile             → Delete file & purge cache
```

Protobuf definitions: `services/common-service/protobuf/media/media.proto`

## Shared Modules

| Module | Purpose |
|---|---|
| `services/common-service/xoxz` | Shared logger, Echo error handler, response model, app errors |
| `services/common-service/monitoring` | OpenTelemetry tracing, metrics, middleware, exporter setup |
| `services/common-service/protobuf` | Shared protobuf definitions & gRPC stubs |

## Observability Stack

- **OpenTelemetry Collector** — Receives traces, metrics, and logs from services
- **Prometheus** — Metrics storage
- **Tempo** — Distributed tracing
- **Loki** — Log aggregation
- **Grafana** — Unified dashboard (auto-provisioned datasources)

## Weverse.io Feature Parity

| Feature | Status | Service |
|---|---|---|
| User accounts & profiles | ✅ Done | Account Service |
| Google OAuth login | ✅ Done | Account Service |
| Artist profiles & groups | ✅ Done | Artist Service |
| Photo/video uploads | ✅ Done | Media Service |
| Live streaming (WebRTC) | ✅ Done | Livestream Service |
| Community feed (posts, comments, likes) | 🚧 Planned | Community Service (new) |
| Push notifications | 🚧 Planned | Notification Service (new) |
| Search (artists, posts, users) | 🚧 Planned | Search Service (new) |
| Shop / merchandise | 🚧 Planned | Shop Service (new) |
| Content moderation | 🚧 Planned | Moderation Service (new) |

## Tooling

```sh
go install github.com/google/wire/cmd/wire@latest
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
