package views

import (
	"image/color"
	"pick-a-bro/internal/commons"
	"pick-a-bro/internal/data"
	"strconv"
	"strings"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func SetRules(window fyne.Window) {
	rules(window)
}

// rules is a function that sets up the rules view in the application window.
// It creates various UI components such as header, select, entry, grid, and buttons,
// and adds them to the content container of the window.
// The function also handles the logic for showing and hiding certain UI components
// based on the selected rule.
//
// Parameters:
// - window: The fyne.Window in which the rules view will be displayed.
//
// Returns: None
func rules(window fyne.Window) {
	headerContainer := createHeaderContainer()

	chancesRule := createSelect()
	chancesLabel := widget.NewLabel(commons.GetTranslation(commons.I18n.ChancesPerPatreon))
	chancesPerUser := createEntry(commons.GetPreferences().IntWithFallback(commons.ChancesPerUser, 1), commons.ChancesPerUser)
	chancesContainer := container.NewHBox(chancesLabel, chancesPerUser)

	membersList := data.GetMembersAndTiers()
	tierEntries := container.NewHBox()
	tierEntries.Add(chancesLabel)
	for _, tier := range membersList.Tiers {
		label := widget.NewLabel(tier.(string))
		entry := createEntry(commons.GetPreferences().IntWithFallback("chances"+tier.(string), 1), "chances"+tier.(string))
		tierEntries.Add(container.NewHBox(label, entry))
	}

	if chancesRule.Selected == commons.GetTranslation(commons.ChancesRules[1]) {
		chancesContainer.Hide()
		tierEntries.Show()
	} else {
		tierEntries.Hide()
	}

	headerGrid, membersGrid := createGrid(membersList)

	rulesViewContainer := container.NewVBox(
		headerContainer,
		chancesRule,
		chancesContainer,
		tierEntries,
		headerGrid,
		membersGrid,
	)

	rect := canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 0, A: 180})
	rect.SetMinSize(fyne.NewSize(50, 20))
	rulesView := container.NewStack(rect, rulesViewContainer)

	confirmButtons := createConfirmButtons(window)

	membersGrid.SetMinSize(fyne.NewSize(window.Canvas().Size().Width,
		window.Canvas().Size().Height-(headerContainer.MinSize().Height+
			chancesRule.MinSize().Height+
			headerGrid.MinSize().Height+
			confirmButtons.MinSize().Height)))
	rulesView.Refresh()

	chancesRule.OnChanged = func(value string) {
		commons.GetPreferences().SetString(commons.ChancesRule, value)
		if value == commons.GetTranslation(commons.ChancesRules[0]) {
			tierEntries.Hide()
			chancesContainer.Show()
			membersGrid.SetMinSize(fyne.NewSize(commons.WindowWidth,
				window.Canvas().Size().Height-(headerContainer.MinSize().Height+
					chancesRule.MinSize().Height+
					chancesContainer.MinSize().Height+
					headerGrid.MinSize().Height+
					confirmButtons.MinSize().Height)))
		} else {
			chancesContainer.Hide()
			tierEntries.Show()
			tierEntries.Refresh()
			membersGrid.SetMinSize(fyne.NewSize(commons.WindowWidth,
				window.Canvas().Size().Height-(headerContainer.MinSize().Height+
					chancesRule.MinSize().Height+
					tierEntries.MinSize().Height+
					headerGrid.MinSize().Height+
					confirmButtons.MinSize().Height)))
		}
	}
	content := container.New(layout.NewStackLayout(), commons.GetBackgroundImage(), rulesView, confirmButtons)

	window.SetContent(content)
}

// createHeaderContainer creates and returns a widget.Check that represents the header container.
// The widget.Check allows the user to exclude winners based on the provided translation.
// The value of the widget.Check is stored in the preferences using the commons.ExcludeWinners key.
func createHeaderContainer() *widget.Check {
	excludeWinners := widget.NewCheck(commons.GetTranslation(commons.I18n.ExcludeWinners), func(value bool) {
		commons.GetPreferences().SetBool(commons.ExcludeWinners, value)
	})
	return excludeWinners
}

