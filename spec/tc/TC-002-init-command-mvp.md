---
id: TC-002
title: Executable checks for init command MVP
status: active
links:
  hls: [HLS-002]
  lls: [LLS-002]
---

# Purpose

Validate that `specguard init` scaffolds required spec directories safely and idempotently.

# Preconditions

- `specguard` binary is available in shell.
- Tests run in isolated temp directories.

# Case 1: Fresh repository initialization

1. In a temp workspace, ensure no `spec/` directory exists.
2. Run `specguard init`.
3. Verify exit code is `0`.
4. Verify `spec/hls`, `spec/lls`, `spec/tc`, and `spec/shared` now exist.
5. Verify output includes `CREATED` lines and success line.

# Case 2: Idempotent re-run

1. In a workspace where all required directories already exist, run `specguard init` twice.
2. Verify both runs exit with `0`.
3. Verify second run outputs `EXISTS` for each required directory.
4. Verify no existing files are changed.

# Case 3: Partial initialization

1. In a temp workspace, create only `spec/hls`.
2. Run `specguard init`.
3. Verify exit code is `0`.
4. Verify missing required directories were created.
5. Verify `spec/hls` remained present and unchanged.

# Case 4: Path conflict

1. In a temp workspace, create file `spec/tc` (not a directory).
2. Run `specguard init`.
3. Verify exit code is `2`.
4. Verify output includes `ERROR spec/tc` and conflict reason.

# Case 5: Compatibility with check

1. In a temp workspace, run `specguard init`.
2. Run `specguard check`.
3. Verify `check` exits with `0` (empty required directories are valid).
