---
name: standup
description: Generate a standup summary (done/doing/blocked). Use for progress check-ins.
---

# Standup Summary

Generate a development progress summary.

## Instructions

1. Check `git log --oneline -20` for recent work
2. Check `git status` for in-progress work
3. Check todo list for planned tasks
4. Review `docs/REQUIREMENTS.md` for remaining scope

## Output Format

### Done (完了)
- What was completed since last check-in
- Reference specific commits/features

### Doing (進行中)
- What is currently being worked on
- Current branch and changes in progress

### Blocked (ブロッカー)
- Any blockers or decisions needed
- Dependencies waiting on

### Next Up (次にやること)
- Priority tasks from backlog
- Estimated focus for next session

### Progress Overview
```
Feature Progress:
[████████░░] 80% - 本の管理
[██████░░░░] 60% - メモ機能
[████░░░░░░] 40% - タグ機能
[██░░░░░░░░] 20% - ポケモンシステム
[░░░░░░░░░░]  0% - 認証(SSO)
```
