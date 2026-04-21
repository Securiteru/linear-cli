package cmd

import (
	"fmt"
	"testing"
)

func TestPriorityLabel(t *testing.T) {
	tests := []struct {
		p    int
		want string
	}{
		{0, "-"},
		{1, "Urgent"},
		{2, "High"},
		{3, "Medium"},
		{4, "Low"},
		{5, "5"},
		{-1, "-1"},
	}
	for _, tt := range tests {
		got := priorityLabel(tt.p)
		if got != tt.want {
			t.Errorf("priorityLabel(%d) = %q, want %q", tt.p, got, tt.want)
		}
	}
}

func TestEscapeGraphQL(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "hello"},
		{`say "hi"`, `say \"hi\"`},
		{"line1\nline2", `line1\nline2`},
		{"tab\there", `tab\there`},
		{"back\\slash", `back\\slash`},
		{"win\rend", `win\rend`},
	}
	for _, tt := range tests {
		got := escapeGraphQL(tt.input)
		if got != tt.want {
			t.Errorf("escapeGraphQL(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestParseFields(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"", nil},
		{"a", []string{"a"}},
		{"a,b,c", []string{"a", "b", "c"}},
	}
	for _, tt := range tests {
		got := parseFields(tt.input)
		if len(got) != len(tt.want) {
			t.Errorf("parseFields(%q) len = %d, want %d", tt.input, len(got), len(tt.want))
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("parseFields(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
			}
		}
	}
}

func TestEffectiveFormat(t *testing.T) {
	tests := []struct {
		name      string
		json      bool
		format    string
		want      string
	}{
		{"json flag", true, "", "json"},
		{"format json", false, "json", "json"},
		{"format tsv", false, "tsv", "tsv"},
		{"format id-only", false, "id-only", "id-only"},
		{"default", false, "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			optJSON = tt.json
			optFormat = tt.format
			got := effectiveFormat()
			if got != tt.want {
				t.Errorf("effectiveFormat() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFieldStr(t *testing.T) {
	tests := []struct {
		name string
		v    any
		want string
	}{
		{"nil", nil, ""},
		{"string", "hello", "hello"},
		{"int float64", float64(42), "42"},
		{"float float64", float64(3.14), "3.14"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"string array", []any{"a", "b"}, "a,b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fieldStr(tt.v)
			if got != tt.want {
				t.Errorf("fieldStr(%v) = %q, want %q", tt.v, got, tt.want)
			}
		})
	}
}

func TestToMap(t *testing.T) {
	type inner struct {
		Name string `json:"name"`
	}
	type outer struct {
		ID    string `json:"id"`
		Label inner  `json:"label"`
	}
	m := toMap(outer{ID: "123", Label: inner{Name: "test"}})
	if m["id"] != "123" {
		t.Errorf("toMap id = %v, want 123", m["id"])
	}
	label, ok := m["label"].(map[string]any)
	if !ok {
		t.Fatal("toMap label is not map")
	}
	if label["name"] != "test" {
		t.Errorf("toMap label.name = %v, want test", label["name"])
	}
}

func TestGetField(t *testing.T) {
	m := map[string]any{
		"id": "ADI-1",
		"state": map[string]any{
			"name": "Todo",
		},
	}
	tests := []struct {
		path string
		want string
	}{
		{"id", "ADI-1"},
		{"state.name", "Todo"},
		{"missing", ""},
		{"state.missing", ""},
	}
	for _, tt := range tests {
		got := getField(m, tt.path)
		var gotStr string
		if got != nil {
			gotStr = fmt.Sprintf("%v", got)
		}
		if gotStr != tt.want {
			t.Errorf("getField(%q) = %v, want %q", tt.path, got, tt.want)
		}
	}
}

func TestToAnySlice(t *testing.T) {
	s := []int{1, 2, 3}
	a := toAnySlice(s)
	if len(a) != 3 {
		t.Fatalf("toAnySlice len = %d, want 3", len(a))
	}
	if a[0].(int) != 1 {
		t.Errorf("toAnySlice[0] = %v, want 1", a[0])
	}
}
