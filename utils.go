package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

func Test() {
	fmt.Println("I am imported")
}

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
