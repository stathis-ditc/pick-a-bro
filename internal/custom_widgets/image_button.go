package custom_widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// ImageTap extends canvas.Image with tap handling
type ImageButton struct {
	widget.BaseWidget
	image *canvas.Image
	onTap func()
}

// NewImageButton creates a new ImageButton widget with the specified resource and onTap function.
// The resource parameter represents the image resource to be displayed on the button.
// The onTap parameter is a function that will be called when the button is tapped.
// The function returns a pointer to the created ImageButton widget.
func NewImageButton(resource fyne.Resource, onTap func()) *ImageButton {
	img := canvas.NewImageFromResource(resource)
	img.FillMode = canvas.ImageFillStretch

	button := &ImageButton{
		image: img,
		onTap: onTap,
	}
	button.ExtendBaseWidget(button)
	return button
}

func (b *ImageButton) Tapped(*fyne.PointEvent) {
	if b.onTap != nil {
		b.onTap()
	}
}

type imageButtonRenderer struct {
	img *canvas.Image
	obj *ImageButton
}

func (r *imageButtonRenderer) MinSize() fyne.Size {
	return fyne.NewSize(300, 200)
}

func (r *imageButtonRenderer) Layout(size fyne.Size) {
	r.img.Resize(size)
}

func (r *imageButtonRenderer) Refresh() {
	canvas.Refresh(r.obj)
}

func (r *imageButtonRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.img}
}

func (r *imageButtonRenderer) Destroy() {}

func (b *ImageButton) CreateRenderer() fyne.WidgetRenderer {
	r := &imageButtonRenderer{obj: b, img: b.image}
	return r
}
