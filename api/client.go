package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const endpoint = "https://api.linear.app/graphql"

var client = &http.Client{Timeout: 15 * time.Second}

type graphResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []graphError    `json:"errors"`
}

type graphError struct {
	Message string `json:"message"`
	Path    []any  `json:"path"`
}

func Query(query string, result any) error {
	key := os.Getenv("LINEAR_API_KEY")
	if key == "" {
		return fmt.Errorf("LINEAR_API_KEY not set")
	}

	body := map[string]string{"query": query}
	buf, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", key)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	var gr graphResponse
	if err := json.Unmarshal(data, &gr); err != nil {
		return fmt.Errorf("parse response: %w", err)
	}

	if len(gr.Errors) > 0 {
		return fmt.Errorf("GraphQL error: %s", gr.Errors[0].Message)
	}

	if gr.Data == nil {
		return fmt.Errorf("no data in response")
	}

	return json.Unmarshal(gr.Data, result)
}
