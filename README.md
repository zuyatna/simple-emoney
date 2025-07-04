# simple-emoney

Roadmap project structure:

```
simple-emoney/
├── config/
│   └── config.go
├── db/
│   └── migrations/
│       ├── 000001_create_users_table.up.sql
│       ├── 000001_create_users_table.down.sql
│       ├── 000002_create_transactions_table.up.sql
│       └── 000002_create_transactions_table.down.sql
├── internal/
│   ├── app/
│   │   ├── handler/
│   │   │   ├── auth_handler.go
│   │   │   ├── transaction_handler.go
│   │   │   └── user_handler.go
│   │   ├── middleware/
│   │   │   └── auth_middleware.go
│   │   ├── repository/
│   │   │   ├── transaction_repository.go
│   │   │   ├── user_repository.go
│   │   │   └── redis_repository.go
│   │   └── service/
│   │       ├── auth_service.go
│   │       ├── transaction_service.go
│   │       └── user_service.go
│   ├── model/
│   │   ├── auth.go
│   │   ├── transaction.go
│   │   └── user.go
│   └── router/
│       └── router.go
├── pkg/
│   ├── database/
│   │   ├── postgres.go
│   │   └── redis.go
│   └── utils/
│       └── jwt.go
├── .env.example
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## Work in progress!!
Project structure may be changed during development.

## Wiki
[https://github.com/zuyatna/simple-emoney.wiki.git](https://github.com/zuyatna/simple-emoney/wiki)
