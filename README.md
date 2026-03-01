# Cleaner

`cleaner` helps you declutter a folder by moving files to a target directory using date-based subfolders. My use case is for removing screenshots from the desktop of my MacOS machine.

It has two modes:

- **Screenshot mode (default):** moves files that start with `Screenshot`
- **Extension mode (`-ext`):** recursively moves files by extension (for example `png` or `.png`)

Moved files are organized as `TARGET/YYYY-MM-DD/`.
If a filename contains a date like `2024-07-30`, cleaner uses that date.
Otherwise, it uses the file's modification time.

## 1) Install

Requires Go `1.23.1+`.

```bash
go install github.com/ezebunandu/cleaner/cmd/cleaner@latest
```

## 2) Learn the command

```text
cleaner [flags] <SOURCE> <TARGET>
```

Flags:

- `-ext string` match extension (example: `png` or `.png`)
- `--dry-run` preview matched files without moving anything

If required arguments are missing, cleaner prints usage.

## 3) Try common workflows

### Move screenshots from Desktop

```bash
cleaner ~/Desktop ~/Pictures/screenshots
```

### Move all PNG files recursively

```bash
cleaner -ext png . ./target
```

### Preview first with dry run

```bash
cleaner -ext png --dry-run . ./target
```

Example dry-run output:

```text
would have moved the following files from . to ./target
./sample.png
```

## 4) Know what to expect

- No matches: `no files to move`
- Success: `moved <N> files to <TARGET>`
- Dry run: prints source, target, and matched files
- Errors (missing paths/permissions): printed to stderr and exit non-zero

## 5) Run tests

```bash
go test ./...
```
