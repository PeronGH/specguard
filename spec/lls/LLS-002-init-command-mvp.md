---
id: LLS-002
title: Init command rules for spec workspace scaffolding
status: active
links:
  hls: [HLS-002]
  tc: [TC-002]
---

# Command contract

- CLI entrypoint: `specguard init`
- Target root: current working directory
- Spec root: `spec/`

# Required directories

`specguard init` ensures all of the following directories exist:

- `spec/hls`
- `spec/lls`
- `spec/tc`
- `spec/shared`

# Behavior rules

- If a required directory does not exist, create it.
- If a required directory already exists, leave it unchanged.
- If a required path exists but is not a directory, report conflict and fail.
- Existing files under `spec/` are never modified by `init`.
- Command is idempotent across repeated runs.

# Output rules

- Print one line per required directory:
  - `CREATED <path>` when newly created
  - `EXISTS <path>` when already present
- On success, print final line `OK spec init complete`.
- On failure, print `ERROR <path> <rule-key>: <message>` for the blocking conflict or error.

# Exit codes

- `0`: initialization completed (with any mix of `CREATED` and `EXISTS`)
- `2`: execution error (filesystem errors, permission failures, path conflicts)

# Library choices (MVP)

- CLI argument parsing: Go standard library `flag` package (required).
- Filesystem operations: Go standard library `os` and `path/filepath` packages (required).
- No additional third-party libraries are required for `init`.
