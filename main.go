package main

import (
	"fmt"
	"os"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "tmux-launcher: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	sessions, err := getSessions()
	if err != nil {
		return err
	}

	// No sessions — create new
	if len(sessions) == 0 {
		return execTmux("new-session")
	}

	// Single detached session — attach
	if len(sessions) == 1 && !sessions[0].Attached {
		return execTmux("attach-session", "-t", sessions[0].Name)
	}

	// Multiple sessions or attached — show picker
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
