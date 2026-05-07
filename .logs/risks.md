# Risk Log
<!-- Tracks risks identified and their mitigations -->
<!-- Format: ### [YYYY-MM-DD HH:MM] SECURITY/PERFORMANCE/DEPENDENCY — Title -->

### [2026-05-06] DEPENDENCY — Kind cluster ephemeral state
- **Status**: OPEN
- **Probability**: High (every machine restart loses cluster)
- **Impact**: Medium (15 min rebuild via `make bootstrap`)
- **Mitigation**: Always run `kind get clusters` at session start; `make bootstrap` rebuilds fully
- **Owner**: DevOps/DevSecOps

### [2026-05-06] DEPENDENCY — IBM MQ AMQP channel configuration
- **Status**: OPEN
- **Probability**: Medium (channel may not start on fresh deploy)
- **Impact**: Medium (enterprise-consumer cannot receive messages)
- **Mitigation**: IBM MQ chart includes MQSC ConfigMap to auto-start AMQP channel; verify with `kubectl -n messaging logs` on IBM MQ pod
- **Owner**: DevOps/DevSecOps
