package defaults

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// CmdRunner abstracts shell command execution for testability.
type CmdRunner interface {
	Run(ctx context.Context, name string, args ...string) ([]byte, error)
}

// RealCmdRunner executes real shell commands.
type RealCmdRunner struct{}

// Run executes a command and returns its stdout.
func (r *RealCmdRunner) Run(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stderr = io.Discard
	return cmd.Output()
}

// Executor defines the interface for interacting with macOS defaults.
type Executor interface {
	Read(domain, key string) (string, error)
	ReadAll(domain string) (map[string]string, error)
	Write(domain, key string, value interface{}, valueType ValueType) error
	Delete(domain, key string) error
	Domains() ([]string, error)
	Find(keyword string) ([]SearchResult, error)
}

// RealExecutor shells out to the `defaults` command.
type RealExecutor struct {
	runner CmdRunner
}

// NewExecutor creates a new RealExecutor with the given command runner.
func NewExecutor(runner CmdRunner) *RealExecutor {
	return &RealExecutor{runner: runner}
}

// Read reads a single value from macOS defaults.
func (e *RealExecutor) Read(domain, key string) (string, error) {
	out, err := e.runner.Run(context.Background(), "defaults", "read", domain, key)
	if err != nil {
		return "", fmt.Errorf("failed to read %s %s: %w", domain, key, err)
	}
	return strings.TrimSpace(string(out)), nil
}

// ReadAll reads all keys from a domain.
func (e *RealExecutor) ReadAll(domain string) (map[string]string, error) {
	out, err := e.runner.Run(context.Background(), "defaults", "read", domain)
	if err != nil {
		return nil, fmt.Errorf("failed to read all from %s: %w", domain, err)
	}

	result := make(map[string]string)
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.TrimSpace(parts[1])
			v = strings.TrimSuffix(v, ";")
			v = strings.TrimSpace(v)
			if k != "" {
				result[k] = v
			}
		}
	}
	return result, nil
}

// Write writes a value to macOS defaults.
func (e *RealExecutor) Write(domain, key string, value interface{}, valueType ValueType) error {
	typeFlag := typeToFlag(valueType)
	valStr := fmt.Sprintf("%v", value)

	_, err := e.runner.Run(context.Background(), "defaults", "write", domain, key, typeFlag, valStr)
	if err != nil {
		return fmt.Errorf("failed to write %s %s: %w", domain, key, err)
	}
	return nil
}

// Delete removes a key from macOS defaults, resetting it to the system default.
func (e *RealExecutor) Delete(domain, key string) error {
	_, err := e.runner.Run(context.Background(), "defaults", "delete", domain, key)
	if err != nil {
		return fmt.Errorf("failed to delete %s %s: %w", domain, key, err)
	}
	return nil
}

// Domains lists all known defaults domains.
func (e *RealExecutor) Domains() ([]string, error) {
	out, err := e.runner.Run(context.Background(), "defaults", "domains")
	if err != nil {
		return nil, fmt.Errorf("failed to list domains: %w", err)
	}

	raw := strings.TrimSpace(string(out))
	parts := strings.Split(raw, ", ")
	domains := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			domains = append(domains, p)
		}
	}
	return domains, nil
}

// Find searches for a keyword across all defaults domains.
func (e *RealExecutor) Find(keyword string) ([]SearchResult, error) {
	out, err := e.runner.Run(context.Background(), "defaults", "find", keyword)
	if err != nil {
		return nil, fmt.Errorf("failed to find %q: %w", keyword, err)
	}

	var results []SearchResult
	var currentDomain string
	lines := strings.Split(string(out), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Found ") && strings.Contains(line, " in domain ") {
			// Extract domain from "Found X keys in domain 'com.apple.dock':"
			parts := strings.SplitN(line, "'", 3)
			if len(parts) >= 2 {
				currentDomain = parts[1]
			}
		} else if currentDomain != "" && strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				k := strings.TrimSpace(parts[0])
				v := strings.TrimSpace(parts[1])
				v = strings.TrimSuffix(v, ";")
				v = strings.TrimSpace(v)
				if k != "" {
					results = append(results, SearchResult{
						Domain: currentDomain,
						Key:    k,
						Value:  v,
					})
				}
			}
		}
	}
	return results, nil
}

// typeToFlag converts a ValueType to the `-type` flag for the defaults command.
func typeToFlag(t ValueType) string {
	switch t {
	case TypeBool:
		return "-bool"
	case TypeInt:
		return "-int"
	case TypeFloat:
		return "-float"
	case TypeString:
		return "-string"
	default:
		return "-string"
	}
}
