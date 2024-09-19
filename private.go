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
	wd, _ := os.Getwd()
	feeds := []*feed{}
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802081", wd + "/data/6802081.nam"})
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802084", wd + "/data/6802084.nam"})
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802085", wd + "/data/6802085.nam"})
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802086", wd + "/data/6802086.nam"})
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802087", wd + "/data/6802087.nam"})
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802089", wd + "/data/6802089.nam"})
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802090", wd + "/data/6802090.nam"})
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802091", wd + "/data/6802091.nam"})
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802092", wd + "/data/6802092.nam"})
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802094", wd + "/data/6802094.nam"})
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802095", wd + "/data/6802095.nam"})
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802097", wd + "/data/6802097.nam"})
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802098", wd + "/data/6802098.nam"})
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802099", wd + "/data/6802099.nam"})
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802100", wd + "/data/6802100.nam"})
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802301", wd + "/data/6802301.nam"})
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/7479572", wd + "/data/7479572.nam"})
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/6802303", wd + "/data/6802303.nam"})
	// feeds = append(feeds, &feed{"https://www.cantabria.es/o/BOC/feed/7293890", wd + "/data/7293890.nam"})
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16413/eid4008216/27148/size20", wd + "/data/27148.nam", 4})
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16413/eid4008216/27147/size20", wd + "/data/27147.nam", 4})
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16413/eid4008216/27141/size20", wd + "/data/27141.nam", 4})
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16413/eid4008216/27139/size20", wd + "/data/27139.nam", 4})
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16413/eid4008216/27146/size20", wd + "/data/27146.nam", 4})
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16413/eid4008216/27140/size20", wd + "/data/27140.nam", 4})
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16413/eid4008216/27145/size20", wd + "/data/27145.nam", 4})
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16413/eid4008216/27149/size20", wd + "/data/27149.nam", 4})
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16413/eid4008216/7720252/size20", wd + "/data/7720252.nam", 4})
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16401/eid4008203/size50", wd + "/data/4008203.nam", 4})
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16401/eid4008202/size50", wd + "/data/4008202.nam", 4})
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16401/eid4008221/size50", wd + "/data/4008221.nam", 4})
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16401/eid4008246/size50", wd + "/data/4008246.nam", 4})
	feeds = append(feeds, &feed{"https://www.cantabria.es/o/GOBIERNO/feed/group16413/eid4008216/size50", wd + "/data/4008216.nam", 4})
	feeds = append(feeds, &feed{"https://empleopublico.cantabria.es/o/GOBIERNO/feed/group16475/inscom_liferay_journal_content_web_portlet_JournalContentPortlet_INSTANCE_6Cx0YFAD8ZVK/size50", wd + "/data/6Cx0YFAD8ZVK.nam", 4})

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
	sender := telegram{"7211500498:AAHDAFhG0CxRxVzYzb9oiOX5y0sc3miyVB8", "-1002415103094"} //-4573799901
	for _, f := range feeds {
		r := reader{f}
		items, err := r.read()
		if err != nil {
			fmt.Println(err)
		}
		err = sender.send(f.messageThreadId, items)
		if err != nil {
			fmt.Println(err)
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
