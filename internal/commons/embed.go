package commons

import (
	"embed"

	"fyne.io/fyne/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var imagesFS *embed.FS
var audioFS *embed.FS
var localesFS *embed.FS
var samplesFS *embed.FS

// EmbedImage embeds an image file located at the specified path and returns a fyne.Resource.
// The imgName parameter is used to set the name of the embedded image resource.
// If the image file cannot be read or an error occurs, the function will log a fatal error.
func EmbedImage(path string, imgName string) fyne.Resource {
	imgData, err := GetImagesFS().ReadFile(path)
	if err != nil {
		GetLogger().Fatal(err)
	}

	return fyne.NewStaticResource(imgName, imgData)
}

// EmbedLocales embeds translation files into the provided i18n.Bundle.
// It reads the translation files from the "locale" directory in the embedded filesystem,
// loads them into the bundle, and parses them.
// If any error occurs during the process, it logs a fatal error.
func EmbedLocales(bundle *i18n.Bundle) {
	files, err := GetLocalesFS().ReadDir("locale")
	if err != nil {
		GetLogger().Fatalf("failed to list translation files: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			// Read the content of the embedded file
			data, err := GetLocalesFS().ReadFile("locale/" + file.Name())
			if err != nil {
				GetLogger().Fatalf("failed to read translation file: %v", err)
			}

			// Load the translation file into the bundle
			_, err = bundle.ParseMessageFileBytes(data, file.Name())
			if err != nil {
				GetLogger().Fatalf("failed to parse translation file: %v", err)
			}
		}
	}
}

func GetImagesFS() *embed.FS {
	return imagesFS
}

func SetImagesFS(fs *embed.FS) {
	imagesFS = fs
}

func GetAudioFS() *embed.FS {
	return audioFS
}

func SetAudioFS(fs *embed.FS) {
	audioFS = fs
}

func GetLocalesFS() *embed.FS {
	return localesFS
}

func SetLocalesFS(fs *embed.FS) {
	localesFS = fs
}

func GetSamplesFS() *embed.FS {
	return samplesFS
}

func SetSamplesFS(fs *embed.FS) {
	samplesFS = fs
}
