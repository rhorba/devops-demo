# Architecture: Industrial IoT Telemetry Platform
**PRD Reference**: docs/prd.md | **Version**: 1.0 | **Date**: 2026-05-02 | **Author**: Tech Lead

## 1. Overview
A multi-broker messaging pipeline deployed on a local Kind (Kubernetes-in-Docker) cluster,
managed via ArgoCD GitOps. Sensor data flows through MQTT → Redpanda Connect → Redpanda/Kafka
→ NATS JetStream → IBM MQ, sinks to LocalStack S3/SQS, and is observed via Prometheus/Grafana/Loki/Thanos.

## 2. Architecture Decision Records

### ADR-1: Local Kubernetes via Kind (not Docker Compose)
- **Context**: Must demonstrate Kubernetes/Helm/ArgoCD proficiency for a senior role
- **Decision**: Kind (Kubernetes-in-Docker) — single-node, reproducible, no VM overhead
- **Alternatives rejected**: Docker Compose (too simple for senior demo), k3d (less familiar tooling)
- **Consequences**: Requires kind-config.yaml, local registry, 8GB RAM minimum

### ADR-2: Redpanda as Kafka drop-in
- **Context**: Job spec lists "Kafka (protocol level, not cluster)" — Redpanda is protocol-compatible
- **Decision**: Redpanda replaces Kafka — same protocol, simpler ops, better local resource usage
- **Alternatives rejected**: Full Kafka+ZK (heavy), Strimzi operator (extra complexity)
- **Consequences**: Use rpk CLI instead of kafka-topics; all Kafka clients work unchanged

### ADR-3: Redpanda Connect as the universal bridge
- **Context**: Job spec lists Redpanda Connect / Wombat / Benthos as Expert/Master — they are the same tool
- **Decision**: Redpanda Connect (formerly Benthos/Wombat) handles all broker-to-broker routing
- **Pipelines**: MQTT→Redpanda, Redpanda→NATS, Redpanda→IBM MQ, Redpanda→S3
- **Consequences**: Single config format (YAML), supports all required protocols natively

### ADR-4: LocalStack for AWS simulation
- **Context**: Must show AWS proficiency without a real account
- **Decision**: LocalStack (community) simulates S3, SQS, Kinesis, Secrets Manager
- **Endpoint**: http://localstack.localstack.svc.cluster.local:4566
- **Consequences**: AWS SDK just needs endpoint override; no credentials required

### ADR-5: ArgoCD App-of-Apps pattern
- **Context**: Must show GitOps discipline — no manual kubectl applies
- **Decision**: One root ArgoCD Application (`argocd/apps/root.yaml`) that bootstraps all others
- **Consequences**: Adding a new service = add an Application YAML to argocd/apps/ and push

### ADR-6: Thanos sidecar (not Thanos Receive)
- **Context**: Must show Thanos for HA and long-term storage
- **Decision**: Thanos sidecar on kube-prometheus-stack, uploads blocks to LocalStack S3
- **Consequences**: No Thanos Receive/Compact needed at this scale; demonstrates the pattern

### ADR-7: Kyverno for admission policies
- **Context**: Must show security hardening — no root containers, resource limits required
- **Decision**: Kyverno (YAML-native policies, easy to show) over OPA/Gatekeeper
- **Consequences**: Policies run at admission time; violating pods are rejected with clear messages

## 3. System Design

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        KIND CLUSTER (local k8s)                             │
│                                                                             │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │  namespace: apps                                                     │   │
│  │                                                                      │   │
│  │  [sensor-simulator (Go)]──MQTT──►[Mosquitto]                        │   │
│  │                                      │                              │   │
│  │                              MQTT subscribe                         │   │
│  │                                      ▼                              │   │
│  │                          [Redpanda Connect]                         │   │
│  │                         /       |        \                          │   │
│  │               Kafka proto    NATS pub   IBM MQ pub                  │   │
│  │                    ▼            ▼            ▼                      │   │
│  │  namespace: messaging           │            │                      │   │
│  │  [Redpanda/Kafka]  [NATS JetStream]    [IBM MQ]                    │   │
│  │         │                │                   │                     │   │
│  │         │           NATS subscribe      MQ get                     │   │
│  │         │                ▼                   ▼                     │   │
│  │  Kafka  │    [api-gateway (Go)]   [enterprise-consumer (Go)]       │   │
│  │  consume│         │REST API               │SQS put                 │   │
│  │         ▼         │ GET /metrics/latest   ▼                        │   │
│  │  [data-sink(Go)]  │           [LocalStack SQS]                     │   │
│  │         │S3 put   │                                                │   │
│  │         ▼         │                                                │   │
│  │  [LocalStack S3]  │  namespace: localstack                         │   │
│  │                   │                                                │   │
│  └───────────────────┼────────────────────────────────────────────────┘   │
│                      │                                                     │
│  ┌───────────────────▼────────────────────────────────────────────────┐   │
│  │  namespace: platform                                               │   │
│  │  [HAProxy Ingress] ◄── routes: grafana, argocd, api, node-red     │   │
│  │  [ArgoCD]          ◄── GitOps: manages all namespaces              │   │
│  │  [cert-manager]    ◄── TLS: self-signed CA for cluster             │   │
│  └────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  namespace: observability                                           │   │
│  │  [Prometheus] ──scrapes all── [Thanos Sidecar] ──S3──► [LocalStack]│   │
│  │  [Loki] ◄── [Promtail (DaemonSet)]                                 │   │
│  │  [Grafana] ◄── datasources: Prometheus, Loki, Thanos               │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  namespace: tools                                                   │   │
│  │  [Node-RED] ── visual flow: NATS subscribe + dashboard              │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  [Kyverno] ── admission: no-root, resource-limits, required-labels         │
│  [NetworkPolicy] ── namespace isolation, allow rules only                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

