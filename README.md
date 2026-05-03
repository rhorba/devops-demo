# DevOps Demo — Industrial IoT Telemetry Platform

> Senior DevOps Engineer showcase for DXC Germany — Expert/Master level messaging, GitOps, cloud-native observability.

## Architecture

```
[sensor-simulator (Go)]──MQTT──►[Mosquitto]
                                     │
                          [Redpanda Connect/Benthos/Wombat]
                         /           │              \
              Kafka proto         NATS pub        IBM MQ put
                   ▼                 ▼                ▼
           [Redpanda/Kafka]  [NATS JetStream]    [IBM MQ]
                │                   │                │
           Kafka consume       NATS subscribe    MQ get
                │                   ▼                ▼
                │       [api-gateway (Go)]  [enterprise-consumer (Go)]
                │         GET /metrics/latest   SQS put
                │                              [LocalStack SQS]
                ▼
         [data-sink (Go)]──S3 put──►[LocalStack S3]

All services → Prometheus → Grafana + Loki + Thanos → LocalStack S3
All deployments → ArgoCD GitOps ← GitHub Actions CI/CD
Admission control: Kyverno (no-root, resource-limits, required-labels)
Ingress: HAProxy → Grafana · ArgoCD · API · Node-RED
Provisioning: Ansible (setup-kind · deploy-infra · configure-mq)
```

## Tech Stack vs Job Spec

| Category | Tool | Spec Level | Demonstrated By |
|---|---|---|---|
| Messaging | IBM MQ | Expert/Master | `enterprise-consumer` → DEV.QUEUE.1 |
| Messaging | MQTT | Expert/Master | `sensor-simulator` → Mosquitto broker |
| Messaging | NATS JetStream | Expert/Master | `api-gateway` subscriber + SSE |
| Messaging | Redpanda Connect/Benthos | Expert/Master | 4 routing pipeline YAMLs |
| Messaging | Kafka (protocol) | Adv. Beginner | Redpanda drop-in, `rpk` CLI |
| Cloud/Platform | ArgoCD | Proficient | App of Apps, auto-sync, self-heal |
| Cloud/Platform | AWS S3/SQS/Kinesis/SecretsMgr | Proficient | LocalStack |
| Cloud/Platform | Docker | Proficient | Multi-stage Dockerfiles (scratch/debian) |
| Cloud/Platform | Kubernetes | Proficient | Kind 3-node cluster |
| Cloud/Platform | Helm | Proficient | 8 charts + 4 service charts |
| Observability | Grafana | Expert/Master | Message Pipeline dashboard |
| Observability | Prometheus | Proficient | Custom metrics + PrometheusRules |
| Observability | Loki | Proficient | Promtail DaemonSet, log correlation |
| Observability | Thanos | Adv. Beginner | Sidecar → LocalStack S3 |
| Scripting/Dev | Go/Golang | Proficient | 4 microservices, ~800 lines |
| Scripting/Dev | Bash | Adv. Beginner | bootstrap/teardown/demo scripts |
| Scripting/Dev | REST APIs | Proficient | api-gateway HTTP + SSE endpoints |
| Systems/Networks | HAProxy | Adv. Beginner | Ingress controller |
| Tools/Methods | Git | Proficient | GitOps source of truth |
| Tools/Methods | Ansible | Adv. Beginner | 3 idempotent playbooks |
| Additional | Microservices | Proficient | 4 decoupled Go services |
| Additional | Node-RED | Proficient | Visual NATS + MQTT live flows |

## Prerequisites

```bash
# Required tools
docker          >= 24.0
kind            >= 0.22
kubectl         >= 1.29
helm            >= 3.14
argocd          >= 2.10   # CLI
ansible         >= 2.15   # optional — bootstrap.sh works without it
golangci-lint   >= 1.57   # for: make lint
```

Add to `C:\Windows\System32\drivers\etc\hosts` (Windows) or `/etc/hosts` (Linux/Mac):
```
127.0.0.1  grafana.demo.local
127.0.0.1  argocd.demo.local
127.0.0.1  api.demo.local
127.0.0.1  nodered.demo.local
```

## Quick Start — One Command

