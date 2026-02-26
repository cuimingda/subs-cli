# subs

`subs` is a command-line utility for batch operations on subtitle files in the **current working directory**.

It works only on `.srt` and `.ass` files and is intentionally non-recursive: only files in the current directory are processed, subdirectories are skipped.

## Features

- List subtitle files in the current directory
- Inspect detected encodings
- Convert subtitle encodings to UTF-8
- Inspect or edit ASS style and dialogue font information

## Installation

```bash
# from repository
$ go install github.com/cuimingda/subs-cli/cmd/subs@latest
```

Or run directly from source:

```bash
go run ./cmd/subs [command]
```

> Requires Go 1.24+

## Usage

```bash
subs [command]
```

Use `--help` on any command to view its subcommands.

### Available Commands

- `list`
- `encoding`
  - `list`
  - `reset`
- `dialogue`
  - `font`
    - `list`
    - `prune`
- `style`
  - `font`
    - `list`
    - `reset`

## Commands

### `subs list`

List all `.srt` and `.ass` files in the current directory, one per line.

```bash
subs list
```

### `subs encoding`

Container command for encoding-related operations.

#### `subs encoding list`

List detected encodings for all subtitle files in the current directory.

```bash
subs encoding list
```

Output format:

```text
file.ext - ENCODING
```

If there are no `.srt`/`.ass` files, it returns `ErrNoSubtitleFiles`.

#### `subs encoding reset`

Convert all subtitle files in the current directory to UTF-8 when needed.

```bash
subs encoding reset
```

Output format:

```text
Total N file(s), updated M file(s)
```

### `subs dialogue`

Container command for ASS dialogue operations.

#### `subs dialogue font`

Shows font-related ASS dialogue operations.

##### `subs dialogue font list`

List font names used by `\fn` tags in each `.ass` file in the current directory.

```bash
subs dialogue font list
```

Output format:

```text
file.ass: font1,font2,font3
```

If no `\fn` font is present in a file, output `None` for that file.

Example:

```text
example.ass: Arial,Times New Roman,None
```

##### `subs dialogue font prune`

Remove all `\fn` font tags from every `.ass` file.

```bash
subs dialogue font prune
```

Output format:

```text
Pruned X font tags in Y files.
```

`Y` is always the number of `.ass` files in the current directory.

### `subs style`

Container command for ASS style operations.

#### `subs style font`

Shows style-font related ASS operations.

##### `subs style font list`

List unique font names referenced by the `Fontname` field in `[V4+ Styles]` section for each `.ass` file.

```bash
subs style font list
```

Output format:

```text
file.ass: font1,font2,font3
```

If no style font can be found in a file, output `None` for that file.

##### `subs style font reset`

Replace all style `Fontname` values in `[V4+ Styles]` with `Microsoft YaHei` for every `.ass` file.

```bash
subs style font reset
```

Output format:

```text
Reset X font names in Y file(s).
```

- `X` is the number of font names replaced.
- `Y` is the number of `.ass` files that were updated.

## Behavior Rules

- `subs list` and `subs encoding` commands are always available without UTF-8 preconditions.
- All other subtitle operations that read/modify files are UTF-8 only:
  - If any file is not UTF-8, the command stops and prints:
    `Please run \`subs encoding reset\` to convert subtitle files to UTF-8 first.`
- Running commands with parent-only arguments (for example `subs dialogue`, `subs dialogue font`, `subs style`, `subs style font`) shows help.

## Tests

Run all tests:

```bash
go test ./...
```
