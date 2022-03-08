package model

// Dimension represents the data for a single dimension
type Dimension struct {
	Options          []string `json:"options"`
	TruncatedOptions []Option `json:"truncated_options"`
	TruncateLink     string   `json:"truncate_link"`
	OptionsCount     int      `json:"options_count"`
	Name             string   `json:"name"`
	EncodedName      string   `json:"encoded_name"`
	URI              string   `json:"uri"`
	IsAreaType       bool     `json:"is_area_type"`
}

// Option represents the data for a single option
type Option struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}
