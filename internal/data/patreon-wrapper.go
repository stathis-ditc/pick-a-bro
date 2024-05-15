package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"pick-a-bro/internal/commons"
	"runtime"
	"time"

	"github.com/austinbspencer/patreon-go-wrapper"
	"golang.org/x/oauth2"
)

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type PatreonMember struct {
	FullName string
	Tier     string
}

var token *oauth2.Token
var client *patreon.Client

// newPatreonClient initializes a new Patreon client and establishes a connection to the API.
// It starts an authentication server, fetches an access token using the received code,
// and creates a Patreon client with the obtained token.
// Returns true if the client is successfully created, otherwise false.
func newPatreonClient() bool {
	code := startAuthServer()
	if code == "" {
		commons.GetLogger().Println("No code received. Unable to establish connection to api")
		return false
	}
	token = fetchToken(code)
	client = createPatreonClient(token)
	return true
}

// startAuthServer starts an HTTP server on port 8080 and waits for a request with a "code" query parameter.
// It opens the authentication URL in the default browser and returns the received code.
func startAuthServer() string {
	done := make(chan bool)
	var code string

	srv := &http.Server{Addr: ":8080"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code = r.URL.Query().Get("code")
		done <- true
	})

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			commons.GetLogger().Println(err)
		}
	}()

	openAuthURLInBrowser()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		commons.GetLogger().Println(err)
	}
	return code
}

// openAuthURLInBrowser opens the authentication URL in the default web browser.
// It constructs the authentication URL using the provided client ID and redirect URI.
// The URL is then opened in the default web browser based on the operating system.
// If an error occurs while opening the URL, it is logged and a fatal error is raised.
func openAuthURLInBrowser() {
	authParams := url.Values{
		"response_type": {"code"},
		"client_id":     {commons.GetPreferences().String(commons.ClientId)},
		"redirect_uri":  {commons.RedirectURI},
	}

	authenticationUrl := fmt.Sprintf("%s?%s", patreon.AuthorizationURL, authParams.Encode())

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", authenticationUrl)
	case "darwin":
		cmd = exec.Command("open", authenticationUrl)
	default:
		cmd = exec.Command("xdg-open", authenticationUrl)
	}

	if err := cmd.Start(); err != nil {
		commons.GetLogger().Fatalf("error: %v", err)
	}
}

// fetchToken fetches an OAuth2 token using the provided authorization code.
// It sends a POST request to the Patreon API to exchange the code for a token.
// The token is then parsed from the response and returned as an oauth2.Token.
// If any error occurs during the process, it will be logged and the function will return nil.
func fetchToken(code string) *oauth2.Token {
	data := url.Values{
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"client_id":     {commons.GetPreferences().String(commons.ClientId)},
		"client_secret": {commons.GetPreferences().String(commons.ClientSecret)},
		"redirect_uri":  {commons.RedirectURI},
	}

	resp, err := http.PostForm(patreon.AccessTokenURL, data)
	if err != nil {
		commons.GetLogger().Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		commons.GetLogger().Fatalf("error: %v", err)
	}

	var respOAuthToken AccessTokenResponse
	if err := json.Unmarshal(body, &respOAuthToken); err != nil {
		commons.GetLogger().Fatalf("error: %v", err)
	}

	return &oauth2.Token{
		AccessToken:  respOAuthToken.AccessToken,
		RefreshToken: respOAuthToken.RefreshToken,
		Expiry:       time.Now().Add(2 * time.Hour),
	}
}

// createPatreonClient creates a new Patreon client using the provided OAuth2 token.
// It returns the initialized Patreon client.
func createPatreonClient(token *oauth2.Token) *patreon.Client {
	config := oauth2.Config{
		ClientID:     commons.GetPreferences().String(commons.ClientId),
		ClientSecret: commons.GetPreferences().String(commons.ClientSecret),
		Endpoint: oauth2.Endpoint{
			AuthURL:  patreon.AuthorizationURL,
			TokenURL: patreon.AccessTokenURL,
		},
		Scopes: patreon.AllScopes,
	}

	tc := config.Client(context.Background(), token)

	return patreon.NewClient(tc)
}

func FetchMembersToLocalStorage() ([]PatreonMember, map[string]interface{}) {
	testMode := commons.GetPreferences().Bool(commons.TestMode)

	if !token.Valid() && !testMode {
		if !newPatreonClient() {
			return nil, nil
		}
	}

	members, tiers, err := getMembers(testMode)
	if err != nil {
		commons.GetLogger().Println(err)
	}

	return members, tiers
}

func getFilePath(testFileName string, realFileName string) string {
	if commons.GetPreferences().Bool(commons.TestMode) {
		return testFileName
	}
	return realFileName
}

func getMembers(testMode bool) ([]PatreonMember, map[string]interface{}, error) {
	useRealData := commons.GetPreferences().Bool(commons.UseRealData)
	membersFilePath := getFilePath(commons.StructuredData.TestDataFileName, commons.StructuredData.RealDataFileName)
	tiersFilePath := getFilePath(commons.StructuredData.TestTiersFileName, commons.StructuredData.RealTiersFileName)
	if testMode && !useRealData {
		return getTestMembers(membersFilePath, tiersFilePath)
	}
	return fetchAndProcessRealMembers(membersFilePath, tiersFilePath)
}

