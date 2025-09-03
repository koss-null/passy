package interactive

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cursor  cursor
	chapter chapter
	options map[chapter][]option
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Final chapter quits immediately
	if m.chapter.isFinal() {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.chapter.isFinal()
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
					cmd = opts[m.cursor].handler()
				}
				m.chapter.change(opts[m.cursor].next)
				m.cursor.setStart()

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

func (m *model) View() string {
	// Define styles
	var sb strings.Builder

	// Title
	addStrL(&sb, title(m.chapter))
	addStrL(&sb, styles().divider.String())
	// Options
	if opts, ok := m.options[m.chapter]; ok {
		for i, opt := range opts {
			m.AddOptionStr(&sb, i, opt)
		}
	}
	// Divider
	addLStrL(&sb, styles().divider.String())
	// Help text
	helpText := styles().help.Render("↑/k: up • ↓/j: down • enter: select • q/ctrl+c: quit")
	addStrL(&sb, helpText)

	return sb.String()
}

// OptionStr adds the line(s) for the sequential option to the strings.Builder
func (m *model) AddOptionStr(sb *strings.Builder, optNum int, opt option) {
	switch opt.optType {
	case OptTypeNextChapter, OptTypeFinish:
		cursorSymbol := "  " // empty spacing for the option item
		style := styles().normal
		if int(m.cursor) == optNum {
			cursorSymbol = "▶ "
			style = styles().cursor
		}
		addStrL(sb, style.Render(cursorSymbol+opt.text))
	case OptTypeUnselectableString:
		addStrL(sb, styles().normal.Render(opt.text))
	case OptTypeUnknown:
		addStrL(sb, styles().normal.Render("option type is Unknown"))
	default:
		addStrL(sb, styles().normal.Render("option not implemented: "+opt.text))
	}
}
