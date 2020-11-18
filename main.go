package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"sort"
	"strings"
	"unicode"
)

// help from https://gist.github.com/dhoss/7532777
// https://medium.com/@kenanbek/golang-html-tokenizer-extract-text-from-a-web-page-kanan-rahimov-8c75704bf8a3
func getHtmlFromPage(url string, c chan []string){ //[]string {
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
			// End of the document, we're done
			c <- res
			return
			//return res
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
type kanjiCount struct {
	Kanji string
	count int
}

func countKanji(s []string) []kanjiCount {
	words := s
	wordCount := make(map[string]int)
	kanjis := make([]kanjiCount, 0)

	for i := range words {
		wordCount[words[i]]++
	}
	for k, v := range wordCount {
		kanjis = append(kanjis, kanjiCount{k,v} )
	}
	return kanjis
}

type byFreq []kanjiCount

func (a byFreq) Len() int {
	return len(a)
}
func (a byFreq) Swap(i, j int) {
	a[i], a[j] =a[j], a[i]
}
func (a byFreq) Less(i, j int) bool { 
	return a[i].count > a[j].count
}




func main() {
	links := []string{
		"https://natgeo.nikkeibp.co.jp/?n_cid=nbpnng_ds99999",	
		"https://mainichi.jp",
		"https://www.nikkei.com",
		"https://www.asahi.com",
		"https://www.yomiuri.co.jp",
		"https://www.kobe-np.co.jp",
		"https://www.kyoto-np.co.jp",

	}

	c := make(chan []string )

	for _, link := range links {
		fmt.Println("firing go routine for ", link)
		go getHtmlFromPage(link, c)
	}

	for Kanji := range c {
		//fmt.Println(Kanji)
		Kanji := countKanji(Kanji)
		sort.Sort(byFreq(Kanji))
	// get Top 10
	fmt.Println(Kanji[:10])
	}
	fmt.Println("Exiting")



	//res := getHtmlFromPage("https://natgeo.nikkeibp.co.jp/?n_cid=nbpnng_ds99999")

	// Kanji := countKanji(res)
	// sort.Sort(byFreq(Kanji))
	// // get Top 10
	// fmt.Println(Kanji[:10])

}
