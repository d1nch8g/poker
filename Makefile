
# Installs required software to work with the project
install:
	go install -tags "postgres,mysql" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/go-bindata/go-bindata/go-bindata@latest
	go install github.com/go-acme/lego/v4/cmd/lego@latest
	go install mvdan.cc/gofumpt@latest

fmt:
	swag fmt
	gofumpt -l -w . 

lint:
	golangci-lint run -v

rundb:
	docker compose up

migrate:
	migrate -path db/migrations/ -database "postgresql://user:password@localhost:5432/db?sslmode=disable" -verbose up

# Generates database related and swagger documentation from request comments (also checks migration schema)
.PHONY: gen
gen:
	sqlc generate
	swag init -o . --ot yaml
	go-bindata -fs -pkg migr -o gen/migr/migr.go -prefix "db/migrations/" db/migrations/...
	# go-bindata -fs -pkg web -o gen/web/web.go -prefix "dist/" dist/...
	# npm run build

# Output in .lego/certificates/chx.su.crt .lego/certificates/chx.su.key
certs:
	lego --email="d1nch8g@gmail.com" --domains="chx.su" --http run

# ip on local network
ip:
	ip route get 1.2.3.4 | grep -oP '(?<=src )\S+'
