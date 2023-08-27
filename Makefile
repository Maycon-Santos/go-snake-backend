include .env.sample
-include .env

GO111MODULE=auto
DOCKER_IMAGE_NAME=$(shell echo $(APP_NAME) | tr A-Z a-z | tr ' ' -)

install-modd-local:
	@echo installing modd
	@go get -v github.com/cortesi/modd/cmd/modd 2>/dev/null || true
	@echo modd installed

test-coverage-local:
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html

setup-local: install-modd-local
	@go mod vendor

setup-dev-local: setup-local install-modd-local
	@go mod tidy

run-dev-local:
	@modd -f ./cmd/server/modd.conf

test-coverage:
	@docker run -ti --rm \
		-v "$(PWD)":/usr/src/app \
		--name $(DOCKER_IMAGE_NAME)-coverage \
		$(DOCKER_IMAGE_NAME) \
		go test ./... -coverprofile=coverage.out && \
    go tool cover -html=coverage.out -o coverage.html

mock:
	@docker run -ti --rm \
		-v "$(PWD)":/usr/src/app \
		--name $(DOCKER_IMAGE_NAME)-coverage \
		$(DOCKER_IMAGE_NAME) \
		mockgen -source=$(source) -destination=$(shell echo $(source) | sed 's/.go$$/_mock.go/') -package=$(package)

setup-dev:
	@echo $(DOCKER_IMAGE_NAME)
	@docker build \
		--target development \
		-t $(DOCKER_IMAGE_NAME) \
		.

run-dev:
	@docker run -ti --rm \
		-v "$(PWD)":/usr/src/app \
    --network="host" \
		--expose $(SERVER_PORT) \
		--name $(DOCKER_IMAGE_NAME)-server \
		$(DOCKER_IMAGE_NAME) \
		modd -f ./cmd/server/modd.conf

migrate-up:
	@docker run \
		-v ${PWD}/db/migrations:/migrations:delegated \
		--name ${DOCKER_IMAGE_NAME}-migrate \
		--rm \
    --network="host" \
		migrate/migrate -verbose -path=/migrations/ -database ${DATABASE_CONN_URI} up

migrate-down:
	@docker run \
		-v ${PWD}/db/migrations:/migrations:delegated \
		--name ${DOCKER_IMAGE_NAME}-migrate \
		--rm \
    --network="host" \
		migrate/migrate -verbose -path=/migrations/ -database ${DATABASE_CONN_URI} down 1

migrate-new:
	@docker run \
		-v ${PWD}/db/migrations:/migrations:delegated \
		--name ${DOCKER_IMAGE_NAME}-migrate \
		-u ${CURRENT_USER}:${CURRENT_GROUP} \
		--rm \
    --network="host" \
		migrate/migrate -verbose -path=/migrations/ -database ${DATABASE_CONN_URI} create -dir ./migrations -ext sql $(FILE)
