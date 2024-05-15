package data

import (
	"encoding/json"
	"image/color"
	"math/rand"
	"os"
	"pick-a-bro/internal/commons"
)

type MembersList struct {
	PatreonMembers []PatreonMember
	Tiers          map[string]interface{}
	ColorCode      map[string]color.Color
}

var list *MembersList

// ExtractDataFromFile reads and extracts data from files based on the test mode and real data preferences.
// It checks if the members file and tiers file exist and are readable, and then generates color codes.
// Returns true if the data extraction is successful, otherwise returns false.
func ExtractDataFromFile() bool {
	testMode := commons.GetPreferences().Bool(commons.TestMode)
	useRealData := commons.GetPreferences().Bool(commons.UseRealData)
	membersFileName := commons.StructuredData.RealDataFileName
	tiersFileName := commons.StructuredData.RealTiersFileName
	if testMode && !useRealData {
		membersFileName = commons.StructuredData.TestDataFileName
		tiersFileName = commons.StructuredData.TestTiersFileName
	}

	list = &MembersList{}
	// Check if the members file exists and is readable
	if !readAndUnmarshal(membersFileName, &list.PatreonMembers) {
		return false
	}

	// Check if the tiers file exists and is readable
	if !readAndUnmarshal(tiersFileName, &list.Tiers) {
		return false
	}

	generateColorCodes()

	return true
}

func GetMembersAndTiers() *MembersList {
	return list
}

func SetMembersList(membersList []PatreonMember) {
	list.PatreonMembers = membersList
}

func SetTiersMap(tiersMap map[string]interface{}) {
	list.Tiers = tiersMap
	generateColorCodes()
}

func readAndUnmarshal(filePath string, v interface{}) bool {
	data, err := os.ReadFile(filePath)
	if err != nil {
		commons.GetLogger().Print(err)
		return false
	}

	if err := json.Unmarshal(data, v); err != nil {
		commons.GetLogger().Print(err)
		return false
	}

	return true
}

// generateColorCodes generates color codes for each tier in the list.
// It assigns predefined colors to tiers and generates random colors for any additional tiers.
// The color codes are stored in the `ColorCode` field of the list.
func generateColorCodes() {
	colors := []color.Color{
		color.RGBA{R: 0, G: 0, B: 255, A: 255},   // Blue
		color.RGBA{R: 0, G: 150, B: 0, A: 255},   // Green
		color.RGBA{R: 100, G: 100, B: 0, A: 255}, // Yellow
		color.RGBA{R: 255, G: 0, B: 0, A: 255},   // Red
	}

	// If there are more tiers than predefined colors, generate random colors for the additional tiers
	if len(list.Tiers) > len(colors) {
		for i := 0; i < len(list.Tiers)-len(colors); i++ {
			colors = append(colors, color.RGBA{R: uint8(rand.Intn(255)), G: uint8(rand.Intn(255)), B: uint8(rand.Intn(255)), A: 255})
		}
	}

	tierColors := make(map[string]color.Color)
	i := 0
	for _, tier := range list.Tiers {
		tierColors[tier.(string)] = colors[i]
		i++
	}

	list.ColorCode = tierColors
}
