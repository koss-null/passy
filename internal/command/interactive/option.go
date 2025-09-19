package interactive

import tea "github.com/charmbracelet/bubbletea"

// option is a common class that contains data for all possible option types
// this may change in future if the class will become too huge, but currently
// consider embed another optType-specific types in here.
type option struct {
	optType

	text       string  // what is printed
	next       chapter // skipped if ""`
	handler    func(*model) (*model, tea.Cmd)
	isSelected bool
}

type optType int

const (
	OptTypeUnknown = iota

	OptTypeUnselectableString = iota
	OptTypeNextChapter        = iota
	OptTypeHandler            = iota
	OptTypeSelectabe          = iota
	// If OptTypeInputWindow is in a chapter, there
	// should be only objects of that type.
	OptTypeInputWindow = iota
	// Currently available one per chapter
	OptTypeOneOf = iota
	// This one should not actually occur, the app about to exit
	// before this value of the chapter will be read.
	OptTypeFinish = iota
)
