package views

import (
	"encoding/json"
	"pick-a-bro/internal/commons"
	"pick-a-bro/internal/custom_widgets"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// SelectLanguage is a function that creates a language selection screen.
// It takes a fyne.Window as a parameter and sets the content of the window to the language selection screen.
func SelectLanguage(window fyne.Window) {
	translations := loadTranslations()
	elBtnImg := commons.EmbedImage(commons.GetAsset(commons.AssetsPaths.ImagesPath, commons.AssetsKeys.ElHandshakeImg), commons.AssetsKeys.ElHandshakeImg)
	enBtnImg := commons.EmbedImage(commons.GetAsset(commons.AssetsPaths.ImagesPath, commons.AssetsKeys.EnHandshakeImg), commons.AssetsKeys.EnHandshakeImg)

	elBtn := createButton(elBtnImg, translations, language.Greek.String(), window)
	enBtn := createButton(enBtnImg, translations, language.AmericanEnglish.String(), window)

	elBtnAlign := alignButton(elBtn)
	enBtnAlign := alignButton(enBtn)

	centerButtonsContainer := container.NewHBox(layout.NewSpacer(), elBtnAlign, layout.NewSpacer(), enBtnAlign, layout.NewSpacer())
	content := container.New(layout.NewStackLayout(), commons.GetBackgroundImage(), centerButtonsContainer)
	window.SetContent(content)
}

// loadTranslations is a function that loads the translations for the application.
// It returns a pointer to an i18n.Bundle that contains the translations.
func loadTranslations() *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	commons.EmbedLocales(bundle)
	return bundle
}

// createButton is a function that creates a custom image button.
// It takes an image resource, a translation bundle, a language string, and a fyne.Window as parameters.
// It returns a pointer to a custom_widgets.ImageButton.
func createButton(img fyne.Resource, bundle *i18n.Bundle, lang string, window fyne.Window) *custom_widgets.ImageButton {
	return custom_widgets.NewImageButton(img, func() {
		commons.SetLocalization(i18n.NewLocalizer(bundle, lang))
		MainMenu(window)
	})
}

// alignButton is a function that aligns a custom image button vertically.
// It takes a pointer to a custom_widgets.ImageButton as a parameter.
// It returns a pointer to a fyne.Container.
func alignButton(button *custom_widgets.ImageButton) *fyne.Container {
	return container.NewVBox(layout.NewSpacer(), widget.NewLabel(""), layout.NewSpacer(), button, layout.NewSpacer())
}
