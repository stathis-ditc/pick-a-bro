package views

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"pick-a-bro/internal/commons"
	"pick-a-bro/internal/data"
	"pick-a-bro/internal/lottery"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

// SetLottery initializes the lottery view and sets up the audio player.
func SetLottery(window fyne.Window) {
	const sampleRate = 44100
	speaker.Init(sampleRate, sampleRate/10)
	lotteryView(window)
}

// lotteryView is a function that creates and displays the lottery view in the application window.
// It takes a fyne.Window as a parameter and initializes the members list for the lottery.
// It creates rectangles for each Patreon member in the members list and adds them to the view.
// It also adds an overlay image to the view.
// The function then sets the content of the window to the created view and starts the lottery process in a separate goroutine.
func lotteryView(window fyne.Window) {
	membersList, _ := lottery.InitMembersList()

	var rectangles []fyne.CanvasObject

	for i := 0; i < len(membersList.PatreonMembers); i++ {
		rect := createRectangle(membersList.ColorCode, membersList.PatreonMembers[i])
		rectangles = append(rectangles, rect)
	}

	overlay := canvas.NewImageFromResource(commons.EmbedImage(
		commons.GetAsset(commons.AssetsPaths.ImagesPath, commons.AssetsKeys.DrawOverlayImg), commons.AssetsKeys.DrawOverlayImg))
	overlay.FillMode = canvas.ImageFillStretch
	rectangles = append(rectangles, overlay)

	columns := int(commons.WindowWidth) / 100
	content := container.NewVScroll(container.NewGridWithColumns(columns, rectangles...))

	go runLottery(rectangles, overlay, membersList.PatreonMembers, window, content)

	window.SetContent(content)
}

// createRectangle creates a rectangle with the specified color based on the member's tier and adds a label with the member's full name.
// It returns a fyne.CanvasObject that contains the rectangle and label.
func createRectangle(colors map[string]color.Color, member data.PatreonMember) fyne.CanvasObject {
	rect := canvas.NewRectangle(colors[member.Tier])
	rect.SetMinSize(fyne.NewSize(50, 20))
	text := widget.NewLabel(member.FullName)
	text.Alignment = fyne.TextAlignCenter
	text.Wrapping = fyne.TextWrapWord

	rectContainer := container.NewStack(rect, text)

	return rectContainer
}

