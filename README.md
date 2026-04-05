# Claude Team Skills (CTS)

**19 AI Specialists Working as One Team — Built for Claude Code**

> Version 1.0 | 19 specialists | 46 files | 6,087 lines | Token-Efficient | YAGNI-Driven

Claude Team Skills is a collection of 19 interconnected AI specialist skills for Claude Code. Instead of generic AI responses, every task is handled by the right specialist — PM, Security Engineer, DBA, Tester, UX Designer, and more — all coordinated by a central Orchestrator.

---

## Why Claude Team Skills?

| | Claude Team Skills | Generic AI |
|---|---|---|
| Specialists | 19 role-specific experts | One-size-fits-all |
| Token usage | ~70% savings via lazy loading | Full context always loaded |
| Project memory | 8 structured log files | No memory between sessions |
| Security | STRIDE, OWASP Top 10, compliance | Ad-hoc |
| Testing | Risk-based ATDD + adversarial review | Basic test generation |
| Workflow | 6-phase structured process | Free-form |

---

## Installation

1. Download and extract `claude-team-skills-final.tar.gz` into your project root.
2. Copy `CLAUDE.md` to your project root (the entry point Claude Code reads first).
3. Open Claude Code and describe what you need.

```
your-project/
├── CLAUDE.md              ← Entry point (copy here)
├── QUICKSTART.md          ← Cheat sheet
├── skills/
│   ├── orchestrator/
│   ├── project-manager/
│   ├── scrum-master/
│   ├── tech-lead/
│   ├── security-engineer/
│   ├── dba/
│   ├── ux-designer/
│   ├── ui-designer/
│   ├── backend-dev/
│   ├── frontend-dev/
│   ├── tester/
│   ├── test-architect/
│   ├── deployment/
│   ├── devops-devsecops/
│   ├── creative-intelligence/
│   ├── digital-marketer/
│   ├── copywriter/
│   ├── content-marketer/
│   └── project-monitor/
└── .logs/                 ← Auto-created on first use
```

**Zero dependencies.** Plain markdown files — no Node.js, no build tools, no npm install.

---

## The Team — 19 Specialists

### Management
| Specialist | What They Do | Use When |
|---|---|---|
| **Project Manager** | Charter, WBS, risk, estimation | New project, scope definition, status reports |
| **Scrum Master** | Sprints, stories, backlog, ceremonies | Sprint planning, task batching, retros |
| **Tech Lead** | ADRs, stack decisions, YAGNI gate | Architecture decisions, code standards |

### Security & Data
| Specialist | What They Do | Use When |
|---|---|---|
| **Security Engineer** | STRIDE, auth design, OWASP, compliance | Threat models, security requirements |
| **DBA** | Schema, queries, migrations, indexes | Database design, optimization |

### Design
| Specialist | What They Do | Use When |
|---|---|---|
| **UX Designer** | Flows, wireframes, IA, usability | User flows, navigation, form design |
| **UI Designer** | Design tokens, colors, typography | Design system, component styling |

### Development
| Specialist | What They Do | Use When |
|---|---|---|
| **Backend Dev** | APIs, services, auth, data layer | Server-side code, business logic |
| **Frontend Dev** | Components, state, accessibility | UI components, client-side code |
| **Tester** | Unit / integration / e2e, QA | Writing tests, bug reports |
| **Test Architect** | Risk-based strategy, ATDD, adversarial review | Test planning, edge case hunting |

### Operations
| Specialist | What They Do | Use When |
|---|---|---|
| **Deployment** | Release, rollback, environments | Shipping code, environment setup |
| **DevOps / DevSecOps** | CI/CD, Docker, K8s, scanning | Pipeline setup, infra, security scans |

### Creative & Marketing
| Specialist | What They Do | Use When |
|---|---|---|
| **Creative Intelligence** | Brainstorming, design thinking, innovation | Ideation, problem reframing, naming |
| **Digital Marketer** | SEO, campaigns, funnels | Launch strategy, analytics |
| **Copywriter** | Headlines, CTAs, UX writing | Ad copy, email copy, UI text |
| **Content Marketer** | Blog, social, content calendar | Blog posts, social media strategy |

---

## How It Works

Every request goes through the **Orchestrator**, which identifies the right specialist, loads only that skill, and passes a 2–3 line handoff summary when switching — keeping token usage minimal.

### 6-Phase Workflow

| Phase | What Happens |
|---|---|
| 1. Understand | Clarify what, where, how big |
| 2. Brainstorm | 2–3 options: simple → complex |
| 3. Plan | Break into batched tasks |
| 4. Execute | One task at a time, interactive checkpoints |
| 5. Verify | Tests + security check |
| 6. Ship | Deploy + document |

### Three Core Principles

**YAGNI** — Always start with the simplest approach. Build only what's needed now.

| Instead of... | Use... |
|---|---|
| Microservices for MVP | Monolith first |
| Kubernetes for 3 containers | Docker Compose |
| Custom auth from scratch | NextAuth / Passport |
| Cache everything upfront | Profile first, cache bottlenecks |

**Token Efficiency** — Load ONE specialist at a time. Handoffs use 2–3 line summaries. Saves ~70% of token budget.

**Project Memory** — Everything logged to `.logs/`. Resume any session exactly where you left off.

---

## Project Monitoring & Logging

Eight structured log files auto-created in `.logs/`:

