# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
# Build for current platform
go build -o cliface .

# Cross-compile (requires CGO)
./build.sh  # Builds darwin-amd64 and windows-amd64
```

## Architecture

cliface is a GUI wrapper for CLI tools built with Fyne. It reads a TOML config file and generates a form-based interface that constructs and executes shell commands.

**Files:**
- `main.go` - Entry point, config loading, window setup
- `config.go` - TOML config structs (Config → Apps → Command/Items)
- `ui.go` - Form generation from config, widget creation, argument building
- `executor.go` - Command execution with three modes: visible (detached), dialog (blocking with output), realtime (streaming)

**Config-driven UI:** Each `[[apps]]` block defines a tab with a command and form items. Items map to CLI arguments via `name`, `type`, and `separator` fields. The `BuildArgs()` method in ui.go constructs the final command line.

**Execution modes:** Set via `command.mode` (visible/hidden) and `command.output` (dialog/realtime/realtime-console).

**Key features:**
- Multi-value fields (`multi: true`) with add/remove buttons
- Positional arguments (`positional: true`) without prefix
- Customizable button text/colors (`run_text`, `run_color`, `debug_text`, `debug_color`)
- Platform-aware quoting (single quotes on Unix, double quotes on Windows)
- Environment variables (`env` in command config)
- Required field validation (`required: true`)
- Regex and range validation (`validate`, `min`, `max`)
- Conditional enable/disable (`condition: "field=value"`)
