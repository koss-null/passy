package interactive

import (
	"strings"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cursor      cursor
	chapter     chapter
	chapterType optType
	options     map[chapter][]option

	context map[contextKey]string
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.chapterType {
	// FIXME: the name of the chapter doesn't fit here;
	// I should split model types with optional types and
	// make a mapper
	case OptTypeNextChapter:
		return handleOptionListInput(m, msg)
	case OptTypeInputWindow:
		return handleOptionInput(m, msg)
	}
	return m, nil
}

func (m *model) View() string {
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

	// Divider + Help text
	addLStrL(&sb, styles().divider.String())
	addStrL(&sb, styles().help.Render("↑/k: up • ↓/j: down • enter: select • q/ctrl+c: quit"))

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
		text := opt.text
		for _, ctxKey := range allContextKeys {
			text = ctxKey.Substitute(text, m.context[ctxKey])
		}
		addStrL(sb, "╔═"+strings.Repeat("═", utf8.RuneCountInString(text))+"═╗")
		// FIXME: need multiple insertions to support multiline strings here
		addStrL(sb, "║ "+styles().normal.Render(text)+" ║")
		addStrL(sb, "╚═"+strings.Repeat("═", utf8.RuneCountInString(text))+"═╝")
	case OptTypeUnknown:
		addStrL(sb, styles().normal.Render("option type is Unknown"))
	default:
		addStrL(sb, styles().normal.Render("option not implemented: "+opt.text))
	}
}
