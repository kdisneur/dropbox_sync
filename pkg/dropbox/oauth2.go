package dropbox

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/kdisneur/dropbox_sync/pkg/dropbox/internal"
	"github.com/pkg/errors"
)

// OAuth2 represents the data needed to connect to Dropbox
type OAuth2 struct {
	API          string
	Site         string
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// NewOAuth2 represents the OAuth2 configuration to connect to Dropbox
func NewOAuth2(clientID string, clientSecret string) OAuth2 {
	return OAuth2{
		API:          "https://api.dropboxapi.com",
		Site:         "https://www.dropbox.com",
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  "http://localhost/dropbox/callback",
	}
}

// AuthorizationURL returns the URL where the user needs to give its consent
func (o OAuth2) AuthorizationURL() *url.URL {
	authorization, _ := url.Parse(o.Site)
	authorization.Path = "/oauth2/authorize"

	query := authorization.Query()
	query.Set("client_id", o.ClientID)
	query.Set("response_type", "code")
	authorization.RawQuery = query.Encode()

	return authorization
}

// AccessTokenURL returns the URL and POST form to exchange an authorization code
// to an access token
func (o OAuth2) AccessTokenURL(code string) (*url.URL, url.Values) {
	token, _ := url.Parse(o.API)
	token.Path = "/oauth2/token"

	form := make(url.Values)
	form.Set("code", code)
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", o.ClientID)
	form.Set("client_secret", o.ClientSecret)

	return token, form
}

// GetAccessToken exchanges an authorization code for an access token
func (o OAuth2) GetAccessToken(code string) (string, error) {
	var client http.Client

	url, form := o.AccessTokenURL(code)

	postBody := strings.NewReader(form.Encode())
	request, err := http.NewRequest("POST", url.String(), postBody)
	if err != nil {
		return "", errors.Wrap(err, "can't build access token request")
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		return "", errors.Wrap(err, "can't execute access token request")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.Wrap(err, "can't read access token response")
	}

	if response.StatusCode != 200 {
		return "", fmt.Errorf("error exchanging token. detail: %s", body)
	}

	var accessToken internal.AccessTokenResponse
	err = json.Unmarshal(body, &accessToken)
	if err != nil {
		return "", errors.Wrap(err, "can't parse access token response")
	}

	return accessToken.Token, nil
}
