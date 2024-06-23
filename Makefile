ifneq (,$(wildcard ./.env))
    include .env
    export
endif
GIT_SHA := $(shell git rev-parse --short HEAD)
ORIGINAL_DISCOVERY_URL=changeme_url
ORIGINAL_SERVICE_ACCOUNT=changeme_email

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

deploy-frontend:
	@echo "Deploying frontend"
	@cd apps/frontend && \
	sed -i '' 's|DISCOVERY_URL: $(ORIGINAL_DISCOVERY_URL)|DISCOVERY_URL: $(DISCOVERY_URL)|' app.yaml && \
	sed -i '' 's|service_account: $(ORIGINAL_SERVICE_ACCOUNT)|service_account: ae-frontend@$(GCP_PROJECT_ID).iam.gserviceaccount.com|' app.yaml && \
	gcloud app deploy app.yaml --quiet && \
	sed -i '' 's|DISCOVERY_URL: $(DISCOVERY_URL)|DISCOVERY_URL: $(ORIGINAL_DISCOVERY_URL)|' app.yaml && \
	sed -i '' 's|service_account: ae-frontend@$(GCP_PROJECT_ID).iam.gserviceaccount.com|service_account: $(ORIGINAL_SERVICE_ACCOUNT)|' app.yaml
