package interactive

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type chapter string

const (
	ChapterMain          = chapter("Main Menu")
	ChapterUnimplemented = chapter("Under Construction")
	ChapterFinal         = chapter("Quitting")
)

type cursor uint16

func (c *cursor) Up() {
	if *c > 0 {
		*c--
	}
}

func (c *cursor) Down(threshold int) {
	if *c < cursor(threshold) {
		*c++
	}
}

type option struct {
	name    string
	next    chapter
	handler func() tea.Cmd
}

type model struct {
	cursor
	chapter
	options  map[chapter][]option
	selected map[chapter]map[int]struct{}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.chapter == ChapterFinal {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.chapter = ChapterFinal
			return m, tea.Quit
		case "up", "k":
			m.Up()
		case "down", "j":
			length := 0
			if opts, ok := m.options[m.chapter]; ok {
				length = len(opts) - 1
			}
			m.Down(length)
		case "enter", " ":
			if opts, ok := m.options[m.chapter]; ok {
				var cmd tea.Cmd
				if opts[m.cursor].handler != nil {
					cmd = opts[m.cursor].handler()
				}
				m.chapter = opts[m.cursor].next
				m.cursor = 0

				// Final chapter quits immediately
				if m.chapter == ChapterFinal {
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
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Padding(0, 1)

	cursorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("250"))

	dividerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		SetString("┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈")

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Italic(true)

	var s strings.Builder

	// Title
	title := titleStyle.Render(string(m.chapter))
	s.WriteString(title + "\n")
	s.WriteString(dividerStyle.String() + "\n")
	// Options
	if opts, ok := m.options[m.chapter]; ok {
		for i := range opts {
			if i == int(m.cursor) {
				s.WriteString(cursorStyle.Render("▶ " + opts[i].name))
			} else {
				s.WriteString(normalStyle.Render("  " + opts[i].name))
			}
			s.WriteString("\n")
		}
	}
	// Divider
	s.WriteString("\n")
	s.WriteString(dividerStyle.String())
	s.WriteString("\n")
	// Help text
	helpText := helpStyle.Render("↑/k: up • ↓/j: down • enter: select • q/ctrl+c: quit")
	s.WriteString(helpText)
	s.WriteString("\n")

	return s.String()
}

func Run(configPath string) error {
	progr := tea.NewProgram(&model{
		cursor:  cursor(0),
		chapter: ChapterMain,
		options: map[chapter][]option{
			ChapterMain: {
				{"Generate Password", ChapterUnimplemented, nil},
				{"Add new password", ChapterUnimplemented, nil},
				{"See passwords", ChapterUnimplemented, nil},
				{"Quit", ChapterFinal, func() tea.Cmd {
					return tea.Quit
				}},
			},
			ChapterUnimplemented: {
				{"Back to Main Menu", ChapterMain, nil},
				{"Quit", ChapterFinal, func() tea.Cmd {
					return tea.Quit
				}},
			},
			ChapterFinal: {
				{"Quit", ChapterFinal, func() tea.Cmd {
					os.Exit(0)
					return nil
				}},
			},
		},
	})

	_, err := progr.Run()
	return err
}