GitHub Actions (CI)
  ├── lint + test Go services
  ├── build Docker images → push to local registry (kind load)
  └── update Helm chart appVersion → ArgoCD detects → auto-sync

Ansible (provisioning)
  ├── playbooks/setup-kind.yml     → create cluster, load images
  ├── playbooks/deploy-infra.yml   → install ArgoCD, bootstrap root app
  └── playbooks/configure-mq.yml  → IBM MQ queue + channel setup
```

## 4. Message Flow (numbered)

```
1. sensor-simulator  → MQTT publish → topic: sensors/{location}/{type}
2. Mosquitto         → broker: receives, retains, fans out
3. Redpanda Connect  → pipeline mqtt-to-redpanda: subscribe MQTT → produce Kafka topic raw.telemetry
4. Redpanda          → topic: raw.telemetry (Kafka protocol, 3 partitions)
5. Redpanda Connect  → pipeline redpanda-to-nats: consume Kafka → publish NATS subject telemetry.>
6. NATS JetStream    → stream: TELEMETRY, subject: telemetry.>
7. api-gateway       → NATS subscribe → serves GET /metrics/latest
8. Redpanda Connect  → pipeline redpanda-to-ibmmq: consume Kafka → put IBM MQ queue DEV.QUEUE.1
9. IBM MQ            → queue: DEV.QUEUE.1
10. enterprise-consumer → MQ get → SQS send → LocalStack SQS queue: enterprise-events
11. Redpanda Connect  → pipeline redpanda-to-s3: consume Kafka → S3 put → LocalStack S3 bucket: telemetry-archive
12. Thanos sidecar   → uploads Prometheus TSDB blocks → LocalStack S3 bucket: thanos-metrics
```

## 5. API Design
| Method | Endpoint | Description |
|---|---|---|
| GET | /health | Service health |
| GET | /metrics/latest | Last 10 telemetry readings from NATS |
| GET | /metrics/stream | SSE stream of live telemetry |

## 6. Security Architecture
- **TLS**: cert-manager self-signed CA; all in-cluster services use TLS where supported
- **RBAC**: Per-service ServiceAccounts with minimal ClusterRole bindings
- **Secrets**: External Secrets Operator pulls from LocalStack Secrets Manager
- **Admission**: Kyverno enforces no-root, no-latest-tag, resource limits, required labels
- **Network**: NetworkPolicy denies all cross-namespace traffic except explicit allow rules
- **IBM MQ**: MQSSLKEYR + channel authentication record enabled

## 7. Infrastructure
- **Cluster**: Kind (1 control-plane + 2 worker nodes)
- **Registry**: Local Docker registry on localhost:5001 (mapped into Kind)
- **GitOps**: ArgoCD v2.x — App of Apps rooted at argocd/apps/root.yaml
- **Ingress**: HAProxy (NodePort 80/443 → cluster services)
- **Secrets**: LocalStack Secrets Manager + External Secrets Operator
- **CI/CD**: GitHub Actions → build → kind load → ArgoCD sync
- **Provisioning**: Ansible 2.15+ (community.kubernetes collection)

## 8. Project Structure
```
devops-demo/
├── .github/workflows/          # CI: ci.yml, release.yml
├── ansible/
│   ├── inventory/hosts.yml
│   ├── playbooks/
│   │   ├── setup-kind.yml
│   │   ├── deploy-infra.yml
│   │   └── configure-mq.yml
│   └── roles/kind/ ibmmq/ localstack/
├── argocd/
│   ├── apps/                   # Application manifests (App of Apps)
│   └── projects/               # AppProject definitions
├── charts/                     # Helm charts (one per service)
│   ├── mosquitto/
│   ├── nats/
│   ├── redpanda/
│   ├── redpanda-connect/
│   ├── ibm-mq/
│   ├── localstack/
│   ├── haproxy/
│   └── node-red/
├── services/                   # Go microservices
│   ├── sensor-simulator/
│   ├── api-gateway/
│   ├── data-sink/
│   └── enterprise-consumer/
├── monitoring/
│   ├── grafana/dashboards/     # Provisioned JSON dashboards
│   └── prometheus/rules/       # Alerting rules
├── security/
│   ├── kyverno-policies/
│   ├── rbac/
│   ├── cert-manager/
│   └── network-policies/
├── scripts/
│   ├── bootstrap.sh
│   ├── demo.sh
│   └── teardown.sh
├── node-red/flows.json
├── kind-config.yaml
├── Makefile
└── README.md
```

## 9. Technical Risks
| Risk | Mitigation | Owner |
|---|---|---|
| IBM MQ image size (~1.5GB) | Pre-pull in bootstrap, document | DevOps |
| RAM pressure (8-10GB total) | Resource limits, optional Thanos disable | Tech Lead |
| LocalStack Secrets Manager (pro feature) | Fall back to k8s Secrets if not available | Backend |
| Redpanda Connect MQTT TLS | Use non-TLS for demo simplicity | DevOps |
