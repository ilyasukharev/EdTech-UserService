### Tests
`TEST_DSN="postgres://login:pwd@localhost:5432/vpn_subscriptions_db?sslmode=disable" go test .
`


### Swagger

1. `go install github.com/swaggo/swag/cmd/swag@v1.8.1`
2. `swag init -g ./cmd/api_service.go --parseInternal`
