package interactive

import tea "github.com/charmbracelet/bubbletea"

type option struct {
	optType
	text    string
	next    chapter
	handler func() tea.Cmd
}

type optType int

const (
	OptTypeUnknown = iota

	OptTypeUnselectableString = iota
	OptTypeNextChapter        = iota
	OptTypeHandler            = iota
	OptTypeSelectabe          = iota
	OptTypeInputWindow        = iota
	OptTypeFinish             = iota
	// currently available one per chapter
	OptTypeOneOf = iota
)
