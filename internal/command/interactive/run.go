package interactive

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

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
