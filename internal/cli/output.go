package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var jsonFlag bool

func printJSON(v any) error {
	return fprintJSON(os.Stdout, v)
}

func fprintJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	return nil
}
