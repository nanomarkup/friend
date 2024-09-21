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
	"go.etcd.io/bbolt"
)

type feed struct {
	Link   string
	Topic  string
	Active bool
}

type telegram struct {
	token  string
	chatId string
}

const (
	dbFileName    string = "feeds.db"
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
	db, err := getDB()
	if err != nil {
		return err
	}
	feeds, err := getFeeds()
	if err != nil {
		return err
	}
	// handle 20 messages per minute
	interval := 3 * time.Second
	currentTime := time.Now()
	lastExecuted := currentTime
	diff := time.Duration(0)

	fp := gofeed.NewParser()
	sender := telegram{"7211500498:AAHDAFhG0CxRxVzYzb9oiOX5y0sc3miyVB8", "-1002415103094"}
	for _, f := range feeds {
		if !f.Active {
			if err = db.Update(
				func(tx *bbolt.Tx) error {
					b := tx.Bucket([]byte("feeds"))
					if b == nil {
						return fmt.Errorf("\"%s\" topic is missing in db", "feeds")
					} else {
						return b.Put([]byte(f.Link), []byte{0})
					}
				}); err != nil {
				fmt.Println(err)
			}
			continue
		}
		// read rss
		feed, err := fp.ParseURL(strings.Trim(f.Link, " "))
		if err != nil {
			fmt.Println(err)
		}
		// check the feed was active
		activated := false
		if err = db.Update(
			func(tx *bbolt.Tx) error {
				feeds := tx.Bucket([]byte("feeds"))
				if feeds == nil {
					return fmt.Errorf("\"%s\" topic is missing in db", "feeds")
				}
				v := feeds.Get([]byte(f.Link))
				if v == nil || v[0] == 0 {
					// add all items to a topic and skip the sending them
					for _, it := range feed.Items {
						topic := tx.Bucket([]byte(f.Topic))
						if topic == nil {
							return fmt.Errorf("\"%s\" topic is missing in db", f.Topic)
						} else if v := topic.Get([]byte(it.Link)); v == nil {
							topic.Put([]byte(it.Link), []byte(time.Now().Format(time.RFC3339)))
						}
					}
					// update the feed
					activated = true
					return feeds.Put([]byte(f.Link), []byte{1})
				}
				return nil
			}); err != nil {
			fmt.Println(err)
			continue
		}
		if activated {
			fmt.Printf("\"%s\" feed is activated\n", f.Link)
			continue
		}
		// iterate all rss
		for _, it := range feed.Items {
			if err = db.Update(
				func(tx *bbolt.Tx) error {
					b := tx.Bucket([]byte(f.Topic))
					if b == nil {
						return fmt.Errorf("\"%s\" topic is missing in db", f.Topic)
					} else if v := b.Get([]byte(it.Link)); v == nil {
						// limit the sending messages
						currentTime = time.Now()
						diff = currentTime.Sub(lastExecuted)
						if diff < interval {
							time.Sleep(interval - diff)
						}
						// send a new message
						err = sender.send(getThreadId(f.Topic), it)
						lastExecuted = currentTime
						if err == nil {
							return b.Put([]byte(it.Link), []byte(time.Now().Format(time.RFC3339)))
						} else {
							return err
						}
					}
					return nil
				}); err != nil {
				fmt.Println(err)
			}
		}
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

func getDB() (*bbolt.DB, error) {
	db, err := bbolt.Open(dbFileName, 0600, nil)
	if err != nil {
		return nil, err
	}
	if err = db.Update(
		func(tx *bbolt.Tx) error {
			if _, err := tx.CreateBucketIfNotExists([]byte("feeds")); err != nil {
				return err
			}
			if _, err := tx.CreateBucketIfNotExists([]byte("cantabria")); err != nil {
				return err
			}
			if _, err := tx.CreateBucketIfNotExists([]byte("travel")); err != nil {
				return err
			}
			if _, err := tx.CreateBucketIfNotExists([]byte("job")); err != nil {
				return err
			}
			return nil
		}); err != nil {
		return nil, err
	}
	return db, nil
}

func getFeeds() ([]feed, error) {
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
	feeds := []feed{}
	err = nanomarkup.Unmarshal(data, &feeds, nil)
	if err != nil {
		return nil, err
	}
	return feeds, nil
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
