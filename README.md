# Claude Team Skills

**19 AI Specialists Working as One Team — Built for Claude Code**

> Version 1.0 · Zero dependencies · Token-efficient · YAGNI-driven

Claude Team Skills is a collection of 19 interconnected specialist skills for Claude Code. Instead of a generic AI assistant, every task is handled by the right expert — PM, Security Engineer, DBA, Test Architect, and more — all routed by a central Orchestrator that loads one specialist at a time.

---

## Installation

### Step 1 — Copy the `skills/` folder into your project root

Extract `claude-team-skills-final.tar.gz` and place its contents inside a `skills/` directory:

```
your-project/
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
│   │   └── references/        ← 5 deep-dive security files
│   ├── creative-intelligence/
│   ├── digital-marketer/
│   ├── copywriter/
│   ├── content-marketer/
│   └── project-monitor/
│       ├── references/
│       └── templates/         ← init-logs.sh, generate-report.py
├── CLAUDE.md                  ← Copy this to your project root
├── QUICKSTART.md              ← Cheat sheet
└── .logs/                     ← Auto-created on first session
```

### Step 2 — Copy `CLAUDE.md` to your project root

This is the entry point. Claude Code reads it first on every session and loads the Orchestrator.

### Step 3 — Start working

Open Claude Code and describe what you need:

```
"I need to add user authentication"
"There's a 500 error on the login endpoint"
"Set up a CI/CD pipeline"
"Write a launch blog post for our new feature"
"Design the schema for a multi-tenant SaaS app"
```

The Orchestrator handles the rest.

**Zero dependencies.** Plain markdown files — no Node.js, no npm, no build step.

---

## How It Works

```
You describe the task
       │
       ▼
CLAUDE.md loads → reads Orchestrator
       │
       ▼
Orchestrator identifies task type
       │
       ▼
Loads ONLY the needed specialist
       │
       ▼
Interactive: options → pick → execute → verify
       │
       ▼
Switches specialist with a 2-3 line handoff summary
```

**Token cost:** CLAUDE.md (~1K) + Orchestrator (~2K) + one specialist (~1-2K) = ~3-5K tokens for routing. Compare to loading all 19 skills at once: 35K+ tokens wasted.

---

## The 6-Phase Workflow

Every session follows this flow. Phases can be skipped when not needed.

| Phase | What Happens | Logged To |
|---|---|---|
| 1. Understand | Clarify what, where, how big | `.logs/communications.md` |
| 2. Brainstorm | 2-3 options ranked simple → complex | `.logs/decisions.md` |
| 3. Plan | Break into 30-60 min task batches | `.logs/activity.md` |
| 4. Execute | One task at a time, interactive checkpoints | `.logs/activity.md` |
| 5. Verify | Tests + security review | `.logs/issues.md` |
| 6. Ship | Deploy, document, announce | `.logs/activity.md` |

Options are always presented as:

```
🟢 SIMPLE:        [lowest effort, fastest]
🟡 BALANCED:      [moderate effort, good tradeoffs]
🔴 COMPREHENSIVE: [highest effort, most robust]
```

---

## Three Core Principles

### 1. YAGNI (You Aren't Gonna Need It)

Before building, designing, or planning anything — ask: *"Is this needed RIGHT NOW?"*

| Instead of... | Use... |
|---|---|
| Microservices for an MVP | Monolith first |
| Kubernetes for 3 containers | Docker Compose |
| Custom auth from scratch | NextAuth / Passport |
| Full design system for 3 pages | Pick a UI framework |
| Cache everything upfront | Profile first, then cache bottlenecks |
| Plan for 10M users on day 1 | Build for current scale |

### 2. Token Efficiency

- Load **one specialist at a time** — never all 19
- Handoffs use a **2-3 line summary**, not full history
- Batch tasks into **30-60 min chunks** with checkpoints
- **Confirm approach before generating code** — prevents expensive rewrites
- Log after the action in the **same tool call** — no extra round-trips

### 3. Project Memory

Everything is appended to `.logs/` — decisions, bugs, risks, scope changes, handoffs. Resume any session exactly where you left off.

---

## The Team — 19 Specialists

### Management
| Specialist | What They Do | Trigger Words |
|---|---|---|
| **Orchestrator** | Routes every request to the right specialist | (first to load, always) |
| **Project Manager** | Charter, WBS, risk matrix, estimation, status reports | project plan, scope, milestone, charter |
| **Scrum Master** | Sprints, user stories, backlog, ceremonies, velocity | sprint, backlog, user story, retro, agile |
| **Tech Lead** | ADRs, architecture, YAGNI gate, code standards | architecture, tech stack, refactor, design pattern |

