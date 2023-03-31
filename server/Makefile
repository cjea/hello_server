.PHONY: all clean build run

LOCAL_DB_USER ?= postgres
LOCAL_DB_PORT ?= 30303
LOCAL_DB_DATABASE ?= postgres
PORT ?= 8080

APP_NAME=hello-api

all: clean build run

clean:
	@echo "Cleaning up..."
	@rm -f ${APP_NAME}

build:
	@echo "Building ${APP_NAME}..."
	@go build -o ${APP_NAME} main.go

run:
	@echo "Starting ${APP_NAME}..."
	@env DATABASE_URL="postgres://$(LOCAL_DB_USER)@127.0.0.1:$(LOCAL_DB_PORT)/$(LOCAL_DB_DATABASE)?sslmode=disable" PORT=$(PORT) ./${APP_NAME}
