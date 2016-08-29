package viceclient

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type Client struct {
	Uri string
}

func NewClient() *Client {
	var c = Client{
		Uri: "http://www.vice.com/rss",
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

func (c *Client) GetFeed() (Rss, error) {
	var url string
	var feed Rss
	url = c.Uri
	rep, err := c.RequestFeed(url)
	if err != nil {
		return feed, err
	}
	xml.Unmarshal(rep, &feed)
	return feed, nil
}
