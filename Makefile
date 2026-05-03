CLUSTER_NAME   := devops-demo
REGISTRY       := localhost:5001
SERVICES       := sensor-simulator api-gateway data-sink enterprise-consumer
ARGOCD_NS      := platform
K8S_VERSION    := v1.29.2

.DEFAULT_GOAL := help

# ─── Cluster lifecycle ────────────────────────────────────────────────────────

.PHONY: bootstrap
bootstrap: ## Full cluster bootstrap (idempotent)
	@bash scripts/bootstrap.sh

.PHONY: teardown
teardown: ## Delete cluster and local registry
	@bash scripts/teardown.sh

.PHONY: cluster
cluster: ## Create kind cluster + local registry only
	@kind create cluster --name $(CLUSTER_NAME) --config kind-config.yaml --image kindest/node:$(K8S_VERSION) || true
	@bash scripts/registry.sh

# ─── Images ───────────────────────────────────────────────────────────────────

.PHONY: images
images: ## Build and push all service images to local registry
	@for svc in $(SERVICES); do \
		echo ">>> Building $$svc"; \
		docker build -t $(REGISTRY)/$$svc:latest services/$$svc; \
		docker push $(REGISTRY)/$$svc:latest; \
	done

.PHONY: image-%
image-%: ## Build and push a single service (e.g. make image-api-gateway)
	docker build -t $(REGISTRY)/$*:latest services/$*
	docker push $(REGISTRY)/$*:latest

# ─── Go tooling ───────────────────────────────────────────────────────────────

.PHONY: tidy
tidy: ## Run go mod tidy in all services
	@for svc in $(SERVICES); do \
		echo ">>> Tidying $$svc"; \
		cd services/$$svc && go mod tidy && cd ../..; \
	done

.PHONY: lint
lint: ## Run golangci-lint on all services
	@for svc in $(SERVICES); do \
		echo ">>> Linting $$svc"; \
		cd services/$$svc && golangci-lint run ./... && cd ../..; \
	done

.PHONY: test
test: ## Run go test on all services
	@for svc in $(SERVICES); do \
		echo ">>> Testing $$svc"; \
		cd services/$$svc && go test -v ./... && cd ../..; \
	done

# ─── ArgoCD ───────────────────────────────────────────────────────────────────

.PHONY: argocd-password
argocd-password: ## Print ArgoCD initial admin password
	@kubectl -n $(ARGOCD_NS) get secret argocd-initial-admin-secret \
		-o jsonpath="{.data.password}" | base64 -d && echo

.PHONY: argocd-sync
argocd-sync: ## Force sync all ArgoCD apps
	@argocd app sync --grpc-web devops-demo-root --prune

# ─── Demo ─────────────────────────────────────────────────────────────────────

.PHONY: demo
demo: ## Print the demo walkthrough script for the interview panel
	@bash scripts/demo.sh

.PHONY: port-forwards
port-forwards: ## Start all useful port-forwards in background
	@kubectl -n observability port-forward svc/grafana 3000:80 &
	@kubectl -n platform port-forward svc/argocd-server 8443:443 &
	@kubectl -n apps port-forward svc/api-gateway 8080:8080 &
	@kubectl -n tools port-forward svc/node-red 1880:1880 &
	@echo "Grafana:  http://localhost:3000  (admin/admin)"
	@echo "ArgoCD:   https://localhost:8443"
	@echo "API:      http://localhost:8080"
	@echo "Node-RED: http://localhost:1880"

# ─── Housekeeping ─────────────────────────────────────────────────────────────

.PHONY: help
help: ## Show available targets
	@echo ""
	@echo "  DevOps Demo — Makefile targets"
	@echo ""
	@grep -E '^[a-zA-Z_%/-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'
	@echo ""
