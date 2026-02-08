---
id: TC-001
title: Executable checks for check-only MVP
status: active
links:
  hls: [HLS-001]
  lls: [LLS-001]
---

# Purpose

Validate that `specguard check` enforces the MVP contract.

# Preconditions

- `specguard` binary is available in shell.
- Tests run from repository root.

# Case 1: Success path

1. Ensure `spec/hls`, `spec/lls`, and `spec/tc` exist.
2. Ensure each directory contains at least one valid `.md` spec file with valid front matter.
3. Run `specguard check`.
4. Verify exit code is `0`.
5. Verify output contains success indication and no `ERROR` lines.

# Case 2: Missing required directory

1. Temporarily rename `spec/tc` to `spec/tc.bak`.
2. Run `specguard check`.
3. Verify exit code is `1`.
4. Verify output includes missing directory violation for `spec/tc`.
5. Restore `spec/tc`.

# Case 3: Filename and ID mismatch

1. Create `spec/lls/LLS-099-temp.md` with front matter `id: LLS-100`.
2. Run `specguard check`.
3. Verify exit code is `1`.
4. Verify output includes rule key `id-mismatch` for that file.
5. Remove temp file.

# Case 4: Invalid status value

1. Create a temp file under `spec/tc` with `status: invalid`.
2. Run `specguard check`.
3. Verify exit code is `1`.
4. Verify output includes invalid status violation.
5. Remove temp file.

# Case 5: Broken link reference

1. Create a temp file with `links` containing a non-existent ID.
2. Run `specguard check`.
3. Verify exit code is `1`.
4. Verify output includes unresolved link violation.
5. Remove temp file.

# Case 6: Missing Gherkin block in HLS

1. Create `spec/hls/HLS-099-temp.md` with valid front matter but no fenced `gherkin` block.
2. Run `specguard check`.
3. Verify exit code is `1`.
4. Verify output includes missing gherkin block violation.
5. Remove temp file.
