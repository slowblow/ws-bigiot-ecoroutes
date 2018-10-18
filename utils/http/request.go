package http

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

const (
	ContentTypeHeaderDefault            = "application/json"
	ContentTypeHeaderXWwwFormUrlencoded = "application/x-www-form-urlencoded"
)

func HttpClientDoRequestBasicAuth(method, urlRequest string, body io.Reader, basicAccessAuthenticationUser, basicAccessAuthenticationPassword string) (*http.Response, error) {
	return HttpClientDoRequestBasicAuthContentTypeHeader(method, urlRequest, body, basicAccessAuthenticationUser, basicAccessAuthenticationPassword, ContentTypeHeaderDefault)
}

func HttpClientDoRequestBasicAuthContentTypeHeader(method, urlRequest string, body io.Reader, basicAccessAuthenticationUser, basicAccessAuthenticationPassword string, contentTypeHeader string) (*http.Response, error) {
	request, err := http.NewRequest(method, urlRequest, body)

	request.Header.Set("Content-Type", contentTypeHeader)

	if err != nil {
		return nil, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		//CheckRedirect: redirectPolicyFunc,
	}

	if strings.HasPrefix(strings.ToLower(urlRequest), "http://") {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			//req.Header.Add("Authorization", "Basic "+basicAuth("cigo-m2m", "dTW7Rqga"))
			SetBasicAccessAuthenticationBySetups(req, basicAccessAuthenticationUser, basicAccessAuthenticationPassword)

			return nil
		}
	} else {
		SetBasicAccessAuthenticationBySetups(request, basicAccessAuthenticationUser, basicAccessAuthenticationPassword)
	}

	//request.SetBasicAuth("cigo-m2m", "dTW7Rqga")
	response, err := client.Do(request)

	return response, err

}

func HttpClientDoRequest(method, urlRequest string, body io.Reader) (*http.Response, error) {
	return HttpClientDoRequestBasicAuth(method, urlRequest, body, "", "")
}

func SetBasicAccessAuthenticationBySetups(request *http.Request, basicAccessAuthenticationUser, basicAccessAuthenticationPassword string) {
	if basicAccessAuthenticationUser != "" && basicAccessAuthenticationPassword != "" {
		request.SetBasicAuth(basicAccessAuthenticationUser, basicAccessAuthenticationPassword)
	}
}

func HttpClientDoRequestWithToken(method, urlRequest, token string) (*http.Response, error) {

	request, err := http.NewRequest(method, urlRequest, nil)

	if err != nil {
		return nil, err
	}

	/*
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	*/

	client := &http.Client{
	//Transport: tr,
	}

	var bearer = "Bearer " + token
	request.Header.Set("Authorization", bearer)

	response, err := client.Do(request)

	return response, err
}

/////////////////////////////////////////////////////
// Method to use as client of our own webservices //
/////////////////////////////////////////////////////

func HttpClientDoRequestBAH(method, urlRequest string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method, urlRequest, body)
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	if err != nil {
		return nil, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
	}

	//add our client BasicAuthenticationHttp (BAH) config
	//setClientBasicAuthenticationHttpByConfig(urlRequest, client, request)
	//previous function doesn't work as we expect
	setClientBasicAuthenticationHttp(request, nil)

	response, err := client.Do(request)

	return response, err
}

func setClientBasicAuthenticationHttpByConfig(urlRequest string, client *http.Client, request *http.Request) {
	if strings.HasPrefix(strings.ToLower(urlRequest), "http://") {
		client.CheckRedirect = setClientBasicAuthenticationHttp
	} else {

	}
}

func setClientBasicAuthenticationHttp(request *http.Request, via []*http.Request) error {

	configUsername := GetClientBasicAuthenticationHttpUsername()

	configPassword := GetClientBasicAuthenticationHttpPassword()

	request.SetBasicAuth(configUsername, configPassword)

	return nil
}

func GetClientBasicAuthenticationHttpUsername() string {
	username := os.Getenv("CLIENT_BASIC_AUTHENTICATION_HTTP_USERNAME")
	if username == "" {
		username = beego.AppConfig.String("ClientBasicAuthenticationHttpUsername")
	}

	return username
}

func GetClientBasicAuthenticationHttpPassword() string {
	password := os.Getenv("CLIENT_BASIC_AUTHENTICATION_HTTP_PASSWORD")
	if password == "" {
		password = beego.AppConfig.String("ClientBasicAuthenticationHttpPassword")
	}

	return password
}

