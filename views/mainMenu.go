package views

import (
	"fmt"
	"pick-a-bro/internal/commons"
	"pick-a-bro/internal/data"
	"pick-a-bro/internal/lottery"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// MainMenu is a function that creates and displays the main menu of the application.
// It takes a `window` parameter of type `fyne.Window` to display the menu.
// The main menu consists of several buttons, including a "New Draw" button, a "Settings" button,
// a "Previous Winners" button, and a "Test Mode" checkbox.
// Clicking the "New Draw" button will either handle the test mode or the normal mode based on the user's preferences.
// Clicking the "Settings" button will open the preferences panel.
// Clicking the "Previous Winners" button will display a list of previous winners and provide an option to clear the winners list.
// Clicking the "Test Mode" checkbox will toggle the test mode on or off based on the user's selection.
// The main menu is displayed within the specified `window`.
func MainMenu(window fyne.Window) {
	newDrawButton := widget.NewButton(commons.GetTranslation(commons.I18n.NewDraw), func() {
		if commons.GetPreferences().Bool(commons.Settings.TestMode) {
			handleTestMode(window)
		} else {
			handleNormalMode(window)
		}
	})

	settingsButton := widget.NewButton(commons.GetTranslation(commons.I18n.Settings), func() {
		preferencesPanel(window)
	})

	previousWinnersButton := widget.NewButton(commons.GetTranslation(commons.I18n.PreviousWinners), func() {
		winners := []fyne.CanvasObject{}

		for _, d := range lottery.GetWinnersList() {
			winners = append(winners, widget.NewLabel(d.FullName), widget.NewLabel(d.DateTime))
		}
		grid := container.NewGridWithColumns(2, winners...)
		scroll := container.NewVScroll(grid)
		scroll.SetMinSize(fyne.NewSize(400, 400))
		clearWinners := widget.NewButton(commons.GetTranslation(commons.I18n.ClearWinners), func() {

			confirmDialog := dialog.NewConfirm(commons.GetTranslation(commons.I18n.ClearWinners),
				commons.GetTranslation(commons.I18n.ConfirmClearWinners),
				func(resp bool) {
					if resp {
						lottery.ClearWinnersList()
						dialog.NewInformation(commons.GetTranslation(commons.I18n.WinnersCleared),
							commons.GetTranslation(commons.I18n.WinnersListCleared), window).Show()
						grid.Hide()
					}
				}, window)
			confirmDialog.SetConfirmText(commons.GetTranslation(commons.I18n.Yes))
			confirmDialog.SetDismissText(commons.GetTranslation(commons.I18n.No))
			confirmDialog.Show()
		})
		memberstable := container.NewVBox(scroll, clearWinners)
		dialogCustom := dialog.NewCustom(commons.GetTranslation(commons.I18n.PreviousWinners),
			commons.GetTranslation(commons.I18n.Close), memberstable, window)
		dialogCustom.Resize(fyne.NewSize(400, 400))
		dialogCustom.Show()
	})

	testModeCheckbox := widget.NewCheck(commons.GetTranslation(commons.I18n.TestMode), func(value bool) {
		commons.GetPreferences().SetBool(commons.Settings.TestMode, value)
	})

	testModeCheckbox.SetChecked(commons.GetPreferences().BoolWithFallback(commons.Settings.TestMode, false))

	mainButtons := container.NewVBox(layout.NewSpacer(), newDrawButton, settingsButton, previousWinnersButton, testModeCheckbox)
	content := container.New(layout.NewStackLayout(), commons.GetBackgroundImage(), mainButtons)
	window.SetContent(content)
}

func handleTestMode(window fyne.Window) {
	testModeWrnLbl := widget.NewLabel(commons.GetTranslation(commons.I18n.TestModeWrn))
	testModeWrnLbl.Wrapping = fyne.TextWrapWord

	var dialogPanel *dialog.CustomDialog
	testData := widget.NewButton(commons.GetTranslation(commons.I18n.TestDummyData), func() {
		commons.GetPreferences().SetBool(commons.Settings.UseRealData, false)
		checkAndGenerateTestData(window)
		SetRules(window)
		dialogPanel.Hide()
	})

	testRealData := widget.NewButton(commons.GetTranslation(commons.I18n.TestRealData), func() {
		commons.GetPreferences().SetBool(commons.Settings.UseRealData, true)
		checkAndGenerateTestData(window)
		SetRules(window)
		dialogPanel.Hide()
	})

	dialogButtons := container.NewHBox(testData, testRealData)
	dialogContent := container.NewVBox(testModeWrnLbl, dialogButtons)

	dialogPanel = dialog.NewCustom(commons.GetTranslation(commons.I18n.TestMode), commons.GetTranslation(commons.I18n.Cancel), dialogContent, window)
	dialogPanel.Resize(fyne.NewSize(200, 200))
	dialogPanel.Show()
}

func handleNormalMode(window fyne.Window) {
	if !data.ExtractDataFromFile() {
		dialog.NewInformation(commons.GetTranslation(commons.I18n.MissingData), commons.GetTranslation(commons.I18n.NoPatreons), window).Show()
		preferencesPanel(window)
	} else {
		dialog.NewCustomConfirm(commons.GetTranslation(commons.I18n.PatreonsList), commons.GetTranslation(commons.I18n.Yes),
			commons.GetTranslation(commons.I18n.No), widget.NewLabel(commons.GetTranslation(commons.I18n.RefreshPatreonsList)), func(resp bool) {
				if resp {
					fetchPatreonsList(window)
				}
				SetRules(window)
			}, window).Show()
	}
}

func checkAndGenerateTestData(window fyne.Window) {
	if !data.ExtractDataFromFile() {
		dialog.NewCustomWithoutButtons(commons.GetTranslation(commons.I18n.TestData),
			widget.NewLabel(commons.GetTranslation(commons.I18n.TestDataGenerated)), window).Show()
		members, tiers := data.FetchMembersToLocalStorage()
		data.SetMembersList(members)
		data.SetTiersMap(tiers)
	}
}

func fetchPatreonsList(window fyne.Window) {
	watingDialog := dialog.NewCustomWithoutButtons(commons.GetTranslation(commons.I18n.FetchingPatreons),
		widget.NewLabel(fmt.Sprintf("%s...", commons.GetTranslation(commons.I18n.FetchingPatreons))), window)
	watingDialog.Show()
	if members, _ := data.FetchMembersToLocalStorage(); members == nil {
		watingDialog.Hide()
		dialog.NewError(fmt.Errorf(commons.GetTranslation(commons.I18n.ErrorFetchingPatreons)), window).Show()
		return
	}
	watingDialog.Hide()
}
