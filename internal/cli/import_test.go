package cli

import (
	"testing"

	"github.com/lu-zhengda/macfig/internal/defaults"
)

func TestParseImportValue(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		vt      defaults.ValueType
		want    interface{}
		wantErr bool
	}{
		{"bool true", "true", defaults.TypeBool, true, false},
		{"bool false", "false", defaults.TypeBool, false, false},
		{"bool 1", "1", defaults.TypeBool, true, false},
		{"bool 0", "0", defaults.TypeBool, false, false},
		{"bool yes", "yes", defaults.TypeBool, true, false},
		{"bool no", "no", defaults.TypeBool, false, false},
		{"bool invalid", "maybe", defaults.TypeBool, false, true},
		{"int valid", "42", defaults.TypeInt, 42, false},
		{"int invalid", "abc", defaults.TypeInt, 0, true},
		{"float valid", "3.14", defaults.TypeFloat, 3.14, false},
		{"float invalid", "abc", defaults.TypeFloat, 0.0, true},
		{"string", "hello", defaults.TypeString, "hello", false},
		{"empty string", "", defaults.TypeString, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseImportValue(tt.raw, tt.vt)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseImportValue(%q, %q) error = %v, wantErr %v",
					tt.raw, tt.vt, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parseImportValue(%q, %q) = %v, want %v",
					tt.raw, tt.vt, got, tt.want)
			}
		})
	}
}

func TestParseBoolImport(t *testing.T) {
	tests := []struct {
		input   string
		want    bool
		wantErr bool
	}{
		{"1", true, false},
		{"0", false, false},
		{"true", true, false},
		{"false", false, false},
		{"yes", true, false},
		{"no", false, false},
		{"TRUE", true, false},
		{"FALSE", false, false},
		{"invalid", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseBoolImport(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseBoolImport(%q) error = %v, wantErr %v",
					tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parseBoolImport(%q) = %v, want %v",
					tt.input, got, tt.want)
			}
		})
	}
}
