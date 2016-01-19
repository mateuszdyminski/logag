package model

// Filter holds info about real-time filter rule.
type Filter struct {
	ID       string   `json:"id,omitempty"`
	Keywords []string `json:"keywords,omitempty"`
	Level    string   `json:"level,omitempty"`
}
