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
				{OptTypeNextChapter, "Generate Password", ChapterUnimplemented, nil},
				{OptTypeNextChapter, "Add new password", ChapterUnimplemented, nil},
				{OptTypeNextChapter, "See passwords", ChapterUnimplemented, nil},
				{OptTypeFinish, "Quit", ChapterFinal, func() tea.Cmd {
					return tea.Quit
				}},
			},
			ChapterUnimplemented: {
				{OptTypeNextChapter, "Back to Main Menu", ChapterMain, nil},
				{OptTypeFinish, "Quit", ChapterFinal, func() tea.Cmd {
					return tea.Quit
				}},
			},
			ChapterFinal: {
				{OptTypeFinish, "Quit", ChapterFinal, func() tea.Cmd {
					os.Exit(0)
					return nil
				}},
			},
		},
	})

	_, err := progr.Run()
	return err
}
