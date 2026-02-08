# Repository Guidelines

This repository is for `sepcguard` which itself is dogfooding a spec-first, agent-native workflow. `sepcguard` aims to be able to init/check/fmt/lint/test the spec directory.

## Workflow

**Order of operations**

1) **Spec first**: update or add spec before changing production code.
2) **Tests next**: add/adjust tests that encode the spec (failing before the fix, passing after).
3) **Code last**: implement the smallest change that makes tests pass.

**Drift rule**

- If implementation or tests reveal a missing/incorrect requirement, **update the spec first** (drift must flow upward).
- If code behavior changes, **spec must reflect it** before merge.

## Spec directory structure

All specs live under `/spec`:

- `/spec/hls/` — High-level spec (Markdown). Behavior-oriented. May contain Gherkin code blocks.
- `/spec/lls/` — Low-level spec (Markdown). Implementation-facing constraints + interfaces + invariants.
- `/spec/tc/` — Test cases as executable Markdown written for an agent; link to code tests when they exist, otherwise include clear step-by-step validation instructions.
- `/spec/shared/` — Shared Markdown fragments referenced by specs/tests (for reusable setup, rules, or procedures).

## Spec conventions

### IDs & filenames

Every spec/test has an ID and filename prefix:

- `HLS-###-*.md`
- `LLS-###-*.md`
- `TC-###-*.md`

### Markdown front matter

Each spec file should have YAML front matter:

- `id`: must match filename prefix (`HLS-010`, `LLS-120`, `TC-501`)
- `title`
- `status`: `draft | active | deprecated`
- `links`: upstream/downstream IDs (when relevant)

Example:

```yaml
---
id: LLS-120
title: Invite acceptance
status: draft
links:
  hls: [HLS-010-S01]
  tc: [TC-500, TC-501]
---
```
