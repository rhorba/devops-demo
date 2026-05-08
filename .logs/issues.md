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
- **Status**: PARTIALLY FIXED (2026-05-08) — TLS passes, auth/queue issues remain
- **Description**: `rpc-redpanda-to-ibmmq` could not POST to `https://ibm-mq.messaging.svc.cluster.local:9443` — IBM MQ self-signed cert only had SAN for `localhost`
- **Fix applied**: Generated new IBM MQ cert (CA:TRUE + correct SANs) → combined with system CA bundle → mounted at /etc/ssl/certs/ca-certificates.crt in pipeline pod. TLS now passes (x509 error gone).
- **Remaining issue 1**: IBM MQ Liberty auth fix is IN-POD ONLY (manual edit to /mnt/mqm/data/web/.../mqwebcontainer.xml). Will be lost on IBM MQ pod restart.
- **Remaining issue 2**: DEV.QUEUE.1 queue is full (5000/5000) — no consumer draining it. Pipeline gets 503 MQRC_Q_FULL.
- **Next step**: Make Liberty auth fix permanent via Helm chart (init-container or ConfigMap override)

### [2026-05-07 09:57] BUG — enterprise-consumer AMQP port 5672 connection refused
- **Status**: PARTIALLY FIXED — new REST API image built but new pod ImagePullBackOff; old AMQP pod still running
- **Description**: enterprise-consumer connecting to IBM MQ via AMQP (port 5672) — IBM MQ Developer image doesn't support AMQP
- **Fix applied**: Rewrote main.go to use REST API (port 9443/HTTPS). Image rebuilt and pushed to localhost:5001.
- **Remaining issue**: New pod (imagePullPolicy: Always) in ImagePullBackOff — Kind worker nodes may not resolve localhost:5001 registry. Old pod (imagePullPolicy: IfNotPresent, cached image) still runs old AMQP binary.
- **Next step**: Fix enterprise-consumer image pull — either investigate Kind containerd mirror config or rebuild with versioned tag (not :latest)

### [2026-05-08 06:10] BLOCKER — IBM MQ Liberty auth ephemeral fix
- **Status**: OPEN
- **Description**: mqwebcontainer.xml in Liberty server directory has plain-text passw0rd hardcoded (manual fix, 2026-05-08). Will be lost if IBM MQ pod restarts.
- **Impact**: IBM MQ REST API will return 401 after pod restart
- **Fix needed**: Helm chart change to mount a ConfigMap at /mnt/mqm/data/web/installations/Installation1/servers/mqweb/mqwebcontainer.xml OR use an init-container to patch it on startup

### [2026-05-08 06:10] BLOCKER — DEV.QUEUE.1 full (5000 messages)
- **Status**: OPEN
- **Description**: IBM MQ DEV.QUEUE.1 is at max depth (5000/5000). Pipeline gets 503 MQRC_Q_FULL. enterprise-consumer not running REST API binary to drain queue.
- **Impact**: rpc-redpanda-to-ibmmq pipeline blocked until queue is drained
- **Fix needed**: Either fix enterprise-consumer to drain the queue, or increase MAXDEPTH, or clear queue manually each session
