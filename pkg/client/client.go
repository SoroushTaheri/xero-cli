package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	RoboEpicsBaseURL   = "https://staging.api.roboepics.com"
	FusionAPIBaseURL   = "https://fusion.roboepics.com/api"
	FusionOAuthBaseURL = "https://fusion.roboepics.com/oauth2"
	FusionClientID     = "7126a051-baea-4fe1-bdf8-fde2fdb31f97"
)

const (
	applicationJSON   string = "application/json"
	wwwFormURLEncoded string = "x-www-form-urlencoded"
	multipartFormData string = "multipart/form-data"
)

type Profile struct {
	Username   string `json:"username"`
	FullName   string `json:"full_name"`
	FullNameEn string `json:"full_name_english"`
	Email      string `json:"email"`
}

type Client struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	ID           string  `json:"user_id"`
	Profile      Profile `json:"profile"`
}

func New() *Client {
	return &Client{}
}

func NewWithToken(accessToken, refreshToken string, updateProfile bool) (*Client, error) {
	client := &Client{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	active, err := client.IsActive()
	if err != nil {
		return nil, fmt.Errorf("failed to check token: %v", err)
	}
	if !active {
		if err := client.RenewAuthSession(); err != nil {
			return nil, fmt.Errorf("failed to renew auth session: %v", err)
		}
	}
	if !updateProfile {
		return client, nil
	}

	incompleteProfile := func() bool {
		result := 1
		v := reflect.ValueOf(client.Profile)
		for i := 0; i < v.NumField(); i++ {
			result *= len(v.Field(i).Interface().(string))
			fmt.Printf("\n%q", v.Field(i).Interface().(string))
		}
		return result == 0
	}()

	if incompleteProfile {
		fmt.Println("Profile is BAD")
		if err := client.UpdateProfile(); err != nil {
			return nil, fmt.Errorf("failed to update user profile: %v", err)
		}
	}

	return client, nil
}

func (c *Client) Creds() (string, string) {
	return c.AccessToken, c.RefreshToken
}

type LoginCredsRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginCredsResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"RefreshToken"`
}

func (c *Client) Login(username, password string) error {
	if len(username)*len(password) == 0 {
		return fmt.Errorf("empty username/password")
	}

	buffer := &bytes.Buffer{}
	if err := json.NewEncoder(buffer).Encode(&LoginCredsRequest{
		Username: username,
		Password: password,
	}); err != nil {
		panic("failed to encode request body")
	}

	response, err := http.Post("https://staging.api.roboepics.com/account/login", "application/json", buffer)
	if err != nil {
		return fmt.Errorf("error while requesting server: %v", err)
	}

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("invalid credentials")
	}

	var loginResponse LoginCredsResponse
	if err := json.NewDecoder(response.Body).Decode(&loginResponse); err != nil {
		return fmt.Errorf("could not parse server response: %v", err)
	}

	c.AccessToken = loginResponse.Token
	c.RefreshToken = loginResponse.RefreshToken
	c.Profile.Username = username

	if err := c.UpdateProfile(); err != nil {
		return fmt.Errorf("could not get user profile: %v", err)
	}

	return nil
}

type FusionAuthGrantRefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (c *Client) RenewAuthSession() error {
	active, err := c.IsActive()
	if err != nil {
		return fmt.Errorf("failed to check token: %v", err)
	}
	if active {
		return nil
	}

	response, err := http.Post(
		fmt.Sprintf("%s/token?grant_type=refresh_token&refresh_token=%s&client_id=%s",
			FusionOAuthBaseURL, c.RefreshToken, FusionClientID),
		"application/x-www-form-urlencoded", nil)
	if err != nil {
		return fmt.Errorf("error while requesting fusionauth server: %v", err)
	}

	// 400: Invalid refresh token
	// 401: Invalid client id
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("server failed to grant new refresh token: %v", err)
	}

	var fusionResponse FusionAuthGrantRefreshTokenResponse
	if err := json.NewDecoder(response.Body).Decode(&fusionResponse); err != nil {
		return fmt.Errorf("could not parse server response: %v", err)
	}

	c.AccessToken = fusionResponse.AccessToken
	c.RefreshToken = fusionResponse.RefreshToken

	return nil
}

type FusionAuthIntrospectTokenResponse struct {
	Active bool `json:"active"`
}

func (c *Client) introspectToken() (bool, error) {
	if c.AccessToken == "" {
		return false, fmt.Errorf("invalid AccessToken")
	}

	// TODO: This endpoint is NOT safe. Token is being transmitted over URL and could easily be sniffed.
	response, err := http.Post(
		fmt.Sprintf("%s/introspect?token=%s&client_id=%s",
			FusionOAuthBaseURL, c.AccessToken, FusionClientID),
		"application/x-www-form-urlencoded", nil)

	if err != nil {
		return false, fmt.Errorf("error while requesting server: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return false, fmt.Errorf("server failed to introspect token: %v", err)
	}

	var fusionResponse FusionAuthIntrospectTokenResponse
	if err := json.NewDecoder(response.Body).Decode(&fusionResponse); err != nil {
		return false, fmt.Errorf("could not parse server response: %v", err)
	}

	return fusionResponse.Active, nil
}

type ProfileResponse struct {
	Username   string `json:"username"`
	FullName   string `json:"full_name"`
	FullNameEn string `json:"full_name_english"`
	Email      string `json:"email"`
	ID         string `json:"fusion_user_id"`
}

// Updates profile fields: email, fullName, fullNameEn, userID. Keep in mind that this method updates more data fields compared to `UpdateProfileFromToken`
func (c *Client) UpdateProfile() error {
	response, err := c.sendProtectedRequest("GET", fmt.Sprintf("%s/account/profile",
		RoboEpicsBaseURL), "application/json", nil)

	if err != nil {
		return fmt.Errorf("error while requesting server: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("server failure: %v", err)
	}

	var profileResponse ProfileResponse
	if err := json.NewDecoder(response.Body).Decode(&profileResponse); err != nil {
		return fmt.Errorf("could not parse server response: %v", err)
	}

	c.Profile.Email = profileResponse.Email
	c.Profile.FullName = profileResponse.FullName
	c.Profile.FullNameEn = profileResponse.FullNameEn
	c.Profile.Username = profileResponse.Username
	c.ID = profileResponse.ID

	return nil
}

// Updates profile fields: email, username, userID via decoding JWT token.
func (c *Client) UpdateProfileFromToken() error {
	claims, err := c.decodeToken()
	if err != nil {
		return fmt.Errorf("failed to parse jwt token: %v", err)
	}

	c.Profile.Email = claims["email"].(string)
	c.Profile.Username = claims["preferred_username"].(string)
	c.ID = claims["aud"].(string)

	return nil
}

func (c *Client) sendProtectedRequest(method string, url string, contentType string, body io.Reader) (*http.Response, error) {
	if c.AccessToken == "" {
		return nil, fmt.Errorf("invalid AccessToken")
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate request")
	}

	req.Header = http.Header{
		"Content-Type":  {contentType},
		"Authorization": {fmt.Sprintf("Bearer %s", c.AccessToken)},
	}

	return http.DefaultClient.Do(req)
}

func (c *Client) IsActive() (bool, error) {
	if !c.IsLoggedIn() {
		return false, fmt.Errorf("user is not logged in")
	}

	claims, err := c.decodeToken()

	if err != nil {
		return false, fmt.Errorf("failed to parse jwt token: %v", err)
	}

	return time.Now().Unix() < int64(claims["exp"].(float64)), nil
}

func (c *Client) decodeToken() (jwt.MapClaims, error) {
	if c.AccessToken == "" {
		return nil, fmt.Errorf("invalid AccessToken")
	}

	claims := jwt.MapClaims{}
	// TODO: it's honestly better to verify but not for now. CTF is 11 days ahead. we don't have time for this shit.
	// +1 by faramin.smile

	// _, err := jwt.ParseWithClaims(c.AccessToken, &claims, func(token *jwt.Token) (interface{}, error) {
	// 	return []byte("123"), nil
	// })
	_, _, err := new(jwt.Parser).ParseUnverified(c.AccessToken, &claims)

	return claims, err
}

func (c *Client) IsLoggedIn() bool {
	return len(c.AccessToken)*len(c.RefreshToken) != 0
}

type CompetitionRetrieveResponse struct {
	TotalParticipants int      `json:"total_participants"`
	TotalSubmissions  int      `json:"total_submissions"`
	Status            int      `json:"status"`
	ParticipantID     int      `json:"participant_id"`
	Tags              []string `json:"tags"`
	Title             string   `json:"title"`
	Subtitle          string   `json:"subtitle"`
	Path              string   `json:"path"`
	Rules             string   `json:"rules"`
	Prize             string   `json:"prize"`
	PhaseSet          []Phase  `json:"phase_set"`
}

type Phase struct {
	ID                  int            `json:"id"`
	Problems            []PhaseProblem `json:"problems"`
	Title               string         `json:"title"`
	Description         string         `json:"description"`
	HideUntilStart      bool           `json:"hide_until_start"`
	ProgressRestriction int            `json:"progress_restriction"`
	EligibleScore       string         `json:"eligible_score"`
	DateCreated         string         `json:"date_created"`
}

type PhaseProblem struct {
	ID            int     `json:"id"`
	Weight        float32 `json:"weight"`
	EligibleScore string  `json:"eligible_score"`
	Problem       Problem `json:"problem"`
}

type Problem struct {
	ID      int    `json:"id"`
	Path    string `json:"path"`
	OwnerID string `json:"owner"`
	Title   string `json:"title"`
}

func (c *Client) GetCompetition(path string) (CompetitionRetrieveResponse, error) {
	var (
		httpResponse *http.Response
		err          error
	)

	fullPath := fmt.Sprintf("%s/competition/%s", RoboEpicsBaseURL, path)

	if c.AccessToken == "" {
		httpResponse, err = http.Get(fullPath)
	} else {
		httpResponse, err = c.sendProtectedRequest("GET", fullPath, applicationJSON, nil)
	}

	if err != nil {
		return CompetitionRetrieveResponse{}, fmt.Errorf("error while requesting server: %v", err)
	}

	if httpResponse.StatusCode != http.StatusOK {
		return CompetitionRetrieveResponse{}, fmt.Errorf("invalid response")
	}

	var response CompetitionRetrieveResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return CompetitionRetrieveResponse{}, fmt.Errorf("could not parse server response: %v", err)
	}

	return response, nil
}

type ProblemRetrieveResponse struct {
	ID          int    `json:"id"`
	Description string `json:"short_description"`
	Title       string `json:"title"`
}

func (c *Client) GetProblem(path string) (ProblemRetrieveResponse, error) {
	fullPath := fmt.Sprintf("%s/problem/%s", RoboEpicsBaseURL, path)

	httpResponse, err := http.Get(fullPath)
	if err != nil {
		return ProblemRetrieveResponse{}, fmt.Errorf("error while requesting server: %v", err)
	}

	if httpResponse.StatusCode != http.StatusOK {
		return ProblemRetrieveResponse{}, fmt.Errorf("invalid response")
	}

	var response ProblemRetrieveResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return ProblemRetrieveResponse{}, fmt.Errorf("could not parse server response: %v", err)
	}

	return response, nil
}

type ProblemTextsListResponse []ProblemText

type ProblemText struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Text        string `json:"text"`
	ContentType int    `json:"content_type"`
	Order       int    `json:"order"`
	Problem     int    `json:"problem"`
}

func (c *Client) GetProblemTexts(path string) (ProblemTextsListResponse, error) {
	fullPath := fmt.Sprintf("%s/problem/%s/texts", RoboEpicsBaseURL, path)

	httpResponse, err := http.Get(fullPath)
	if err != nil {
		return ProblemTextsListResponse{}, fmt.Errorf("error while requesting server: %v", err)
	}

	if httpResponse.StatusCode != http.StatusOK {
		return ProblemTextsListResponse{}, fmt.Errorf("invalid response")
	}

	var response ProblemTextsListResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return ProblemTextsListResponse{}, fmt.Errorf("could not parse server response: %v", err)
	}

	return response, nil
}

func (c *Client) GetProblemText(id string) (ProblemText, error) {
	fullPath := fmt.Sprintf("%s/problem/text/%s", RoboEpicsBaseURL, id)

	httpResponse, err := http.Get(fullPath)
	if err != nil {
		return ProblemText{}, fmt.Errorf("error while requesting server: %v", err)
	}

	if httpResponse.StatusCode != http.StatusOK {
		return ProblemText{}, fmt.Errorf("invalid response")
	}

	var response ProblemText
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return ProblemText{}, fmt.Errorf("could not parse server response: %v", err)
	}

	return response, nil
}

type AddParticipantResponse struct {
	ParticipantID int `json:"id"`
}

func (c *Client) AddParticipant(id string) (AddParticipantResponse, error) {
	fullPath := fmt.Sprintf("%s/competition/%s/add-participant", RoboEpicsBaseURL, id)

	httpResponse, err := http.Get(fullPath)
	if err != nil {
		return AddParticipantResponse{}, fmt.Errorf("error while requesting server: %v", err)
	}

	if httpResponse.StatusCode != http.StatusOK {
		return AddParticipantResponse{}, fmt.Errorf("invalid response: %d", httpResponse.StatusCode)
	}

	var response AddParticipantResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return AddParticipantResponse{}, fmt.Errorf("could not parse server response: %v", err)
	}

	return response, nil
}

type LeaderboardRow struct {
	TeamName        string   `json:"team_name"`
	CreatorFullName string   `json:"creator_full_name"`
	Score           float64  `json:"score"`
	Individual      bool     `json:"individual"`
	LastSubmission  string   `json:"last_submission"`
	SubmissionDates []string `json:"submission_dates"`
}

type GetLeaderboardResponse struct {
	Results       []LeaderboardRow `json:"results"`
	Total         int              `json:"total"`
	EligibleScore float64          `json:"eligible_score"`
}

func (c *Client) GetLeaderboard(problemPath string, start, end int) (GetLeaderboardResponse, error) {
	fullPath := fmt.Sprintf("%s/problem/%s/leaderboard?start=%d&end=%d", RoboEpicsBaseURL, problemPath, start, end)

	httpResponse, err := http.Get(fullPath)
	if err != nil {
		return GetLeaderboardResponse{}, fmt.Errorf("error while requesting server: %v", err)
	}

	if httpResponse.StatusCode != http.StatusOK {
		responseBody, _ := ioutil.ReadAll(httpResponse.Body)
		return GetLeaderboardResponse{}, fmt.Errorf("invalid response: %d %s", httpResponse.StatusCode, responseBody)
	}

	var response GetLeaderboardResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return GetLeaderboardResponse{}, fmt.Errorf("could not parse server response: %v", err)
	}

	return response, nil
}
