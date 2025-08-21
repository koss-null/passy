package interactive

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/koss-null/passy/internal/command/impl"
)

func Run(configPath string) error {
	app := tview.NewApplication()

	// Main menu
	list := tview.NewList().
		AddItem("Add Password", "Add new password entry", 'a', func() {
			showAddPasswordForm(app, configPath)
		}).
		AddItem("Get Password", "Retrieve stored password", 'g', func() {
			showGetPasswordForm(app, configPath)
		}).
		AddItem("Generate Password", "Create random password", 'n', func() {
			showPasswordGenerationOptions(app, configPath)
		}).
		AddItem("List Keys", "Show all stored keys", 'l', func() {
			showListKeysOptions(app, configPath)
		}).
		AddItem("Delete Password", "Remove password entry", 'd', func() {
			showDeletePasswordForm(app, configPath)
		}).
		AddItem("Edit Configuration", "Modify config file", 'c', func() {
			go executeWithLoading(app, configPath, "Editing config...", func() error {
				return impl.HandleEditConfig(configPath)
			})
		}).
		AddItem("Generate Key", "Create encryption key", 'k', func() {
			showGenerateKeyForm(app, configPath)
		}).
		AddItem("Quit", "Exit application", 'q', func() {
			app.Stop()
		})

	list.SetBorder(true).SetTitle("═══   Passy - Interactive Mode   ═══   exit with ctrl+c   ═══").SetTitleAlign(tview.AlignLeft)
	list.SetMainTextColor(tcell.ColorWhite).
		SetSecondaryTextColor(tcell.ColorLightBlue).
		SetShortcutColor(tcell.ColorBlue)

	return app.SetRoot(list, true).SetFocus(list).Run()
}

func showAddPasswordForm(app *tview.Application, configPath string) {
	form := tview.NewForm().
		AddInputField("Key Path (e.g., email/gmail):", "", 40, nil, nil).
		AddPasswordField("Password (leave empty to generate):", "", 40, '*', nil).
		AddCheckbox("Generate Readable Password", false, nil).
		AddCheckbox("Generate Safe Password", true, nil).
		AddCheckbox("Generate Strong Password", false, nil)

	form.SetBorder(true).SetTitle(" Add Password ").SetTitleAlign(tview.AlignLeft)
	form.SetButtonsAlign(tview.AlignCenter)

	form.AddButton("Save", func() {
		key := form.GetFormItem(0).(*tview.InputField).GetText()
		password := form.GetFormItem(1).(*tview.InputField).GetText()
		readable := form.GetFormItem(2).(*tview.Checkbox).IsChecked()
		safe := form.GetFormItem(3).(*tview.Checkbox).IsChecked()
		strong := form.GetFormItem(4).(*tview.Checkbox).IsChecked()

		if key == "" {
			showErrorModal(app, configPath, "Key path is required")
			return
		}

		go executeWithLoading(app, configPath, "Adding password...", func() error {
			return impl.HandleAddPassword(configPath, key, password, readable, safe, strong)
		})
	})

	form.AddButton("Cancel", func() {
		app.SetRoot(createMainMenu(configPath), true)
	})

	app.SetRoot(form, true).SetFocus(form)
}

func showGetPasswordForm(app *tview.Application, configPath string) {
	form := tview.NewForm().
		AddInputField("Key Path to retrieve:", "", 40, nil, nil)

	form.SetBorder(true).SetTitle(" Get Password ").SetTitleAlign(tview.AlignLeft)

	form.AddButton("Retrieve", func() {
		key := form.GetFormItem(0).(*tview.InputField).GetText()
		if key == "" {
			showErrorModal(app, configPath, "Key path is required")
			return
		}

		go executeWithLoading(app, configPath, "Retrieving password...", func() error {
			return impl.HandleGetPass(configPath, key)
		})
	})

	form.AddButton("Cancel", func() {
		app.SetRoot(createMainMenu(configPath), true)
	})

	app.SetRoot(form, true).SetFocus(form)
}

