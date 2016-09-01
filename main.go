package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bluele/slack"
)

var (
	timeToExpire time.Duration
	webHook      *slack.WebHook
)

func main() {
	// Parse port from command-line parameters
	port := flag.String("port", "8080", "HTTP Port to listen on")
	flag.Parse()

	// declare variables
	timeToExpire = 10 * time.Minute
	webHook = slack.NewWebHook(os.Getenv("WEBHOOK_URL"))

	// Start our Server
	log.Println("Starting Server on", *port)
	http.HandleFunc("/", index)
	http.HandleFunc("/news", news)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Slacker News, a slack integration to provide current news. Check it out at www.github.com/munrocape/slacker-news"))
}

func news(w http.ResponseWriter, r *http.Request) {
	newsSource := r.URL.Query().Get("text")
	log.Println(newsSource)
	tokens := strings.Split(strings.ToLower(newsSource), " ")
	var source, argument string
	if len(tokens) == 2 {
		source, argument = tokens[0], tokens[1]
	} else {
		source, argument = tokens[0], ""
	}

	userName := formatUserName(r)
	channel := getChannel(r)
	switch source {
	case "hn":
		stories, err := GetHnTop10()
		if err == nil {
			postToSlack(userName+stories, channel)
		} else {
			w.Write([]byte("Server Error - Firebase could not be reached"))
		}
		return
	case "ph":
		posts, err := GetPhTop10()
		if err == nil {
			postToSlack(userName+posts, channel)
		} else {
			w.Write([]byte("Server Error - Product Hunt could not be reached"))
		}
		return
	case "vice":
		articles, err := GetViceTop10()
		if err == nil {
			postToSlack(userName+articles, channel)
		} else {
			w.Write([]byte("Server Error - Vice News could not be reached"))
		}
		return
	case "bbc":
		articles, err := GetBbcTop10(argument)
		if err == nil {
			postToSlack(userName+articles, channel)
		} else {
			if strings.Contains(err.Error(), "Invalid feed selection") {
				response := fmt.Sprintf("That is an invalid BBC category: %s\nTry `/news help` to view all sources.", argument)
				w.Write([]byte(response))
			} else {
				w.Write([]byte("Server Error - the BBC could not be reached"))
			}
		}
		return
	case "538":
		articles, err := GetFteTop10(argument)
		if err == nil {
			postToSlack(userName+articles, channel)
		} else {
			if strings.Contains(err.Error(), "Invalid feed selection") {
				response := fmt.Sprintf("That is an invalid FiveThirtyEight category: %s\nTry `/news help` to view all sources.", argument)
				w.Write([]byte(response))
			} else {
				w.Write([]byte("Server Error - FiveThirtyEight could not be reached"))
			}
		}
		return
	case "dn":
		stories, err := GetDnArgument(argument)
		if err == nil {
			postToSlack(userName+stories, channel)
		} else {
			if strings.Contains(err.Error(), "Invalid argument") {
				w.Write([]byte("Invalid argument - try `/news help` to view valid selections."))
			} else {
				w.Write([]byte("Server Error - Designer News could not be reached"))
			}
		}
		return
	case "help":
		w.Write([]byte(GetSources()))
		return
	case "":
		w.Write([]byte(GetSources()))
		return
	}
	userArgument := fmt.Sprintf("%s %s", source, argument)
	w.Write([]byte("Hmm.. I can't figure out what news you are looking for :( I received \"" + strings.TrimSpace(userArgument) + "\"\nTry `/news help` to view all sources."))
}

func GetSources() string {
	hn := "Hacker News: hn\n"
	ph := "Product Hunt: ph\n"
	vice := "Vice News: vice\n"
	designer := "Designer News: dn [stories, motd]\n"
	fte := "FiveThirtyEight: 538 [" + GetFteSources() + "]\n"
	bbc := "BBC: bbc [" + GetBbcSources() + "]\n"
	return fmt.Sprintf("%s%s%s%s%s%s", hn, ph, vice, designer, fte, bbc)
}

func ExpiredResponse(timestamp time.Time) bool {
	timeSinceLast := time.Since(timestamp)
	if timeSinceLast > timeToExpire {
		return true
	}
	return false
}

func formatUserName(r *http.Request) string {
	return fmt.Sprintf("// %s\n", r.URL.Query().Get("user_name"))
}

func getChannel(r *http.Request) string {
	channelName := r.URL.Query().Get("channel_name")
	switch channelName {
	case "directmessage":
		return "@" + r.URL.Query().Get("user_name")
	case "privategroup":
		return r.URL.Query().Get("channel_id")
	default:
		return "#" + channelName
	}
}

func postToSlack(text string, channel string) {
	err := webHook.PostMessage(&slack.WebHookPostPayload{
		Text:    text,
		Channel: channel,
	})
	if err != nil {
		panic(err)
	}
}
