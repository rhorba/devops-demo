# Stories: Industrial IoT Telemetry Platform
**PRD**: docs/prd.md | **Architecture**: docs/architecture.md

---

## Epic 1: Project Scaffold
*Foundational structure that all other epics depend on*

### Story 1.1: Kind cluster config + local registry
**Priority**: Must | **Size**: S | **Specialist**: DevOps
- kind-config.yaml (1 control-plane + 2 workers + extraPortMappings)
- Local Docker registry on localhost:5001
- Registry patch ConfigMap for Kind

### Story 1.2: Makefile with all top-level targets
**Priority**: Must | **Size**: S | **Specialist**: DevOps
- Targets: bootstrap, teardown, demo, images, lint, test, argocd-password
- One-command setup: `make bootstrap`

### Story 1.3: .gitignore + repo initialization
**Priority**: Must | **Size**: XS | **Specialist**: DevOps

---

## Epic 2: Go Microservices
*The application layer — 4 services in Go*

### Story 2.1: sensor-simulator
**Priority**: Must | **Size**: M | **Specialist**: Backend Dev
- Publishes simulated telemetry (temperature, pressure, humidity) via MQTT
- Topics: `sensors/{zone}/{type}` (e.g. sensors/zone-a/temperature)
- Rate: 1 msg/sec per sensor, configurable via env
- Prometheus metrics: messages_published_total, publish_errors_total
- Dockerfile (scratch base, non-root)

### Story 2.2: api-gateway
**Priority**: Must | **Size**: M | **Specialist**: Backend Dev
- NATS JetStream subscriber on `telemetry.>`
- REST endpoints: GET /health, GET /metrics/latest, GET /metrics/stream (SSE)
- In-memory ring buffer of last 100 messages
- Prometheus metrics: http_requests_total, nats_messages_received_total
- Dockerfile (scratch base, non-root)

### Story 2.3: data-sink
**Priority**: Must | **Size**: M | **Specialist**: Backend Dev
- Kafka consumer (Redpanda) on topic `raw.telemetry`
- Writes JSON batches to LocalStack S3 bucket `telemetry-archive`
- Key: `{date}/{hour}/{uuid}.json`
- Prometheus metrics: records_sinked_total, s3_put_errors_total
- Dockerfile (scratch base, non-root)

### Story 2.4: enterprise-consumer
**Priority**: Must | **Size**: M | **Specialist**: Backend Dev
- IBM MQ consumer on queue `DEV.QUEUE.1`
- Forwards messages to LocalStack SQS queue `enterprise-events`
- Prometheus metrics: mq_messages_consumed_total, sqs_send_errors_total
- Dockerfile (scratch base, non-root)

---

## Epic 3: Helm Charts — Messaging Layer
*Packaging all messaging components*

### Story 3.1: Mosquitto chart
**Priority**: Must | **Size**: S | **Specialist**: DevOps
- MQTT broker, port 1883 (internal) + 9001 (WebSocket)
- ConfigMap for mosquitto.conf
- No anonymous access for production values

### Story 3.2: NATS chart (values overlay)
**Priority**: Must | **Size**: S | **Specialist**: DevOps
- nats.io Helm chart with JetStream enabled
- Stream: TELEMETRY, subjects: telemetry.>, retention: limits, maxAge: 24h

### Story 3.3: Redpanda chart (values overlay)
**Priority**: Must | **Size**: S | **Specialist**: DevOps
- Redpanda single-broker values
- Topic: raw.telemetry (3 partitions, 1 replica)
- Schema registry enabled (for demo)

### Story 3.4: Redpanda Connect chart + pipeline configs
**Priority**: Must | **Size**: M | **Specialist**: DevOps
- 4 pipeline ConfigMaps:
  - mqtt-to-redpanda.yaml
  - redpanda-to-nats.yaml
  - redpanda-to-ibmmq.yaml
  - redpanda-to-s3.yaml
- Deployment with pipeline selector via env

### Story 3.5: IBM MQ chart
**Priority**: Must | **Size**: S | **Specialist**: DevOps
- icr.io/ibm-messaging/mq:latest (Developer edition — free)
- MQ_QMGR_NAME: QM1, MQ_APP_PASSWORD set via secret
- Queue DEV.QUEUE.1 pre-created via MQ config

---

## Epic 4: Platform & GitOps
*ArgoCD, HAProxy, LocalStack, cert-manager*

### Story 4.1: LocalStack chart
**Priority**: Must | **Size**: S | **Specialist**: DevOps
- localstack/localstack Helm chart
- Services: s3, sqs, kinesis, secretsmanager
- Init script creates buckets + queues on startup

### Story 4.2: HAProxy ingress chart
**Priority**: Must | **Size**: S | **Specialist**: DevOps
- Routes: grafana.demo.local, argocd.demo.local, api.demo.local, nodered.demo.local
- /etc/hosts entry in README

### Story 4.3: ArgoCD App-of-Apps
**Priority**: Must | **Size**: M | **Specialist**: DevOps
- argocd/projects/demo.yaml (AppProject)
- argocd/apps/root.yaml (root Application)
- argocd/apps/{mosquitto,nats,redpanda,ibmmq,redpanda-connect,localstack,
                haproxy,apps,observability,security,tools}.yaml