func showDeletePasswordForm(app *tview.Application, configPath string) {
	form := tview.NewForm().
		AddInputField("Key Path to delete:", "", 40, nil, nil)

	form.SetBorder(true).SetTitle(" Delete Password ").SetTitleAlign(tview.AlignLeft)

	form.AddButton("Delete", func() {
		key := form.GetFormItem(0).(*tview.InputField).GetText()
		if key == "" {
			showErrorModal(app, configPath, "Key path is required")
			return
		}

		// Show confirmation modal
		showConfirmationModal(
			app,
			configPath,
			fmt.Sprintf("Delete '%s' permanently?", key),
			func() {
				go executeWithLoading(app, configPath, "Deleting...", func() error {
					return impl.HandleDeletePassword(configPath, key)
				})
			},
		)
	})

	form.AddButton("Cancel", func() {
		app.SetRoot(createMainMenu(configPath), true)
	})

	app.SetRoot(form, true).SetFocus(form)
}

func showPasswordGenerationOptions(app *tview.Application, configPath string) {
	form := tview.NewForm().
		AddCheckbox("Readable (memorable)", false, nil).
		AddCheckbox("Safe (balanced)", true, nil).
		AddCheckbox("Strong (maximum security)", false, nil)

	form.SetBorder(true).SetTitle(" Generate Password ").SetTitleAlign(tview.AlignLeft)

	form.AddButton("Generate", func() {
		readable := form.GetFormItem(0).(*tview.Checkbox).IsChecked()
		safe := form.GetFormItem(1).(*tview.Checkbox).IsChecked()
		strong := form.GetFormItem(2).(*tview.Checkbox).IsChecked()

		go executeWithLoading(app, configPath, "Generating password...", func() error {
			return impl.HandlePasswordComposition(readable, safe, strong)
		})
	})

	form.AddButton("Cancel", func() {
		app.SetRoot(createMainMenu(configPath), true)
	})

	app.SetRoot(form, true).SetFocus(form)
}

func showListKeysOptions(app *tview.Application, configPath string) {
	form := tview.NewForm().
		AddCheckbox("Show passwords (reveal all)", false, nil)

	form.SetBorder(true).SetTitle(" List Keys ").SetTitleAlign(tview.AlignLeft)

	form.AddButton("List", func() {
		showAll := form.GetFormItem(0).(*tview.Checkbox).IsChecked()

		go executeWithLoading(app, configPath, "Loading keys...", func() error {
			return impl.HandleShowKeys(configPath, showAll)
		})
	})

	form.AddButton("Cancel", func() {
		app.SetRoot(createMainMenu(configPath), true)
	})

	app.SetRoot(form, true).SetFocus(form)
}

func showGenerateKeyForm(app *tview.Application, configPath string) {
	form := tview.NewForm().
		AddInputField("Key file path:", "", 40, nil, nil)

	form.SetBorder(true).SetTitle(" Generate Encryption Key ").SetTitleAlign(tview.AlignLeft)

	form.AddButton("Generate", func() {
		keyPath := form.GetFormItem(0).(*tview.InputField).GetText()
		if keyPath == "" {
			showErrorModal(app, configPath, "Key file path is required")
			return
		}

		go executeWithLoading(app, configPath, "Generating key...", func() error {
			return impl.HandleKeyGeneration(keyPath)
		})
	})

	form.AddButton("Cancel", func() {
		app.SetRoot(createMainMenu(configPath), true)
	})

	app.SetRoot(form, true).SetFocus(form)
}

func executeWithLoading(app *tview.Application, configPath, message string, task func() error) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons(nil) // No buttons for loading

	app.QueueUpdateDraw(func() {
		app.SetRoot(modal, false)
	})

	err := task()

	app.QueueUpdateDraw(func() {
		if err != nil {
			showErrorModal(app, configPath, err.Error())
		} else {
			app.SetRoot(createMainMenu(configPath), true)
		}
	})
}

func showErrorModal(app *tview.Application, configPath, message string) {
	modal := tview.NewModal().
		SetText("Error: " + message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(createMainMenu(configPath), true)
		})
	app.SetRoot(modal, false)
}

func showConfirmationModal(app *tview.Application, configPath, message string, confirmFunc func()) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				confirmFunc()
			} else {
				app.SetRoot(createMainMenu(configPath), true)
			}
		})
	app.SetRoot(modal, false)
}

func createMainMenu(configPath string) tview.Primitive {
	list := tview.NewList().
		AddItem("Add Password", "Add new password entry", 'a', func() {
			showAddPasswordForm(tview.NewApplication(), configPath)
		}).
		AddItem("Quit", "Exit application", 'q', func() {
			tview.NewApplication().Stop()
		})

	list.SetBorder(true).SetTitle(" Passy - Interactive Mode ")
	return list
}
