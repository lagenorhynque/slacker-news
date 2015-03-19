package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

var (
	timeToExpire time.Duration
)

func main() {
	// Parse port from command-line parameters
	port := flag.String("port", "8080", "HTTP Port to listen on")
	flag.Parse()

	// declare variables
	timeToExpire = 10 * time.Minute

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
	news_source := r.URL.Query().Get("text")
	log.Println(news_source)
	switch {
	case news_source == "hn":
		stories, err := getHnTop10()
		if err == nil {
			w.Write([]byte(stories))
		} else {
			w.Write([]byte("Server Error - Firebase could not be reached"))
		}
		return
	case news_source == "ph":
		posts, err := GetPhTop10()
		if err == nil {
			w.Write([]byte(posts))
		} else {
			w.Write([]byte("Server Error - Product Hunt could not be reached"))
		}
		return
	}
	w.Write([]byte("can't find that one! " + news_source))
}
