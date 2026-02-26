# subs

`subs` is a focused command-line utility for batch operations on subtitle files in the **current working directory**.

It is designed for `.srt` and `.ass` files and is intentionally non-recursive: only files in the current directory are processed, subdirectories are skipped.

## Goal

The tool is intended to grow into a small toolkit for common subtitle maintenance tasks, such as:

- listing subtitle files
- inspecting encoding information
- normalizing encoding to UTF-8
- processing ASS `\fn` font tags

## Installation

```bash
# from repository
$ go install github.com/cuimingda/subs-cli/cmd/subs@latest
```

You can also run it directly from source:

```bash
go run .
```

> Requirements: Go 1.24+

## Usage

```bash
subs [command]
```

Use `--help` on any command to view its subcommands.

### Available Commands

- `list`
- `encoding`
- `dialogue`

## Commands

### `subs list`

List all `.srt` and `.ass` files in the current directory.

```bash
subs list
```

### `subs encoding`

Container command for encoding-related operations.

#### `subs encoding list`

List each subtitle file and the detected encoding in the current directory.

```bash
subs encoding list
```

Output format:

```text
file.ext - ENCODING
```

If no `.srt`/`.ass` file exists, the command returns `ErrNoSubtitleFiles`.

#### `subs encoding reset`

Convert all subtitle files in the current directory to UTF-8 when needed.

```bash
subs encoding reset
```

Output:

```text
Total N file(s), updated M file(s)
```

(`N` is total subtitle files, `M` is how many files were changed.)

### `subs dialogue`

Container command for dialogue-related ASS operations.

#### `subs dialogue font`

Shows font-related ASS operations.

#### `subs dialogue font list`

List font names used by `\fn` tags in every `.ass` file. Output is one line per file:

```text
file.ass: font1,font2,font3
```

If no `\fn` font appears in a file, `None` is shown.

```bash
subs dialogue font list
```

#### `subs dialogue font prune`

Remove `\fn` font tags from all `.ass` files in the current directory.

```bash
subs dialogue font prune
```

Output format:

```text
Pruned X font tags in Y files.
```

`Y` is always the number of `.ass` files in the current directory; repeated runs are supported.

## Behavior Notes

- Commands with no extra arguments generally show help for the parent command.
- All core logic lives under `internal/subtitles` and is covered by tests.
- Current scope is one directory level only; recursion is intentionally not used.

## Tests

Run all tests:

```bash
go test ./...
```
