package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type action int

const (
	actionNone action = iota
	actionAttach
	actionNew
	actionQuit
)

type model struct {
	list     list.Model
	action   action
	selected *Session
}

type newSessionItem struct{}

func (n newSessionItem) FilterValue() string { return "new" }
func (n newSessionItem) Title() string       { return "+ Create new session" }
func (n newSessionItem) Description() string { return "Start a fresh tmux session" }

func newTUI(sessions []Session) model {
	items := make([]list.Item, len(sessions))
	for i, s := range sessions {
		items[i] = s
	}
	items = append(items, newSessionItem{})

	accent := lipgloss.Color("#7D56F4")

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(accent).
		BorderForeground(accent)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(accent).
		BorderForeground(accent)

	l := list.New(items, delegate, 0, 0)
	l.Title = "tmux sessions"
	l.Styles.Title = lipgloss.NewStyle().
		Background(accent).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1).
		Bold(true)
	l.SetFilteringEnabled(true)
	l.SetShowStatusBar(true)

	return model{list: l}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		// Don't intercept keys while filtering
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {
		case "q", "ctrl+c":
			m.action = actionQuit
			return m, tea.Quit

		case "enter":
			switch item := m.list.SelectedItem().(type) {
			case Session:
				m.selected = &item
				m.action = actionAttach
			case newSessionItem:
				m.action = actionNew
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.action == actionQuit {
		return ""
	}
	return m.list.View()
}
