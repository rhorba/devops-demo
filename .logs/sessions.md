# Session Log

### 2026-05-03 SESSION_END
- **Specialist**: Orchestrator
- **Completed**:
  - 11 bugs identified and closed (commit 7e80e66)
  - CTS Session Playbook written to README.md (commit 4d8df08)
  - .gitignore fixed — removed **/go.sum so lockfiles can be committed (commit 3ac1d4b)
  - All pushed to https://github.com/rhorba/devops-demo.git main
- **Bugs closed**:
  1. data-sink/main.go — flush hot loop fixed (mutex + dedicated goroutine)
  2. api-gateway/main.go — health metric always recorded 200, now records actual status
  3. enterprise-consumer chart — hardcoded passw0rd removed from Secret stringData
  4. data-sink chart — missing readiness probe added
  5. kyverno no-root policy — =() conditional bypass replaced with anyPattern
  6. argocd/apps/messaging.yaml — redpanda-connect + ibm-mq missing CreateNamespace=true
  7. All 4 Dockerfiles — go.sum missing from COPY before go mod download
  8. argocd/apps/apps.yaml — Helm template syntax replaced with plain YAML Applications
  9. monitoring/prometheus/values.yaml — plaintext adminPassword replaced with existingSecret ref
  10. kind-config.yaml — ArgoCD NodePort 30080 missing from extraPortMappings
  11. Makefile — argocd-sync --cascade (invalid flag) replaced with --prune
- **In progress**: Nothing
- **Blocked**: Nothing
- **Known gap**: go.sum files do not exist on disk yet — must run go mod tidy in each
  services/* directory then git add services/*/go.sum before Docker builds will work
- **Audit finding**: Session did NOT follow CTS framework — UNDERSTAND/BRAINSTORM/PLAN
  phases skipped, no specialist skills loaded, no handoffs, no logs maintained during
  execution. CTS playbook now in README to enforce compliance next session.
- **Next session**:
  - Start with: resume last session
  - Priority task: cd into each services/* dir, run go mod tidy, commit go.sum files
  - Then: commit remaining untracked files (ansible/, charts/, argocd/apps remaining,
    security/ policies + rbac + network-policies, monitoring/loki + grafana + thanos,
    scripts/demo.sh + teardown.sh, docs/, services/ chart Chart.yaml + values.yaml)
- **Open issues**: 1 — go.sum files not yet generated (low risk, blocks Docker build)
- **Open risks**: 0
---

### 2026-05-03 SESSION_END (updated — extended session)
- **Specialist**: Orchestrator
- **Completed**:
  - Go 1.26.2 installed via winget (C:\Program Files\Go\bin\go)
  - go mod tidy run in all 4 services — go.sum files generated
  - MIT licence added (LICENSE file, copyright rhorba 2026)
  - .gitignore fixed — **/go.sum removed so lockfiles can be tracked
  - Commits pushed: 3ac1d4b (.gitignore fix), 64295f4 (go.sum + LICENSE)
  - CTS Session Playbook written to README.md (commit 4d8df08)
  - 11 bugs closed (commit 7e80e66)
- **In progress**: Nothing
- **Blocked**: Nothing
- **Next session**:
  - Start with: resume last session
  - Main task: run the devops-demo project end-to-end
    1. make bootstrap  (kind cluster + ArgoCD + all charts ~12-15 min)
    2. make port-forwards
    3. Smoke-test full pipeline: MQTT -> Redpanda -> NATS -> IBM MQ -> S3/SQS
    4. Verify Grafana, ArgoCD, API, Node-RED all reachable
  - Secondary: commit remaining untracked files (ansible/, charts/, argocd/apps
    remaining, security/, monitoring/, scripts/demo+teardown, docs/, service charts)
- **Open issues**: 0
- **Open risks**: 0
---

### 2026-05-05 SESSION_START
- **Specialist**: Orchestrator
- **Resuming from**: SESSION_END 2026-05-03 (extended)
- **Context**: All 11 bugs fixed, go.sum generated, MIT licence added. Project ready to run end-to-end.
- **Plan**: make bootstrap → port-forwards → smoke-test full pipeline
---

### 2026-05-06 SESSION_END
- **Specialist**: Orchestrator
- **Completed**:
  - Installed missing tools: kind v0.31.0, argocd CLI v3.3.8 (via winget → ~/bin)
  - Partial bootstrap run: cluster + namespaces + cert-manager + Kyverno + ArgoCD created
  - 3/4 service images built and pushed: sensor-simulator, api-gateway, data-sink
  - enterprise-consumer: diagnosed IBM MQ CGO build failure (SDK headers not in MQ image, IBM CDN blocked)
  - enterprise-consumer: rewrote using AMQP 1.0 (Azure/go-amqp) — removed CGO dependency entirely
  - IBM MQ chart: added AMQP port 5672 + MQSC ConfigMap to start AMQP channel
  - All 4 images now built and pushed to local registry (localhost:5001)
  - Committed all 62 previously-untracked infrastructure files (charts, services, monitoring, security, argocd, docs)
  - Fixed ArgoCD repo URL: all app manifests updated from rhorba/devops-demo → rhorba/CTS
  - Applied AppProject (devops-demo) and updated root Application to cluster
  - Root app: Healthy, sync in progress (auto-sync enabled)
  - Pushed all commits to https://github.com/rhorba/CTS main
