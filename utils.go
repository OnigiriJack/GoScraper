package main

import (
	"log"
	"net/http"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

func ScrapeHtmlFromPage(url string) []string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	textTags := []string{
		"a",
		"p", "span", "em", "string", "blockquote", "q", "cite",
		"h1", "h2", "h3", "h4", "h5", "h6", "title",
	}
	tag := ""
	enter := false
	res := make([]string, 0)
	z := html.NewTokenizer(resp.Body)

	for {
		tt := z.Next()
		t := z.Token()
		switch {
		case tt == html.ErrorToken:
			return res
		case tt == html.StartTagToken, tt == html.SelfClosingTagToken:
			enter = false
			tag = t.Data
			for _, ttt := range textTags {
				if tag == ttt {
					enter = true
					break
				}
			}
		case tt == html.TextToken:
			if enter {
				data := strings.TrimSpace(t.Data)
				if len(data) > 0 {
					for _, jchar := range data {
						if unicode.Is(unicode.Han, jchar) {
							res = append(res, string(jchar))
						}
					}
				}
			}
		}
	}
}

//help from https://gist.github.com/dhoss/7532777
//https://medium.com/@kenanbek/golang-html-tokenizer-extract-text-from-a-web-page-kanan-rahimov-8c75704bf8a3
func ScrapeHtmlFromPageConcurrent(url string, c chan []string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	textTags := []string{
		"a",
		"p", "span", "em", "string", "blockquote", "q", "cite",
		"h1", "h2", "h3", "h4", "h5", "h6", "title",
	}
	tag := ""
	enter := false
	res := make([]string, 0)
	z := html.NewTokenizer(resp.Body)

	for {
		tt := z.Next()
		t := z.Token()
		switch {
		case tt == html.ErrorToken:
			// End of document
			res = append(res, url)
			c <- res
			return
		case tt == html.StartTagToken, tt == html.SelfClosingTagToken:
			enter = false
			tag = t.Data
			for _, ttt := range textTags {
				if tag == ttt {
					enter = true
					break
				}
			}
		case tt == html.TextToken:
			if enter {
				data := strings.TrimSpace(t.Data)
				if len(data) > 0 {
					for _, jchar := range data {
						if unicode.Is(unicode.Han, jchar) {
							res = append(res, string(jchar))
						}
					}
				}
			}
		}
	}
}

////////////TWITTER ABOVE/////////////

// links := []string{
// 	"https://natgeo.nikkeibp.co.jp/?n_cid=nbpnng_ds99999",
// 	"https://mainichi.jp",
// 	"https://www.nikkei.com",
// 	"https://www.asahi.com",
// }

// kanjiChannel := make(chan []string)

// for _, link := range links {
// 	fmt.Println("firing go routine for ", link)
// 	go getHtmlFromPage(link, kanjiChannel)
// }

// for Kanji := range kanjiChannel {
// 	urlFromSite := Kanji[len(Kanji)-1]
// 	Kanji := countKanji(Kanji)
// 	sort.Sort(byFreq(Kanji))
// 	// get Top 10
// 	fmt.Println(Kanji[:10])

// 	ten := Kanji[:10]
// 	var s []string
// 	for _, v := range ten {
// 		s = append(s, v.Kanji+": "+strconv.Itoa(v.count)+" ")
// 	}
// 	kanjis := strings.Join(s, ",")
// 	// tweet, resp, err := client.Statuses.Update("今日世界 top TEN KANJI from "+urlFromSite+" "+kanjis+"", nil)
// 	// if err != nil {
// 	// 	log.Println("error making tweet", err)
// 	// }
// 	// log.Printf("%+v\n", resp)
// 	// log.Printf("%+v\n", tweet)

// }
