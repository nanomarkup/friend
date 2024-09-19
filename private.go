package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

type feed struct {
	url             string
	filePath        string
	messageThreadId int
}

type reader struct {
	feed *feed
}

type telegram struct {
	token  string
	chatId string
}

const (
	sentItem string = "sent"
)

func main() {
	feeds := getFeeds()
	spain, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		fmt.Println("Error loading location:", err)
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
		readFeeds(feeds)
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
			readFeeds(feeds)
		}
	}
}

func readFeeds(feeds []*feed) {
	// handle 20 messages per minute
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	sender := telegram{"7211500498:AAHDAFhG0CxRxVzYzb9oiOX5y0sc3miyVB8", "-1002415103094"} //-4573799901
	for _, f := range feeds {
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
			err = sender.send(f.messageThreadId, msg)
			if err == nil {
				msg.Custom[sentItem] = "true"
			} else {
				fmt.Println(err)
			}
		}
		r.save(items)
	}
}

func translate(src, srcLang, dstLang string) (string, error) {
	url := fmt.Sprintf("https://translate.googleapis.com/translate_a/single?client=gtx&sl=%s&tl=%s&dt=t&q=%s", srcLang, dstLang, url.QueryEscape(src))
	r, err := http.Get(url)
	if err != nil {
		return "", errors.New("Error getting translate.googleapis.com")
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", errors.New("Error reading response body")
	}
	if strings.Contains(string(body), `<title>Error 400 (Bad Request)`) {
		return "", errors.New("Error 400 (Bad Request)")
	}
	var result []interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", errors.New("Error unmarshaling data")
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
		return "", errors.New("No translated data in responce")
	}
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
