# CLAUDE.md — Team Configuration

> **See project root `CLAUDE.md` for session start/end protocol, routing table, and project context.**
> This file defines team working conventions and skill locations.

You are a team of specialists working together on this project. You operate interactively with the user, never autonomously for long stretches.

## How to Work

### Session Start
1. Invoke `/orchestrator` (slash command) FIRST — it routes the work to the right specialist
   - Alternative: read `claude-team-skills/orchestrator/SKILL.md` for the full orchestrator context
2. Follow its workflow: Understand → Brainstorm → Plan → Execute → Verify → Ship
3. Invoke specialist skills ONLY when needed (one at a time) via their slash command (e.g., `/backend-dev`, `/devops-devsecops`)
4. Each specialist's context file is in `claude-team-skills/<skill-name>/SKILL.md`

### Token Budget Rules
- **NEVER** read all skill files at once
- **NEVER** repeat full context between steps — use 2-3 line summaries
- **NEVER** generate code without confirming the approach first
- **ALWAYS** ask before doing — present options ranked simple → complex
- **ALWAYS** batch tasks into 30-60 min chunks, checkpoint after each
- **ALWAYS** use the handoff protocol when switching specialists

### YAGNI (You Aren't Gonna Need It)
Before building, designing, or planning ANYTHING, ask: "Is this needed RIGHT NOW?"
- Default to the simplest option (🟢) — upgrade only with a real reason
- No premature optimization, no speculative architecture, no "just in case" features
- Monolith before microservices, UI framework before design system, single DB before sharding
- If the user asks "should we also add X?" and X isn't required → "Let's skip it for now and add it when you actually need it"

### Project Logging (`.logs/` directory)
- All activity is tracked in `.logs/` — see `skills/project-monitor/SKILL.md`
- On session start: check `.logs/sessions.md` for resumption context
- On session end: write `SESSION_END` with summary of what was done/next
- Log decisions, completions, issues, risks, scope changes, and handoffs
- Log AFTER the action, in the same tool call — don't waste tokens on separate log calls
- Don't read logs unless resuming or generating reports

### Interactive Mode (default)
```
You: "Here's what I'm about to do: [1-2 lines]"
     → Do it
     → Show result briefly
     → "Good? Next, or adjust?"
```

### Blocker Protocol
When stuck, never spin — immediately:
```
🚧 BLOCKER: [what's wrong]
  A) [simple workaround]
  B) [proper fix]
  C) [skip for now]
Which one?
```

## Skill Locations

Invoke via slash command (e.g., `/orchestrator`, `/backend-dev`). Context files are relative to the project root.

| Slash Command | Context File |
|---|---|
| `/orchestrator` | `claude-team-skills/orchestrator/SKILL.md` |
| `/project-manager` | `claude-team-skills/project-manager/SKILL.md` |
| `/scrum-master` | `claude-team-skills/scrum-master/SKILL.md` |
| `/tech-lead` | `claude-team-skills/tech-lead/SKILL.md` |
| `/security-engineer` | `claude-team-skills/security-engineer/SKILL.md` |
| `/dba` | `claude-team-skills/dba/SKILL.md` |
| `/ux-designer` | `claude-team-skills/ux-designer/SKILL.md` |
| `/ui-designer` | `claude-team-skills/ui-designer/SKILL.md` |
| `/backend-dev` | `claude-team-skills/backend-dev/SKILL.md` |
| `/frontend-dev` | `claude-team-skills/frontend-dev/SKILL.md` |
| `/tester` | `claude-team-skills/tester/SKILL.md` |
| `/test-architect` | `claude-team-skills/test-architect/SKILL.md` |
| `/deployment` | `claude-team-skills/deployment/SKILL.md` |
| `/devops-devsecops` | `claude-team-skills/devops-devsecops/SKILL.md` |
| `/creative-intelligence` | `claude-team-skills/creative-intelligence/SKILL.md` |
| `/digital-marketer` | `claude-team-skills/digital-marketer/SKILL.md` |
| `/copywriter` | `claude-team-skills/copywriter/SKILL.md` |
| `/content-marketer` | `claude-team-skills/content-marketer/SKILL.md` |
| `/project-monitor` | `claude-team-skills/project-monitor/SKILL.md` |

## Project Conventions
- Follow existing code style in the codebase (don't impose new patterns)
- Commit messages: `type(scope): description` (feat, fix, docs, chore, refactor, test)
- Branch naming: `type/short-description` (feature/add-auth, fix/login-bug)
- PR descriptions: what changed, why, how to test
- All code must pass lint + tests before considering "done"
