PROJECT_NAME=ymir
PROJECT_GIT_URL=github.com/svartlfheim/ymir
DOCKER_COMPOSE=docker compose -f ./.docker/docker-compose.yml --project-name="$(PROJECT_NAME)_dev"
DOCKER_COMPOSE_RUN=$(DOCKER_COMPOSE) run --rm
DOCKER_COMPOSE_RUN_TF=$(DOCKER_COMPOSE_RUN) --entrypoint /bin/terraform tf

ifndef ARGS
	ARGS=-count=1 ./...
endif

UPSERT_PLAN_NAME=local-test.tfplan
DESTROY_PLAN_NAME=destroy-local-test.tfplan

ifdef BUILD_HASH
	LDFLAGS_VAL :=$(LDFLAGS_VAL) -X '$(PROJECT_GIT_URL)/cmd/ymir.versionHashRef=$(BUILD_HASH)'
endif

ifdef BUILD_VERSION
	LDFLAGS_VAL :=$(LDFLAGS_VAL) -X '$(PROJECT_GIT_URL)/cmd/ymir.version=$(BUILD_VERSION)'
endif

LDFLAGS=-ldflags="$(LDFLAGS_VAL)"

# This is a combination of the following suggestions:
# https://gist.github.com/prwhite/8168133#gistcomment-1420062
help: ## This help dialog.
	@IFS=$$'\n' ; \
	help_lines=(`fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##/:/'`); \
	printf "%-30s %s\n" "target" "help" ; \
	printf "%-30s %s\n" "------" "----" ; \
	for help_line in $${help_lines[@]}; do \
			IFS=$$':' ; \
			help_split=($$help_line) ; \
			help_command=`echo $${help_split[0]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
			help_info=`echo $${help_split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
			printf '\033[36m'; \
			printf "%-30s %s" $$help_command ; \
			printf '\033[0m'; \
			printf "%s\n" $$help_info; \
	done

init: install-hosts hosts-dns app-env tls-trust-ca link-ca-for-tf ## Initialise the project

hosts-dns: ## Setup /etc/hosts entries for decvelopment
	sudo ./bin/hosts add 127.0.0.1 ymir.local || true
	sudo ./bin/hosts add 127.0.0.1 pgadmin.ymir.local || true

app-env: ## Generate the app env file if it does not exist
	if [ ! -f ./.docker/.env ]; then cp ./.docker/.env.example ./.docker/.env; fi;

up: ## Start the docker-compose development environment
	$(DOCKER_COMPOSE) up -d

build-images:
	$(DOCKER_COMPOSE) build --no-cache

down: ## Destroy the docker-compose development environment
	$(DOCKER_COMPOSE) down

restart: destroy up## Restart the docker-compose development environment; destroy && up

restart-ymir: ## Restart the ymir container only
	$(DOCKER_COMPOSE) restart ymir

destroy-all: destroy ## Destroy the docker-compose development environment and built images
	docker rmi ymir_dev_ymir

exec-ls: ## Exec into a shell in the localstack runtime
	$(DOCKER_COMPOSE) exec -it lstack bash

exec-tf: ## Exec into a shell in the terraform runtime
	$(DOCKER_COMPOSE_RUN) --entrypoint /bin/bash tf

exec-ymir: ## Exec into a shell in the ymir runtime
	$(DOCKER_COMPOSE) exec ymir /bin/bash

exec-pgsql: ## Exec into a shell in the postgres runtime
	$(DOCKER_COMPOSE) exec postgres /bin/bash

exec-pgadmin: ## Exec into a shell in the pgadmin runtime
	$(DOCKER_COMPOSE) exec pgadmin /bin/sh

tf-init: ## Run the terraform init command
	$(DOCKER_COMPOSE_RUN_TF) init

tf-plan: ## Run the terraform plan command
	$(DOCKER_COMPOSE_RUN_TF) plan -out $(UPSERT_PLAN_NAME)

tf-apply: tf-plan ## Run the terraform apply command (runs plan first)
	$(DOCKER_COMPOSE_RUN_TF) apply $(UPSERT_PLAN_NAME)

tf-plan-destroy: ## Run the terraform plan command to destroy the resources
	$(DOCKER_COMPOSE_RUN_TF) plan -destroy -out $(DESTROY_PLAN_NAME)

tf-destroy: tf-plan-destroy ## Run the terraform apply command to destroy all resources (runs a plan for destroy first)
	$(DOCKER_COMPOSE_RUN_TF) apply $(DESTROY_PLAN_NAME)

tail-ls: ## Tail logs from localstack
	$(DOCKER_COMPOSE) logs -f localstack

tail-ymir: ## Tail logs from ymir
	$(DOCKER_COMPOSE) logs -f ymir

tail-pgadmin: ## Tail logs from pgadmin
	$(DOCKER_COMPOSE) logs -f pgadmin

tail-reverse-proxy: ## Tail logs from reverse-proxy
	$(DOCKER_COMPOSE) logs -f reverse-proxy

install-hosts: ## Installs the hosts cli utility in ./bin
	sudo curl -L https://raw.github.com/xwmx/hosts/master/hosts -o ./bin/hosts
	sudo chmod +x ./bin/hosts

tls-trust-ca: ## Trust the self-signed HTTPS certification
	sudo security add-trusted-cert -d -r trustRoot -k "/Library/Keychains/System.keychain" "./.docker/certs/minica.pem"

link-ca-for-tf: ## Link CA into the tf build context in ./docker/tf
	ln -s $(shell pwd)/.docker/certs/minica.pem $(shell pwd)/.docker/tf/minica.pem || true

build: ## Build the ymir binary
	go mod tidy
	go build $(LDFLAGS) -o /go/bin/$(PROJECT_NAME) cmd/main.go

.PHONY: fmt
fmt: ## Format code using go fmt ./...
	go fmt ./...

.PHONY: lint
lint: ## Lint with golangci-lint run ./...
	golangci-lint run ./...

.PHONY: unit-test
unit-test: ## Run the unit tests in ymir environment
	@$(DOCKER_COMPOSE) exec ymir /bin/bash -c "go test -cover $(ARGS)"

.PHONY: test
test: ## Run unit and integration tests ymir environment
	@$(DOCKER_COMPOSE) exec ymir /bin/bash -c "YMIR_CI_INTEGRATION_TESTS_ENABLED=true go test -coverprofile test/coverage.out -cover $(ARGS)"
	@$(DOCKER_COMPOSE) exec ymir /bin/bash -c "go tool cover -o=test/coverage.html -html=test/coverage.out"

serve: build ## Run the serve command of ymir (runs the http server)
	ymir serve