package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mmcdole/gofeed"
)

func (v *telegram) send(threadId int, item *gofeed.Item) error {
	cats := item.Custom["categorias"]
	if len(item.Categories) > 0 {
		cats = item.Categories[0]
	}
	mes, err := v.getMessage(item.Title, cats, item.Link)
	if err != nil {
		return err
	}
	// send
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", v.token)
	payload := map[string]interface{}{
		"chat_id":           v.chatId,
		"message_thread_id": threadId,
		"parse_mode":        "HTML",
		"text":              mes,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	} else {
		return fmt.Errorf("failed to send message. Status code: %d", resp.StatusCode)
	}
}

func (v *telegram) getMessage(title string, cats string, link string) (string, error) {
	website, err := isWebsite(link)
	if err != nil {
		return "", err
	}
	mes := title + "</a>\n\n"
	if cats != "" {
		items := strings.Split(cats, "|")
		if len(items) > 1 {
			cats = items[0]
			if len(cats) < 10 {
				cats += " | " + items[1]
			}
		}
		mes += fmt.Sprintf("Categorias: %s\n\n", cats)
	}
	mes, err = translate(mes, "es", "uk")
	if err != nil {
		return "", err
	}
	if website {
		mes = fmt.Sprintf("<a href=\"%s\">%s</a>\n\n<a href=\"http://translate.google.com/translate?sl=es&tl=uk&u=%s&client=webapp/\">%s", link, title, link, mes)
	} else {
		mes = fmt.Sprintf("<a href=\"%s\">%s\n\n%s", link, title, mes)
	}
	return mes, nil
}
