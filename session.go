package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Session struct {
	Name     string
	Windows  int
	Created  time.Time
	Attached bool
}

func (s Session) FilterValue() string { return s.Name }

func (s Session) Title() string {
	status := "detached"
	if s.Attached {
		status = "attached"
	}
	return fmt.Sprintf("%s (%s)", s.Name, status)
}

func (s Session) Description() string {
	word := "window"
	if s.Windows != 1 {
		word = "windows"
	}
	return fmt.Sprintf("%d %s · created %s", s.Windows, word, s.Created.Format("Jan 2 15:04"))
}

func getSessions() ([]Session, error) {
	tmuxBin, err := findTmux()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(tmuxBin,
		"list-sessions", "-F",
		"#{session_name}|#{session_windows}|#{session_created}|#{session_attached}")

	output, err := cmd.CombinedOutput()
	if err != nil {
		out := string(output)
		if strings.Contains(out, "no server running") ||
			strings.Contains(out, "no sessions") ||
			strings.Contains(out, "error connecting") {
			return nil, nil
		}
		return nil, fmt.Errorf("tmux list-sessions: %s", strings.TrimSpace(out))
	}

	return parseSessions(string(output))
}

func parseSessions(output string) ([]Session, error) {
	var sessions []Session
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 4)
		if len(parts) != 4 {
			continue
		}

		windows, _ := strconv.Atoi(parts[1])
		ts, _ := strconv.ParseInt(parts[2], 10, 64)
		attached := parts[3] == "1"

		sessions = append(sessions, Session{
			Name:     parts[0],
			Windows:  windows,
			Created:  time.Unix(ts, 0),
			Attached: attached,
		})
	}
	return sessions, nil
}

func findTmux() (string, error) {
	path, err := exec.LookPath("tmux")
	if err != nil {
		return "", fmt.Errorf("tmux not found in PATH: %w", err)
	}
	return path, nil
}
