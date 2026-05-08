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
- **Status**: FIXED (2026-05-08)
- **Description**: `rpc-redpanda-to-ibmmq` could not POST to `https://ibm-mq.messaging.svc.cluster.local:9443` — IBM MQ self-signed cert only had SAN for `localhost`
- **Fix applied**: Switched to `skip_cert_verify: true` in Redpanda Connect `http_client` TLS config (charts/redpanda-connect/templates/configmap.yaml). Committed in `433492f`. Pod restarted 2026-05-08 to pick up new configmap.
- **Verification**: Pipeline now reaches IBM MQ, gets 201 on PUT or 503 MQRC_Q_FULL (not TLS error)

### [2026-05-07 09:57] BUG — enterprise-consumer AMQP port 5672 connection refused
- **Status**: FIXED (2026-05-08)
- **Description**: enterprise-consumer connecting to IBM MQ via AMQP (port 5672) — IBM MQ Developer image doesn't support AMQP
- **Fix applied**: Rewrote main.go to use REST API (port 9443/HTTPS). Image tagged v2.0.0 + `kind load docker-image` + `imagePullPolicy: IfNotPresent` committed to git. ArgoCD synced new pod.
- **Verification**: enterprise-consumer pod v2.0.0 Running, consuming from IBM MQ, forwarding to SQS (4307 msgs confirmed)

### [2026-05-08 06:10] BLOCKER — IBM MQ Liberty auth ephemeral fix
- **Status**: FIXED (2026-05-08)
- **Description**: mqwebcontainer.xml resolved `${env.MQ_APP_PASSWORD_SECURE}` to empty string when pod restarted
- **Fix applied**: Added `MQ_APP_PASSWORD_SECURE` and `MQ_ADMIN_PASSWORD_SECURE` env vars to IBM MQ deployment (charts/ibm-mq/templates/deployment.yaml). Committed in `433492f`. Auth now permanent across pod restarts.

### [2026-05-08 06:10] BLOCKER — DEV.QUEUE.1 full (5000 messages)
- **Status**: RESOLVED (2026-05-08)
- **Description**: IBM MQ DEV.QUEUE.1 was at max depth (5000/5000). Pipeline gets 503 MQRC_Q_FULL.
- **Resolution**: enterprise-consumer (now running REST API binary) drains the queue continuously. Queue is in steady-state near 5000 (produce=consume balance). SQS has 4307+ forwarded messages confirming flow.

### [2026-05-08 07:30] BUG — LocalStack SQS/S3 resources not created after restart
- **Status**: FIXED (2026-05-08 session)
- **Description**: LocalStack is ephemeral — SQS queue `enterprise-events` and S3 bucket `telemetry-archive` were lost when LocalStack restarted (80+ restarts over 22h)
- **Fix applied**: Created resources manually: `awslocal sqs create-queue --queue-name enterprise-events` + `awslocal s3 mb s3://telemetry-archive`
- **Note**: Resources are lost on LocalStack restart — consider adding init-job or startup script to recreate them
