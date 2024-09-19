package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/mmcdole/gofeed"
)

func (v *telegram) send(threadId int, items []*gofeed.Item) error {
	mes := ""
	var err error
	for _, it := range items {
		// it.Custom[sentItem] = "true"
		// continue
		if it.Custom[sentItem] == "true" {
			continue
		}
		mes, err = v.getMessage(it.Title, it.Custom["categorias"], it.Link)
		if err != nil {
			return err
		}
		// send
		params := url.Values{}
		params.Add("chat_id", v.chatId)
		if threadId > 0 {
			params.Add("message_thread_id", strconv.Itoa(threadId))
		}
		params.Add("parse_mode", "HTML")
		params.Add("text", mes)
		resp, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?%s", v.token, params.Encode()))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		// check
		if resp.StatusCode == 200 {
			it.Custom[sentItem] = "true"
		}
	}
	return nil
}

func (v *telegram) getMessage(title string, cats string, link string) (string, error) {
	website, err := isWebsite(link)
	if err != nil {
		return "", err
	}
	mes := title + "</a>\n\n"
	if cats != "" {
		items := strings.Split(cats, "|")
		if len(items) > 0 {
			cats = items[0]
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
