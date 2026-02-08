# Repository Guidelines

This repository is for `specguard` which itself is dogfooding a spec-first, agent-native workflow. `specguard` aims to be able to init/check/fmt/test the spec directory.

## Workflow

**Order of operations**

1) **Spec first**: update or add spec before changing production code.
2) **Tests next**: add/adjust tests that encode the spec (failing before the fix, passing after).
3) **Code last**: implement the spec fully.

**Drift rule**

- If implementation or tests reveal a missing/incorrect requirement, **update the spec first** (drift must flow upward).
- If code behavior changes, **spec must reflect it** before merge.

## Spec directory structure

All specs live under `/spec`:

- `/spec/hls/`
  - High-level spec in Markdown.
  - Behavior-oriented.
  - Must contain Gherkin code blocks.
  - One HLS can correspond to multiple LLS files.
- `/spec/lls/`
  - Must be implementable in one go; otherwise split into multiple LLS files.
  - Low-level spec in Markdown.
  - Implementation-facing.
  - Defines constraints, interfaces, and invariants.
  - One LLS can have multiple TC files.
- `/spec/tc/`
  - Test cases as executable Markdown written for an agent.
  - Each TC should use one shared precondition.
  - Link to code tests when they exist.
  - Otherwise include clear step-by-step validation instructions.
  - All cases in a TC must share same test environments; otherwise split into multiple TC files.
- `/spec/shared/`
  - Shared Markdown fragments referenced by specs/tests.
  - Use for reusable setup, rules, or procedures.

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
  hls: [HLS-010]
  tc: [TC-500, TC-501]
---
```

## Commit conventions

- Commit messages must follow `type(scope): summary`.
- Use `Specs:` trailer for any commit that changes spec files, spec-driven tests, or implementation behavior; optional for unrelated repo/meta changes.
- Keep `Specs:` minimal:
  - include `LLS-...` when implementing behavior,
  - include `TC-...` when tests/manual TC steps are run,
  - include `HLS-...` only when editing HLS directly or writing LLS from HLS.
- Trailer format example:

```text
Specs: LLS-001, TC-001
```
