# Security Best Practices Review Report

## Executive summary
Reviewed the Go CLI codebase (`github.com/cuimingda/subs-cli`) against Go security guidance for CLI applications. No obvious remote-exposure vulnerabilities (no network/external execution/auth/database/templating) were found. The primary risks are local integrity/availability hardening gaps in file-processing paths.

## Scope
- Language/framework identified: Go 1.23.4 CLI with Cobra.
- Reviewed entrypoints in `cmd/*` and file processing logic in `internal/subtitles/*`.

## Critical findings
None identified.

## High severity

### 1) Unbounded in-memory file processing (availability)
- Status: **Fixed**
- Location: `internal/subtitles/encoding.go:53`, `internal/subtitles/encoding.go:81`, `internal/subtitles/file_limits.go`, `internal/subtitles/reset.go:46`, `internal/subtitles/style_fonts.go:181`, `internal/subtitles/dialogue_fonts.go:57`
- Impact: The CLI reads full subtitle files into memory before validation or conversion (`os.ReadFile`, `io.ReadAll`). A very large or intentionally crafted file can consume excessive memory and cause slowdowns or OOM termination, including when running in automation pipelines.
- Fix:
  - Add a maximum input size policy before processing (`os.Stat` + hard cap).
  - Added input size enforcement in `validateSubtitleFileSize` (`internal/subtitles/file_limits.go`) with an environment override (`SUBS_MAX_SUBTITLE_FILE_BYTES`).
  - Applied size checks before file reads in all subtitle read/convert paths.
- Mitigation: Document operational constraints and enforce worker-level resource limits in CI/local runners.
- False-positive note: This is currently local-file only (no network/file upload endpoint), but damage potential is still realistic in shared/CI contexts.

## Medium severity

### 2) File writes use a fixed `0644` mode, potentially widening file permissions
- Status: **Fixed**
- Location: `internal/subtitles/file_write.go`, `internal/subtitles/reset.go:69`, `internal/subtitles/style_fonts.go:253`, `internal/subtitles/dialogue_fonts.go:67`
- Impact: Files are rewritten with `os.WriteFile(..., 0o644)` regardless of original mode. This can silently broaden existing permissions and make previously restricted subtitle files more accessible after conversion/reset operations.
- Fix:
  - Added shared helper `writeFilePreserveMode` to reuse original file permissions (`os.Stat(path).Mode().Perm()`) when rewriting.
  - Replaced direct write calls in mutation paths with permission-preserving helper.
  - Avoid forcing broad world-readable permissions in user directories containing sensitive content.
- Mitigation: At minimum, document that the tool may alter permission bits and provide a `--preserve-perms` option.

### 3) Non-atomic in-place rewrites increase corruption risk
- Location: `internal/subtitles/reset.go:69`, `internal/subtitles/style_fonts.go:253`, `internal/subtitles/dialogue_fonts.go:67`
- Impact: The code writes converted content directly back to the original path. Process interruption (crash/power loss) can leave files partially written, which can result in data loss. In hostile local environments, this also increases overwrite risk.
- Fix:
  - Write to a temp file in the same directory, `fsync`, then atomically `os.Rename`.
  - Keep a `.bak` or checksum strategy for recoverability when running bulk operations.
- Mitigation: Provide dry-run and backup options before bulk edits.

## Low severity

### 4) Supply-chain/toolchain hardening is not explicitly enforced in repository configuration
- Location: `go.mod` and repository root configuration (no CI config files found)
- Impact: No visible CI policy enforces toolchain minimum versions, dependency scanning, or `govulncheck`.
- Fix:
  - Add CI checks for `go version` policy and `govulncheck`.
  - Track dependency updates in a periodic policy.
- False-positive note: Absence in repo does not prove production CI does not exist elsewhere, but risk is not observable from checked-in files.

## Additional notes
- No HTTP handlers, SQL, subprocess execution, OS command calls, shell usage, or template execution were found in the scanned codebase.
- No hard-coded secrets or credential-like values were found.

## Recommended fix order
1. Add size limits before file parsing/conversion (High)
2. Preserve file permissions and switch to atomic writes (Medium)
3. Add vuln checks/toolchain policy in CI (Low)
