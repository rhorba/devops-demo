#!/usr/bin/env bash
# demo.sh — Interview walkthrough script for the technical panel
set -euo pipefail

CYAN='\033[0;36m'; YELLOW='\033[1;33m'; GREEN='\033[0;32m'; BOLD='\033[1m'; NC='\033[0m'

step() { echo -e "\n${BOLD}${CYAN}▶ STEP $1: $2${NC}"; }
show() { echo -e "  ${YELLOW}$*${NC}"; }
note() { echo -e "  ${GREEN}→ $*${NC}"; }

clear
echo -e "${BOLD}╔══════════════════════════════════════════════════════════════════════╗"
echo -e "║       Industrial IoT Telemetry Platform — Demo Walkthrough          ║"
echo -e "║       DXC Germany — Senior DevOps Interview                         ║"
echo -e "╚══════════════════════════════════════════════════════════════════════╝${NC}"

step 1 "Cluster & GitOps — ArgoCD"
note "Everything was deployed via GitOps — zero manual kubectl apply"
show "kubectl get nodes"
kubectl get nodes
show "kubectl get ns"
kubectl get ns
show "argocd app list (or open http://localhost:30080)"
kubectl -n platform get applications.argoproj.io 2>/dev/null || echo "(argocd CLI not logged in — check UI)"

step 2 "Messaging — Multi-broker pipeline"
note "MQTT → Redpanda Connect → Redpanda(Kafka) → NATS → IBM MQ"
show "kubectl -n messaging get pods"
kubectl -n messaging get pods
echo ""
note "Live messages arriving via MQTT sensor simulator:"
show "kubectl -n apps logs -l app=sensor-simulator --tail=5"
kubectl -n apps logs -l app=sensor-simulator --tail=5 2>/dev/null || echo "(pod not ready)"

step 3 "IBM MQ — Enterprise queue"
note "DEV.QUEUE.1 receives bridged messages from Redpanda Connect"
show "kubectl -n messaging exec -it deploy/ibm-mq -- bash -c 'echo \"DIS QL(DEV.QUEUE.1) CURDEPTH\" | runmqsc QM1'"
kubectl -n messaging exec deploy/ibm-mq -- bash -c \
    'echo "DIS QL(DEV.QUEUE.1) CURDEPTH" | runmqsc QM1' 2>/dev/null || echo "(IBM MQ pod not ready)"

step 4 "REST API Gateway — NATS subscriber"
note "api-gateway subscribes to NATS JetStream and exposes REST"
show "curl http://localhost:8080/health"
curl -sf http://localhost:8080/health 2>/dev/null || echo "(port-forward needed: make port-forwards)"
show "curl http://localhost:8080/metrics/latest | jq ."
curl -sf http://localhost:8080/metrics/latest 2>/dev/null | python3 -m json.tool 2>/dev/null || echo "(run: make port-forwards)"

step 5 "LocalStack — AWS services"
note "S3 archival + SQS enterprise queue — no real AWS account needed"
show "aws --endpoint-url=http://localhost:4566 s3 ls s3://telemetry-archive/"
kubectl -n localstack exec deploy/localstack -- \
    awslocal s3 ls s3://telemetry-archive/ 2>/dev/null | tail -5 || echo "(LocalStack pod not ready)"

step 6 "Observability — Prometheus + Grafana + Loki + Thanos"
note "Open http://localhost:3000 — dashboards: 'Message Pipeline' and 'Service Health'"
show "kubectl -n observability get pods"
kubectl -n observability get pods
note "Thanos is uploading Prometheus TSDB blocks to LocalStack S3 (thanos-metrics bucket)"
kubectl -n localstack exec deploy/localstack -- \
    awslocal s3 ls s3://thanos-metrics/ 2>/dev/null | head -5 || true

step 7 "Security — Kyverno admission policies"
note "All pods must be non-root, have resource limits, and carry required labels"
show "kubectl get clusterpolicies"
kubectl -n security get clusterpolicies 2>/dev/null || kubectl get clusterpolicies 2>/dev/null || echo "(Kyverno not ready)"
note "Test — trying to deploy a root container (should be REJECTED):"
show "kubectl run bad-pod --image=nginx:latest -n apps"
kubectl run bad-pod --image=nginx:latest -n apps 2>&1 || echo "↑ Correctly rejected by Kyverno"

step 8 "CI/CD — GitHub Actions"
note "Push to main → lint+test → build images → push to registry → ArgoCD auto-syncs"
show "cat .github/workflows/ci.yml | head -30"
head -30 .github/workflows/ci.yml 2>/dev/null || echo "(see .github/workflows/ci.yml)"

step 9 "Ansible — Infrastructure provisioning"
note "The entire cluster was bootstrapped with Ansible playbooks"
show "ansible-playbook ansible/playbooks/setup-kind.yml --check"
echo "(run in dry-run mode to show idempotent provisioning)"

echo ""
echo -e "${GREEN}${BOLD}Demo complete. Stack covered:${NC}"
echo "  Messaging:      IBM MQ · MQTT/Mosquitto · NATS JetStream · Redpanda (Kafka) · Redpanda Connect"
echo "  Platform:       Kubernetes (Kind) · Helm · ArgoCD GitOps · HAProxy · cert-manager"
echo "  AWS (Local):    S3 · SQS · Kinesis · Secrets Manager (LocalStack)"
echo "  Observability:  Prometheus · Grafana · Loki · Thanos"
echo "  Security:       Kyverno · RBAC · NetworkPolicy · TLS"
echo "  IaC:            Ansible · Bash · Go microservices"
echo "  CI/CD:          GitHub Actions"
echo ""