- Auto-sync + self-heal enabled

### Story 4.4: cert-manager setup
**Priority**: Should | **Size**: S | **Specialist**: DevOps
- ClusterIssuer: selfsigned-ca
- Certificate resources per namespace

---

## Epic 5: Observability Stack
*Prometheus + Grafana + Loki + Thanos*

### Story 5.1: kube-prometheus-stack values
**Priority**: Must | **Size**: M | **Specialist**: DevOps
- Thanos sidecar on Prometheus
- Thanos object store config → LocalStack S3 bucket: thanos-metrics
- ServiceMonitor for all custom services
- Alertmanager rules: high message lag, pod restarts

### Story 5.2: Loki + Promtail values
**Priority**: Must | **Size**: S | **Specialist**: DevOps
- Loki single-binary mode
- Promtail DaemonSet scraping /var/log/pods

### Story 5.3: Grafana dashboards
**Priority**: Must | **Size**: M | **Specialist**: DevOps
- Dashboard 1: Message Pipeline Overview
  - Panels: msg/sec per broker, end-to-end latency, error rate
- Dashboard 2: Service Health
  - Panels: pod CPU/RAM, restarts, HTTP req rate
- Dashboard 3: Infrastructure
  - Panels: node resources, PVC usage, network

---

## Epic 6: Security Hardening
*Kyverno, RBAC, NetworkPolicy*

### Story 6.1: Kyverno policies
**Priority**: Must | **Size**: M | **Specialist**: Security
- policy-no-root.yaml: all containers must runAsNonRoot
- policy-resource-limits.yaml: all containers must have limits
- policy-no-latest-tag.yaml: image tags must not be latest
- policy-required-labels.yaml: app + version labels required

### Story 6.2: RBAC per service
**Priority**: Must | **Size**: S | **Specialist**: Security
- ServiceAccount per Go service
- Role/ClusterRole with minimal permissions (no wildcard)

### Story 6.3: Network Policies
**Priority**: Must | **Size**: S | **Specialist**: Security
- Default deny-all ingress per namespace
- Explicit allow rules: apps→messaging, apps→localstack, observability→all

---

## Epic 7: Ansible Provisioning
*Infrastructure-as-code for cluster setup*

### Story 7.1: setup-kind.yml playbook
**Priority**: Must | **Size**: M | **Specialist**: DevOps
- Tasks: check prerequisites, create kind cluster, start local registry, load base images
- Idempotent (safe to re-run)

### Story 7.2: deploy-infra.yml playbook
**Priority**: Must | **Size**: M | **Specialist**: DevOps
- Tasks: install ArgoCD via Helm, wait for ArgoCD ready, apply root Application, wait for sync
- Tags: argocd, bootstrap

### Story 7.3: configure-mq.yml playbook
**Priority**: Should | **Size**: S | **Specialist**: DevOps
- Tasks: exec into MQ pod, run runmqsc to verify queues, set channel auth

---

## Epic 8: CI/CD Pipeline
*GitHub Actions*

### Story 8.1: ci.yml — build + test
**Priority**: Must | **Size**: M | **Specialist**: DevOps
- Trigger: push to main, PR
- Jobs: lint (golangci-lint), test (go test), build images, push to GHCR
- Matrix: all 4 Go services

### Story 8.2: release.yml — tag + update charts
**Priority**: Should | **Size**: S | **Specialist**: DevOps
- Trigger: tag v*.*.*
- Updates appVersion in Chart.yaml
- ArgoCD detects new image tag → auto-syncs

---

## Epic 9: Developer Experience
*Scripts, Node-RED, README*

### Story 9.1: bootstrap.sh + demo.sh + teardown.sh
**Priority**: Must | **Size**: S | **Specialist**: DevOps
- bootstrap.sh: wraps `make bootstrap` with pre-flight checks
- demo.sh: step-by-step demo narrative (print what to show the panel)
- teardown.sh: `make teardown` with cluster delete

### Story 9.2: Node-RED flows
**Priority**: Should | **Size**: S | **Specialist**: DevOps
- Flow 1: NATS subscribe → debug output (live message viewer)
- Flow 2: HTTP inject → MQTT publish (manual trigger for demo)

### Story 9.3: README
**Priority**: Must | **Size**: M | **Specialist**: DevOps
- Architecture diagram (ASCII)
- Prerequisites + install steps
- make bootstrap usage
- Demo walkthrough script
- Tech stack table mapped to job spec

---

## Sprint Allocation
| Sprint | Epics | Stories | Est. |
|---|---|---|---|
| Sprint 1 (Day 1) | 1, 2 | 1.1–1.3, 2.1–2.4 | 6h |
| Sprint 2 (Day 2 AM) | 3, 4 | 3.1–3.5, 4.1–4.4 | 5h |
| Sprint 3 (Day 2 PM) | 5, 6 | 5.1–5.3, 6.1–6.3 | 4h |
| Sprint 4 (Day 3) | 7, 8, 9 | 7.1–7.3, 8.1–8.2, 9.1–9.3 | 5h |
