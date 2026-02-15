package defaults

// ValueType represents the type of a defaults value.
type ValueType string

const (
	TypeBool   ValueType = "bool"
	TypeInt    ValueType = "int"
	TypeFloat  ValueType = "float"
	TypeString ValueType = "string"
)

// SearchResult represents a single match from `defaults find`.
type SearchResult struct {
	Domain string `json:"domain"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}
