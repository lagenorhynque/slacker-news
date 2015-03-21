package main

import (
	"fmt"
	bbc  "github.com/munrocape/bbc/bbcclient"
	"strings"
	"time"
)

var (
	BbcClient           *bbc.Client
	currentBbcResponses  map[string]string
	bbcTimestamps map[string]time.Time
)

func GetBbcTop10(category string) (string, error) {
	var err error
	var rep string
	if expiredBbcResponse(bbcTimestamps[category]) {
		rep, err = generateNewBbcResponse(category)
		currentBbcResponses[category] = rep
	}
	return currentBbcResponses[category], err
}

func expiredBbcResponse(timestamp time.Time) bool {
	timeSinceLast := time.Since(timestamp)
	if timeSinceLast > timeToExpire {
		return true
	}
	return false
}

func generateNewBbcResponse(category string) (string, error) {
	c := getBbcClient()
	rss, err := c.GetFeed(category)
	if err != nil {
		return "", err
	}

	var urls [11]string
	urls[0] = "Top Stories from BBC " + category
	items := rss.Channel.Items
	for index, element := range items {
		index = index + 1
		if (index < 11){
			urls[index] = fmt.Sprintf("<%s|%d. %s>\n\t%s", element.Link, index, element.Title, element.Description)
		}
	}

	response := strings.Join(urls[:], "\n")
	bbcTimestamps[category] = time.Now().Local()
	return response, nil
}

func GetBbcSources() string {
	c := getBbcClient()
	res := ""
	first := true
	for k, _ := range c.NewsCategories {
		if (first){
			res = res + k
			first = false
		} else {
			res = res + ", " + k
		}
		
	}
	for k, _ := range c.SportsCategories {
		res = res + ", " + k
	}
	return res
}

func getBbcClient() *bbc.Client {
	if BbcClient == nil {
		BbcClient = bbc.NewClient()
		currentBbcResponses = make(map[string]string)
		bbcTimestamps = make(map[string]time.Time)
		initializeTimestampMap(BbcClient)
		initializeResponseMap(BbcClient)
	}
	return BbcClient
}

func initializeTimestampMap(c *bbc.Client) {
	for k, _ := range c.NewsCategories {
		bbcTimestamps[k] = time.Now().Local().AddDate(0, 0, -11)
	}
	for k, _ := range c.SportsCategories {
		bbcTimestamps[k] = time.Now().Local().AddDate(0, 0, -11)
	}
}

func initializeResponseMap(c *bbc.Client) {
	for k, _ := range c.NewsCategories {
		currentBbcResponses[k] = ""
	}
	for k, _ := range c.SportsCategories {
		currentBbcResponses[k] = ""
	}
}