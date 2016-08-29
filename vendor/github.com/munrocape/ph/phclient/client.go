package phclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Client struct {
	BaseUrl              string
	CurrentToken         Token
	Credential           Credential
	MarshalledCredential []byte
	TokenSuffix          string
	PostsSuffix          string
}

type Credential struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
	CreatedAt   int64  `json:"created_at"`
}

func NewClient() *Client {

	var c = Client{
		BaseUrl:     "https://api.producthunt.com/v1/",
		TokenSuffix: "oauth/token",
		PostsSuffix: "posts",
	}
	var cred, err = c.GetCredential()
	if err != nil {
		fmt.Errorf("The correct environment variables could not be found.")
	}
	marshalled, _ := json.Marshal(cred)
	c.Credential = cred
	c.MarshalledCredential = marshalled
	token, err := c.GenerateToken()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	c.CurrentToken = token
	return &c
}

func (c *Client) Get(url string, params url.Values) ([]byte, error) {
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(c.MarshalledCredential))

	// set the headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	token := c.GetToken().AccessToken
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Host", "api.producthunt.com")

	// set the params
	values := req.URL.Query()
	for key, val := range params {
		values.Add(key, val[0])
	}
	req.URL.RawQuery = values.Encode()

	// make the request
	req.Close = true
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return response, err
}

func (c *Client) Post(url string) ([]byte, error) {
	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(c.MarshalledCredential))
	req.Header.Set("Content-Type", "application/json")
	req.Close = true
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return response, err
}

func (c *Client) GetCredential() (Credential, error) {
	id := os.Getenv("PH_CLIENT_ID")
	secret := os.Getenv("PH_CLIENT_SECRET")
	var cred Credential
	if (id == "") || (secret == "") {
		return cred, fmt.Errorf("Missing environment key for the id or secret. Set PH_CLIENT_ID and PH_CLIENT_SECRET accordingly.")
		os.Exit(1)
	}
	cred = Credential{
		ClientId:     id,
		ClientSecret: secret,
		GrantType:    "client_credentials",
	}
	return cred, nil
}

// Returns a valid client authorization token.
// Returns the current token if valid
// If the current token is not valid, it will generate, assign, and return a new one
func (c *Client) GetToken() Token {
	now := time.Now().Unix()
	currentToken := c.CurrentToken
	var err error
	if currentToken.CreatedAt+currentToken.ExpiresIn <= now {
		c.CurrentToken, err = c.GenerateToken()
		if err != nil {
			fmt.Errorf("Could not generate a new token. Reason: %s\n", err)
			os.Exit(1)
		}
	}
	return c.CurrentToken
}

// Generate and returns a new authorization token
func (c *Client) GenerateToken() (Token, error) {
	url := c.BaseUrl + c.TokenSuffix
	rep, err := c.Post(url)

	var token Token
	if err != nil {
		return token, err
	}

	err = json.Unmarshal(rep, &token)
	token.CreatedAt = time.Now().Unix()
	return token, err
}

func (c *Client) GetPostsToday() (PostsResponse, error) {
	url := c.BaseUrl + c.PostsSuffix
	rep, err := c.Get(url, nil)
	var posts PostsResponse
	if err != nil {
		return posts, err
	}

	err = json.Unmarshal(rep, &posts)
	return posts, err
}

func (c *Client) GetPostsOffset(n int) (PostsResponse, error) {
	endpoint := c.BaseUrl + c.PostsSuffix
	params := url.Values{}
	params.Set("days_ago", strconv.Itoa(n))
	rep, err := c.Get(endpoint, params)
	var posts PostsResponse
	if err != nil {
		return posts, err
	}

	err = json.Unmarshal(rep, &posts)
	return posts, err

}
