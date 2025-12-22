# cliface

A lightweight GUI wrapper for command-line tools. Define forms via TOML config to generate graphical interfaces that construct and execute shell commands.

## Features

- Config-driven UI generation from TOML files
- Multiple input types: string, number, boolean, choice
- File and directory pickers
- Multi-value fields with add/remove buttons
- Field validation (required, regex, range)
- Conditional field visibility
- Multiple execution modes: visible window, dialog output, realtime streaming
- Cross-platform (macOS, Windows, Linux)

[中文文档](README_zh.md)

## Repository

https://github.com/listeng/cliface

## Install

```bash
go build -o cliface .
```

## Usage

```bash
cliface -c examples/curl.toml
```

Without `-c`, reads `config.toml` from the executable's directory.

## Config Example

```toml
title = "My Tools"
width = 500
height = 400

[[apps]]
[apps.command]
path = "/usr/bin/ffmpeg"
name = "Video Convert"
args = ["-y"]
mode = "hidden"      # hidden | visible
output = "dialog"    # dialog | realtime
debug = true

[[apps.items]]
text = "Convert video to MP4"  # label only

[[apps.items]]
name = "i"
short = true
type = "string"
label = "Input"
picker = "file"
separator = " "

[[apps.items]]
name = "crf"
type = "number"
label = "Quality"
default = 23

[[apps.items]]
name = "preset"
type = "choice"
label = "Speed"
choices = ["ultrafast", "fast", "medium", "slow"]
default = "medium"
```

## Config Reference

### Global

| Field | Description | Default |
|-------|-------------|---------|
| title | Window title | command name or "cliface" |
| width | Window width | 400 |
| height | Window height | 300 |

### Command

| Field | Description |
|-------|-------------|
| path | Executable path |
| name | Display name (tab title for multiple apps) |
| args | Fixed arguments |
| mode | `hidden` or `visible` window |
| output | `dialog` (show after completion), `realtime` (streaming window), or `realtime-console` (streaming to terminal) |
| debug | Show "Show Command" button |
| run_text / run_color | Run button text and color (high/danger/warning/success/low) |
| debug_text / debug_color | Debug button text and color |
| env | Environment variables as key-value pairs |

### Item

| Field | Description |
|-------|-------------|
| text | Label text (ignores other fields if set) |
| name | Argument name |
| type | `string` / `number` / `bool` / `choice` |
| short | Use single dash `-name` if true |
| positional | Positional argument (no prefix) if true |
| label | Display label |
| description | Field description |
| default | Default value |
| choices | Options for choice type |
| picker | `file` or `directory` picker |
| picker_text | Custom picker button text |
| separator | Arg separator: `" "` for space, `"none"` for no separator, default `=` |
| multi | Allow multiple values with add/remove buttons |
| required | Field must have a value before running |
| validate | Regex pattern for validation |
| min / max | Number range validation |
| condition | Show/enable based on another field (e.g., `field=value` or `field!=value`) |

## License

MIT