// runLottery runs the lottery process by animating the countdown, selecting random rectangles,
// playing a beep sound, and displaying the winner dialog.
// It takes the following parameters:
// - rectangles: a slice of fyne.CanvasObject representing the rectangles to select from.
// - overlay: a pointer to a canvas.Image representing the overlay image.
// - membersList: a slice of data.PatreonMember representing the list of members.
// - window: a fyne.Window representing the application window.
// - content: a pointer to a container.Scroll representing the scrollable content.
// The function does not return any value.
func runLottery(rectangles []fyne.CanvasObject, overlay *canvas.Image, membersList []data.PatreonMember, window fyne.Window, content *container.Scroll) {
	buffer1, _, err := loadMP3ToBuffer(commons.GetAsset(commons.AssetsPaths.AudioPath, commons.AssetsKeys.BeepAudio))
	if err != nil {
		fmt.Println(err)
	}

	animatedImage := canvas.NewImageFromResource(commons.GetCoundownImages()[2])
	animatedImage.SetMinSize(fyne.NewSize(300, 300))
	countdown := dialog.NewCustomWithoutButtons(commons.GetTranslation(commons.I18n.Ready), animatedImage, window)
	countdown.Show()
	var wg sync.WaitGroup // Declare a WaitGroup

	wg.Add(1) // Increment the WaitGroup counter.
	go animateImage(&wg, animatedImage)
	wg.Wait()
	countdown.Hide()

	var randomNumber int

	for i := 0; i < 40; i++ {
		randomNumber = rand.Intn(len(rectangles))
		if rectangles[randomNumber].Position().Y > content.Offset.Y+600 || rectangles[randomNumber].Position().Y < content.Offset.Y {
			content.Offset = fyne.NewPos(rectangles[randomNumber].Position().X, rectangles[randomNumber].Position().Y-float32(rand.Intn(300)))
			content.Refresh()
		}
		overlay.Resize(fyne.NewSize(rectangles[randomNumber].Size().Width, rectangles[randomNumber].Size().Height))
		overlay.Move(fyne.NewPos(rectangles[randomNumber].Position().X, rectangles[randomNumber].Position().Y))
		beepStream := buffer1.Streamer(0, buffer1.Len())
		done := make(chan bool)
		speaker.Play(beep.Seq(beepStream, beep.Callback(func() {
			done <- true
		})))
		<-done
		content.Refresh()
		switch {
		case i >= 20 && i < 30:
			time.Sleep(time.Millisecond * 200)
		case i >= 30 && i < 35:
			time.Sleep(time.Millisecond * 400)
		case i >= 35 && i < 37:
			time.Sleep(time.Millisecond * 600)
		case i >= 37 && i < 38:
			time.Sleep(time.Millisecond * 800)
		case i >= 38:
			time.Sleep(time.Millisecond * 900)
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}

	showWinnerDialog(membersList[randomNumber].FullName, window)
}

// showWinnerDialog displays a dialog box to congratulate the winner and play a winner audio.
// It takes the winner's name and a fyne.Window as parameters.
// The function loads an MP3 audio file, plays the audio, and creates a dialog box with a congratulatory message.
// The dialog box is then shown to the user.
// After the dialog box is closed, the function adds the winner's name to the winners list (if not in test mode)
// and returns to the main menu.
func showWinnerDialog(winnerName string, window fyne.Window) {
	buffer, _, err := loadMP3ToBuffer(commons.GetAsset(commons.AssetsPaths.AudioPath, commons.AssetsKeys.WinnerAudio))
	if err != nil {
		commons.GetLogger().Fatal(err)
	}

	winnerStream := buffer.Streamer(0, buffer.Len())
	speaker.Play(winnerStream)
	congratsLabel := widget.NewLabel(fmt.Sprintf(commons.GetTranslation(commons.I18n.Congrats), winnerName))
	dialogContent := container.NewHBox(congratsLabel)

	winnersDialog := dialog.NewCustom(commons.GetTranslation(commons.I18n.Winner), "Done", dialogContent, window)
	winnersDialog.Resize(fyne.NewSize(200, 200))
	winnersDialog.Show()
	winnersDialog.SetOnClosed(func() {
		if !commons.GetPreferences().Bool(commons.Settings.TestMode) {
			lottery.AddToWinnersList(winnerName)
		}
		MainMenu(window)
	})
}

// loadMP3ToBuffer loads an MP3 file from the specified filePath and returns a buffer, format, and error.
func loadMP3ToBuffer(filePath string) (*beep.Buffer, beep.Format, error) {
	audioFS, err := commons.GetAudioFS().Open(filePath)
	if err != nil {
		log.Fatalf("failed to open embedded audio file: %v", err)
	}
	defer audioFS.Close()

	streamer, format, err := mp3.Decode(audioFS)
	if err != nil {
		return nil, beep.Format{}, fmt.Errorf("failed to decode file: %w", err)
	}
	defer streamer.Close()

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)

	return buffer, format, nil
}

// animateImage animates an image by cycling through a list of images and playing a beep sound.
// It takes a wait group `wg` and an image `img` as parameters.
// The function loads an MP3 audio file, plays a beep sound, updates the image resource, refreshes the image, and then waits for a second before cycling to the next image.
// The function uses synchronization with the wait group to indicate when it has finished.
func animateImage(wg *sync.WaitGroup, img *canvas.Image) {
	// List of images to cycle through
	defer wg.Done()

	buffer1, _, _ := loadMP3ToBuffer(commons.GetAsset(commons.AssetsPaths.AudioPath, commons.AssetsKeys.BeepAudio))

	for i := 2; i >= 0; i-- {
		beepStream := buffer1.Streamer(0, buffer1.Len())
		done := make(chan bool)
		speaker.Play(beep.Seq(beepStream, beep.Callback(func() {
			done <- true
		})))
		<-done

		img.Resource = commons.GetCoundownImages()[i]
		img.Refresh()
		time.Sleep(time.Second)
	}
}
