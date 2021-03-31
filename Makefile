include .env.sample
-include .env

GO111MODULE=auto
DOCKER_IMAGE_NAME=$(shell echo $(APP_NAME) | tr A-Z a-z | tr ' ' -)

install-modd-local:
	@echo installing modd
	@go get -v github.com/cortesi/modd/cmd/modd 2>/dev/null || true
	@echo modd installed

setup-local:
	@go mod vendor

setup-dev-local: setup-local install-modd-local
	@go mod tidy

run-dev-local:
	@modd -f ./cmd/server/modd.conf

setup-dev:
	@echo $(DOCKER_IMAGE_NAME)
	@docker build \
		--target development \
		-t $(DOCKER_IMAGE_NAME) \
		.

run-dev:
	@docker run -ti --rm \
		-v "$(PWD)":/usr/src/app \
		--expose $(SERVER_PORT) \
		--name $(DOCKER_IMAGE_NAME)-server \
		$(DOCKER_IMAGE_NAME) \
		modd -f ./cmd/server/modd.conf
