# Pick A Bro

Pick a Bro is a draw app, primarily created for the youtube show broscar. The app fetches Patreon members of a connected acoount and performs a draw, randomly selecting a winner

The app is written in Go with support of Fyne library for GUI

## Features
- Fetch Patreon Members using Patron's API
- Dynamically design draw rectangles
- Customize draw settings
- Available in Greek and English

## Requirements
To build the app Go 1.21+ is required. 

## Run the app
To run the app, navigate to `cmd/pick-a-bro` and run ```go run main.go```

## Build
To build the app navoigate to cmd/pick-a-bro and run ```go build```. That is enough for mac/linux machines
If running from a windows machine or the build must be an exe file for windows run ```CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -v -o pickabro.exe``` 
or ```CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -v -o pickabro.exe``` if building for a mac/linux machine

## License
This project is licensed under the MIT License - see the [LICENSE](https://github.com/devs-in-the-cloud/pick-a-bro/LICENSE) file for details.

## Commercial Use
Commercial use of this software requires explicit written permission from the author. Please contact the author if you wish to use this software commercially.

