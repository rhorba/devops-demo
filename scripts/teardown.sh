#!/usr/bin/env bash
set -euo pipefail
CLUSTER_NAME="devops-demo"
REGISTRY_NAME="kind-registry"

echo "Deleting kind cluster '${CLUSTER_NAME}'..."
kind delete cluster --name "${CLUSTER_NAME}" 2>/dev/null || echo "Cluster not found"

echo "Stopping local registry..."
docker rm -f "${REGISTRY_NAME}" 2>/dev/null || echo "Registry not running"

echo "Teardown complete."
