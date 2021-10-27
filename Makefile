DEFAULT_ENV_FILE='.env'
DEV_ENV_FILE='.env.local.dev'
ENV_FILE := $(shell if [ ! -f $(DEFAULT_ENV_FILE) ]; then echo $(DEV_ENV_FILE) ; else echo $(DEFAULT_ENV_FILE) ; fi   )
include $(ENV_FILE)
export

MAIN_APP_FILE=./cmd/chat/main.go
MAIN_APP_DIR:= $(shell dirname $(MAIN_APP_FILE))
PROJECT_NAME=nixchat_backend

## ----------------------------------------------------------------------
## 		A little manual for using this Makefile.
## ----------------------------------------------------------------------


.PHONY: build
build:	## Compile the code into an executable application
	go build -v -o ./bin/main ${MAIN_APP_FILE}


.PHONY: docker-build
docker-build:	## Build docker image
	docker build -t ${PROJECT_NAME} .


.PHONY: run
run:	## Run application
	go run ${MAIN_APP_FILE}


MOCKS_DESTINATION=mocks
.PHONY: mocks
mocks: ## Generate mocks
	@echo "Generating mocks..."
	go generate ./...

.PHONY: test
test: mocks ## Run golang tests
	go test -race  -coverprofile=coverage.out -cover `go list ./... | grep -v mocks `


.PHONY: migrate-up
migrate-up:	## run db migration scripts
	migrate -path ./migrations/ -database 'postgres://${PG_USERNAME}:${PG_PASSWORD}@${PG_HOST}:${PG_PORT}/${PG_NAME}?sslmode=${PG_SSLMODE}' up

.PHONY: migrate-down
migrate-down:	## rollback db migrations
	migrate -path ./migrations/ -database 'postgres://${PG_USERNAME}:${PG_PASSWORD}@${PG_HOST}:${PG_PORT}/${PG_NAME}?sslmode=${PG_SSLMODE}' down

.PHONY: linter
linter:	## Run linter for *.go files
	revive -config .linter.toml  -exclude ./vendor/... -formatter unix ./...


.PHONY: docker-compose-up
docker-compose-up:	## Run application and app environment in docker
	docker-compose --env-file ${ENV_FILE} up db redis backend-app migrate


.PHONY: docker-compose-dev-up
docker-compose-dev-up:	## Run local dev environment
	docker-compose --env-file ${ENV_FILE} up db redis migrate


.PHONY: docker-compose-dev-down
docker-compose-dev-down:	## Stop local dev environment
	docker-compose --env-file ${ENV_FILE} down


.PHONY: swagger
swagger:	## Generate swagger api specs
	swag init --output ./api --dir ${MAIN_APP_DIR},./internal --parseInternal true


.PHONY: help
help:     ## Show this help.
	echo $(REDIS_PORT)
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)


.DEFAULT_GOAL := build