```bash
make bootstrap
# Takes ~12-15 min (IBM MQ image pull is the slowest step)
```

## Or Step-by-Step with Ansible

```bash
ansible-playbook -i ansible/inventory/hosts.yml ansible/playbooks/setup-kind.yml
ansible-playbook -i ansible/inventory/hosts.yml ansible/playbooks/deploy-infra.yml
ansible-playbook -i ansible/inventory/hosts.yml ansible/playbooks/configure-mq.yml
```

## Access Services

```bash
make port-forwards   # start all port-forwards in background

# Open:
#   Grafana:   http://localhost:3000   (admin / admin)
#   ArgoCD:    http://localhost:30080  (admin / make argocd-password)
#   API:       http://localhost:8080
#   Node-RED:  http://localhost:1880
```

## Demo Walkthrough

```bash
make demo    # prints step-by-step interview script
```

## Key API Endpoints

| Endpoint | Description |
|---|---|
| `GET /health` | Service health + NATS connection status |
| `GET /metrics/latest` | Last 10 telemetry readings (JSON) |
| `GET /metrics/stream` | Live SSE stream of sensor data |

## Verify the Message Pipeline

```bash
# 1. Watch MQTT sensor data publishing
kubectl -n apps logs -f deploy/sensor-simulator

# 2. Confirm Kafka/Redpanda has messages
kubectl -n messaging exec -it deploy/redpanda-0 -- \
  rpk topic consume raw.telemetry --num=5

# 3. Check IBM MQ queue depth
kubectl -n messaging exec -it deploy/ibm-mq -- bash -c \
  "echo 'DISPLAY QL(DEV.QUEUE.1) CURDEPTH' | runmqsc QM1"

# 4. See S3 archival in LocalStack
kubectl -n localstack exec deploy/localstack -- \
  awslocal s3 ls s3://telemetry-archive/ --recursive | tail -10

# 5. SQS message count
kubectl -n localstack exec deploy/localstack -- \
  awslocal sqs get-queue-attributes \
  --queue-url http://localhost:4566/000000000000/enterprise-events \
  --attribute-names ApproximateNumberOfMessages

# 6. Live API stream
curl -N http://localhost:8080/metrics/stream
```

## CI/CD

```bash
git push origin main         # → lint + test + build + push to GHCR → ArgoCD syncs
git tag v1.0.0 && git push   # → versioned images + chart bump + GitHub Release
```

## Kyverno Security Policies

```bash
# Test: try to run a root container — should be REJECTED
kubectl run bad-pod --image=nginx:latest -n apps

# List all policies
kubectl get clusterpolicies
```

## Teardown

```bash
make teardown   # deletes kind cluster + local registry
```

## Project Structure

```
devops-demo/
├── .github/workflows/          # CI (lint+test+push) + Release (tag+chart bump)
├── ansible/                    # Idempotent provisioning
│   ├── inventory/hosts.yml
│   └── playbooks/              # setup-kind | deploy-infra | configure-mq
├── argocd/
│   ├── apps/                   # App of Apps (root + messaging + platform + security + tools)
│   └── projects/demo.yaml      # AppProject with source/destination restrictions
├── charts/                     # Helm charts
│   ├── mosquitto/              # MQTT broker
│   ├── nats/                   # NATS JetStream
│   ├── redpanda/               # Kafka-compatible event streaming + console
│   ├── redpanda-connect/       # 4 routing pipelines (MQTT→Kafka→NATS→IBM MQ→S3)
│   ├── ibm-mq/                 # IBM MQ Developer edition (icr.io/ibm-messaging/mq)
│   ├── localstack/             # AWS S3 + SQS + Kinesis + Secrets Manager
│   ├── haproxy/                # Ingress controller
│   └── node-red/               # Visual NATS + MQTT flow editor
├── services/                   # Go microservices
│   ├── sensor-simulator/       # MQTT publisher (6 zones × 3 types), Prometheus metrics
│   ├── api-gateway/            # NATS JetStream subscriber + REST + SSE + ring buffer
│   ├── data-sink/              # Kafka consumer (franz-go) → LocalStack S3 archival
│   └── enterprise-consumer/    # IBM MQ consumer (ibmmq CGO) → LocalStack SQS
├── monitoring/
│   ├── prometheus/             # kube-prometheus-stack + Thanos sidecar + alerting rules
│   ├── loki/                   # Loki + Promtail DaemonSet
│   └── grafana/dashboards/     # Message Pipeline Overview dashboard (JSON)
├── security/
│   ├── kyverno-policies/       # no-root | resource-limits | no-latest | required-labels
│   ├── rbac/                   # Per-service ServiceAccounts + minimal Roles
│   └── network-policies/       # Default deny + explicit allow rules
├── scripts/
│   ├── bootstrap.sh            # Full idempotent cluster setup
│   ├── demo.sh                 # Interview walkthrough script
│   └── teardown.sh             # Clean cluster deletion
├── docs/
│   ├── prd.md                  # Product Requirements Document
│   ├── architecture.md         # 9 ADRs + system design diagram + message flow
│   └── stories.md              # 9 Epics · 28 Stories · 4 Sprints
├── kind-config.yaml            # 3-node Kind cluster (1 control-plane + 2 workers)
└── Makefile                    # bootstrap | teardown | images | lint | test | demo
```

