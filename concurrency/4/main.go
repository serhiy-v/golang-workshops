package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const stringToSearch = "concurrency"

var sites = []string{
	"https://google.com",
	"https://itc.ua/",
	"https://twitter.com/concurrencyinc",
	"https://twitter.com/",
	"https://github.com/bradtraversy/go_restapi/blob/master/main.go",
	"https://www.youtube.com/",
	"https://postman-echo.com/get",
	"https://en.wikipedia.org/wiki/Concurrency_(computer_science)#:~:text=In%20computer%20science%2C%20concurrency%20is,without%20affecting%20the%20final%20outcome.",
}

type SiteData struct {
	data []byte
	uri  string
}

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	resultsCh := make(chan SiteData, len(sites))

	Worker(ctx, resultsCh, sites)
	Reader(ctx, stringToSearch, resultsCh)

	// give one second to validate if all other goroutines are closed
	time.Sleep(time.Second)
}

func Worker(ctx context.Context, ch chan SiteData, sites []string) {
	for _, site := range sites {
		go func(site string) {
			fmt.Println("Starting sending request to", site)
			res, err := request(ctx, site)
			if err != nil {
				return
			}
			ch <- SiteData{data: res, uri: site}
		}(site)
	}
}

func Reader(ctx context.Context, stringToSearch string, ch chan SiteData) {
	for {
		data := <-ch
		if findString(data.data, stringToSearch) {
			fmt.Println(fmt.Sprintf("'%v' string is found in: %v", stringToSearch, data.uri))
			return
		}
	}
}

func findString(data []byte, substr string) bool {
	return strings.Index(string(data), substr) != -1
}

// TODO implement function that will execute request function, will validate the output and cancel all other requests when needed page is found
// and will listen to cancellation signal from context and will exit from the func when will receive it

// TODO implement function that will perfrom request using the example under

// TODO hint request function code:
func request(ctx context.Context, uri string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return bodyBytes, nil
}
