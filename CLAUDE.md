# tmux-launcher

A smart tmux wrapper written in Go that automatically manages session attachment.

## Build

```bash
go build -o tmux-launcher .
```

Or use the Makefile:

```bash
make build      # build binary
make install    # build and copy to ~/bin/
make clean      # remove binary
```

## Project Structure

- `main.go` — Entry point, decision logic, nesting guard, `execTmux` helper using `syscall.Exec`
- `config.go` — Config struct, layered loading via koanf (defaults → YAML file → CLI flags)
- `session.go` — `Session` struct, tmux querying via `list-sessions -F`, output parsing
- `tui.go` — Bubble Tea TUI model using `bubbles/list` for session picker

## Configuration

Config file: `~/.config/tmux-launcher/config.yaml` (optional)

CLI flags `--allow-nested` and `--no-auto-attach` override config file values.

## Dependencies

- [bubbletea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [bubbles](https://github.com/charmbracelet/bubbles) — List component
- [lipgloss](https://github.com/charmbracelet/lipgloss) — Styling
- [koanf](https://github.com/knadh/koanf) — Config loading (defaults, YAML file, flags)
- [pflag](https://github.com/spf13/pflag) — POSIX-style CLI flag parsing

## Conventions

- Flat package structure (all files in `package main`)
- tmux binary discovered via `exec.LookPath("tmux")`, not hardcoded
- `syscall.Exec` replaces the process when launching tmux (no child process)
