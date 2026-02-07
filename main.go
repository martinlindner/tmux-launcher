package main

import (
	"fmt"
	"os"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
)

var version = "dev"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "tmux-launcher: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	// Nesting guard
	if _, nested := os.LookupEnv("TMUX"); nested && !cfg.AllowNested {
		return fmt.Errorf("already inside a tmux session (use --allow-nested to override)")
	}

	sessions, err := getSessions()
	if err != nil {
		return err
	}

	// Auto-create a new session when none exist
	if len(sessions) == 0 && cfg.AutoNewSession {
		return execTmux("new-session")
	}

	// Auto-attach to a single detached session
	if cfg.AutoAttach {
		var detached []Session
		for _, s := range sessions {
			if !s.Attached {
				detached = append(detached, s)
			}
		}
		if len(detached) == 1 {
			return execTmux("attach-session", "-t", detached[0].Name)
		}
	}

	// Show TUI picker
	p := tea.NewProgram(newTUI(sessions), tea.WithAltScreen())
	result, err := p.Run()
	if err != nil {
		return fmt.Errorf("TUI: %w", err)
	}

	m := result.(model)
	switch m.action {
	case actionAttach:
		if m.selected != nil {
			return execTmux("attach-session", "-t", m.selected.Name)
		}
	case actionNew:
		return execTmux("new-session")
	}

	return nil
}

func execTmux(args ...string) error {
	tmuxBin, err := findTmux()
	if err != nil {
		return err
	}

	argv := append([]string{"tmux"}, args...)
	return syscall.Exec(tmuxBin, argv, os.Environ())
}
