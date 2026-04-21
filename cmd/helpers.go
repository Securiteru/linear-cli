package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/securiter/linear-cli/api"
)

var linearURLRe = regexp.MustCompile(`linear\.app/[^/]+/issue/([A-Za-z]+-\d+)`)

func parseIssueIdentifier(input string) string {
	input = strings.TrimSpace(input)
	if m := linearURLRe.FindStringSubmatch(input); len(m) > 1 {
		return m[1]
	}
	return input
}

func getViewerID() (string, error) {
	q := `query { viewer { id } }`
	var res struct {
		Viewer struct {
			ID string `json:"id"`
		} `json:"viewer"`
	}
	if err := api.Query(q, &res); err != nil {
		return "", err
	}
	return res.Viewer.ID, nil
}

func getViewerName() (string, error) {
	q := `query { viewer { name } }`
	var res struct {
		Viewer struct {
			Name string `json:"name"`
		} `json:"viewer"`
	}
	if err := api.Query(q, &res); err != nil {
		return "", err
	}
	return res.Viewer.Name, nil
}

func fetchTeamStates(teamID string) ([]struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}, error) {
	q := fmt.Sprintf(`query { workflowStates(filter: { team: { id: { eq: "%s" } } }) { nodes { id name } } }`, teamID)
	var res struct {
		WorkflowStates struct {
			Nodes []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"nodes"`
		} `json:"workflowStates"`
	}
	if err := api.Query(q, &res); err != nil {
		return nil, err
	}
	return res.WorkflowStates.Nodes, nil
}

func fuzzyMatchState(states []struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}, input string) (string, string, error) {
	lower := strings.ToLower(input)
	var exactMatch struct {
		ID   string
		Name string
	}
	var partialMatches []struct {
		ID   string
		Name string
	}

	for _, s := range states {
		sl := strings.ToLower(s.Name)
		if sl == lower {
			exactMatch.ID = s.ID
			exactMatch.Name = s.Name
		}
		if strings.Contains(sl, lower) {
			partialMatches = append(partialMatches, struct {
				ID   string
				Name string
			}{s.ID, s.Name})
		}
	}

	if exactMatch.ID != "" {
		return exactMatch.ID, exactMatch.Name, nil
	}

	if len(partialMatches) == 1 {
		return partialMatches[0].ID, partialMatches[0].Name, nil
	}

	if len(partialMatches) > 1 {
		names := make([]string, len(partialMatches))
		for i, m := range partialMatches {
			names[i] = m.Name
		}
		return "", "", fmt.Errorf("status %q matched multiple states: %s", input, strings.Join(names, ", "))
	}

	names := make([]string, len(states))
	for i, s := range states {
		names[i] = s.Name
	}
	return "", "", fmt.Errorf("status %q not found. Available: %s", input, strings.Join(names, ", "))
}

func formatAvailableStates(states []struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}) string {
	names := make([]string, len(states))
	for i, s := range states {
		names[i] = s.Name
	}
	return strings.Join(names, ", ")
}
