package commons

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

// GetBackgroundImage returns a fyne.Container with a background image.
// The background image is retrieved using the GetAsset function from the specified AssetsPaths.ImagesPath
// and the AssetsKeys.BackgroundImg key.
// The background image is then displayed using a canvas.Image and set to fill the container using ImageFillStretch mode.
// The resulting container is created with a stack layout and contains the background image.
func GetBackgroundImage() *fyne.Container {

	backgroundImg := EmbedImage(GetAsset(AssetsPaths.ImagesPath, AssetsKeys.BackgroundImg), AssetsKeys.BackgroundImg)

	backgroundResource := canvas.NewImageFromResource(backgroundImg)
	backgroundResource.FillMode = canvas.ImageFillStretch
	return container.New(layout.NewStackLayout(), backgroundResource)
}

func GetCoundownImages() []fyne.Resource {
	return []fyne.Resource{
		EmbedImage(GetAsset(AssetsPaths.ImagesPath, AssetsKeys.Countdown1Img), AssetsKeys.Countdown1Img),
		EmbedImage(GetAsset(AssetsPaths.ImagesPath, AssetsKeys.Countdown2Img), AssetsKeys.Countdown2Img),
		EmbedImage(GetAsset(AssetsPaths.ImagesPath, AssetsKeys.Countdown3Img), AssetsKeys.Countdown3Img),
	}
}

// GetAsset returns the path of the specified asset based on the asset type and filename.
func GetAsset(assetType string, filename string) string {
	return fmt.Sprintf("%s%s%s", AssetsPaths.FilePath, assetType, AssetsPaths.Files[filename])
}
