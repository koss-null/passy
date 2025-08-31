package interactive

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cursor   cursor
	chapter  chapter
	options  map[chapter][]option
	selected map[chapter]map[int]struct{}
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
	var s strings.Builder

	// Title
	title := Styles().title.Render(string(m.chapter))
	s.WriteString(title + "\n")
	s.WriteString(Styles().divider.String() + "\n")
	// Options
	if opts, ok := m.options[m.chapter]; ok {
		for i := range opts {
			if i == int(m.cursor) {
				s.WriteString(Styles().cursor.Render("▶ " + opts[i].name))
			} else {
				s.WriteString(Styles().normal.Render("  " + opts[i].name))
			}
			s.WriteString("\n")
		}
	}
	// Divider
	s.WriteString("\n")
	s.WriteString(Styles().divider.String())
	s.WriteString("\n")
	// Help text
	helpText := Styles().help.Render("↑/k: up • ↓/j: down • enter: select • q/ctrl+c: quit")
	s.WriteString(helpText)
	s.WriteString("\n")

	return s.String()
}

type chapter string

const (
	ChapterMain          = chapter("Main Menu")
	ChapterUnimplemented = chapter("Under Construction")
	ChapterFinal         = chapter("Quitting")
)

func (c *chapter) isFinal() bool {
	return *c == ChapterFinal
}

func (c *chapter) change(nextChap chapter) bool {
	// TODO: should add some state machine here
	*c = nextChap
	return true
}

type option struct {
	name    string
	next    chapter
	handler func() tea.Cmd
}
