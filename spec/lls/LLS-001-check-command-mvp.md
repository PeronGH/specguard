---
id: LLS-001
title: Validation rules for check-only MVP
status: active
links:
  hls: [HLS-001]
  tc: [TC-001]
---

# Command contract

- CLI entrypoint: `specguard check`
- Default scan root: `spec/`
- The command is read-only.

# Library choices (MVP)

- YAML parser: `gopkg.in/yaml.v3` (required).
- CLI argument parsing: Go standard library `flag` package (required).
- Markdown parsing: no third-party Markdown parser in MVP.

Implementation notes:

- Front matter delimiter extraction uses simple text scanning for leading `---` blocks.
- HLS Gherkin validation uses text scanning for fenced blocks starting with ```` ```gherkin ````.
- Introducing additional libraries (for example Cobra or Goldmark) is out of MVP scope.

# Required directory rules

The following directories must exist:

- `spec/hls`
- `spec/lls`
- `spec/tc`

Directory presence is required, but directories are allowed to be empty.

# File discovery rules

- Only `.md` files under `spec/hls`, `spec/lls`, and `spec/tc` are considered.
- Expected filename prefixes by directory:
  - `spec/hls`: `HLS-`
  - `spec/lls`: `LLS-`
  - `spec/tc`: `TC-`
- Filename pattern must match `<PREFIX><3 digits>-<slug>.md`.

# Front matter rules

Each discovered spec file must contain YAML front matter with:

- `id` (string, required)
- `title` (string, required)
- `status` (enum: `draft`, `active`, `deprecated`)
- `links` (object, optional)

Validation constraints:

- `id` must match filename ID prefix (example: `HLS-001` in `HLS-001-*.md`).
- `title` must be non-empty.
- `status` must be one of the allowed enum values.
- If `links` is present, each value must be a list of IDs.

# Cross-reference rules

- All IDs are globally unique across `spec/hls`, `spec/lls`, and `spec/tc`.
- Each linked ID in `links` must resolve to an existing file ID.

# HLS-specific rules

- Every file under `spec/hls` must contain at least one fenced Gherkin code block:
  - opening fence: ```` ```gherkin ````
  - closing fence: ```` ``` ````

# Exit codes

- `0`: no violations
- `1`: one or more validation violations
- `2`: execution error (example: cannot read filesystem)

# Diagnostic format

- One violation per line.
- Each line includes:
  - severity (`ERROR`)
  - file path or logical target (for missing directory)
  - short rule key
  - human-readable message

Example:

`ERROR spec/hls/HLS-001-foo.md id-mismatch: front matter id HLS-002 does not match filename HLS-001`
