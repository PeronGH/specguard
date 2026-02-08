---
id: HLS-001
title: Check-only MVP for specguard
status: active
links:
  lls: [LLS-001]
  tc: [TC-001]
---

# Outcome

`specguard` MVP provides a single command, `check`, that validates the spec workspace and reports whether it is compliant.

# User-facing behavior

- The tool reads spec files under `spec/`.
- The tool validates required structure and metadata.
- Empty required spec directories are allowed and do not fail validation.
- The tool returns a non-zero exit code when violations exist.
- The tool prints actionable diagnostics with file path context.
- The tool does not modify files in MVP.

# Non-goals for MVP

- No `init`, `fmt`, or `test` commands.
- No auto-fix behavior.
- No schema customization.

# Acceptance scenarios

```gherkin
Feature: Validate spec workspace with a check-only MVP

  Scenario: Workspace is fully compliant
    Given a repository with spec/hls, spec/lls, and spec/tc directories
    And each present spec file has valid naming and front matter
    When I run "specguard check"
    Then the process exits with code 0
    And output indicates success

  Scenario: Required directories are empty
    Given a repository with spec/hls, spec/lls, and spec/tc directories
    And those directories contain no spec files
    When I run "specguard check"
    Then the process exits with code 0
    And output indicates success

  Scenario: Workspace has violations
    Given a repository with one or more invalid spec files
    When I run "specguard check"
    Then the process exits with a non-zero code
    And output lists each violation with file path and reason

  Scenario: Missing required spec directories
    Given a repository without one or more required spec directories
    When I run "specguard check"
    Then the process exits with a non-zero code
    And output includes each missing directory
```
