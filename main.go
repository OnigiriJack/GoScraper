package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
	"golang.org/x/net/html"
)

type kanjiCount struct {
	Kanji string
	count int
}

type byFreq []kanjiCount

func (a byFreq) Len() int {
	return len(a)
}
func (a byFreq) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a byFreq) Less(i, j int) bool {
	return a[i].count > a[j].count
}

// help from https://gist.github.com/dhoss/7532777
// https://medium.com/@kenanbek/golang-html-tokenizer-extract-text-from-a-web-page-kanan-rahimov-8c75704bf8a3
// func scrapeHtmlFromPageConcurrent(url string, c chan []string) {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()
// 	textTags := []string{
// 		"a",
// 		"p", "span", "em", "string", "blockquote", "q", "cite",
// 		"h1", "h2", "h3", "h4", "h5", "h6", "title",
// 	}
// 	tag := ""
// 	enter := false
// 	res := make([]string, 0)
// 	z := html.NewTokenizer(resp.Body)

// 	for {
// 		tt := z.Next()
// 		t := z.Token()
// 		switch {
// 		case tt == html.ErrorToken:
// 			// End of document
// 			res = append(res, url)
// 			c <- res
// 			return
// 		case tt == html.StartTagToken, tt == html.SelfClosingTagToken:
// 			enter = false
// 			tag = t.Data
// 			for _, ttt := range textTags {
// 				if tag == ttt {
// 					enter = true
// 					break
// 				}
// 			}
// 		case tt == html.TextToken:
// 			if enter {
// 				data := strings.TrimSpace(t.Data)
// 				if len(data) > 0 {
// 					for _, jchar := range data {
// 						if unicode.Is(unicode.Han, jchar) {
// 							res = append(res, string(jchar))
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// }

func scrapeHtmlFromPage(url string) []string {
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

func countKanji(s []string) []kanjiCount {
	words := s
	wordCount := make(map[string]int)
	kanjis := make([]kanjiCount, 0)

	for i := range words {
		wordCount[words[i]]++
	}
	for k, v := range wordCount {
		kanjis = append(kanjis, kanjiCount{k, v})
	}
	return kanjis
}

func formatCountForTwitter(allKanji []string, url string) string {
	Kanji := countKanji(allKanji)
	sort.Sort(byFreq(Kanji))
	// get Top 10
	ten := Kanji[:10]
	//fmt.Println(Kanji[:10])
	var s []string
	for _, v := range ten {
		s = append(s, v.Kanji+": "+strconv.Itoa(v.count)+" ")
	}
	kanjis := strings.Join(s, ",")
	kanjis = kanjis + url
	return kanjis
}

func twitterSend(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/scrape" {
		http.Error(w, "404 Not Found.", http.StatusNotFound)
		return
	}
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "parseform() err: %v", err)
		return
	}

	enverr := godotenv.Load()
	if enverr != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Fprintf(w, "Post from website r.PostForm = %v\n", r.PostForm)
	url := r.FormValue("url")

	allKanji := ScrapeHtmlFromPage(url)
	countForTwitter := formatCountForTwitter(allKanji, url)

	///////////TWITTER CONfIGS//////////////
	config := oauth1.NewConfig(os.Getenv("TWITTER_CONSUMER_KEY"), os.Getenv("TWITTER_CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	tweet, resp, err := client.Statuses.Update(countForTwitter, nil)
	if err != nil {
		log.Println("error making tweet", err)
	}
	log.Printf("%+v\n", resp)
	log.Printf("%+v\n", tweet)
	http.Redirect(w, r, "/", 301)
}

func main() {
	Test()
	enverr := godotenv.Load()
	if enverr != nil {
		log.Fatal("Error loading .env file")
	}
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/scrape", twitterSend)
	fmt.Printf(" hosting at %s ", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))

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

}
