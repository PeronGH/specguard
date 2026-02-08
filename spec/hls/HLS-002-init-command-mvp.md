---
id: HLS-002
title: Init command MVP for spec workspace scaffolding
status: active
links:
  lls: [LLS-002]
  tc: [TC-002]
---

# Outcome

`specguard` provides an `init` command that prepares a repository for spec-first workflow by creating the required spec directory structure.

# User-facing behavior

- The command creates missing spec directories under `spec/`.
- The command is idempotent: running it multiple times is safe.
- The command does not delete or overwrite existing files in MVP.
- The command reports what it created and what already existed.

# Non-goals for MVP

- No template/spec-file generation.
- No migration of existing custom layouts.
- No destructive cleanup behavior.

# Acceptance scenarios

```gherkin
Feature: Initialize spec workspace layout

  Scenario: Fresh repository without spec directory
    Given a repository with no spec directory
    When I run "specguard init"
    Then directories spec/hls, spec/lls, spec/tc, and spec/shared are created
    And the process exits with code 0
    And output indicates initialization success

  Scenario: Repository already initialized
    Given a repository where spec/hls, spec/lls, spec/tc, and spec/shared already exist
    When I run "specguard init"
    Then no existing content is modified
    And the process exits with code 0
    And output indicates each directory already exists

  Scenario: Partially initialized repository
    Given a repository where only some required spec directories exist
    When I run "specguard init"
    Then the missing required directories are created
    And existing directories remain unchanged
    And the process exits with code 0

  Scenario: Path conflict blocks initialization
    Given a repository where a required directory path is occupied by a file
    When I run "specguard init"
    Then the process exits with a non-zero code
    And output identifies the conflicting path
```