| File | Tracks |
|---|---|
| `activity.md` | Completed tasks, milestones |
| `decisions.md` | Architecture choices with rationale |
| `issues.md` | Bugs, blockers, errors |
| `risks.md` | Risks + mitigation status |
| `corrections.md` | Scope changes, pivots |
| `communications.md` | Preferences, handoffs |
| `sessions.md` | Session start/end snapshots |
| `metrics.md` | KPI snapshots |

### Report Commands

```
"status report"      → Full project status with KPIs
"what did we do?"    → Activity log
"what's broken?"     → Open issues and blockers
"what did we decide?" → Decision log
"retro"              → Retrospective analysis
"KPIs"               → Project metrics
```

---

## Document-First Development

For new projects and large features (3+ days of work):

```
PRD  →  Architecture Doc  →  Epics & Stories  →  Implementation
 PM      Tech Lead + DBA       Scrum Master +       Dev team
                                Test Architect
```

**YAGNI scaling:**
- Small task (< 1 day): skip the chain
- Medium feature: mini-PRD + ADR only
- Large feature / new project: full chain

---

## Test Architecture & Adversarial Review

The **Test Architect** provides:

- **Risk-based strategy** — Score components by impact × frequency × complexity; allocate testing effort accordingly
- **ATDD / Gherkin** — Write acceptance criteria before coding; they become the tests
- **Adversarial review** — Input abuse (SQLi, XSS, traversal), auth bypass (IDOR, CSRF, escalation), race conditions, business logic exploits
- **Traceability matrix** — Every requirement → story → test → code
- **Release gates** — Go/no-go criteria for deployment

---

## Security Coverage

- STRIDE threat modeling
- OWASP Top 10 — full coverage with detection + remediation per vulnerability
- Compliance frameworks: SOC2, PCI-DSS, HIPAA, ISO 27001, NIST, GDPR
- CI/CD, container, and cloud security scanning references

---

## Creative Intelligence Suite

7 structured ideation frameworks:

| Framework | Best For |
|---|---|
| Brainstorming | Generating lots of ideas fast |
| SCAMPER | Improving existing solutions |
| Design Thinking | User-centered innovation (5 phases) |
| Storytelling | Pitches and product narratives |
| First Principles | Breaking assumptions, rebuilding from facts |
| Problem Reframing | Unsticking hard problems |
| Name Generation | Product / feature naming with scoring |

All sessions follow: **diverge → cluster → converge**, scored with an Impact × Confidence / Effort matrix.

---

## Quick Command Reference

| Command | What It Does |
|---|---|
| `"simpler"` | Switch to a less complex approach |
| `"YAGNI"` | Remind to keep it simple |
| `"option A/B/C"` | Pick from presented options |
| `"good, next"` | Approve current step and move on |
| `"skip this"` | Move to next task |
| `"go back"` | Revisit a previous decision |
| `"what's the plan?"` | Show current batch/task list |
| `"ship it"` | Move to deployment |
| `"test this"` | Run tests on current work |
| `"secure this"` | Run security review |
| `"adversarial review"` | Run adversarial security checklist |
| `"test strategy"` | Design risk-based test plan |
| `"brainstorm this"` | Start creative ideation session |
| `"status report"` | Generate project status |
| `"retro"` | Generate retrospective |

---

## Token Budget by Session Type

| Session Type | Expected Tokens | Skills Loaded |
|---|---|---|
| Quick bug fix | 3–8K | Orchestrator + 1 dev + tester |
| Small feature | 8–15K | Orchestrator + tech lead + 1–2 devs |
| Large feature | 15–30K | Multiple specialists across batches |
| New project scaffold | 15–25K | PM + tech lead + DBA + devs |
| UX/UI design | 5–12K | UX + UI + frontend + copywriter |
| Security audit | 5–15K | Security + DevSecOps |
| Brainstorm session | 3–8K | Creative Intelligence + PM/Tech Lead |
| Test strategy | 5–10K | Test Architect + Tester + Security |

---

## Customization

- **Edit `CLAUDE.md`** — Add your tech stack, naming conventions, deployment URLs, team contacts.
- **Disable unused skills** — Remove rows for specialists you don't need from the CLAUDE.md skill table.
- **Add new specialists** — Create a folder in `skills/` with a `SKILL.md` (YAML frontmatter + role + templates + handoff points), then register it in the orchestrator routing table and `CLAUDE.md`.

---

## Troubleshooting

| Problem | Solution |
|---|---|
| Claude doesn't use skills | Ensure `CLAUDE.md` is in the project root |
| Wrong specialist loaded | Say "switch to [specialist]" or be more specific |
| Logs not created | Run `bash skills/project-monitor/templates/init-logs.sh` |
| Session not resuming | Check `.logs/sessions.md` for SESSION_END entries |
| Too many tokens used | Say "YAGNI" to simplify, or "skip" unnecessary steps |
| Skill not triggering | Check trigger words in the skill's `SKILL.md` description |

---

## vs. BMAD Method

CTS intentionally trades some BMAD features for simplicity and zero-dependency setup:

| CTS Advantage | BMAD Advantage |
|---|---|
| 19 specialists (Security, DBA, UI, Test Arch, Marketing) | Named agent personas |
| Zero dependencies — plain markdown, no Node.js | NPM installer |
| 8 structured log files + 10 KPIs | Party mode (multi-agent debate) |
| OWASP Top 10 + compliance frameworks | 10+ IDE support |
| Blocker protocol with A/B/C options | Modular expansion packs |
| Adversarial review checklist | Scale-adaptive intelligence |
| Session resumption with snapshots | Web bundles (ChatGPT/Gemini) |

---

*Claude Team Skills v1.0 — 19 specialists, 46 files, 6,087 lines*
