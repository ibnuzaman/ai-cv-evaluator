## Change these variables as necessary.
include .env

binary_name = backend
db_user = $(DB_USER)
db_password = $(DB_PASSWORD)
db_host = $(DB_HOST)
db_port = $(DB_PORT)
db_name = $(DB_NAME)
migration_dir = database/migrations
migration_database = postgres://${db_user}:${db_password}@${db_host}:${db_port}/${db_name}?sslmode=disable

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	@test -z "$(git status --porcelain)"
    
## audit: run quality control checks
.PHONY: audit
audit: test
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)" 
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## tidy: tidy modfiles and format .go files
.PHONY: tidy
tidy:
	go mod tidy -v
	go mod vendor
	go fmt ./...

## build: build the application
.PHONY: build
build:
	go build -o /tmp/${binary_name} .

## build/run: build and run the application
.PHONY: build/run
build/run: build run

## run: run the application
.PHONY: run
run:
# 	go run main.go
	export ENV=development && /tmp/${binary_name}

## run: run the application with environment set to production
.PHONY: run/production
run/production: build
	export ENV=production && /tmp/${binary_name}

## run: run the application with environment set to staging
.PHONY: run/staging
run/staging: build
	export ENV=staging && /tmp/${binary_name}

## dev: run the application with reloading on file changes
.PHONY: dev
dev:
	export ENV=development && air -c .air.toml

## dev/run: run the docs, lint, and build run
.PHONY: dev/run
dev/run: docs lint build run

## dev/install: install air for reloading on file changes and golangci-lint for linting
.PHONY: dev/install
dev/install:
# 	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/swaggo/swag/cmd/swag@latest

## lint: run linters on the code
.PHONY: lint
lint:
	go fmt ./...
	go vet ./...
	swag fmt -d .
	golangci-lint run -v --fast

## docs: generate swagger documentation
.PHONY: docs
docs:
	swag init .
ifeq ($(shell uname), Darwin)
	gsed -i "s/x-nullable/nullable/g" ./docs/docs.go
	gsed -i "s/x-omitempty/omitempty/g" ./docs/docs.go
	gsed -i "s/x-nullable/nullable/g" ./docs/swagger.json
	gsed -i "s/x-omitempty/omitempty/g" ./docs/swagger.json
	gsed -i "s/x-nullable/nullable/g" ./docs/swagger.yaml
	gsed -i "s/x-omitempty/omitempty/g" ./docs/swagger.yaml
else
	sed -i "s/x-nullable/nullable/g" ./docs/docs.go
	sed -i "s/x-omitempty/omitempty/g" ./docs/docs.go
	sed -i "s/x-nullable/nullable/g" ./docs/swagger.json
	sed -i "s/x-omitempty/omitempty/g" ./docs/swagger.json
	sed -i "s/x-nullable/nullable/g" ./docs/swagger.yaml
	sed -i "s/x-omitempty/omitempty/g" ./docs/swagger.yaml
endif

## clean: clean up the project
.PHONY: clean
clean:
	rm -f /tmp/${binary_name}
	rm -f /tmp/coverage.out

## migrate/create: create a new migration with argument name. The name is snake case, for example: name=create_users_table
.PHONY: migrate/create
migrate/create:
	migrate create -ext sql -dir $(migration_dir) -seq $(name)

## migrate/up: run all migrations. Optional many argument can be passed to migrate to a specific many, for example: many=1
.PHONY: migrate/up
migrate/up:
	migrate -path ${migration_dir} -database ${migration_database} up $(many)

## migrate/down: rollback the last migration. Optional many argument can be passed to rollback to a specific many, for example: many=1
.PHONY: migrate/down
migrate/down:
	migrate -path ${migration_dir} -database ${migration_database} down $(many)

## migrate/drop: drop all tables
.PHONY: migrate/drop
migrate/drop:
	migrate -path ${migration_dir} -database ${migration_database} drop

## migrate/version: show the current migration version
.PHONY: migrate/version
migrate/version:
	migrate -path ${migration_dir} -database ${migration_database} version

## migrate/force: force a specific version. The version is an integer, for example: version=1
.PHONY: migrate/force
migrate/force:
	migrate -path ${migration_dir} -database ${migration_database} force $(version)

## migrate/goto: migrate to a specific version. The version is an integer, for example: version=1
.PHONY: migrate/goto
migrate/goto:
	migrate -path ${migration_dir} -database ${migration_database} goto $(version)