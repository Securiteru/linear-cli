package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var (
	optQuiet   bool
	optJSON    bool
	optCompact bool
	optFormat  string
	optFields  string
)

func IsCompact() bool { return optCompact }

func effectiveFormat() string {
	if optJSON {
		return "json"
	}
	switch optFormat {
	case "json", "tsv", "id-only":
		return optFormat
	}
	return ""
}

func writeJSON(v any) error {
	return json.NewEncoder(os.Stdout).Encode(v)
}

func toMap(v any) map[string]any {
	b, _ := json.Marshal(v)
	var m map[string]any
	json.Unmarshal(b, &m)
	return m
}

func getField(m map[string]any, path string) any {
	parts := strings.Split(path, ".")
	var cur any = m
	for _, p := range parts {
		cm, ok := cur.(map[string]any)
		if !ok {
			return nil
		}
		cur = cm[p]
		if cur == nil {
			return nil
		}
	}
	return cur
}

func fieldStr(v any) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%g", val)
	case bool:
		return fmt.Sprintf("%t", val)
	case []any:
		names := make([]string, 0, len(val))
		for _, item := range val {
			if obj, ok := item.(map[string]any); ok {
				if name, ok := obj["name"].(string); ok {
					names = append(names, name)
					continue
				}
			}
			names = append(names, fieldStr(item))
		}
		return strings.Join(names, ",")
	default:
		b, _ := json.Marshal(val)
		return strings.Trim(string(b), "\"")
	}
}

func tsvPrint(fields ...string) {
	fmt.Println(strings.Join(fields, "\t"))
}

func parseFields(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

func outputListItems(items []any, idFn func(any) string, defaultFields []string, defaultFn func()) error {
	f := effectiveFormat()
	if f == "json" {
		return writeJSON(items)
	}
	if f == "id-only" {
		for _, item := range items {
			fmt.Println(idFn(item))
		}
		return nil
	}
	if f == "tsv" || optFields != "" {
		fields := parseFields(optFields)
		if len(fields) == 0 {
			fields = defaultFields
		}
		tsvPrint(fields...)
		for _, item := range items {
			m := toMap(item)
			row := make([]string, len(fields))
			for i, fi := range fields {
				row[i] = fieldStr(getField(m, fi))
			}
			tsvPrint(row...)
		}
		return nil
	}
	if optQuiet {
		for _, item := range items {
			fmt.Println(idFn(item))
		}
		return nil
	}
	defaultFn()
	return nil
}
