package gomyob

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
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

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = tokenURL
	urlStr := fmt.Sprintf("%v", u)

	request := fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&client_secret=%s&scope=CompanyFile", v.StoreCode, v.RedirectURI, v.ClientID, v.ClientSecret)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(request)))

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, _ := client.Do(r)

	rawResBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", "", time.Now(), fmt.Errorf("%v", string(rawResBody))
	}

	if res.StatusCode == 200 {
		resp := &TokenResponse{}
		if err := json.Unmarshal(rawResBody, resp); err != nil {
			fmt.Println("3", string(rawResBody), err)
			return "", "", time.Now(), err
		}
		fmt.Println("4", string(rawResBody))
		return resp.AccessToken, resp.RefreshToken, time.Now().Add(time.Duration(1200) * time.Millisecond), nil
	}

	return "", "", time.Now(), fmt.Errorf("Failed to get access token: %s", res.Status)
}

// RefreshToken will get a new refresg token
func (v *MYOB) RefreshToken(refreshtoken string) (string, string, time.Time, error) {

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = tokenURL
	urlStr := fmt.Sprintf("%v", u)

	request := fmt.Sprintf("grant_type=refresh_token&refresh_token=%s&redirect_uri=%s&client_id=%s&client_secret=%s", refreshtoken, v.RedirectURI, v.ClientID, v.ClientSecret)

	fmt.Println("--------------------------")
	fmt.Println("Calling Reckon", urlStr)
	fmt.Println(string(request))

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(request)))

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	fmt.Println("--------------------------")

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
		return resp.AccessToken, resp.RefreshToken, time.Now().Add(time.Duration(1200) * time.Millisecond), nil
	}

	fmt.Println(string(rawResBody))

	return "", "", time.Now(), fmt.Errorf("Failed to get refresh token: %s", res.Status)
}
