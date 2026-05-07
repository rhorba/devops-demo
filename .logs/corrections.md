# Corrections & Plan Changes
<!-- Tracks scope changes, pivots, plan adjustments, mistakes corrected -->
<!-- Format: ### [YYYY-MM-DD HH:MM] SCOPE_CHANGE/PIVOT/REPLAN/CORRECTION — Title -->

### [2026-05-03] CORRECTION — enterprise-consumer IBM MQ CGO build failure
- **Mistake**: enterprise-consumer was built with CGO dependency on IBM MQ C SDK headers, which are not available in the MQ image and IBM CDN was blocked
- **Correction**: Rewrote enterprise-consumer using AMQP 1.0 (Azure/go-amqp library) — removes all CGO dependency
- **Prevention**: Prefer pure-Go libraries for production microservices; avoid CGO in containerized builds

### [2026-05-03] CORRECTION — ArgoCD repo URL mismatch
- **Mistake**: All ArgoCD app manifests pointed to rhorba/CTS instead of rhorba/devops-demo
- **Correction**: Updated all manifests in argocd/apps/*.yaml to correct repo URL
- **Prevention**: Validate repo URLs in ArgoCD manifests before first sync

### [2026-05-03] CORRECTION — 11 bugs in sprint 4 output
- **Mistake**: Sprint 4 produced code without following CTS UNDERSTAND/BRAINSTORM/PLAN phases — resulted in 11 bugs
- **Correction**: All 11 bugs fixed in post-session audit (commit 7e80e66)
- **Prevention**: ALWAYS follow the /orchestrator workflow before implementing

### [2026-05-06] CORRECTION — Accidental commits to wrong repo (rhorba/CTS)
- **Mistake**: 62 infrastructure files were committed to rhorba/CTS instead of rhorba/devops-demo
- **Correction**: Reverted rhorba/CTS to b985dc2; all files re-committed to correct repo
- **Prevention**: Verify `git remote -v` before committing large changesets
