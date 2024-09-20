package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/mmcdole/gofeed"
	"github.com/nanomarkup/nanomarkup.go"
)

func (v *reader) read() ([]*gofeed.Item, error) {
	// get items from an online resource
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(strings.Trim(v.feed.Link, " "))
	if err != nil {
		return nil, err
	}
	// get items from a file
	_, err = os.Stat(v.feed.File)
	if err != nil {
		if os.IsNotExist(err) {
			return feed.Items, nil
		} else {
			return nil, err
		}
	}
	items := []gofeed.Item{}
	data, err := os.ReadFile(v.feed.File)
	if err != nil {
		return nil, err
	}
	err = nanomarkup.Unmarshal(data, &items, nil)
	if err != nil {
		return nil, err
	}
	// synchronize items
	for _, it := range feed.Items {
		for i, d := range items {
			if d.Link == it.Link {
				val, ok := d.Custom[sentItem]
				if ok {
					it.Custom[sentItem] = val
				} else {
					it.Custom[sentItem] = "false"
				}
				items = append(items[:i], items[i+1:]...)
				break
			}
		}
	}
	return feed.Items, nil
}

func (v *reader) save(items []*gofeed.Item) error {
	_, err := os.Stat(v.feed.File)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(v.feed.File), os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	file, err := os.Create(v.feed.File)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := nanomarkup.MarshalIndent(items, "", "    ")
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	_, err = writer.Write(data)
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}
