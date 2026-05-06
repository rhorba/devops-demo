# PRD: Senior DevOps Demo — Industrial IoT Telemetry Platform
**Version**: 1.0 | **Date**: 2026-05-02 | **Author**: PM | **Status**: Approved

## 1. Problem Statement
A senior technical panel at DXC Germany needs to evaluate a candidate's Expert/Master-level competency
across messaging middleware, cloud-native infrastructure, GitOps, and observability. The demo must show
*depth* — not just awareness — of the full stack listed in the job spec, running as a coherent system.

## 2. Goals & Success Metrics
| Goal | Metric | Target |
|---|---|---|
| Demonstrate messaging expertise | Brokers running & exchanging messages | IBM MQ + MQTT + NATS + Redpanda all active |
| Show GitOps discipline | ArgoCD managing all deployments | Zero manual kubectl apply |
| Prove observability skills | Grafana dashboards live | Metrics + Logs + Thanos HA |
| Evidence security hardening | Kyverno policies enforced | All pods pass admission |
| Run fully locally | No cloud account needed | Kind cluster on laptop |

## 3. Demo Narrative
**"Industrial IoT Telemetry Pipeline"** — a factory floor sensor network where:
- Go-based sensor simulators publish telemetry via **MQTT**
- **Redpanda Connect (Benthos/Wombat)** routes and transforms events across brokers
- Data flows through **NATS JetStream** (real-time), **Redpanda/Kafka** (streaming), and **IBM MQ** (enterprise)
- A REST API gateway exposes processed data
- Everything sinks to **LocalStack S3** (archival) and **SQS** (enterprise queue)
- Full observability: Prometheus + Grafana + Loki + Thanos

## 4. Scope
### In Scope
- Kind (local k8s) cluster with all services
- IBM MQ, Eclipse Mosquitto, NATS JetStream, Redpanda (Kafka), Redpanda Connect
- LocalStack simulating: S3, SQS, Kinesis, Secrets Manager
- ArgoCD GitOps (App-of-Apps pattern)
- Helm charts for every service
- 4 Go microservices (sensor-sim, api-gateway, data-sink, enterprise-consumer)
- Prometheus + Grafana + Loki + Thanos observability stack
- HAProxy ingress
- Ansible provisioning playbooks
- Kyverno admission policies + RBAC + Network Policies + cert-manager
- GitHub Actions CI/CD (build + push to local registry)
- Node-RED visual flow editor
- Makefile one-command bootstrap

### Out of Scope
- Real AWS account (LocalStack only)
- Production HA multi-node cluster
- Real IBM MQ license (use Developer edition free image)
- Frontend UI (REST API + Grafana is sufficient)
- Authentication layer on API (not in spec)

## 5. Requirements
### Functional
- FR-1: Sensor simulator publishes 1 msg/sec per topic via MQTT
- FR-2: Redpanda Connect bridges MQTT → Redpanda → NATS → IBM MQ
- FR-3: api-gateway exposes GET /metrics/latest and GET /health
- FR-4: data-sink consumes from Redpanda and writes to LocalStack S3
- FR-5: enterprise-consumer reads IBM MQ queue, forwards to LocalStack SQS
- FR-6: ArgoCD syncs all apps from Git (no manual deploys)
- FR-7: Grafana shows message throughput, latency, and service health
- FR-8: Loki aggregates all service logs, queryable in Grafana
- FR-9: Thanos stores Prometheus metrics in LocalStack S3

### Non-Functional
- NFR-1: All services start within `make bootstrap` (≤15 min cold start)
- NFR-2: No container runs as root (Kyverno enforced)
- NFR-3: All inter-service traffic uses TLS (cert-manager self-signed CA)
- NFR-4: Every service has resource limits defined in Helm values

## 6. Constraints & Assumptions
- Windows 11 host with Docker Desktop + kind installed
- kubectl, helm, argocd CLI available
- 16GB RAM minimum (full stack is ~8-10GB)
- GitHub repo is the GitOps source of truth

## 7. Risks
| Risk | Probability | Impact | Mitigation |
|---|---|---|---|
| IBM MQ image pull slow | M | L | Pre-pull in bootstrap |
| Kind cluster RAM pressure | M | H | Resource limits on all pods |
| LocalStack cold start | L | M | Health-check in bootstrap |
| Redpanda Connect MQTT auth | L | M | Config tested in unit mode |

## 8. Timeline
| Milestone | Target |
|---|---|
| PRD + Architecture + Stories | Day 1 morning |
| Scaffold + Go services | Day 1 afternoon |
| Helm charts + ArgoCD | Day 2 morning |
| Observability + Security | Day 2 afternoon |
| Ansible + CI/CD + README | Day 3 |
| Full demo run-through | Day 3 evening |
