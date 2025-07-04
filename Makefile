include .env
export

CONFIG_DIR := config
GEN_DIR := gen

.PHONY: run
run:
	go run cmd/main.go

.PHONY: gencert
gencert:
	go tool cfssl gencert \
		-initca $(CONFIG_DIR)/ca-csr.json | go tool cfssljson -bare $(GEN_DIR)/ca
	go tool cfssl gencert \
		-ca $(GEN_DIR)/ca.pem \
		-ca-key $(GEN_DIR)/ca-key.pem \
		-config $(CONFIG_DIR)/ca-config.json \
		-profile server \
		$(CONFIG_DIR)/server-csr.json | go tool cfssljson -bare $(GEN_DIR)/server

.PHONY: test
test:
	go test -race ./...

.PHONY: deploy
deploy:
	helm upgrade --install gokube ./charts/gokube \
		--namespace gokube \
		--create-namespace \
		--values env.yaml

.PHONY: undeploy
undeploy:
	helm uninstall gokube --namespace gokube

.PHONY: delete
delete:
	helm uninstall gokube --namespace gokube
	kubectl delete namespace gokube

.PHONY: port-forward
port-forward:
	kubectl port-forward svc/gokube 8080:8080 --namespace gokube

.PHONY: help
help:
	@echo "Usage: make <target>"
	@echo "Targets:"
	@echo "  run - Run the application"
	@echo "  test - Run tests"
	@echo "  deploy - Deploy the application"