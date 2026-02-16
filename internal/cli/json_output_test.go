package cli

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestFprintJSON(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]string{"key": "value"}
	if err := fprintJSON(&buf, data); err != nil {
		t.Fatalf("fprintJSON returned error: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if result["key"] != "value" {
		t.Errorf("expected key=value, got key=%s", result["key"])
	}
}

func TestFprintJSONIndented(t *testing.T) {
	var buf bytes.Buffer
	data := struct {
		Name string `json:"name"`
	}{"Alice"}
	if err := fprintJSON(&buf, data); err != nil {
		t.Fatalf("fprintJSON returned error: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("\n  ")) {
		t.Error("expected indented output")
	}
}

func TestSetResultJSON(t *testing.T) {
	tests := []struct {
		name   string
		result setResult
		checks map[string]interface{}
	}{
		{
			name: "basic set",
			result: setResult{
				OK:     true,
				Domain: "com.apple.dock",
				Key:    "autohide",
				Value:  "true",
			},
			checks: map[string]interface{}{
				"ok":     true,
				"domain": "com.apple.dock",
				"key":    "autohide",
				"value":  "true",
			},
		},
		{
			name: "set with restart",
			result: setResult{
				OK:        true,
				Domain:    "com.apple.dock",
				Key:       "tilesize",
				Value:     "48",
				Restarted: "Dock",
			},
			checks: map[string]interface{}{
				"ok":        true,
				"restarted": "Dock",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := fprintJSON(&buf, tt.result); err != nil {
				t.Fatalf("fprintJSON error: %v", err)
			}

			var got map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
				t.Fatalf("invalid JSON: %v", err)
			}

			for k, want := range tt.checks {
				if got[k] != want {
					t.Errorf("%s: got %v, want %v", k, got[k], want)
				}
			}
		})
	}
}

func TestSetResultOmitsEmptyRestarted(t *testing.T) {
	var buf bytes.Buffer
	r := setResult{OK: true, Domain: "d", Key: "k", Value: "v"}
	if err := fprintJSON(&buf, r); err != nil {
		t.Fatalf("fprintJSON error: %v", err)
	}
	if bytes.Contains(buf.Bytes(), []byte("restarted")) {
		t.Error("expected restarted to be omitted when empty")
	}
}

func TestResetResultJSON(t *testing.T) {
	tests := []struct {
		name   string
		result resetResult
		checks map[string]interface{}
	}{
		{
			name: "basic reset",
			result: resetResult{
				OK:     true,
				Domain: "com.apple.dock",
				Key:    "autohide",
			},
			checks: map[string]interface{}{
				"ok":     true,
				"domain": "com.apple.dock",
				"key":    "autohide",
			},
		},
		{
			name: "reset with restart",
			result: resetResult{
				OK:        true,
				Domain:    "com.apple.finder",
				Key:       "ShowPathbar",
				Restarted: "Finder",
			},
			checks: map[string]interface{}{
				"ok":        true,
				"restarted": "Finder",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := fprintJSON(&buf, tt.result); err != nil {
				t.Fatalf("fprintJSON error: %v", err)
			}

			var got map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
				t.Fatalf("invalid JSON: %v", err)
			}

			for k, want := range tt.checks {
				if got[k] != want {
					t.Errorf("%s: got %v, want %v", k, got[k], want)
				}
			}
		})
	}
}

func TestResetResultOmitsEmptyRestarted(t *testing.T) {
	var buf bytes.Buffer
	r := resetResult{OK: true, Domain: "d", Key: "k"}
	if err := fprintJSON(&buf, r); err != nil {
		t.Fatalf("fprintJSON error: %v", err)
	}
	if bytes.Contains(buf.Bytes(), []byte("restarted")) {
		t.Error("expected restarted to be omitted when empty")
	}
}

func TestRestoreResultJSON(t *testing.T) {
	tests := []struct {
		name   string
		result restoreResult
		checks map[string]interface{}
	}{
		{
			name: "all restored",
			result: restoreResult{
				OK:       true,
				Path:     "/tmp/backup.json",
				Restored: 10,
				Failed:   0,
			},
			checks: map[string]interface{}{
				"ok":       true,
				"path":     "/tmp/backup.json",
				"restored": float64(10),
				"failed":   float64(0),
			},
		},
		{
			name: "partial failure",
			result: restoreResult{
				OK:       false,
				Path:     "/tmp/backup.json",
				Restored: 8,
				Failed:   2,
			},
			checks: map[string]interface{}{
				"ok":       false,
				"restored": float64(8),
				"failed":   float64(2),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := fprintJSON(&buf, tt.result); err != nil {
				t.Fatalf("fprintJSON error: %v", err)
			}

			var got map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
				t.Fatalf("invalid JSON: %v", err)
			}

			for k, want := range tt.checks {
				if got[k] != want {
					t.Errorf("%s: got %v, want %v", k, got[k], want)
				}
			}
		})
	}
}
