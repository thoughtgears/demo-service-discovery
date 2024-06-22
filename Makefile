ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: infra deploy-api deploy-frontend

infra:
	@echo "Creating infra"
	terraform -chdir=terraform init
	terraform -chdir=terraform plan -out=tf.plan
	terraform -chdir=terraform apply tf.plan

deploy-api:
	@echo "Deploying API"
	cd functions/item-api && \
	gcloud functions deploy item-api \
	--runtime go121 \
	--trigger-http \
	--region=$(GCP_REGION) \
	--project=$(GCP_PROJECT_ID) \
	--no-allow-unauthenticated \
	--entry-point=app \
	--no-gen2 \
	--ingress-settings=internal-only \
	--service-account=cf-item-api@$(GCP_PROJECT_ID).iam.gserviceaccount.com

deploy-frontend:
	@echo "Deploying frontend"
	cd functions/item-frontend && \
	gcloud functions deploy item-frontend \
	--runtime go121 \
	--trigger-http \
	--region=$(GCP_REGION) \
	--project=$(GCP_PROJECT_ID) \
	--no-allow-unauthenticated \
	--entry-point=app \
	--no-gen2 \
	--set-env-vars=ENVIRONMENT=dev,DISCOVERY_URL=$(DISCOVERY_URL) \
	--service-account=cf-item-bff@$(GCP_PROJECT_ID).iam.gserviceaccount.com
