package msal

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	// DefaultContentType defined as per https://docs.microsoft.com/en-us/graph/auth-v2-service
	DefaultContentType = "application/x-www-form-urlencoded"
	// RequiredGrantType must be client_credentials
	RequiredGrantType = "client_credentials"
)

// Token contains the access token and expiration
type Token struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	TokenType    string `json:"token_type"`
}

// ClientApplication represents the client application
type ClientApplication struct {
	TenantID     string
	ClientID     string
	clientSecret string
	Scope        string
	BaseURL      string
	httpClient   *http.Client
}

// NewClientApplication returns and initialized OAuth client
func NewClientApplication(tenantID string, clientID string, clientSecret string, scope string, hc *http.Client) (*ClientApplication, error) {

	if tenantID == "" || clientID == "" || clientSecret == "" || scope == "" {
		return nil, errors.New("Missing TenantID/ClientID/ClientSecret")
	}

	var httpClient *http.Client
	if hc == nil {
		httpClient = &http.Client{}
	}

	app := ClientApplication{
		ClientID:     clientID,
		clientSecret: clientSecret,
		Scope:        scope,
		BaseURL:      "https://login.microsoftonline.com/" + tenantID + "/oauth2/v2.0/token",
		httpClient:   httpClient,
	}

	return &app, nil
}

// AcquireTokenForClient gets authentication token with client id and client secert
func (c *ClientApplication) AcquireTokenForClient() (*Token, error) {

	payload := strings.NewReader("grant_type=" + RequiredGrantType + "&client_id=" + c.ClientID + "&client_secret=" + c.clientSecret + "&scope=" + c.Scope)
	client := c.httpClient
	req, err := http.NewRequest("POST", c.BaseURL, payload)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", DefaultContentType)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	token := Token{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
