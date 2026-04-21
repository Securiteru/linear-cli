package cmd

import (
	"testing"
)

func TestParseIssueIdentifier(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"ADI-36", "ADI-36"},
		{"  ADI-36  ", "ADI-36"},
		{"https://linear.app/myteam/issue/ADI-36/test-title", "ADI-36"},
		{"http://linear.app/myteam/issue/FOO-1", "FOO-1"},
		{"https://linear.app/team/issue/BUG-123/some-slug?comment=1", "BUG-123"},
		{"not-a-url", "not-a-url"},
		{"", ""},
	}
	for _, tt := range tests {
		got := parseIssueIdentifier(tt.input)
		if got != tt.want {
			t.Errorf("parseIssueIdentifier(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

type stateEntry = struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func TestFuzzyMatchState(t *testing.T) {
	states := []stateEntry{
		{"s1", "Todo"},
		{"s2", "In Progress"},
		{"s3", "In Review"},
		{"s4", "Done"},
		{"s5", "Cancelled"},
	}

	tests := []struct {
		name    string
		input   string
		wantID  string
		wantErr bool
	}{
		{"exact match", "Todo", "s1", false},
		{"exact case insensitive", "todo", "s1", false},
		{"exact case insensitive uppercase", "DONE", "s4", false},
		{"partial unique", "prog", "s2", false},
		{"partial unique rev", "review", "s3", false},
		{"partial ambiguous", "In", "", true},
		{"not found", "Deployed", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, gotName, err := fuzzyMatchState(states, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("fuzzyMatchState(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && gotID != tt.wantID {
				t.Errorf("fuzzyMatchState(%q) id = %q, want %q", tt.input, gotID, tt.wantID)
			}
			if !tt.wantErr && gotName == "" {
				t.Errorf("fuzzyMatchState(%q) name should not be empty", tt.input)
			}
		})
	}
}

func TestFuzzyMatchState_Ambiguous(t *testing.T) {
	states := []stateEntry{
		{"s1", "In Progress"},
		{"s2", "In Review"},
	}
	_, _, err := fuzzyMatchState(states, "in")
	if err == nil {
		t.Fatal("expected error for ambiguous partial match")
	}
	t.Logf("ambiguous error: %v", err)
}

func TestFormatAvailableStates(t *testing.T) {
	states := []stateEntry{
		{"s1", "Todo"},
		{"s2", "Done"},
	}
	got := formatAvailableStates(states)
	if got != "Todo, Done" {
		t.Errorf("formatAvailableStates() = %q, want %q", got, "Todo, Done")
	}
}
