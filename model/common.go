package model

/* SearchOutput represents the presentable data required to display search output section
HasNoResults is a bool which displays messaging if there are no search results
SearchResults is an array of search results
Selections is an array of previously added selections
Language is the user set language */
type SearchOutput struct {
	HasNoResults    bool                `json:"has_no_results"`
	SearchResults   []SelectableElement `json:"search_results"`
	Selections      []SelectableElement `json:"selections"`
	SelectionsTitle string              `json:"selections_title"`
	Language        string              `json:"language"`
}

/* SelectableElement represents the data required for a selectable element.
Text is the human readable label.
InnerText is human readable inner text within the element.
Value is the value sent to the server.
Name is the name attribute.
IsSelected is a boolean representing whether the element is selected.
IsDisabled is a boolean representing whether the element is disabled */
type SelectableElement struct {
	Text       string `json:"text"`
	InnerText  string `json:"inner_text"`
	Value      string `json:"value"`
	Name       string `json:"name"`
	IsSelected bool   `json:"is_selected"`
	IsDisabled bool   `json:"is_disabled"`
}

// SearchField represents the data required to populate the search input partial
type SearchField struct {
	Value    string `json:"value"`
	Name     string `json:"name"`
	ID       string `json:"id"`
	Language string `json:"language"`
}
