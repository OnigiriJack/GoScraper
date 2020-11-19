This was created during my time as a student at [Code Chrysalis](https://www.codechrysalis.io/)

# GoScraper

![スクリーンショット 2020-11-19 20 24 42](https://user-images.githubusercontent.com/35797565/99660476-c94a6b00-2aa5-11eb-9f87-50a9aa855b21.png)

## About
This [project](https://kanji-counter-twitter.herokuapp.com/) scrapes websites in GOlang and posts results of the Top Ten Kanji 
of the scraped website on Twitter onto my twitter via a post Request. 
Building this application in Go allows one to utilize the power of concurrency that Go provides.

**Features**
- Scrape a webapage and rank the Kanji on the site
- Post the results on twitter
- Scrape many sites at once with (Goroutines)[https://blog.golang.org/context]

## Tech
Golang, [go-twitter](https://github.com/dghubble/go-twitter), HTML, [TwitterAPI](https://developer.twitter.com/en), Heroku

### How to Run

Retrieve dependencies
```
go get ./...
```

Compile Program
```
go build
```

Run Program
```
go run main.go
```



