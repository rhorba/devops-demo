# Project Metrics
<!-- KPI snapshots over time -->
<!-- Format: ### [YYYY-MM-DD HH:MM] SPRINT_SNAPSHOT/DAILY_SNAPSHOT — Title -->

### [2026-05-03] SPRINT_SNAPSHOT — Project code-complete (Sprint 4 done)
- **Services built**: 4/4 (sensor-simulator, api-gateway, data-sink, enterprise-consumer)
- **Helm charts**: 12 (8 infrastructure + 4 service)
- **ArgoCD apps**: 7 (root + 6 child apps)
- **Kyverno policies**: 4 (no-root, resource-limits, no-latest-tag, require-labels)
- **Bugs found/fixed**: 11
- **Files committed**: 62 infrastructure files
- **CI/CD**: GitHub Actions (lint + test + build + push, matrix for 4 services)

### [2026-05-06] DAILY_SNAPSHOT — Partial bootstrap
- **Images pushed to local registry**: 3/4 (missing enterprise-consumer initially, then rebuilt with AMQP)
- **ArgoCD apps synced**: Root app healthy; child apps syncing
- **Kind cluster**: 3-node cluster running (1 control-plane + 2 workers)
- **Tools installed**: kind v0.31.0, argocd CLI v3.3.8 (via winget)