// createSelect creates and returns a new widget.Select with options populated from commons.ChancesRules.
// It sets the selected option based on the user's preferences.
func createSelect() *widget.Select {
	options := make([]string, len(commons.ChancesRules))

	copy(options, commons.ChancesRules)
	for i, option := range options {
		options[i] = commons.GetTranslation(option)
	}
	selectWidget := widget.NewSelect(options, nil)
	selectWidget.SetSelected(commons.GetPreferences().StringWithFallback(commons.ChancesRule, commons.GetTranslation(commons.ChancesRules[0])))
	return selectWidget
}

// createEntry creates a new widget.Entry with the specified default value and preference key.
// The default value is converted to a string and set as the initial text of the entry.
// The entry's OnChanged event is set to a function that updates the entry's text based on user input,
// filters out non-digit characters, and updates the corresponding preference value.
// If the entry's text is empty, it is set to "1" as the default value.
// The filtered text is then converted to an integer and stored in the preference with the specified key.
// The created entry is returned.
func createEntry(defaultValue int, preferenceKey string) *widget.Entry {
	entry := widget.NewEntry()
	entry.Text = strconv.Itoa(defaultValue)
	entry.OnChanged = func(value string) {
		if value == "" {
			entry.Text = "1"
			return
		}

		filtered := strings.Map(func(r rune) rune {
			if unicode.IsDigit(r) {
				return r
			}
			return -1 // Discard this rune
		}, value)

		// Prevent infinite loop by checking if filtering is necessary
		if filtered != value {
			entry.SetText(filtered)
		}
		chances, _ := strconv.Atoi(filtered)
		commons.GetPreferences().SetInt(preferenceKey, chances)
	}
	return entry
}

// createGridCells creates grid cells for each member in the given membersList.
// It takes a pointer to a MembersList struct and returns a slice of fyne.CanvasObject.
func createGridCells(membersList *data.MembersList) []fyne.CanvasObject {
	members := []fyne.CanvasObject{}
	for _, d := range membersList.PatreonMembers {
		color := membersList.ColorCode[d.Tier]
		members = append(members, makeCellWithBackground(widget.NewLabel(d.FullName).Text, color),
			makeCellWithBackground(widget.NewLabel(d.Tier).Text, color))
	}
	return members
}

// Function to create a widget with a background color
// makeCellWithBackground creates a fyne.CanvasObject that consists of a label with the specified text and a background color.
// The label displays the given text, and the background is a rectangle filled with the specified color.
// The function returns a fyne.CanvasObject that can be added to a UI layout.
func makeCellWithBackground(text string, color color.Color) fyne.CanvasObject {
	label := widget.NewLabel(text)
	background := canvas.NewRectangle(color)
	return container.NewStack(background, label)
}

// createGrid creates a grid layout containing the header and members grid.
// It takes a pointer to a MembersList and returns a Container and Scroll widget.
func createGrid(membersList *data.MembersList) (*fyne.Container, *container.Scroll) {
	cells := createGridCells(membersList)

	headerGrid := container.NewHBox(
		widget.NewLabel(commons.Fellowship),
	)
	membersGrid := container.NewVScroll(container.NewGridWithColumns(2, cells...))
	return headerGrid, membersGrid
}

// createConfirmButtons creates and returns a container with confirm buttons for the window.
// It takes a fyne.Window as input and returns a *fyne.Container.
func createConfirmButtons(window fyne.Window) *fyne.Container {
	lotteryBtn := widget.NewButton(commons.PickABro, func() {
		SetLottery(window)
	})

	cancelDraw := widget.NewButton(commons.GetTranslation(commons.I18n.Cancel), func() {
		MainMenu(window)
	})

	mainButtons := container.NewVBox(layout.NewSpacer(), lotteryBtn, cancelDraw)
	return mainButtons
}
