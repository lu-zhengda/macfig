package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func TestPrintJSON(t *testing.T) {
	// Capture stdout.
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	data := map[string]string{"key": "value"}
	err := printJSON(data)

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("printJSON returned error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)

	var result map[string]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if result["key"] != "value" {
		t.Errorf("expected key=value, got key=%s", result["key"])
	}
}

func TestPrintJSONIndented(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	data := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{"Alice", 30}

	err := printJSON(data)

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("printJSON returned error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)

	output := buf.String()

	// Verify it's indented (has newline + spaces).
	if !bytes.Contains([]byte(output), []byte("\n  ")) {
		t.Error("expected indented output")
	}
}
