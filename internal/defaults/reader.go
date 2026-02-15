package defaults

import (
	"fmt"
	"strconv"
	"strings"
)

// ReadBool reads a boolean value from defaults.
func ReadBool(exec Executor, domain, key string) (bool, error) {
	val, err := exec.Read(domain, key)
	if err != nil {
		return false, err
	}
	return parseBool(val)
}

// ReadInt reads an integer value from defaults.
func ReadInt(exec Executor, domain, key string) (int, error) {
	val, err := exec.Read(domain, key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(val))
}

// ReadFloat reads a float value from defaults.
func ReadFloat(exec Executor, domain, key string) (float64, error) {
	val, err := exec.Read(domain, key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(strings.TrimSpace(val), 64)
}

// ReadString reads a string value from defaults.
func ReadString(exec Executor, domain, key string) (string, error) {
	return exec.Read(domain, key)
}

// ReadTyped reads a value and returns it as the appropriate Go type based on ValueType.
func ReadTyped(exec Executor, domain, key string, vt ValueType) (interface{}, error) {
	switch vt {
	case TypeBool:
		return ReadBool(exec, domain, key)
	case TypeInt:
		return ReadInt(exec, domain, key)
	case TypeFloat:
		return ReadFloat(exec, domain, key)
	case TypeString:
		return ReadString(exec, domain, key)
	default:
		return ReadString(exec, domain, key)
	}
}

// parseBool handles macOS defaults boolean representations:
// 1/0, true/false, yes/no.
func parseBool(s string) (bool, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	switch s {
	case "1", "true", "yes":
		return true, nil
	case "0", "false", "no":
		return false, nil
	default:
		return false, fmt.Errorf("cannot parse %q as bool", s)
	}
}
