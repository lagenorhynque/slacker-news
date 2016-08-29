package fivethirtyeightclient

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	Url string
	Categories   map[string]Category
}

type Category struct {
	Uri    string
	Pretty string
}

func NewClient() *Client {
	var categories = map[string]Category{
		"all":           Category{Uri: "all", Pretty: "All Stories"},
		"features":         Category{Uri: "features", Pretty: "Features"},
		"data":       Category{Uri: "datalab", Pretty: "DataLab Posts"},
		"politics":          Category{Uri: "politics", Pretty: "Politics"},
		"economics":            Category{Uri: "economics", Pretty: "Economics"},
		"sports":      Category{Uri: "sports", Pretty: "Sports"},
		"science":      Category{Uri: "science", Pretty: "Science"},
		"life":        Category{Uri: "life", Pretty: "Life"},
	}
	var c = Client{
		Url:          "http://fivethirtyeight.com/%s/feed",
		Categories:   categories,
	}
	return &c
}

func (c *Client) RequestFeed(url string) ([]byte, error) {
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
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

func (c *Client) GetFeed(category string) (Rss, error) {
	var url string
	var feed Rss
	if val, ok := c.Categories[category]; ok {
		url = fmt.Sprintf(c.Url, val.Uri)
	} else {
		return feed, fmt.Errorf("Invalid feed selection: %s\n", category)
	}

	rep, err := c.RequestFeed(url)
	if err != nil {
		return feed, err
	}
	xml.Unmarshal(rep, &feed)
	return feed, nil
}

func (c *Client) GetPretty(category string) string {
	if val, ok := c.Categories[category]; ok {
		return val.Pretty
	} else {
		return ""
	}
}

func (c *Client) GetUri(category string) string {
	if val, ok := c.Categories[category]; ok {
		return val.Uri
	} else {
		return ""
	}
}
