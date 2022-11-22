CURRENT_DIR=$(shell pwd)
APP=template
APP_CMD_DIR=./cmd

build:
	CGO_ENABLED=0 GOOS=darwin go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go

proto-gen:
	./scripts/gen-proto.sh	${CURRENT_DIR}

run:
	go run cmd/main.go

lint: ## Run golangci-lint with printing to stdout
	golangci-lint -c .golangci.yaml run --build-tags "musl" ./...

pull_submodule:
	git submodule update --init --recursive

update_submodule:
	git submodule update --remote --merge

create_submodule:
	git submodule add gitpathToProtos

migrate_up:
	migrate -path migrations/ -database postgres://postgres:1234@localhost:5432/review up

migrate_down:
	migrate -path migrations/ -database postgres://postgres:1234@localhost:5432/review down

migrate_force:
	migrate -path migrations/ -database postgres://postgres:1234@localhost:5432/review force 1