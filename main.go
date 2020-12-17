package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
)

type kanjiCount struct {
	Kanji string
	count int
}

type kanjiData struct {
	url    string
	kanjis kanjiCount
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
	ten := Kanji[:10]
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

func goRoutine(w http.ResponseWriter, r *http.Request) {

	tmplt := template.New("goroutine.html")       // create a new template with some name
	tmplt, _ = tmplt.ParseFiles("goroutine.html") // parse some content and generate a template, which is an internal representation

	links := []string{
		"https://natgeo.nikkeibp.co.jp/?n_cid=nbpnng_ds99999",
		"https://mainichi.jp",
		"https://www.nikkei.com",
		"https://www.asahi.com",
	}

	kanjiChannel := make(chan []string)

	for _, link := range links {
		fmt.Println("firing go routine for ", link)
		go ScrapeHtmlFromPageConcurrent(link, kanjiChannel)
	}
	result := make([]kanjiData, 0)
	for Kanji := range kanjiChannel {
		urlFromSite := Kanji[len(Kanji)-1]
		Kanji := countKanji(Kanji)
		sort.Sort(byFreq(Kanji))
		result = append(result, kanjiData{urlFromSite, Kanji})
		fmt.Println(Kanji[:10])
	}
	tmplt.Execute(w, result)
}

func main() {

	enverr := godotenv.Load()
	if enverr != nil {
		log.Fatal("Error loading .env file")
	}
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/go", goRoutine)
	http.HandleFunc("/scrape", twitterSend)
	fmt.Printf(" hosting at %s ", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
