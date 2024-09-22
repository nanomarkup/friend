package main

import (
	"bytes"
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

func (v *telegram) send(threadId int, item *gofeed.Item) error {
	cats := item.Custom["categorias"]
	if len(item.Categories) > 0 {
		cats = item.Categories[0]
	}
	mes, err := v.getMessage(item.Title, cats, item.Link)
	if err != nil {
		return err
	}
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
	website, err := v.isWebsite(link)
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
	mes, err = v.translateMessage(mes, "es", "uk")
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

func (v *telegram) translateMessage(mes, srcLang, dstLang string) (string, error) {
	url := fmt.Sprintf("https://translate.googleapis.com/translate_a/single?client=gtx&sl=%s&tl=%s&dt=t&q=%s", srcLang, dstLang, url.QueryEscape(mes))
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

func (v *telegram) isWebsite(url string) (bool, error) {
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