// getTestMembers reads a sample data file containing Patreon members and their tier information,
// processes this data, and writes the results to specified files.
//
// Parameters:
// - membersFilePath: The path to the file where the processed list of members will be saved.
// - tiersFilePath: The path to the file where the processed tiers' details will be saved.
//
// Returns:
//   - A slice of PatreonMember, representing the members parsed from the sample data.
//   - A map[string]interface{}, representing the tiers' details as parsed from the sample data.
//     The structure of this map depends on the processing logic of the getTiersMap function.
//   - An error, which is non-nil if any errors occurred during the execution of the function.
func getTestMembers(membersFilePath string, tiersFilePath string) ([]PatreonMember, map[string]interface{}, error) {
	data, err := commons.GetSamplesFS().ReadFile("tests/samples/patreons.json")
	if err != nil {
		commons.GetLogger().Fatalf("failed to read samples file: %v", err)
	}

	var membersResp *patreon.MembersResponse
	if err := json.Unmarshal(data, &membersResp); err != nil {
		return nil, nil, err
	}
	commons.GetLogger().Print("Members test data generated")

	tiersMap := getTiersMap(membersResp)
	writeToFile(tiersFilePath, tiersMap)

	membersList := getMembersList(membersResp, tiersMap)
	writeToFile(membersFilePath, membersList)

	return membersList, tiersMap, nil
}

// fetchAndProcessRealMembers fetches members from a Patreon campaign and processes their information,
// including tiers they are entitled to. It writes the members' information and tiers' details to specified files.
//
// Parameters:
// - membersFilePath: The path to the file where the list of members will be saved.
// - tiersFilePath: The path to the file where the tiers' details will be saved.
//
// Returns:
//   - A slice of PatreonMember, representing the campaign members fetched from Patreon.
//   - A map[string]interface{}, representing the tiers' details. The exact structure of the map depends on
//     the structure of the tiers information returned by the Patreon API and processed by getTiersMap function.
//   - An error, which is non-nil if any errors occurred during the function's execution.
func fetchAndProcessRealMembers(membersFilePath string, tiersFilePath string) ([]PatreonMember, map[string]interface{}, error) {
	membersList := []PatreonMember{}
	var nextCursor string
	var tiersMap map[string]interface{}

	for {
		membersResp, err := client.FetchCampaignMembers(commons.GetPreferences().String(commons.CampaignId),
			patreon.WithIncludes("currently_entitled_tiers"),
			patreon.WithFields("member", commons.MemberFields...),
			patreon.WithFields("tier", commons.TierFields...),
			patreon.WithCursor(nextCursor),
		)

		if err != nil {
			return nil, nil, err
		}

		if tiersMap == nil {
			tiersMap = getTiersMap(membersResp)
			writeToFile(tiersFilePath, tiersMap)
		}

		membersList = append(membersList, getMembersList(membersResp, tiersMap)...)

		nextCursor = membersResp.Meta.Pagination.Cursors.Next
		if nextCursor == "" {
			break
		}
	}
	writeToFile(membersFilePath, membersList)
	return membersList, tiersMap, nil
}

// getTiersMap returns a map of Patreon tier IDs to their corresponding titles.
func getTiersMap(members *patreon.MembersResponse) map[string]interface{} {
	tiersMap := make(map[string]interface{})
	for _, item := range members.Included.Items {
		if tier, ok := item.(*patreon.Tier); ok {
			tiersMap[tier.ID] = tier.Attributes.Title
		}
	}
	return tiersMap
}

// getMembersList retrieves a list of active and paid Patreon members from the given `members` response
// and maps them to a slice of `PatreonMember` structs. It filters out members who are not active patrons
// or whose last charge status is not "Paid".
//
// Parameters:
// - members: A pointer to a `patreon.MembersResponse` object containing the Patreon members data.
// - tiersMap: A map of Patreon tier IDs to their corresponding names.
//
// Returns:
// - A slice of `PatreonMember` structs representing the active and paid Patreon members.
func getMembersList(members *patreon.MembersResponse, tiersMap map[string]interface{}) []PatreonMember {
	membersList := []PatreonMember{}
	for _, member := range members.Data {
		if member.Attributes.PatronStatus == "active_patron" && member.Attributes.LastChargeStatus == "Paid" {
			membersList = append(membersList, PatreonMember{
				FullName: member.Attributes.FullName,
				Tier:     tiersMap[member.Relationships.CurrentlyEntitledTiers.Data[0].ID].(string),
			})
		}
	}
	return membersList
}

func writeToFile(filePath string, data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		commons.GetLogger().Fatalf("error %v", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		commons.GetLogger().Fatalf("error %v", err)
	}
	commons.GetLogger().Printf("%s generated", filePath)
}
