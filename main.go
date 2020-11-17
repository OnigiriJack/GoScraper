package main

import (
	"fmt"
	"net/http"
)

func fetchUrl(url string, chFailedUrls chan string, chIsFinished chan bool){
	client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("If-None-Match", `W/"wyzzy"`)
	resp, err := client.Do(req) //sends a request and gets a response

	defer func() {
		chIsFinished <- true
	}()

	if err!= nil || resp.StatusCode != 200 {
	 chFailedUrls <- url
	 return
 }
}


func main() {

	urlsList := [10]string{
			 "http://example1.com",
			 "http://example2.com",
			 "http://example3.com",
			 "http://example4.com",
			 "http://example5.com",
			 "http://example10.com",
			 "http://example20.com",
			 "http://example30.com",
			 "http://example40.com",
			 "http://example50.com",
	 }

	 chFailedUrls := make(chan string)
	 chIsFinished:= make(chan bool)

	 for _, url := range urlsList {
		 go fetchUrl(url, chFailedUrls, chIsFinished)
	 }

	 failedUrls := make([]string, 0)
	 for i := 0 ; i < len(urlsList); {
		 select {
		 case url := <-chFailedUrls:
			 failedUrls = append(failedUrls, url)
		 case <- chIsFinished:
			 i++
		 }
	 }
	 fmt.Println("could not fetch these urls: ", failedUrls)

}
