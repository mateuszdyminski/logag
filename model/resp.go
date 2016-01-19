package model

// Response holds information about queried data and total number of hits.
type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Total int64       `json:"total,omitempty"`
}