### Security & Data
| Specialist | What They Do | Trigger Words |
|---|---|---|
| **Security Engineer** | STRIDE, auth design, OWASP, compliance (SOC2/HIPAA/PCI-DSS) | threat model, security review, OWASP, auth design |
| **DBA** | Schema, queries, indexes, migrations, backup strategy | database, schema, slow query, migration, index |

### Design
| Specialist | What They Do | Trigger Words |
|---|---|---|
| **UX Designer** | User flows, wireframes, IA, usability | user flow, wireframe, UX, usability |
| **UI Designer** | Design tokens, colors, typography, visual polish | design tokens, UI, color scheme, typography |

### Development
| Specialist | What They Do | Trigger Words |
|---|---|---|
| **Backend Dev** | APIs, services, auth, ORM, business logic | API, endpoint, backend, REST, GraphQL, service |
| **Frontend Dev** | Components, state, accessibility, responsive | component, React, Vue, Next, frontend, CSS |
| **Tester** | Unit / integration / e2e tests, QA, bug reports | test, unit test, e2e, QA, coverage, bug report |
| **Test Architect** | Risk-based strategy, ATDD/Gherkin, adversarial review | test strategy, ATDD, adversarial review, edge case |

### Operations
| Specialist | What They Do | Trigger Words |
|---|---|---|
| **Deployment** | Release, rollback, blue-green, environment management | deploy, release, rollback, staging, production |
| **DevOps/DevSecOps** | CI/CD, Docker, K8s, Terraform, scanning (Semgrep/Trivy/Checkov) | CI/CD, pipeline, Docker, Kubernetes, security scan |

### Creative & Marketing
| Specialist | What They Do | Trigger Words |
|---|---|---|
| **Creative Intelligence** | Brainstorming, SCAMPER, Design Thinking, storytelling, naming | brainstorm, ideate, design thinking, name ideas |
| **Digital Marketer** | SEO, campaigns, funnels, analytics | SEO, campaign, launch, analytics, growth |
| **Copywriter** | Headlines, CTAs, email copy, UX writing | copy, headline, CTA, landing page, email |
| **Content Marketer** | Blog posts, social media, content calendar | blog post, content, social media, newsletter |
| **Project Monitor** | Logs, KPIs, status reports, retrospectives, session resumption | status report, KPIs, retro, what did we do |

---

## Pre-Built Workflows

The Orchestrator recognizes these scenarios and uses an optimised fast path:

| Workflow | Duration | Specialists |
|---|---|---|
| 🐛 Bug Fix | 15-60 min | Tech Lead → Backend/Frontend → Tester |
| ✨ New Feature | 1-4 hrs | PM → UX → Tech Lead → DBA → Devs → Tester → Security |
| 🚀 New Project | 2-6 hrs | PM → UX → Tech Lead → DBA → Security → all devs → DevOps |
| 📋 Document-First Build | 1-2 hrs planning + execution | PM → Tech Lead → DBA → Security → Scrum → Test Architect |
| 🎨 UX/UI Design | 30-120 min | UX → UI → Frontend → Copywriter |
| 🗄️ Database Design | 15-60 min | DBA → Tech Lead → Backend Dev |
| 🔒 Security Audit | 30-90 min | Security Engineer → DevSecOps → DBA |
| 🔐 Security Review | 20-60 min | Security Engineer → DevSecOps → Backend/Frontend |
| 🔄 Refactor | 30-120 min | Tech Lead → Backend/Frontend → Tester |
| 📝 Documentation | 15-60 min | Tech Lead → relevant specialist |
| 📢 Marketing Launch | 1-3 hrs | Digital Marketer → Copywriter → Content Marketer |
| 💡 Brainstorm/Ideation | 15-45 min | Creative Intelligence → PM or Tech Lead |
| 🧪 Test Strategy | 20-60 min | Test Architect → Tester → Security Engineer |

---

## Project Logging

Eight append-only log files auto-created in `.logs/`:

| File | What Gets Logged |
|---|---|
| `activity.md` | Completed tasks, milestones, skill changes |
| `decisions.md` | Architecture choices, approach selections (mini-ADRs) |
| `issues.md` | Bugs, blockers (with severity and status) |
| `risks.md` | Risks, probability, mitigation, status |
| `corrections.md` | Scope changes, pivots, plan corrections |
| `communications.md` | User preferences, handoffs, inter-skill consultations |
| `sessions.md` | SESSION_START / SESSION_END snapshots |
| `metrics.md` | Sprint/batch KPI snapshots |

### Log Entry Format
```
### [YYYY-MM-DD HH:MM] [CATEGORY] — [Short Title]
- Specialist: [who logged this]
- Summary: [1-2 sentences]
- Status: open | in-progress | resolved
- Impact: low | medium | high | critical
```

