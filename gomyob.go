package gomyob

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"golang.org/x/oauth2"
)

const (
	baseURL    = "https://secure.myob.com"
	tokenURL   = "oauth2/v1/authorize"
	refreshURL = "oauth2/v1/authorize"
)

var (
	defaultSendTimeout = time.Second * 30
)

// MYOB The main struct of this package
type MYOB struct {
	StoreCode    string
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Timeout      time.Duration
}

// NewClient will create a MYOB client with default values
func NewClient(code string, clientID string, clientSecret string, redirectURI string) *MYOB {
	return &MYOB{
		StoreCode:    code,
		Timeout:      defaultSendTimeout,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
	}
}

// AccessToken will get a new access token
func (v *MYOB) AccessToken() (string, string, time.Time, error) {
	/*
		u, _ := url.ParseRequestURI(baseURL)
		u.Path = tokenURL
		urlStr := fmt.Sprintf("%v", u)

		data := TokenRequest{
			ClientID:     v.ClientID,
			ClientSecret: v.ClientSecret,
			Code:         v.StoreCode,
			Scope:        "CompanyFile",
			RedirectURI:  v.RedirectURI,
			GrantType:    "authorization_code",
		}
		//request, _ := json.Marshal(data)
		*
		data := url.Values{}
		data.Set("client_id", v.ClientID)
		data.Set("client_secret", v.ClientSecret)
		data.Set("code", v.StoreCode)
		data.Set("scope", "CompanyFile")
		data.Set("redirect_uri", v.RedirectURI)
		data.Set("grant_type", "authorization_code")

		fmt.Println(urlStr)

		client := &http.Client{}
		//r, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer(request))
		r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))

		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Add("Accept", "application/json")
		r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

		res, _ := client.Do(r)

		rawResBody, err := ioutil.ReadAll(res.Body)

		fmt.Println(string(rawResBody))

		if err != nil {
			return "", "", time.Now(), fmt.Errorf("%v", string(rawResBody))
		}

		if res.StatusCode == 200 {
			resp := &TokenResponse{}
			if err := json.Unmarshal(rawResBody, resp); err != nil {
				return "", "", time.Now(), err
			}
			return resp.AccessToken, resp.RefreshToken, time.Now().Add(time.Duration(resp.ExpiresIn) * time.Millisecond), nil
		}

		fmt.Println(string(rawResBody))

		return "", "", time.Now(), fmt.Errorf("Failed to get access token: %s", res.Status)
	*/

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = tokenURL
	urlStr := fmt.Sprintf("%v", u)

	fmt.Println(urlStr)

	config := oauth2.Config{
		ClientID:     v.ClientID,
		ClientSecret: v.ClientSecret,
		Scopes:       []string{"CompanyFile"},
		RedirectURL:  v.RedirectURI,
		Endpoint: oauth2.Endpoint{
			AuthURL:  urlStr,
			TokenURL: urlStr,
		},
	}

	fmt.Println(config)

	token, err := config.Exchange(oauth2.NoContext, v.StoreCode)

	if err != nil {
		return "", "", time.Now(), err
	}

	return token.AccessToken, token.RefreshToken, time.Now(), err

}

// RefreshToken will get a new refresg token
func (v *MYOB) RefreshToken(refreshtoken string) (string, string, time.Time, error) {

	data := TokenRequest{
		ClientID:     v.ClientID,
		ClientSecret: v.ClientSecret,
		RefreshToken: refreshtoken,
		GrantType:    "refresh_token",
	}

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = refreshURL
	urlStr := fmt.Sprintf("%v", u)

	request, _ := json.Marshal(data)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer(request))

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Length", strconv.Itoa(len(request)))

	res, _ := client.Do(r)

	rawResBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", "", time.Now(), fmt.Errorf("%v", string(rawResBody))
	}

	if res.StatusCode == 200 {
		resp := &TokenResponse{}
		if err := json.Unmarshal(rawResBody, resp); err != nil {
			return "", "", time.Now(), err
		}
		return resp.AccessToken, resp.RefreshToken, time.Now().Add(time.Duration(resp.ExpiresIn) * time.Millisecond), nil
	}

	return "", "", time.Now(), fmt.Errorf("Failed to get access token: %s", res.Status)
}
