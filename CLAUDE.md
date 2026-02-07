# tmux-launcher

A smart tmux wrapper/entrypoint for WSL2 written in Go.

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

- `main.go` — Entry point, decision logic (no sessions → new, single detached → attach, otherwise → TUI), `execTmux` helper using `syscall.Exec`
- `session.go` — `Session` struct, tmux querying via `list-sessions -F`, output parsing
- `tui.go` — Bubble Tea TUI model using `bubbles/list` for session picker

## Dependencies

- [bubbletea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [bubbles](https://github.com/charmbracelet/bubbles) — List component
- [lipgloss](https://github.com/charmbracelet/lipgloss) — Styling

## Conventions

- Flat package structure (all files in `package main`)
- tmux binary discovered via `exec.LookPath("tmux")`, not hardcoded
- `syscall.Exec` replaces the process when launching tmux (no child process)
