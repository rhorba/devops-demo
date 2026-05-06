#!/usr/bin/env bash
# registry.sh — start local Docker registry and wire it into the kind cluster
set -euo pipefail

REGISTRY_NAME="kind-registry"
REGISTRY_PORT="5001"

RED='\033[0;31m'; GREEN='\033[0;32m'; CYAN='\033[0;36m'; NC='\033[0m'
log() { echo -e "${CYAN}[registry]${NC} $*"; }
ok()  { echo -e "${GREEN}[✓]${NC} $*"; }
die() { echo -e "${RED}[✗]${NC} $*"; exit 1; }

# ─── Start registry container ─────────────────────────────────────────────────
if ! docker ps --format '{{.Names}}' | grep -q "^${REGISTRY_NAME}$"; then
    log "Starting local registry on port ${REGISTRY_PORT}..."
    docker run -d --restart=always \
        -p "127.0.0.1:${REGISTRY_PORT}:5000" \
        --network bridge \
        --name "${REGISTRY_NAME}" \
        registry:2
    ok "Registry started"
else
    ok "Registry '${REGISTRY_NAME}' already running"
fi

# ─── Connect registry to kind network ────────────────────────────────────────
if docker network inspect kind &>/dev/null; then
    if ! docker network inspect kind 2>/dev/null | grep -q "${REGISTRY_NAME}"; then
        log "Connecting registry to kind network..."
        docker network connect kind "${REGISTRY_NAME}" 2>/dev/null || true
        ok "Registry connected to kind network"
    else
        ok "Registry already connected to kind network"
    fi
else
    log "kind network not found — cluster may not be running yet"
fi

# ─── Apply registry ConfigMap (containerd discovery) ─────────────────────────
if kubectl cluster-info &>/dev/null 2>&1; then
    log "Applying registry ConfigMap..."
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
else
    log "Cluster not reachable — skipping ConfigMap (run after cluster is up)"
fi