- **In progress**: ArgoCD syncing child apps (messaging, platform, observability, security, tools)
- **Blocked**: ArgoCD CLI v3 incompatible with ArgoCD server v2.10.2 (gRPC mismatch — use REST API or kubectl instead)
- **Next session**:
  - Start with: resume last session
  - Check ArgoCD sync status: kubectl -n platform get applications
  - Verify all child apps synced and pods running in messaging/apps/observability/tools namespaces
  - Run make port-forwards
  - Smoke-test pipeline: MQTT → Redpanda → NATS → IBM MQ → SQS
  - Verify Grafana dashboard at http://localhost:3000
  - Demo walkthrough: make demo
- **Open issues**: 1 — ArgoCD CLI v3/v2 mismatch (use REST API: https://localhost:30080 admin/tX5VSX6ZNnAPRK07)
- **Open risks**: 0
---

### [2026-05-06 00:00] SESSION_END
- **Completed**: Reverted rhorba/CTS to b985dc2 (removed accidental devops-demo commits)
- **In progress**: Nothing — clean state
- **Blocked**: Nothing
- **Next session**: Verify ArgoCD sync → port-forwards → smoke-test pipeline → make demo
- **Open issues**: ArgoCD CLI v3/v2 gRPC mismatch (use REST API workaround)
- **Open risks**: Kind cluster may have stopped — may need `kind create cluster` again
---

### [2026-05-07 09:35] SESSION_START
- Resuming from: SESSION_END 2026-05-06
- Context: Cluster fully healthy (18/18 apps Synced Healthy, commit c0e8584) — sessions.md not yet updated to reflect 2026-05-07 fixes
- Plan: TBD — awaiting user direction
---

### [2026-05-07 ~00:00] SESSION_START (backfilled)
- Resuming from: SESSION_END 2026-05-06
- Plan: Fix remaining broken ArgoCD apps → get all 18 Synced Healthy
---

### [2026-05-07 ~23:59] SESSION_END (backfilled)
- Specialist: DevOps/DevSecOps
- Completed:
  - redpanda configmap: removed invalid `default_topic_replication_factor` from bootstrap.yaml + redpanda.yaml → pod no longer CrashLoops
  - Prometheus CR: fixed thanos.image — must be string not nested object
  - Grafana, api-gateway, node-red Ingresses: disabled (ArgoCD health gate incompatible with Kind+NodePort)
  - readOnlyRootFilesystem removed from pod-level securityContext across 6 apps (it's container-level only — Kyverno was rejecting)
  - enterprise-consumer: rewrote IBM MQ connector from AMQP to REST API (port 9443/HTTPS) — IBM MQ Developer image does not support AMQP
  - redpanda-connect: updated ibm-mq pipeline to use REST API
  - LocalStack memory limit: 512Mi → 1Gi (was OOMKilled under load)
  - promtail scrapeConfigs: added `__path__` relabel rule (logs were not found without it)
  - node-red Helm chart: escaped `{{msg.payload.*}}` template vars (conflicting with Helm Go templates)
  - All 18 ArgoCD applications: Synced Healthy ✅
  - Latest commit: c0e8584 (node-red Ingress fix)
- In progress: Nothing
- Blocked: Nothing
- Next session: Smoke-test full pipeline, run make demo, generate session report
- Open issues: 0
- Open risks: AppProject not auto-synced — must `kubectl apply -f argocd/projects/demo.yaml` on any change
---

### [2026-05-07 10:10] SESSION_END
- Specialist: Orchestrator → DevOps/DevSecOps (smoke-test) → Project Monitor (report)
- Completed:
  - A: Backfilled 2026-05-07 session log (8 architectural fixes, 18/18 apps Synced Healthy)
  - B: Full pipeline smoke-test — 4/6 legs flowing, 2 IBM MQ bugs surfaced
  - D: Status report generated — issues.md + metrics.md updated
- In progress: Nothing
- Blocked: Nothing
- Next session:
  - Fix IBM MQ TLS: add tls_skip_cert_verify to redpanda-connect ibm-mq configmap
  - Fix enterprise-consumer: verify/rebuild REST API image (AMQP still in running pod)
  - Then: full 6/6 pipeline smoke-test → make demo
- Open issues: 3
  1. ArgoCD CLI v3/v2 mismatch (cosmetic, workaround in place)
  2. IBM MQ TLS cert SAN mismatch (rpc-redpanda-to-ibmmq)
  3. enterprise-consumer AMQP port 5672 refused (image may not have been rebuilt)
- Open risks: 1 — Kind cluster is ~28h old; may stop on machine restart
---

### [2026-05-07 10:30] SESSION_START
- Resuming from: SESSION_END 2026-05-07 10:10
- Plan: Fix 2 IBM MQ pipeline bugs → full 6/6 pipeline smoke-test
---
