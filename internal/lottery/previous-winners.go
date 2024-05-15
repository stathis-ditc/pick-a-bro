package lottery

import (
	"encoding/json"
	"os"
	"pick-a-bro/internal/commons"
	"time"
)

type Winner struct {
	FullName string
	DateTime string
}

type Winners struct {
	Winners []Winner `json:"winners"`
}

// AddToWinnersList adds a new winner to the list of previous winners.
// It takes the name of the winner as a parameter and appends it to the list.
// The updated list is then written back to the file.
func AddToWinnersList(winner string) {
	newWinner := createWinner(winner)
	winners := readWinnersFromFile()

	winners.Winners = append(winners.Winners, newWinner)

	writeWinnersToFile(winners)
}

func GetWinnersList() []Winner {
	winners := readWinnersFromFile()

	return winners.Winners
}

func createWinner(name string) Winner {
	return Winner{FullName: name, DateTime: time.Now().Format("02/01/2006 15:04:05")}
}

// readWinnersFromFile reads the previous winners from a file and returns them.
// If the file does not exist or there is an error reading the file, an empty
// Winners struct is returned.
func readWinnersFromFile() Winners {
	filename := commons.StructuredData.WinnersFileName

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		commons.GetLogger().Fatalf("Failed opening file: %s", err)
	}
	defer file.Close()

	var winners Winners
	err = json.NewDecoder(file).Decode(&winners)
	if err != nil {
		winners.Winners = []Winner{}
	}

	return winners
}

// writeWinnersToFile writes the given winners data to a file.
// It takes a parameter `winners` of type `Winners`, which represents the winners data to be written.
// The function marshals the winners data into JSON format and writes it to the file specified by the `WinnersFileName` field in the `StructuredData` struct.
// If any error occurs during the marshaling or writing process, the function logs a fatal error using the logger from the `commons` package.
func writeWinnersToFile(winners Winners) {
	filename := commons.StructuredData.WinnersFileName

	jsonData, err := json.Marshal(winners)
	if err != nil {
		commons.GetLogger().Fatalf("JSON marshaling failed: %s", err)
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		commons.GetLogger().Fatalf("Failed writing to file: %s", err)
	}
}

// ClearWinnersList removes the file containing the list of previous winners.
func ClearWinnersList() {
	filename := commons.StructuredData.WinnersFileName
	err := os.Remove(filename)
	if err != nil {
		commons.GetLogger().Fatalf("Failed removing file: %s", err)
	}
}
