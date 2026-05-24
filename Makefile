include .env
export

# Директория, в которой хранятся исполняемые
# файлы проекта и зависимости, необходимые для сборки.
LOCAL_BIN := $(CURDIR)/bin

# Шорткат для golangci-lint
GOLANGCI_BIN := $(LOCAL_BIN)/golangci-lint

# Шорткаты для создания миграций
migrations_dir := $(CURDIR)/migrations

# Переменная пути к сгенерированным swagger.json
SWAGGER_SOURCES := $(shell find  pkg/** -name '*.json')

# Переменная со списком бинарников для сборки,
# по умолчанию только основной бинарник сервиса.
BUILD_PATHS ?= ./cmd/main

# Дополнительные параметры для сборки приложения.
BUILD_ENVPARMS ?= CGO_ENABLED=0

# Определяем операционную систему
PLATFORM := $(shell uname)

PROTOC_VERSION := 31.1
PROTOC_ZIP := protoc-$(PROTOC_VERSION)-linux-x86_64.zip
ifeq ($(PLATFORM),Darwin)
	PROTOC_ZIP = protoc-$(PROTOC_VERSION)-osx-x86_64.zip
endif
	PROTOC_URL := https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/$(PROTOC_ZIP)


ifeq ($(PLATFORM),Darwin)
    SED_COMMAND = sed -i ''
else ifeq ($(PLATFORM),Linux)
    SED_COMMAND = sed -i
endif

# Установить зависимости
bin-deps:
	$(info #Installing project binary dependencies...)
	curl -sSL $(PROTOC_URL) -o /tmp/$(PROTOC_ZIP)
	unzip -o /tmp/$(PROTOC_ZIP) -d /tmp/protoc
	chmod u+w /tmp/protoc/bin/protoc
	cp /tmp/protoc/bin/protoc $(LOCAL_BIN)/
	cp -r /tmp/protoc/include $(LOCAL_BIN)/include
	rm -rf /tmp/$(PROTOC_ZIP) /tmp/protoc
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.18.0
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.2.1
	GOBIN=$(LOCAL_BIN) go install github.com/easyp-tech/easyp/cmd/easyp@v0.7.15

# Генерация контрактов
generate-proto:
	$(info Starting proto generation)
	PATH=$(LOCAL_BIN):$$PATH easyp generate

# Запустить линтер
lint:
	$(info Running lint against all project files...)
	$(GOLANGCI_BIN) run --config=.golangci.yml ./...

# Сборка приложения
build:
	$(info Build application)
	$(BUILD_ENVPARMS) go build -o $(LOCAL_BIN) $(BUILD_PATHS)

# Создать миграцию
migration:
	mkdir -p $(migrations_dir)
	$(LOCAL_BIN)/goose -dir $(migrations_dir) create $(shell bash -c 'read -p "Migration name: " migration_name; echo $$migration_name') sql

migration-up:
	$(LOCAL_BIN)/goose $(opts) -allow-missing -dir ./migrations postgres "host=$$POSTGRES_HOST port=$$POSTGRES_PORT user=$$POSTGRES_USER password=$$POSTGRES_PASSWORD dbname=$$POSTGRES_DB sslmode=disable" up

migration-down:
	$(LOCAL_BIN)/goose $(opts) -dir ./migrations postgres "host=$$POSTGRES_HOST port=$$POSTGRES_PORT user=$$POSTGRES_USER password=$$POSTGRES_PASSWORD dbname=$$POSTGRES_DB sslmode=disable" down

migration-status:
	$(LOCAL_BIN)/goose $(opts) -dir ./migrations postgres "host=$$POSTGRES_HOST port=$$POSTGRES_PORT user=$$POSTGRES_USER password=$$POSTGRES_PASSWORD dbname=$$POSTGRES_DB sslmode=disable" status

.PHONY:
	bin-deps \
 	generate-proto \
 	lint \
 	build \
 	migration \
 	migration-up \
 	migration-down \
 	migration-status