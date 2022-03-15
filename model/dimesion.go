package model

// Dimension represents the data for a single dimension
type Dimension struct {
	Options      []string `json:"options"`
	IsTruncated  bool     `json:"is_truncated"`
	TruncateLink string   `json:"truncate_link"`
	OptionsCount int      `json:"options_count"`
	Name         string   `json:"name"`
	EncodedName  string   `json:"encoded_name"`
	URI          string   `json:"uri"`
	IsAreaType   bool     `json:"is_area_type"`
}
