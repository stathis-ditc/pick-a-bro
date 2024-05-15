package internal

import (
	"embed"
	"log"
	"os"
	"pick-a-bro/internal/commons"
	"pick-a-bro/views"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

//go:embed assets/images/*.png
var imagesFS embed.FS

//go:embed assets/audio/*.mp3
var audioFS embed.FS

//go:embed locale/*.json
var localeFS embed.FS

//go:embed tests/samples/*.json
var samplesFS embed.FS

// RunApp initializes and runs the Pick a Bro application.
func RunApp() {
	// Create a new instance of the Pick a Bro application
	pickABro := app.NewWithID("cloud.devsinthe.pick-a-bro")

	// Create the main window for the application
	mainPanel := pickABro.NewWindow("Pick a Bro")

	// Resize the main window to a specific size
	mainPanel.Resize(fyne.NewSize(commons.WindowWidth, commons.WindowHeight))

	// Set the application preferences
	commons.SetPreferences(pickABro.Preferences())

	// Show the language selection view
	views.SelectLanguage(mainPanel)

	// Show and run the main window
	mainPanel.ShowAndRun()
}

func init() {
	// Open a file for writing logs
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new logger that writes to the file
	logger := log.New(file, "", log.LstdFlags)

	// Set the logger for the application
	commons.SetLogger(logger)

	// Set the embedded file systems for images, audio, locales, and samples
	commons.SetImagesFS(&imagesFS)
	commons.SetAudioFS(&audioFS)
	commons.SetLocalesFS(&localeFS)
	commons.SetSamplesFS(&samplesFS)
}
