package interactive

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type tuiStyles struct {
	title   lipgloss.Style
	cursor  lipgloss.Style
	normal  lipgloss.Style
	divider lipgloss.Style
	help    lipgloss.Style
}

func styles() tuiStyles {
	return tuiStyles{
		title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Padding(0, 1),

		cursor: lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Bold(true),

		normal: lipgloss.NewStyle().
			Foreground(lipgloss.Color("250")),

		divider: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			SetString("┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈"),

		help: lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true),
	}
}

func title(cpt chapter) string {
	return styles().title.Render(string(cpt))
}

func addStrL(sb *strings.Builder, str string) {
	sb.WriteString(str + "\n")
}

func addLStrL(sb *strings.Builder, str string) {
	sb.WriteString("\n" + str + "\n")
}
