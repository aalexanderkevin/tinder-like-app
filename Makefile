APP_NAME=tinder-like-app
VERSION_VAR=main.Version
VERSION=$(shell git describe --tags)

dep:
	@echo ">> Downloading Dependencies"
	@go mod download

dep-lint:
	@echo ">> Downloading golangci-lint"
	@curl -sSfL "https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh" | sh -s -- -b $$(go env GOPATH)/bin v1.47.2

lint:
	@echo ">> Running Linter"
	@$$(go env GOPATH)/bin/golangci-lint run

lint-new:
	@echo ">> Running Linter on new/updated files"
	@$$(go env GOPATH)/bin/golangci-lint run --new-from-rev=HEAD~1

hooks:
	@echo ">> Installing git hooks"
	git config core.hooksPath .githooks

build: dep
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-X ${VERSION_VAR}=${VERSION}" -a -installsuffix nocgo -o ./bin ./...

docker:
	@echo ">> Building Docker Image"
	@docker build -t ${APP_NAME}:latest .

run-server:
	env $$(cat .env | xargs) go run tinder-like-app/cmd server
	
migrate:
	eval $$(egrep -v '^#' .env | xargs -0) go run tinder-like-app/cmd migrate

test-all: test-unit test-integration-with-infra

test-unit: dep
	@echo ">> Running Unit Test"
	@env $$(cat .env.testing | xargs) go test -tags=unit -failfast -cover -covermode=atomic ./...

test-integration: dep
	@echo ">> Running Integration Test"
	@env $$(cat .env.testing | xargs) env POSTGRES_MIGRATION_PATH=$$(pwd)/database/migrations go test -tags=integration -failfast -cover -covermode=atomic ./...

test-integration-with-infra: test-infra-up test-integration test-infra-down

test-infra-up:
	$(MAKE) test-infra-down
	@echo ">> Starting Test DB"
	docker run -d --rm --name test-postgres -p 5431:5432 --env-file .env.testing postgres:12

test-infra-down:
	@echo ">> Shutting Down Test DB"
	@-docker kill test-postgres