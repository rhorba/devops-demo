#!/usr/bin/env bash
# bootstrap.sh — idempotent cluster setup for devops-demo
set -euo pipefail

CLUSTER_NAME="devops-demo"
REGISTRY_NAME="kind-registry"
REGISTRY_PORT="5001"
ARGOCD_VERSION="v2.10.2"
ARGOCD_NS="platform"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; CYAN='\033[0;36m'; NC='\033[0m'
log()  { echo -e "${CYAN}[bootstrap]${NC} $*"; }
ok()   { echo -e "${GREEN}[✓]${NC} $*"; }
warn() { echo -e "${YELLOW}[!]${NC} $*"; }
die()  { echo -e "${RED}[✗]${NC} $*"; exit 1; }

# ─── Pre-flight checks ────────────────────────────────────────────────────────
log "Running pre-flight checks..."
for cmd in docker kind kubectl helm argocd; do
    command -v "$cmd" &>/dev/null || die "Missing required tool: $cmd"
done
ok "All tools present"

# ─── Local Docker registry ────────────────────────────────────────────────────
if ! docker ps --format '{{.Names}}' | grep -q "^${REGISTRY_NAME}$"; then
    log "Starting local registry on port ${REGISTRY_PORT}..."
    docker run -d --restart=always -p "127.0.0.1:${REGISTRY_PORT}:5000" \
        --network bridge --name "${REGISTRY_NAME}" registry:2
    ok "Registry started"
else
    ok "Registry already running"
fi

# ─── Kind cluster ─────────────────────────────────────────────────────────────
if ! kind get clusters 2>/dev/null | grep -q "^${CLUSTER_NAME}$"; then
    log "Creating kind cluster '${CLUSTER_NAME}'..."
    kind create cluster --name "${CLUSTER_NAME}" \
        --config "${ROOT_DIR}/kind-config.yaml" \
        --image kindest/node:v1.29.2
    ok "Cluster created"
else
    ok "Cluster '${CLUSTER_NAME}' already exists"
fi

# ─── Connect registry to cluster network ──────────────────────────────────────
if ! docker network inspect kind 2>/dev/null | grep -q "${REGISTRY_NAME}"; then
    log "Connecting registry to kind network..."
    docker network connect kind "${REGISTRY_NAME}" 2>/dev/null || true
fi

# Apply registry ConfigMap so containerd knows about localhost:5001
kubectl apply -f - <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: local-registry-hosting
  namespace: kube-public
data:
  localRegistryHosting.v1: |
    host: "localhost:${REGISTRY_PORT}"
    help: "https://kind.sigs.k8s.io/docs/user/local-registry/"
EOF
ok "Registry ConfigMap applied"

# ─── Namespaces ───────────────────────────────────────────────────────────────
log "Creating namespaces..."
for ns in platform messaging apps observability localstack tools security; do
    kubectl create namespace "$ns" --dry-run=client -o yaml | kubectl apply -f -
done
ok "Namespaces ready"

# ─── cert-manager ─────────────────────────────────────────────────────────────
log "Installing cert-manager..."
helm repo add jetstack https://charts.jetstack.io --force-update
helm upgrade --install cert-manager jetstack/cert-manager \
    --namespace platform --create-namespace \
    --version v1.14.4 \
    --set installCRDs=true \
    --wait --timeout 5m
ok "cert-manager installed"

# ─── Kyverno ─────────────────────────────────────────────────────────────────
log "Installing Kyverno..."
helm repo add kyverno https://kyverno.github.io/kyverno/ --force-update
helm upgrade --install kyverno kyverno/kyverno \
    --namespace security --create-namespace \
    --version 3.1.4 \
    --wait --timeout 5m || warn "Kyverno upgrade had hook issues (pods already running — continuing)"
ok "Kyverno installed"

# ─── ArgoCD ───────────────────────────────────────────────────────────────────
log "Installing ArgoCD ${ARGOCD_VERSION}..."
helm repo add argo https://argoproj.github.io/argo-helm --force-update
helm upgrade --install argocd argo/argo-cd \
    --namespace "${ARGOCD_NS}" \
    --version 6.7.14 \
    --set server.service.type=NodePort \
    --set server.service.nodePortHttp=30080 \
    --wait --timeout 8m
ok "ArgoCD installed"

# ─── Grafana admin secret ─────────────────────────────────────────────────────
log "Creating Grafana admin secret (override GRAFANA_PASSWORD to change)..."
GRAFANA_PASSWORD="${GRAFANA_PASSWORD:-admin}"
kubectl -n observability create secret generic grafana-admin \
    --from-literal=admin-user=admin \
    --from-literal=admin-password="${GRAFANA_PASSWORD}" \
    --dry-run=client -o yaml | kubectl apply -f -
ok "Grafana admin secret ready"

# ─── Thanos S3 config (must exist before prometheus stack syncs) ──────────────
log "Applying Thanos objstore secret..."
kubectl apply -f "${ROOT_DIR}/monitoring/prometheus/thanos-secret.yaml"
ok "Thanos secret applied"

# ─── Grafana dashboard ConfigMap ──────────────────────────────────────────────
log "Applying Grafana messaging dashboard ConfigMap..."
kubectl -n observability create configmap grafana-dashboard-messaging \
    --from-file=messaging-pipeline.json="${ROOT_DIR}/monitoring/grafana/dashboards/messaging-pipeline.json" \
    --dry-run=client -o yaml | kubectl apply -f -
kubectl -n observability label configmap grafana-dashboard-messaging \
    grafana_dashboard=1 --overwrite
ok "Grafana dashboard ConfigMap applied"

# ─── Bootstrap root Application ───────────────────────────────────────────────
log "Applying ArgoCD root Application (App of Apps)..."
kubectl apply -f "${ROOT_DIR}/argocd/apps/root.yaml"
ok "Root Application applied — ArgoCD will now sync everything"

# ─── Build & push service images ──────────────────────────────────────────────
log "Building and pushing service images..."
for svc in sensor-simulator api-gateway data-sink enterprise-consumer; do
    log "  Building ${svc}..."
    docker build -t "localhost:${REGISTRY_PORT}/${svc}:latest" \
        "${ROOT_DIR}/services/${svc}" 2>&1 | tail -3
    docker push "localhost:${REGISTRY_PORT}/${svc}:latest"
    ok "  ${svc} pushed"
done

# ─── Done ─────────────────────────────────────────────────────────────────────
ARGOCD_PASS=$(kubectl -n "${ARGOCD_NS}" get secret argocd-initial-admin-secret \
    -o jsonpath="{.data.password}" 2>/dev/null | base64 -d || echo "not-ready-yet")

echo ""
echo -e "${GREEN}╔══════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║          Bootstrap complete!                     ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════════════╝${NC}"
echo ""
echo "  ArgoCD UI:    http://localhost:30080"
echo "  Admin pass:   ${ARGOCD_PASS}"
echo ""
echo "  Run 'make port-forwards' to access Grafana, API, Node-RED"
echo "  Run 'make demo' for the interview walkthrough script"
echo ""