---

## CTS Session Playbook

This project uses the **Claude Team Skills (CTS)** framework for all AI-assisted work.
Follow this playbook every session — do not skip steps.

### One-time setup (run once per project)

```
initialize project logs
```

Creates `.logs/` with all 7 tracking files (activity, decisions, issues, risks, corrections, communications, sessions).

---

### Every session — opening sequence

```
1.  resume last session
```
I read the last `SESSION_END` entry and tell you what was done, what is next, and any open blockers.
If no logs exist I will say so and start fresh.

```
2.  [confirm or redirect]
    e.g. "yes continue"  /  "skip that, today we are doing X instead"
```

---

### Starting a task

```
3.  [describe the work in plain language]
    e.g. "fix the 11 bugs"
         "add auth to the api"
         "security audit the infra"
         "plan sprint 2"
```

I run **UNDERSTAND** — ask you 3 questions (what / where / how big) then route to the right specialists.

```
4.  answer my 3 understand questions
```

---

### Brainstorm + plan approval — do NOT skip

```
5.  [pick an option after I present them]
    e.g. "go with balanced"  /  "option 2"  /  "mix 1 and 3"
```

I show you `🟢 SIMPLE` / `🟡 BALANCED` / `🔴 COMPREHENSIVE` options. You pick one. I log the decision.

```
6.  approve the plan
    e.g. "plan looks good, start batch 1"
         "adjust batch 2 — drop task X"
```

I show the full batch plan with estimates. You approve or change it **before any code is written**.

---

### During execution — batch-by-batch gates

```
7.  After each fix / task:   "good, next"   or   "adjust: [what to change]"
8.  After each batch:        "verify batch [N]"
9.  If blocked:              "show options"   (I give A / B / C choices)
```

---

### Specialist override

Force a specific skill at any time:

```
"load tech lead"
"load security engineer"
"load tester"
"load devops"
"load backend dev"
```

---

### Logging and status on demand

```
"show project status"          → full status report from .logs/
"what did we decide"           → reads .logs/decisions.md
"show open issues"             → reads .logs/issues.md
"show risks"                   → reads .logs/risks.md
"what changed from the plan"   → reads .logs/corrections.md
"run retrospective"            → reads all logs, generates retro
```

---

### Closing a session — ALWAYS end with this

```
10.  end session
```

I write `SESSION_END` to `.logs/sessions.md` with: what was completed, what is in progress, blockers, and next steps.
**Without this, resumption next session is blind.**

---

### Full strict sequence at a glance

```
SESSION START
─────────────────────────────────────────────────
1.   resume last session
2.   confirm / redirect
3.   [describe the task]
4.   answer understand questions
5.   pick brainstorm option  (green / yellow / red)
6.   approve the batch plan

EXECUTION  (repeat per batch)
─────────────────────────────────────────────────
7.   "good, next"       after each task
8.   "verify batch N"   after each batch
9.   "show options"     if blocked

SESSION END
─────────────────────────────────────────────────
10.  end session
```
