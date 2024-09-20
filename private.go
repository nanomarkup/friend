package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/nanomarkup/nanomarkup.go"
)

type feed struct {
	Link   string
	File   string
	Topic  string
	Active bool
}

type reader struct {
	feed *feed
}

type telegram struct {
	token  string
	chatId string
}

const (
	sentItem      string = "sent"
	feedsFileName string = "feeds.nano"
)

func main() {
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
		err = readFeeds()
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
			hour = now.Hour()
			minute = now.Minute()
			fmt.Printf("[%d:%d] Reading feeds...\n", hour, minute)
			err = readFeeds()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func readFeeds() error {
	feeds, err := getFeeds()
	if err != nil {
		return err
	}
	// handle 20 messages per minute
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	sender := telegram{"7211500498:AAHDAFhG0CxRxVzYzb9oiOX5y0sc3miyVB8", "-1002415103094"}
	for _, f := range feeds {
		if !f.Active {
			continue
		}
		// read rss
		r := reader{f}
		items, err := r.read()
		if err != nil {
			fmt.Println(err)
		}
		// get items to send
		messages := make(chan *gofeed.Item)
		go func() {
			for _, it := range items {
				if it.Custom == nil {
					it.Custom = map[string]string{}
				}
				// it.Custom[sentItem] = "true"
				// continue
				if it.Custom[sentItem] != "true" {
					messages <- it
				}
			}
			close(messages)
		}()
		// send items using limitation in time
		for msg := range messages {
			<-ticker.C
			err = sender.send(getThreadId(f.Topic), msg)
			if err == nil {
				msg.Custom[sentItem] = "true"
			} else {
				fmt.Println(err)
			}
		}
		r.save(items)
	}
	return nil
}

func translate(src, srcLang, dstLang string) (string, error) {
	url := fmt.Sprintf("https://translate.googleapis.com/translate_a/single?client=gtx&sl=%s&tl=%s&dt=t&q=%s", srcLang, dstLang, url.QueryEscape(src))
	r, err := http.Get(url)
	if err != nil {
		return "", errors.New("error getting translate.googleapis.com")
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", errors.New("error reading response body")
	}
	if strings.Contains(string(body), `<title>Error 400 (Bad Request)`) {
		return "", errors.New("error 400 (Bad Request)")
	}
	var result []interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", errors.New("error unmarshaling data")
	}
	if len(result) > 0 {
		var text []string
		inner := result[0]
		for _, slice := range inner.([]interface{}) {
			for _, translatedText := range slice.([]interface{}) {
				text = append(text, fmt.Sprintf("%v", translatedText))
				break
			}
		}
		return strings.Join(text, ""), nil
	} else {
		return "", errors.New("no translated data in responce")
	}
}

func getThreadId(topic string) int {
	switch topic {
	case "job":
		return 6
	case "travel":
		return 11
	default:
		return 4 // cantabria
	}
}

func getFeeds() ([]*feed, error) {
	// get feeds from a file
	wd, _ := os.Getwd()
	filePath := fmt.Sprintf("%s/%s", wd, feedsFileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("\"%s\" does not exist", filePath)
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	feeds := []*feed{}
	err = nanomarkup.Unmarshal(data, &feeds, nil)
	if err != nil {
		return nil, err
	}
	// update path to files
	for _, it := range feeds {
		it.File = fmt.Sprintf("%s/%s", wd, it.File)
	}
	return feeds, nil
	// 4, 6, 11
}

func isWebsite(url string) (bool, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return false, err
	}
	// optionally set User-Agent to mimic a browser (some servers block non-browser requests)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	contentType := resp.Header.Get("Content-Type")
	types := strings.Split(contentType, ";")
	for _, it := range types {
		if it == "text/html" || it == "text/xhtml" {
			return true, nil
		}
	}
	return false, nil
}