### Report Commands
```
"status report"        → Full project status with KPIs
"what did we do?"      → Activity log
"what's broken?"       → Open issues and blockers
"what did we decide?"  → Decision log
"what changed?"        → Plan corrections
"retro"                → Retrospective analysis
"KPIs"                 → Project metrics
```

### Session Resumption
On session start, the Orchestrator reads only the last `SESSION_END` entry from `sessions.md` and asks:
```
📋 Last session summary:
- Completed: [what was done]
- In progress: [what was started]
- Blocked: [any blockers]
- Next: [what was planned]

"Continue from here, or start something new?"
```

---

## KPIs Tracked

| KPI | Formula | Target |
|---|---|---|
| Velocity | Tasks completed per sprint | Stable or increasing |
| Completion Rate | Completed / Planned × 100 | ≥ 80% |
| Bug Rate | Bugs found / Tasks completed | < 0.5 |
| Bug Fix Rate | Bugs fixed / Bugs found × 100 | ≥ 90% |
| Scope Creep | Scope changes / Original tasks × 100 | < 20% |
| Estimation Accuracy | Actual / Estimated time | 0.8 – 1.2 |
| Blocker Duration | Time from open → resolved | < 1 session |
| Security Issue Rate | Security issues / Tasks completed | < 0.2 |
| Rework Rate | Corrections / Completed × 100 | < 15% |

---

## Document-First Development

For new projects and large features (3+ days of work):

```
PRD → Architecture Doc → Epics & Stories → Implementation
 PM    Tech Lead + DBA +    Scrum Master +     Dev team
        Security Eng         Test Architect
```

Each artifact is saved to `docs/`, logged as a MILESTONE, and requires user approval before proceeding.

**YAGNI scaling:**
- Small task (< 1 day): skip the chain entirely
- Medium feature: mini-PRD + ADR only
- Large feature / new project: full chain

---

## Security Coverage

The **Security Engineer** provides:
- **STRIDE threat modeling** — fast 5-minute threat model or full rigor when needed
- **Auth design decision tree** — API keys / sessions / JWT / mTLS / SSO based on context
- **Auth checklist** — bcrypt/argon2id, MFA, HttpOnly cookies, JWT claim validation, account lockout
- **Authorization patterns** — Simple roles / RBAC / ABAC / ReBAC
- **Security review checklist** — input/output, auth, data, infra
- **Compliance frameworks** — SOC2, PCI-DSS, HIPAA, ISO 27001, NIST, GDPR
- **Incident response template** — Detect → Assess → Contain → Fix → Recover → Learn

The **DevOps/DevSecOps** engineer runs automated scanning:
```bash
semgrep ci --config p/owasp-top-ten    # SAST
trivy fs --severity CRITICAL,HIGH .    # SCA
trivy image --severity CRITICAL,HIGH   # Container scan
checkov -d ./terraform                 # IaC scan
gitleaks detect --source .             # Secrets scan
```

---

## Test Architecture

The **Test Architect** (distinct from the basic Tester) designs the test system:

### Risk-Based Strategy
Score components by Impact + Frequency + Complexity:
- **13-15**: Maximum (unit + integration + e2e + adversarial + load)
- **9-12**: High (unit + integration + e2e)
- **5-8**: Standard (unit + integration)
- **1-4**: Minimal (smoke tests only)

### ATDD / Gherkin
Acceptance criteria written *before* coding — they become the regression suite.

### Adversarial Review
Automatically triggered for auth, payments, and data handling:
- **Input abuse** — SQLi, XSS, command injection, path traversal, type confusion
- **Auth & access abuse** — IDOR, CSRF, expired tokens, privilege escalation
- **Race conditions** — double-submit, TOCTOU, concurrent updates
- **Business logic abuse** — negative quantities, coupon stacking, skipping required steps

### Edge Case Hunter
Systematic discovery across: boundaries, types, time (DST, leap year), state (empty/full/mid-migration), and network (timeout, partial response, retry storm).

---

## Creative Intelligence Suite

7 structured ideation frameworks, matched to the problem:

| Framework | Best For | Time |
|---|---|---|
| Structured Brainstorming | Generating 10+ ideas fast (diverge → cluster → converge) | 5-10 min |
| SCAMPER | Improving existing solutions (Substitute/Combine/Adapt/Modify/Put to use/Eliminate/Reverse) | 5 min |
| Design Thinking | User-centered innovation (Empathize → Define → Ideate → Prototype → Test) | 15-30 min |
| Storytelling (Hero's Journey) | Pitches and product narratives | 10-15 min |
| First Principles | Breaking assumptions, rebuilding from facts | 10 min |
| Problem Reframing | Unsticking hard problems (4 reframe angles) | 5 min |
| Name Generation | Product/feature naming (descriptive → metaphorical → invented → compound) | 5 min |

All sessions score ideas with: **Impact × Confidence / Effort**

---

## Quick Command Reference

| Say this | What it does |
|---|---|
| `"simpler"` | Switch to a less complex approach |
| `"YAGNI"` | Remind to keep it simple |
| `"option A/B/C"` | Pick from presented options |
| `"good, next"` | Approve current step, move on |
| `"skip this"` | Jump to next task |
| `"go back"` | Revisit a previous decision |
| `"stop"` | Pause and discuss |
| `"what's the plan?"` | Show current batch/task list |
| `"ship it"` | Move to deployment |
| `"test this"` | Run tests on current work |
| `"secure this"` | Run security review |
| `"adversarial review"` | Run adversarial security checklist |
| `"test strategy"` | Design risk-based test plan |
| `"brainstorm this"` | Start creative ideation session |
| `"status report"` | Generate project status with KPIs |
| `"retro"` | Generate retrospective from logs |

### Blocker Protocol
When stuck, specialists never spin — they immediately present:
```
🚧 BLOCKER: [what's wrong]
  A) [simple workaround]
  B) [proper fix]
  C) [skip for now]
Which one?
```

---

## Token Budget by Session Type

| Session Type | Expected Tokens | Skills Loaded |
|---|---|---|
| Quick bug fix | 3-8K | Orchestrator + 1 dev + Tester |
| Small feature | 8-15K | Orchestrator + Tech Lead + 1-2 devs + Tester |
| Large feature | 15-30K | Multiple specialists across batches |
| New project scaffold | 15-25K | PM + Tech Lead + DBA + devs + DevOps |
| Database design | 5-10K | DBA + Backend Dev |
| UX/UI design | 5-12K | UX + UI + Frontend + Copywriter |
| Security audit | 5-15K | Security Eng + DevSecOps + references |
| Brainstorm session | 3-8K | Creative Intelligence + PM or Tech Lead |
| Test strategy | 5-10K | Test Architect + Tester + Security Eng |
| Document-first build | 10-20K | PM + Tech Lead + Scrum + devs |
| Sprint planning | 3-8K | PM + Scrum Master |

---

## Customization

**Edit `CLAUDE.md`** — Add your tech stack, naming conventions, deployment URLs, team contacts, and project-specific rules.

**Disable unused skills** — Remove rows from the skill locations table in `CLAUDE.md`. Backend-only project? Remove `frontend-dev`, `copywriter`, `content-marketer`, `digital-marketer`.

**Add new specialists** — Create `skills/your-skill/SKILL.md` with:
1. YAML frontmatter: `name` + `description` + trigger words
2. Role definition
3. Key templates / checklists
4. Handoff points to other specialists

Then add the skill to the Orchestrator's routing table and the `CLAUDE.md` skill locations table.

---

## Troubleshooting

| Problem | Solution |
|---|---|
| Claude doesn't use skills | Ensure `CLAUDE.md` is in the project root |
| Wrong specialist loaded | Say `"switch to [specialist]"` or describe the task more specifically |
| Logs not created | Run `bash skills/project-monitor/templates/init-logs.sh` |
| Session not resuming | Check `.logs/sessions.md` for `SESSION_END` entries |
| Too many tokens used | Say `"YAGNI"` or `"simpler"` to reduce scope |
| Skill not triggering | Check the trigger words in the skill's `SKILL.md` frontmatter |

**Commit message convention:** `type(scope): description` — e.g. `feat(auth): add refresh token rotation`

**Branch naming:** `type/short-description` — e.g. `feature/add-auth`, `fix/login-bug`

---

## Intentional Trade-offs (vs BMAD Method)

CTS skips these BMAD features by design to stay **lightweight and zero-dependency**:

| BMAD Feature Skipped | Why |
|---|---|
| NPM installer | CTS is zero-dependency — copy markdown, no Node.js |
| Multi-IDE support (10+) | Optimised for Claude Code specifically |
| Named agent personas | Role-based skills use fewer tokens and are easier to maintain |
| Party Mode (multi-agent debate) | Sequential handoffs save ~70% tokens |
| Modular expansion packs | Single bundle — one download, everything included |
| Scale-adaptive intelligence | Explicit YAGNI gates — you pick complexity, not auto-detected |
| Agent-as-Code (YAML compilation) | Plain markdown — editable in any text editor |
| Web bundles (ChatGPT/Gemini) | Claude Code native — one platform done well |

---

*Claude Team Skills v1.0 — 19 specialists, 46 files, 6,087 lines*
