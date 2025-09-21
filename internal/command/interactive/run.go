package interactive

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/koss-null/passy/internal/passgen"
)

func Run(configPath string) error {
	progr := tea.NewProgram(&model{
		cursor:      cursor(0),
		chapter:     ChapterMain,
		chapterType: OptTypeNextChapter,
		options: map[chapter][]option{
			ChapterMain: {
				{
					optType: OptTypeNextChapter,
					text:    "Generate Password",
					next:    ChapterPasswordGen,
					handler: nil,
				},
				{
					optType: OptTypeNextChapter,
					text:    "Add new password",
					next:    ChapterUnimplemented,
					handler: nil,
				},
				{
					optType: OptTypeNextChapter,
					text:    "See passwords",
					next:    ChapterUnimplemented,
					handler: nil,
				},
				{
					optType: OptTypeFinish,
					text:    "Quit",
					next:    ChapterFinal,
					handler: func(m *model) (*model, tea.Cmd) {
						return m, tea.Quit
					},
				},
			},
			ChapterPasswordGen: {
				{
					optType: OptTypeNextChapter,
					text:    "Readable",
					next:    ChapterShowGeneratedPass,
					handler: func(m *model) (*model, tea.Cmd) {
						gen, err := passgen.New()
						if err != nil {
							m.context[ContextKeyPass] = gen.GenReadablePass()
							return m, tea.Quit
						}
						m.context[ContextKeyPass] = gen.GenReadablePass()
						m.context[ContextKeyLastPassStrength] = "Readable"
						return m, nil
					},
				},
				{
					optType: OptTypeNextChapter,
					text:    "Safe",
					next:    ChapterShowGeneratedPass,
					handler: func(m *model) (*model, tea.Cmd) {
						gen, err := passgen.New()
						if err != nil {
							m.context[ContextKeyPass] = gen.GenReadablePass()
							return m, tea.Quit
						}
						m.context[ContextKeyPass] = gen.GenSafePass()
						m.context[ContextKeyLastPassStrength] = "Safe"
						return m, nil
					},
				},
				{
					optType: OptTypeNextChapter,
					text:    "Insane",
					next:    ChapterShowGeneratedPass,
					handler: func(m *model) (*model, tea.Cmd) {
						gen, err := passgen.New()
						if err != nil {
							m.context[ContextKeyPass] = gen.GenReadablePass()
							return m, tea.Quit
						}
						m.context[ContextKeyPass] = gen.GenInsanePass()
						m.context[ContextKeyLastPassStrength] = "Insane"
						return m, nil
					},
				},
			},
			ChapterShowGeneratedPass: {
				{
					optType: OptTypeUnselectableString,
					text:    ContextKeyPass.Template(),
					next:    "",
					handler: nil,
				},
				{
					optType: OptTypeNextChapter,
					text:    "One more",
					next:    ChapterShowGeneratedPass,
					handler: func(m *model) (*model, tea.Cmd) {
						gen, err := passgen.New()
						if err != nil {
							m.context[ContextKeyPass] = gen.GenReadablePass()
							return m, tea.Quit
						}
						switch m.context[ContextKeyLastPassStrength] {
						case "Readable":
							m.context[ContextKeyPass] = gen.GenReadablePass()
						case "Safe":
							m.context[ContextKeyPass] = gen.GenSafePass()
						case "Insane":
							m.context[ContextKeyPass] = gen.GenInsanePass()
						}
						return m, nil
					},
				},
				{
					optType: OptTypeNextChapter,
					text:    "Back",
					next:    ChapterPasswordGen,
				},
				{
					optType: OptTypeNextChapter,
					text:    "Main menu",
					next:    ChapterMain,
				},
			},
			ChapterUnimplemented: {
				{
					optType: OptTypeNextChapter,
					text:    "Back to Main Menu",
					next:    ChapterMain,
					handler: nil,
				},
				{
					optType: OptTypeFinish,
					text:    "Quit",
					next:    ChapterFinal,
					handler: func(m *model) (*model, tea.Cmd) {
						return m, tea.Quit
					},
				},
			},
			ChapterFinal: {
				{
					optType: OptTypeFinish,
					text:    "Quit",
					next:    ChapterFinal,
					handler: func(m *model) (*model, tea.Cmd) {
						return m, tea.Quit
					},
				},
			},
		},
		context: make(map[contextKey]string),
	})

	_, err := progr.Run()
	return err
}
