.PHONY: infra dev-api

infra:
	@echo "Creating infra"
	terraform -chdir=terraform init
	terraform -chdir=terraform plan -out=tf.plan
	terraform -chdir=terraform apply tf.plan

dev-api:
	@echo "Starting API"
	cd functions/item-api && \
	go mod tidy && \
	FUNCTION_TARGET=app LOCAL_ONLY=true DISCOVERY_URL=http://discovery go run cmd/main.go
