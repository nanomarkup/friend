package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/nanomarkup/nanomarkup.go"
	"go.etcd.io/bbolt"
)

func (v *app) getFeeds() ([]feed, error) {
	// get feeds from a file
	wd, _ := os.Getwd()
	filePath := v.feedsFileName
	if wd != "" {
		filePath = fmt.Sprintf("%s/%s", wd, v.feedsFileName)
	}
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

func (v *app) processFeeds(telegramToken string, telegramChat string) error {
	// read feeds
	db, err := bbolt.Open(v.dbFileName, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()
	feeds, err := v.getFeeds()
	if err != nil {
		return err
	}
	fp := gofeed.NewParser()
	// iterate all items
	sender := telegram{telegramToken, telegramChat}
	diff := time.Duration(0)
	interval := 3 * time.Second // handle 20 messages per minute
	currentTime := time.Now()
	lastExecuted := currentTime
	for _, f := range feeds {
		if !f.Active {
			if err = db.Update(
				func(tx *bbolt.Tx) error {
					b := tx.Bucket([]byte(v.feedsBucket))
					if b == nil {
						return fmt.Errorf("\"%s\" topic is missing in db", v.feedsBucket)
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
				feeds := tx.Bucket([]byte(v.feedsBucket))
				if feeds == nil {
					return fmt.Errorf("\"%s\" topic is missing in db", v.feedsBucket)
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
					} else if d := b.Get([]byte(it.Link)); d == nil {
						// limit the sending messages
						currentTime = time.Now()
						diff = currentTime.Sub(lastExecuted)
						if diff < interval {
							time.Sleep(interval - diff)
						}
						// send a new message
						err = sender.send(v.buckets[f.Topic], it)
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

func (v *app) updateDB() error {
	db, err := bbolt.Open(v.dbFileName, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Update(func(tx *bbolt.Tx) error {
		for name := range v.buckets {
			if _, err := tx.CreateBucketIfNotExists([]byte(name)); err != nil {
				return err
			}
		}
		return nil
	})
}
