package interactive

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func handleOptionListInput(m *model, msg tea.Msg) (tea.Model, tea.Cmd) {
	// NPE protection
	if m == nil {
		fmt.Println("[Unexpected Exception] handleOptionListInput is called with nil model")
		return m, tea.Quit
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			m.cursor.up()
		case "down", "j":
			length := 0
			if opts, ok := m.options[m.chapter]; ok {
				length = len(opts) - 1
			}
			m.cursor.down(length)
		case "enter", " ":
			if opts, ok := m.options[m.chapter]; ok {
				var cmd tea.Cmd
				if opts[m.cursor].handler != nil {
					m, cmd = opts[m.cursor].handler(m)
					if cmd != nil {
						return m, cmd
					}
				}
				if !opts[m.cursor].next.isNil() { // if has next
					// TODO: move this trash to the chapter functionality
					m.chapter.change(opts[m.cursor].next)
					// this is 50% bullshit
					m.chapterType = calculateOpts(m.options[m.chapter])
					m.cursor.setStart(m.options[m.chapter])
				}

				// Final chapter quits immediately
				if m.chapter.isFinal() {
					return m, tea.Quit
				}
				return m, cmd
			}
		}
	}
	return m, nil
}

func calculateOpts(opts []option) optType {
	resultType := OptTypeUnknown
	for _, o := range opts {
		switch o.optType {
		case OptTypeUnselectableString, OptTypeNextChapter, OptTypeOneOf:
			if resultType == OptTypeUnselectableString ||
				resultType == OptTypeNextChapter ||
				resultType == OptTypeOneOf ||
				resultType == OptTypeUnknown {
				resultType = OptTypeNextChapter
				continue
			}
			return OptTypeUnknown
		case OptTypeInputWindow:
			if resultType == OptTypeInputWindow ||
				resultType == OptTypeUnknown {
				resultType = OptTypeInputWindow
				break
			}
			return OptTypeUnknown
		}
	}
	// other types are unsupported yet
	return optType(resultType)
}

func handleOptionInput(m *model, _ tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: implement
	return m, nil
}
