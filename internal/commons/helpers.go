package commons

import (
	"log"

	"fyne.io/fyne/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var pickABro fyne.App
var preferences fyne.Preferences
var loc *i18n.Localizer
var logger *log.Logger

func SetApplications(app fyne.App) {
	pickABro = app
}

func GetApplication() fyne.App {
	return pickABro
}

func SetPreferences(prefs fyne.Preferences) {
	preferences = prefs
}

func GetPreferences() fyne.Preferences {
	return preferences
}

func SetLocalization(l *i18n.Localizer) {
	loc = l
}

func GetLogger() *log.Logger {
	return logger
}

func SetLogger(l *log.Logger) {
	logger = l
}

// GetTranslation retrieves the translation for the given key.
// It uses the loc.Localize function from the i18n package to perform the translation.
// If an error occurs during the translation process, an empty string is returned.
func GetTranslation(key string) string {
	translation, err := loc.Localize(&i18n.LocalizeConfig{MessageID: key})
	if err != nil {
		return ""
	}
	return translation
}
