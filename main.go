package main

import (
	"fmt"
	"golang.org/x/net/html"
	//"io"
	"log"
	"net/http"
	"strings"
)
// help from https://gist.github.com/dhoss/7532777
 // https://medium.com/@kenanbek/golang-html-tokenizer-extract-text-from-a-web-page-kanan-rahimov-8c75704bf8a3
func getHtmlFromPage(url string) []string {
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
			return res
		case tt == html.StartTagToken, tt ==html.SelfClosingTagToken:
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
				//fmt.Println("I am a jchar  ",jchar )
				if len(data) > 0 {
					jchar := strings.Split(data,"") 
					for _, kanji := range jchar {
						//fmt.Println("I am a jchar  ",kanji )
						res = append(res, kanji)
					}
				}
			}

		}
	}
}


func wordCount(s []string) map[string]int {
    words := s
    wordCount := make(map[string]int)
    for i := range words {
        wordCount[words[i]]++
    }
        
    return wordCount
}

func main() {

   res := getHtmlFromPage("https://www.nikkei.com/")

   count := wordCount(res)
   fmt.Println(count)

    // fmt.Println("result html",res )
	//  for i,v := range strings.Split(res[0],""){
 	// fmt.Printf("element info: %d %d\n", i, v)
	//   }
}

