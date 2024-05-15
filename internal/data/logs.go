package data

import (
	"log"
	"os"
)

// GetLogs retrieves the contents of the log file and returns it as a string.
// If an error occurs while reading the file, it returns an empty string and the error.
func GetLogs() (string, error) {
	data, err := os.ReadFile("log.txt")
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return string(data), nil
}
