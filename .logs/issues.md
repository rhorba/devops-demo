# Issue Log
<!-- Tracks bugs, errors, blockers with status -->
<!-- Format: ### [YYYY-MM-DD HH:MM] BUG/BLOCKER/ERROR — Title -->

### [2026-05-06] BLOCKER — ArgoCD CLI v3/v2 gRPC incompatibility
- **Status**: OPEN (workaround in place)
- **Description**: ArgoCD CLI v3.3.8 is incompatible with ArgoCD server v2.10.2 due to gRPC protocol mismatch
- **Impact**: Cannot use `argocd app sync`, `argocd app get` etc.
- **Workaround**: Use `kubectl -n platform get/describe application` or REST API at https://localhost:30080
- **Resolution**: Will resolve when ArgoCD is upgraded to v2.13+ or CLI downgraded to v2.x
