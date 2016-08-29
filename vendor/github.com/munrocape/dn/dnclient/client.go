package dnclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type Client struct {
	GetUrl  string
	Stories string
	Motd    string
	Search  string
}

func NewClient() *Client {
	client_id := os.Getenv("DN_CLIENT_ID")
	if client_id == "" {
		fmt.Printf("ERROR: Could not find Designer News Client Id\n")
		os.Exit(1)
	}
	var c = Client{
		GetUrl:  "https://api-news.layervault.com/api/v1/%sclient_id=" + client_id,
		Stories: "stories?",
		Motd:    "motd?",
		Search:  "stories/search?query=%s&",
	}
	return &c
}

func (c *Client) Get(url string, params url.Values) ([]byte, error) {
	httpClient := &http.Client{}
	fmt.Printf("%s\n", url)
	req, err := http.NewRequest("GET", url, nil)

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

func (c *Client) GetStories() (Stories, error) {
	url := fmt.Sprintf(c.GetUrl, c.Stories)
	rep, err := c.Get(url, nil)
	var stories Stories
	if err != nil {
		return stories, err
	}
	err = json.Unmarshal(rep, &stories)
	return stories, err
}

func (c *Client) GetMotd() (Motd, error) {
	url := fmt.Sprintf(c.GetUrl, c.Motd)
	rep, err := c.Get(url, nil)
	var wrapper MotdWrapper
	var motd Motd
	if err != nil {
		return motd, err
	}
	err = json.Unmarshal(rep, &wrapper)
	motd = wrapper.Motd
	return motd, err
}

func (c *Client) GetSearch(search string) (Stories, error) {
	query := url.QueryEscape(search)
	value := fmt.Sprintf(c.Search, query)
	url := fmt.Sprintf(c.GetUrl, value)
	rep, err := c.Get(url, nil)
	var stories Stories
	if err != nil {
		return stories, err
	}
	err = json.Unmarshal(rep, &stories)
	return stories, err
}

// func main() {
// 	c := NewClient()
// 	s, _ := c.GetStories()
// 	for _, element := range s.Stories {
// 		fmt.Printf("%d %s %s %s %d\n", element.VoteCount, element.Title, element.Url, element.SiteUrl, len(element.Comments))
// 	}
// 	m, _ := c.GetMotd()
// 	fmt.Printf("%+v\n", m)
// 	s, _ = c.GetSearch("Product Hunt")
// 	for _, element := range s.Stories {
// 		fmt.Printf("%d %s %s %s %d %+v\n", element.VoteCount, element.Title, element.Url, element.SiteUrl, element.CommentCount, element.Comments)
// 	}
// }
