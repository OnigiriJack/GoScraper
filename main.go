package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//chGoodUrls chan string,
func fetchUrl(url string, chFailedUrls chan string, chIsFinished chan bool) {
	client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("If-None-Match", `W/"wyzzy"`)
	resp, err := client.Do(req)
	//sends a request and gets a response
	// executed when fetch url outer scope is finished
	defer func() {
		chIsFinished <- true
	}()

	if err != nil || resp.StatusCode != 200 {
		chFailedUrls <- url
		return
	}

	// 	fmt.Println("here")
	// 	chGoodUrls <- url

	//}
}

func getHtmlFromPage(url string) {
	resp, _ := http.Get(url)
	bytes, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("HTML:\n\n", string(bytes), bytes)
	resp.Body.close()
}
func main() {

	urlsList := [10]string{
		"https://natgeo.nikkeibp.co.jp",
		"http://example2.com",
		"http://example3.com",
		"http://example4.com",
		// "http://example5.com",
		// "http://example10.com",
		// "http://example20.com",
		// "http://example30.com",
		"http://example40.com",
		// "http://example50.com",
	}

	chFailedUrls := make(chan string)
	//chGoodUrls := make(chan string)
	chIsFinished := make(chan bool)
	// do each fetch in the list concurrently
	for _, url := range urlsList {
		fmt.Println(url)
		go fetchUrl(url, chFailedUrls, chIsFinished)
	}

	failedUrls := make([]string, 0)
	//goodUrls := make([]string, 0)

	for i := 0; i < len(urlsList); {
		select {
		case url := <-chFailedUrls:
			failedUrls = append(failedUrls, url)
		//case url := <-chGoodUrls:
		//	goodUrls = append(goodUrls, url)
		case <-chIsFinished:
			i++
		}

	}
	fmt.Println("could not fetch these urls: ", failedUrls)
	fmt.Println("could fetch these urls: ", failedUrls)

}
