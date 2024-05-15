package commons

// Window
var WindowWidth = float32(800)
var WindowHeight = float32(600)

// Labels
var Fellowship = "The Fellowship"
var PickABro = "Pick A Bro"

// List of constants used both for labels and preferences saving
var ClientId = "Client ID"
var ClientSecret = "Client Secret"
var AccessToken = "Access Token"
var ExcludeWinners = "Exclude previous winners"
var RefreshToken = "Refresh Token"
var CampaignId = "Campaign Id"

// Patreon variables
var RedirectURI = "http://localhost:8080"
var WithIncludes = "currently_entitled_tiers"
var MemberFields = []string{"full_name", "patron_status", "last_charge_status"}
var TierFields = []string{"title"}

// Preferences keys
var ChancesRule = "chancesRule"
var ChancesPerUser = "chancesPerUser"
var NumberOfWinners = "numberOfWinners"
var TestMode = "testMode"
var UseRealData = "useRealData"

// Lists
var ChancesRules = []string{I18n.AllEqualChances, I18n.ChancesByTier}

// JSON file names
var StructuredData = struct {
	OutputPath        string
	RealDataFileName  string
	TestDataFileName  string
	RealTiersFileName string
	TestTiersFileName string
	WinnersFileName   string
}{
	OutputPath:        "structured_data/",
	RealDataFileName:  "eligle_patreons.json",
	TestDataFileName:  "eligle_patreons_test.json",
	RealTiersFileName: "tiers.json",
	TestTiersFileName: "tiers_test.json",
	WinnersFileName:   "winners.json",
}

// Assets
var AssetsPaths = struct {
	FilePath   string
	AudioPath  string
	ImagesPath string
	Files      map[string]string
}{
	FilePath:   "assets/",
	AudioPath:  "audio/",
	ImagesPath: "images/",
	Files: map[string]string{
		AssetsKeys.Countdown1Img:  "1.png",
		AssetsKeys.Countdown2Img:  "2.png",
		AssetsKeys.Countdown3Img:  "3.png",
		AssetsKeys.BeepAudio:      "beep.mp3",
		AssetsKeys.WinnerAudio:    "winner.mp3",
		AssetsKeys.BackgroundImg:  "background.png",
		AssetsKeys.ElHandshakeImg: "el_handshake.png",
		AssetsKeys.EnHandshakeImg: "en_handshake.png",
		AssetsKeys.DrawOverlayImg: "overlay.png",
	},
}

var AssetsKeys = struct {
	BackgroundImg  string
	BeepAudio      string
	Countdown1Img  string
	Countdown2Img  string
	Countdown3Img  string
	ElHandshakeImg string
	EnHandshakeImg string
	DrawOverlayImg string
	WinnerAudio    string
}{
	Countdown1Img:  "1",
	Countdown2Img:  "2",
	Countdown3Img:  "3",
	BeepAudio:      "beep",
	WinnerAudio:    "winner",
	BackgroundImg:  "background",
	ElHandshakeImg: "el_handshake",
	EnHandshakeImg: "en_handshake",
	DrawOverlayImg: "overlay",
}

var LocalePaths = struct {
	Greek   string
	English string
}{
	Greek:   "locale/el-GR.json",
	English: "locale/en-US.json",
}

// Settings keys
var Settings = struct {
	TestMode    string
	UseRealData string
}{
	TestMode:    "testMode",
	UseRealData: "useRealData",
}

// Translation keys

var I18n = struct {
	AllEqualChances       string
	Cancel                string
	ChancesByTier         string
	ChancesPerPatreon     string
	ClearWinners          string
	Close                 string
	ConfirmClearWinners   string
	Congrats              string
	ExcludeWinners        string
	ErrorFetchingPatreons string
	FetchingPatreons      string
	MissingData           string
	NewDraw               string
	No                    string
	NoPatreons            string
	PatreonsList          string
	PreviousWinners       string
	ReadLogs              string
	Ready                 string
	RefreshPatreonsList   string
	Settings              string
	Success               string
	SuccessfulReceive     string
	TestData              string
	TestDataGenerated     string
	TestDummyData         string
	TestMode              string
	TestModeWrn           string
	TestRealData          string
	Yes                   string
	Winner                string
	WinnersCleared        string
	WinnersListCleared    string
}{
	AllEqualChances:       "all_equal_chances",
	Cancel:                "cancel",
	ChancesByTier:         "chances_by_tier",
	ChancesPerPatreon:     "chances_per_patreon",
	ClearWinners:          "clear_winners",
	Close:                 "close",
	ConfirmClearWinners:   "confirm_clear_winners",
	Congrats:              "congratulations",
	ExcludeWinners:        "exclude_winners",
	ErrorFetchingPatreons: "error_fetching_patreons",
	FetchingPatreons:      "fetching_patreons",
	MissingData:           "missing_data",
	NewDraw:               "new_draw",
	No:                    "no",
	NoPatreons:            "no_patreons_found",
	PatreonsList:          "patreons_list",
	PreviousWinners:       "previous_winners",
	ReadLogs:              "read_logs",
	Ready:                 "ready",
	RefreshPatreonsList:   "refresh_patreons_list",
	Settings:              "settings",
	Success:               "success",
	SuccessfulReceive:     "succsfull_received_patreons",
	TestData:              "test_data",
	TestDataGenerated:     "test_data_generated",
	TestDummyData:         "test_dummy_data",
	TestMode:              "test_mode",
	TestModeWrn:           "test_mode_warning",
	TestRealData:          "test_real_data",
	Yes:                   "yes",
	Winner:                "winner",
	WinnersCleared:        "winners_cleared",
	WinnersListCleared:    "winners_list_cleared",
}
