```text
services/
└── name-service/
    ├── cmd/
    │   └── main.go
    │
    ├── configs/
    │   └── config.yaml
    │
    ├── docs/
    │   ├── docs.go
    │   ├── swagger.json
    │   └── swagger.yaml
    │
    ├── internal/
    │   ├── adapter/
    │   ├── bootstrap/
    │   ├── config/
    │   ├── const/
    │   ├── controller/
    │   ├── di/
    │   ├── domain/
    │   │   ├── entity/
    │   │   ├── repository/
    |   |       ├── repo_impl/
    │   │   └── seeder/
    │   │
    │   ├── logic/
    │   ├── middleware/
    │   ├── model/
    │   ├── service/
    │   └── util/
    │
    ├── migrations/
    │
    ├── .dockerignore
    ├── .env
    ├── .gitignore
    ├── Dockerfile
    ├── go.mod
    ├── go.sum
    └── Makefile
```