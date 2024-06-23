ifneq (,$(wildcard ./.env))
    include .env
    export
endif
GIT_SHA := $(shell git rev-parse --short HEAD)

.PHONY: all infra build push deploy-api deploy-frontend

infra:
	@echo "Creating infra"
	terraform -chdir=terraform init
	terraform -chdir=terraform plan -out=tf.plan
	terraform -chdir=terraform apply tf.plan

build:
	docker build --platform=linux \
				-t $(GCP_REGION)-docker.pkg.dev/$(GCP_PROJECT_ID)/demos/item-api:latest \
 				-t $(GCP_REGION)-docker.pkg.dev/$(GCP_PROJECT_ID)/demos/item-api:$(GIT_SHA) \
 				-f apps/items-api/Dockerfile ./apps/items-api

	docker build --platform=linux \
				-t $(GCP_REGION)-docker.pkg.dev/$(GCP_PROJECT_ID)/demos/store-bff:latest \
				-t $(GCP_REGION)-docker.pkg.dev/$(GCP_PROJECT_ID)/demos/store-bff:$(GIT_SHA) \
				-f apps/items-api/Dockerfile ./apps/store-bff

push:
	docker push $(GCP_REGION)-docker.pkg.dev/$(GCP_PROJECT_ID)/demos/item-api:latest
	docker push $(GCP_REGION)-docker.pkg.dev/$(GCP_PROJECT_ID)/demos/item-api:$(GIT_SHA)
	docker push $(GCP_REGION)-docker.pkg.dev/$(GCP_PROJECT_ID)/demos/store-bff:latest
	docker push $(GCP_REGION)-docker.pkg.dev/$(GCP_PROJECT_ID)/demos/store-bff:$(GIT_SHA)
