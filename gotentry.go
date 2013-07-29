package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type RSS struct {
	Channel Channel `xml:"channel"`
	Item    []Item  `xml:"item"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
}

type Item struct {
	Title     string `xml:"title"`
	Link      string `xml:"link"`
	Bookmarks int    `xml:"bookmarkcount"`
}

func hotentryUrl(keyword string, threshold int) string {
	tmpl := `http://b.hatena.ne.jp/search/tag?q=%s&users=%d&mode=rss`
	ret := fmt.Sprintf(tmpl, keyword, threshold)

	return ret
}

func main() {
	threshold := flag.Int("t", 3, "threshold of bookmarks")
	limit := flag.Int("l", 10, "limit of printing entries")
	flag.Parse()

	key := flag.Arg(0)
	if key == "" {
		fmt.Printf("Please specified 'keyword'!!\n")
		os.Exit(1)
	}

	url := hotentryUrl(key, *threshold)
	response, err := http.Get(url)
	if err != nil {
		log.Printf("Can't download '%s'", url)
		os.Exit(1)
	}
	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Can't read response: %s", err.Error())
		os.Exit(1)
	}

	var rss RSS
	if err := xml.Unmarshal(bytes, &rss); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	if *limit > len(rss.Item) {
		fmt.Printf("Limit '%s' is too long", *limit)
		*limit = len(rss.Item)
	}

	for i, item := range rss.Item[:*limit] {
		fmt.Printf("%2d: %s [%d]\n", i+1, item.Title, item.Bookmarks)
	}
}
