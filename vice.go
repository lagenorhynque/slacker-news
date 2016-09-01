package main

import (
	"fmt"
	"strings"
	"time"

	vice "github.com/munrocape/vice/viceclient"
)

var (
	ViceClient           *vice.Client
	currentViceResponse  string
	currentViceTimestamp time.Time
)

func GetViceTop10() (string, error) {
	var err error
	if ExpiredResponse(currentViceTimestamp) {
		currentViceResponse, err = generateNewViceResponse()
	}
	return currentViceResponse, err
}

func generateNewViceResponse() (string, error) {
	c := getViceClient()
	var rss vice.Rss
	rss, err := c.GetFeed()
	if err != nil {
		return "", err
	}

	var urls [11]string
	urls[0] = "Top Stories from <https://www.vice.com|Vice News>"
	items := rss.Channel.Items
	for index, element := range items {
		index = index + 1
		if index < 11 {
			urls[index] = fmt.Sprintf("%d. <%s|%s>", index, element.Link, element.Title)

		}
	}

	response := strings.Join(urls[:], "\n")
	currentViceTimestamp = time.Now().Local()
	return response, nil
}

func getViceClient() *vice.Client {
	if ViceClient == nil {
		ViceClient = vice.NewClient()
	}
	return ViceClient
}
