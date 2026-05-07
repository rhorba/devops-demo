# Issue Log
<!-- Tracks bugs, errors, blockers with status -->
<!-- Format: ### [YYYY-MM-DD HH:MM] BUG/BLOCKER/ERROR — Title -->

### [2026-05-06] BLOCKER — ArgoCD CLI v3/v2 gRPC incompatibility
- **Status**: OPEN (workaround in place)
- **Description**: ArgoCD CLI v3.3.8 is incompatible with ArgoCD server v2.10.2 due to gRPC protocol mismatch
- **Impact**: Cannot use `argocd app sync`, `argocd app get` etc.
- **Workaround**: Use `kubectl -n platform get/describe application` or REST API at https://localhost:30080
- **Resolution**: Will resolve when ArgoCD is upgraded to v2.13+ or CLI downgraded to v2.x

### [2026-05-07 09:57] BUG — IBM MQ TLS cert SAN mismatch (redpanda-connect → IBM MQ)
- **Status**: OPEN
- **Description**: `rpc-redpanda-to-ibmmq` cannot POST to `https://ibm-mq.messaging.svc.cluster.local:9443` — IBM MQ self-signed cert only has SAN for `localhost`, not the K8s service DNS name
- **Impact**: Redpanda → IBM MQ pipeline leg broken; messages not reaching IBM MQ queue
- **Fix options**:
  - A) Add `tls_skip_cert_verify: true` to redpanda-connect ibm-mq pipeline configmap (fast, acceptable for demo)
  - B) Generate IBM MQ cert with correct SAN (proper fix, complex)
- **Recommended**: Option A for demo

### [2026-05-07 09:57] BUG — enterprise-consumer AMQP port 5672 connection refused
- **Status**: OPEN
- **Description**: enterprise-consumer is connecting to IBM MQ via AMQP (port 5672), but IBM MQ Developer image does not expose AMQP — connection refused
- **Impact**: enterprise-consumer → LocalStack SQS leg broken; no SQS messages
- **Fix**: enterprise-consumer must use IBM MQ REST API (port 9443) not AMQP; verify the REST API rewrite is in the running image
- **Note**: Memory says REST API rewrite was done 2026-05-07, but running pod still shows AMQP errors — image may not have been rebuilt and pushed
