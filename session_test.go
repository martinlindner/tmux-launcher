package main

import (
	"testing"
	"time"
)

func TestParseSessions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Session
	}{
		{
			name:     "empty output",
			input:    "",
			expected: nil,
		},
		{
			name:  "single detached session",
			input: "dev|3|1738950234|0\n",
			expected: []Session{
				{Name: "dev", Windows: 3, Created: time.Unix(1738950234, 0), Attached: false},
			},
		},
		{
			name:  "single attached session",
			input: "work|5|1738949021|1\n",
			expected: []Session{
				{Name: "work", Windows: 5, Created: time.Unix(1738949021, 0), Attached: true},
			},
		},
		{
			name:  "multiple sessions",
			input: "dev|3|1738950234|0\nwork|5|1738949021|1\n",
			expected: []Session{
				{Name: "dev", Windows: 3, Created: time.Unix(1738950234, 0), Attached: false},
				{Name: "work", Windows: 5, Created: time.Unix(1738949021, 0), Attached: true},
			},
		},
		{
			name:     "malformed line skipped",
			input:    "bad-line\ndev|3|1738950234|0\n",
			expected: []Session{
				{Name: "dev", Windows: 3, Created: time.Unix(1738950234, 0), Attached: false},
			},
		},
		{
			name:     "blank lines skipped",
			input:    "\ndev|3|1738950234|0\n\n",
			expected: []Session{
				{Name: "dev", Windows: 3, Created: time.Unix(1738950234, 0), Attached: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessions, err := parseSessions(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(sessions) != len(tt.expected) {
				t.Fatalf("got %d sessions, want %d", len(sessions), len(tt.expected))
			}
			for i, s := range sessions {
				exp := tt.expected[i]
				if s.Name != exp.Name {
					t.Errorf("session[%d].Name = %q, want %q", i, s.Name, exp.Name)
				}
				if s.Windows != exp.Windows {
					t.Errorf("session[%d].Windows = %d, want %d", i, s.Windows, exp.Windows)
				}
				if !s.Created.Equal(exp.Created) {
					t.Errorf("session[%d].Created = %v, want %v", i, s.Created, exp.Created)
				}
				if s.Attached != exp.Attached {
					t.Errorf("session[%d].Attached = %v, want %v", i, s.Attached, exp.Attached)
				}
			}
		})
	}
}

func TestSessionTitle(t *testing.T) {
	s := Session{Name: "dev", Attached: false}
	if got := s.Title(); got != "dev (detached)" {
		t.Errorf("Title() = %q, want %q", got, "dev (detached)")
	}

	s.Attached = true
	if got := s.Title(); got != "dev (attached)" {
		t.Errorf("Title() = %q, want %q", got, "dev (attached)")
	}
}

func TestSessionDescription(t *testing.T) {
	ts := time.Date(2025, 2, 7, 14, 30, 0, 0, time.UTC)

	s := Session{Windows: 1, Created: ts}
	if got := s.Description(); got != "1 window · created Feb 7 14:30" {
		t.Errorf("Description() = %q, want %q", got, "1 window · created Feb 7 14:30")
	}

	s.Windows = 3
	if got := s.Description(); got != "3 windows · created Feb 7 14:30" {
		t.Errorf("Description() = %q, want %q", got, "3 windows · created Feb 7 14:30")
	}
}

func TestSessionFilterValue(t *testing.T) {
	s := Session{Name: "mySession"}
	if got := s.FilterValue(); got != "mySession" {
		t.Errorf("FilterValue() = %q, want %q", got, "mySession")
	}
}
