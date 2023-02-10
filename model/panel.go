package model

type PanelType int

const (
	Info PanelType = iota
	Pending
	Success
	Error
)

// FuncGetPanelType returns the panel type as a string
func (p Panel) FuncGetPanelType() (panelType string) {
	switch p.Type {
	case Info:
		return "info"
	case Pending:
		return "pending"
	case Success:
		return "success"
	case Error:
		return "error"
	}
	return panelType
}

// Panel contains the data required to populate a panel UI component
type Panel struct {
	Type       PanelType `json:"type"`
	CssClasses []string  `json:"css_classes"`
	Body       string    `json:"body"`
	SafeHTML   []string  `json:"safe_html"`
	Language   string    `json:"language"`
}
