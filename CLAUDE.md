# devops-demo — CLAUDE.md

> Auto-loaded every session. All instructions below are **mandatory** — no exceptions.

---

## SESSION START — DO THIS FIRST (every session, no skipping)

1. **Read** `.logs/sessions.md` (last 30 lines)
2. **Report** in one paragraph: what was done last session, what is in progress, blockers, next steps
3. **Ask**: *"Resume from last session, or start a new task?"*
4. **Log** a `SESSION_START` entry to `.logs/sessions.md` before any work begins:

```
### [YYYY-MM-DD HH:MM] SESSION_START
- Resuming from: [last SESSION_END date or "fresh start"]
- Plan: [today's goal in one line]
---
```

---

## MANDATORY WORKFLOW — Every Task

All tasks route through this chain. Never skip it:

```
/orchestrator  →  [specialist skill]  →  /project-monitor (log)
```

- **New task or unclear scope** → always start with `/orchestrator`
- **Specialist switch** → log the handoff in `.logs/sessions.md` before switching
- **Never implement** without first routing and confirming the approach

---

## SKILL ROUTING TABLE (devops-demo)

| When asked to... | Invoke |
|---|---|
| Bootstrap, `make bootstrap`, run end-to-end | `/devops-devsecops` + `/deployment` |
| Fix Helm chart, values.yaml, ArgoCD app manifest | `/devops-devsecops` |
| Fix or update Ansible playbook | `/devops-devsecops` |
| Fix or update GitHub Actions workflow | `/devops-devsecops` |
| Add/fix Go microservice (api-gateway, data-sink, etc.) | `/backend-dev` |
| Write or fix Go unit/integration tests | `/tester` |
| Design test strategy for the pipeline | `/test-architect` |
| Kyverno policy, RBAC, NetworkPolicy, cert-manager | `/security-engineer` |
| Security check before git push | `/security-review` |
| Prometheus alert rule, Grafana panel, Loki config | `/devops-devsecops` |
| Kafka topic/partition design, LocalStack S3/SQS schema | `/dba` |
| Sprint planning, backlog management, velocity | `/scrum-master` |
| Architecture decision, ADR, tech stack choice | `/tech-lead` |
| PRD update, epics, scope definition | `/project-manager` |
| Session summary, KPIs, status report | `/project-monitor` |
| README, docs, tutorials | `/content-marketer` |
| Brainstorm architecture options | `/creative-intelligence` |
| Code quality, simplify, refactor review | `/simplify` |
| Ambiguous or multi-domain request | `/orchestrator` |

---

## SESSION END — DO BEFORE EVERY CLOSE (mandatory)

Write a `SESSION_END` entry to `.logs/sessions.md`:

```
### [YYYY-MM-DD HH:MM] SESSION_END
- Specialist: [last skill invoked]
- Completed: [bullet list of what was done this session]
- In progress: [description or "Nothing"]
- Blocked: [description or "Nothing"]
- Next session: [exact first command + ordered priority task list]
- Open issues: [N — brief descriptions]
- Open risks: [N — brief descriptions]
---
```

Then commit the logs:
```bash
git add .logs/ && git commit -m "chore(logs): session end $(date +%Y-%m-%d)" && git push
```

---

## PROJECT CONTEXT

| | |
|---|---|
| **Repo** | https://github.com/rhorba/devops-demo |
| **Branch** | main |
| **Registry** | localhost:5001 (Kind-integrated local) |
| **ArgoCD** | https://localhost:30080 — check `.logs/sessions.md` for credentials |
| **Go services** | sensor-simulator, api-gateway, data-sink, enterprise-consumer |
| **Messaging** | MQTT (Mosquitto) → Redpanda (Kafka) → NATS / IBM MQ → LocalStack S3/SQS |
| **Observability** | Prometheus + Grafana + Loki + Thanos |
| **Security** | Kyverno + RBAC + NetworkPolicies + cert-manager |
| **GitOps** | ArgoCD App-of-Apps (root → messaging, platform, observability, security, tools, apps) |

**Permanent constraints** (never forget):
- ArgoCD CLI v3 is incompatible with ArgoCD server v2.10.2 — use `kubectl` or the REST API, not `argocd` CLI
- Kind cluster is ephemeral — always run `kind get clusters` at session start before assuming it exists

**First commands at every session start**:
```bash
kind get clusters                                   # cluster alive?
kubectl -n platform get applications                # ArgoCD sync status
make port-forwards                                  # expose services locally
```

---

## LOG FILES REFERENCE

| File | What to log there |
|---|---|
| `.logs/sessions.md` | SESSION_START / SESSION_END per session |
| `.logs/activity.md` | Completions, milestones, deployments |
| `.logs/decisions.md` | Architecture/tech decisions (ADR-worthy choices) |
| `.logs/issues.md` | Bugs, blockers found or resolved |
| `.logs/risks.md` | Risks identified and their mitigations |
| `.logs/corrections.md` | Mistakes made and how they were corrected |
| `.logs/communications.md` | PRs opened, releases cut, external comms |
| `.logs/metrics.md` | Build times, test coverage, deploy counts |

**Logging rules**:
- Log AFTER the action, in the same response — no separate log calls
- Log decisions to `decisions.md` whenever an ADR-worthy choice is made
- Log issues to `issues.md` whenever a blocker or bug is found or resolved
- Log corrections to `corrections.md` whenever a mistake is caught and fixed

---

## COMMIT PROTOCOL

Before every `git push`:
1. Run `/security-review` (or manually check: no hardcoded secrets, no credentials in code)
2. Run `go test ./...` in any changed Go service directory
3. Use conventional commits: `type(scope): description`
   - Types: `feat`, `fix`, `docs`, `chore`, `refactor`, `test`, `security`, `infra`
4. `git push origin main`

---

## TEAM SKILLS REFERENCE

Specialists are invoked via slash commands. Context files are in `claude-team-skills/`.

| Slash command | Context file |
|---|---|
| `/orchestrator` | `claude-team-skills/orchestrator/SKILL.md` |
| `/project-manager` | `claude-team-skills/project-manager/SKILL.md` |
| `/scrum-master` | `claude-team-skills/scrum-master/SKILL.md` |
| `/tech-lead` | `claude-team-skills/tech-lead/SKILL.md` |
| `/security-engineer` | `claude-team-skills/security-engineer/SKILL.md` |
| `/dba` | `claude-team-skills/dba/SKILL.md` |
| `/backend-dev` | `claude-team-skills/backend-dev/SKILL.md` |
| `/tester` | `claude-team-skills/tester/SKILL.md` |
| `/test-architect` | `claude-team-skills/test-architect/SKILL.md` |
| `/deployment` | `claude-team-skills/deployment/SKILL.md` |
| `/devops-devsecops` | `claude-team-skills/devops-devsecops/SKILL.md` |
| `/project-monitor` | `claude-team-skills/project-monitor/SKILL.md` |
| `/creative-intelligence` | `claude-team-skills/creative-intelligence/SKILL.md` |
| `/content-marketer` | `claude-team-skills/content-marketer/SKILL.md` |

Full team configuration and working protocols: `claude-team-skills/CLAUDE.md`
