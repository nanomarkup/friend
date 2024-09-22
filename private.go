package main

import (
	"fmt"
	"time"
)

type app struct {
	buckets       map[string]int
	feedsBucket   string
	dbFileName    string
	feedsFileName string
}

type feed struct {
	Link     string
	Topic    string
	Language string
	Active   bool
}

type telegram struct {
	token  string
	chatId string
}

const (
	dbFileName    string = "feeds.db"
	feedsFileName string = "feeds.nano"
	feedsBucket   string = "feeds"
	telegramToken string = "7211500498:AAHDAFhG0CxRxVzYzb9oiOX5y0sc3miyVB8"
	telegramChat  string = "-1002415103094"
)

var buckets = map[string]int{
	feedsBucket:   0,
	"immigration": 2,
	"charity":     3,
	"cantabria":   4,
	"school":      5,
	"job":         6,
	"transport":   7,
	"health":      8,
	"market":      9,
	"home":        10,
	"travel":      11,
	"events":      12,
	"galicia":     13,
	"asturias":    14,
	"basque":      16,
	"larioja":     17,
	"navarra":     18,
	"leon":        19,
	"children":    55,
	"spam":        57,
	"sport":       73,
	"ukraine":     301,
}

func main() {
	p := app{
		buckets,
		feedsBucket,
		dbFileName,
		feedsFileName,
	}
	if err := p.updateDB(); err != nil {
		fmt.Println(err)
		return
	}
	spain, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		fmt.Printf("Error loading location: %s\n", err)
		return
	}

	startHour := 9
	endHour := 21
	wait := 0
	now := time.Now().In(spain)
	hour := now.Hour()
	minute := now.Minute()
	if hour >= startHour && hour <= endHour {
		fmt.Printf("[%d:%d] Reading feeds...\n", hour, minute)
		err = p.processFeeds(telegramToken, telegramChat)
		if err != nil {
			fmt.Println(err)
		}
	}
	for {
		now = time.Now().In(spain)
		hour = now.Hour()
		minute = now.Minute()
		wait = 60 - minute
		fmt.Printf("[%d:%d] Waiting for %d minutes...\n", hour, minute, wait)
		time.Sleep(time.Duration(wait) * time.Minute)
		if hour >= startHour && hour <= endHour {
			now = time.Now().In(spain)
			fmt.Printf("[%d:%d] Reading feeds...\n", now.Hour(), now.Minute())
			err = p.processFeeds(telegramToken, telegramChat)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
