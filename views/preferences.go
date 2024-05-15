package views

import (
	"fmt"
	"image/color"
	"pick-a-bro/internal/commons"
	"pick-a-bro/internal/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// preferencesPanel is a function that creates and displays the preferences panel in the application window.
// It takes a fyne.Window as a parameter and sets the content of the window to the preferences panel.
func preferencesPanel(window fyne.Window) {
	formItems := createFormItems()
	preferencesProcessing(formItems, "onLoad")

	form := createForm(window, formItems)

	discardPreferences := createDiscardPreferencesButton(window)
	readLogs := createReadLogsButton(window)

	mainButtons := container.NewVBox(layout.NewSpacer(), discardPreferences, readLogs)
	background := createBackground()
	rulesView := container.NewStack(background, form)

	content := container.New(layout.NewStackLayout(), commons.GetBackgroundImage(), rulesView, mainButtons)
	window.SetContent(content)
}

// createFormItems creates and returns a slice of *widget.FormItem.
// Each *widget.FormItem consists of a label and a password entry widget.
// The labels are predefined constants from the commons package.
// The password entry widgets are created using widget.NewPasswordEntry().
func createFormItems() []*widget.FormItem {
	clientIdEntry := widget.NewPasswordEntry()
	clientSecretEntry := widget.NewPasswordEntry()
	campaignIdEntry := widget.NewPasswordEntry()

	clientIdForm := widget.NewFormItem(commons.ClientId, clientIdEntry)
	clientSecretForm := widget.NewFormItem(commons.ClientSecret, clientSecretEntry)
	campaignIdForm := widget.NewFormItem(commons.CampaignId, campaignIdEntry)

	return []*widget.FormItem{clientIdForm, clientSecretForm, campaignIdForm}
}

// createForm creates a new widget.Form with the specified formItems and submit handler.
// It returns the created widget.Form.
func createForm(window fyne.Window, formItems []*widget.FormItem) *widget.Form {
	return &widget.Form{
		Items: formItems,
		OnSubmit: func() {
			handleFormSubmit(window, formItems)
		},
		SubmitText: "You, Fetch!",
	}
}

// handleFormSubmit handles the form submission in the preferences view.
// It shows a waiting dialog while processing the form items, then updates the preferences accordingly.
// If the test mode is enabled, it toggles it off temporarily and restores it afterwards.
// It fetches members from local storage and shows a success dialog if successful.
// If there is an error fetching the members, it shows an error dialog.
func handleFormSubmit(window fyne.Window, formItems []*widget.FormItem) {
	waitingDialog := dialog.NewCustomWithoutButtons(commons.GetTranslation(commons.I18n.FetchingPatreons),
		widget.NewLabel(fmt.Sprintf("%s...", commons.GetTranslation(commons.I18n.FetchingPatreons))), window)
	waitingDialog.Show()
	preferencesProcessing(formItems, "onSave")
	testMode := commons.GetPreferences().BoolWithFallback(commons.TestMode, false)
	testModeToogled := false
	if testMode {
		commons.GetPreferences().SetBool(commons.TestMode, false)
		testModeToogled = true
	}
	if members, _ := data.FetchMembersToLocalStorage(); members == nil {
		waitingDialog.Hide()
		dialog.NewError(fmt.Errorf(commons.GetTranslation(commons.I18n.ErrorFetchingPatreons)), window).Show()
	} else {
		waitingDialog.Hide()
		showSuccessDialog(window)
	}

	if testModeToogled {
		commons.GetPreferences().SetBool(commons.TestMode, true)
	}
}

// showSuccessDialog displays a success dialog with a custom message.
// It takes a fyne.Window as a parameter and shows the dialog on that window.
// After the dialog is closed, it calls the MainMenu function to return to the main menu.
func showSuccessDialog(window fyne.Window) {
	dialogCustom := dialog.NewCustom(commons.GetTranslation(commons.I18n.Success),
		commons.GetTranslation(commons.I18n.Close), widget.NewLabel(commons.GetTranslation(commons.I18n.SuccessfulReceive)), window)
	dialogCustom.SetOnClosed(func() {
		MainMenu(window)
	})
	dialogCustom.Show()
}

// createDiscardPreferencesButton creates a button that discards the preferences and returns to the main menu.
func createDiscardPreferencesButton(window fyne.Window) *widget.Button {
	return widget.NewButton(commons.GetTranslation(commons.I18n.Cancel), func() {
		MainMenu(window)
	})
}

// createReadLogsButton creates a button that, when clicked, shows the logs in a window.
func createReadLogsButton(window fyne.Window) *widget.Button {
	return widget.NewButton(commons.GetTranslation(commons.I18n.ReadLogs), func() {
		showLogs(window)
	})
}

// showLogs displays the logs in a dialog window.
// It retrieves the logs using the GetLogs function from the data package,
// sets them as the content of a multi-line entry widget, and displays the content in a scrollable container.
// The logs are shown in a dialog window with a custom title and a close button.
// The window parameter is used to determine the parent window for the dialog.
func showLogs(window fyne.Window) {
	logs, _ := data.GetLogs()

	content := widget.NewEntry()
	content.SetText(logs)
	content.MultiLine = true
	content.Disable()
	scroll := container.NewVScroll(content)

	dialogCustom := dialog.NewCustom("Logs", commons.GetTranslation(commons.I18n.Close), scroll, window)
	dialogCustom.Resize(fyne.NewSize(400, 400))
	dialogCustom.Show()
}

// createBackground creates and returns a new canvas rectangle with a specified color and size.
func createBackground() *canvas.Rectangle {
	rect := canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 0, A: 180})
	rect.SetMinSize(fyne.NewSize(50, 20))
	return rect
}

// preferencesProcessing processes the form items based on the specified action.
// It iterates over the form items and performs different actions on the widget.Entry items
// based on the given action.
//
// Parameters:
// - formItems: A slice of form items.
// - action: The action to perform on the form items.
//
// Possible actions:
// - "onLoad": Sets the text of each widget.Entry item to the corresponding value from the preferences.
// - "onSave": Sets the value of each widget.Entry item in the preferences to the entered text.
func preferencesProcessing(formItems []*widget.FormItem, action string) {
	for _, item := range formItems {
		entry, ok := item.Widget.(*widget.Entry)
		if ok {
			switch action {
			case "onLoad":
				entry.SetText(commons.GetPreferences().String(item.Text))
			case "onSave":
				commons.GetPreferences().SetString(item.Text, entry.Text)
			}
		}
	}
}