func GetOAuth2ClientID() string {
	oAuth2ClientID := os.Getenv("OAUTH2_CLIENT_ID")
	//fmt.Println("OAUTH2_CLIENT_ID", oAuth2ClientID)
	if oAuth2ClientID == "" {
		oAuth2ClientID = beego.AppConfig.String("OAuth2ClientID")
		//fmt.Println("OAuth2ClientID", oAuth2ClientID)
	}

	return oAuth2ClientID
}

func GetOAuth2ClientSecret() string {
	oAuth2ClientSecret := os.Getenv("OAUTH2_CLIENT_SECRET")
	//fmt.Println("OAUTH2_CLIENT_SECRET", oAuth2ClientSecret)
	if oAuth2ClientSecret == "" {
		oAuth2ClientSecret = beego.AppConfig.String("OAuth2ClientSecret")
		//fmt.Println("OAuth2ClientSecret", oAuth2ClientSecret)
	}

	return oAuth2ClientSecret
}

func GetOAuth2Scopes() []string {
	oAuth2Scopes := os.Getenv("OAUTH2_SCOPES")
	//fmt.Println("OAUTH2_SCOPES", oAuth2Scopes)
	if oAuth2Scopes == "" {
		oAuth2Scopes = beego.AppConfig.String("OAuth2Scopes")
		//fmt.Println("OAuth2Scopes", oAuth2Scopes)
	}

	return strings.Split(oAuth2Scopes, " ")
}

func GetOAuth2AccessTokenURL() string {
	oAuth2AccessTokenURL := os.Getenv("OAUTH2_ACCESS_TOKEN_URL")
	//fmt.Println("OAUTH2_ACCESS_TOKEN_URL", oAuth2AccessTokenURL)
	if oAuth2AccessTokenURL == "" {
		oAuth2AccessTokenURL = beego.AppConfig.String("OAuth2AccessTokenURL")
		//fmt.Println("OAuth2AccessTokenURL", oAuth2AccessTokenURL)
	}

	return oAuth2AccessTokenURL
}

type oauth2TokenJSON struct {
	AccessToken string `json:"access_token"` //"5c712150-7f12-46e9-9caf-2e335166b88d"
	TokenType   string `json:"token_type"`   //"bearer"
	ExpiresIn   int32  `json:"expires_in"`   //75816 (seconds to expire)
	Scope       string `json:"scope"`        //"vehicle:action vehicle:max-battery"
}

type OAuth2Token struct {
	AccessToken string    `json:"access_token"` //"5c712150-7f12-46e9-9caf-2e335166b88d"
	TokenType   string    `json:"token_type"`   //"bearer"
	ExpiresIn   time.Time `json:"expires_in"`   // time.Now() + 75816 seconds
	Scope       string    `json:"scope"`        //"vehicle:action vehicle:max-battery"
}

func (e *oauth2TokenJSON) expiry() (t time.Time) {
	if v := e.ExpiresIn; v != 0 {
		return time.Now().Add(time.Duration(v) * time.Second)
	}

	return
}

func GetOAuth2Token() (*OAuth2Token, error) {

	clientID := GetOAuth2ClientID()
	clientSecret := GetOAuth2ClientSecret()
	tokenURL := GetOAuth2AccessTokenURL()

	body := strings.NewReader("grant_type=client_credentials")

	response, err := HttpClientDoRequestBasicAuthContentTypeHeader(http.MethodPost, tokenURL, body, clientID, clientSecret, ContentTypeHeaderXWwwFormUrlencoded)

	token := &oauth2TokenJSON{}

	if err == nil {
		if response.StatusCode == http.StatusOK {
			var responseBody []byte
			responseBody, err = ioutil.ReadAll(io.LimitReader(response.Body, 1048576))

			if err == nil {
				err = json.Unmarshal(responseBody, token)
				if err == nil {
					return &OAuth2Token{
						AccessToken: token.AccessToken,
						TokenType:   token.TokenType,
						ExpiresIn:   token.expiry(),
						Scope:       token.Scope,
					}, nil
				}
			}
		} else {
			err = errors.New("GetOAuth2Token (" + strconv.Itoa(response.StatusCode) + ") " + response.Status)
		}
	}

	return nil, err
}
