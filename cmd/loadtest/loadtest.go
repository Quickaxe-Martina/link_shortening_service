package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

// JSONGenerateURLRequest model
// generate:reset
type JSONGenerateURLRequest struct {
	URL string `json:"url"`
}

// JSONGenerateURLResponse model
// generate:reset
type JSONGenerateURLResponse struct {
	Result string `json:"result"`
}

func main() {
	var (
		target      string
		requests    int
		concurrency int
	)

	flag.StringVar(&target, "target", "http://localhost:8080", "Server URL")
	flag.IntVar(&requests, "n", 10000, "Requests per generator")
	flag.IntVar(&concurrency, "c", 50, "Workers per endpoint")
	flag.Parse()

	client := resty.New()
	client.SetRedirectPolicy(resty.NoRedirectPolicy())

	codes := make(chan string, requests*2)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		runLoad(requests, concurrency, func() {
			doGenerate(client, target, codes)
		})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		runLoad(requests, concurrency, func() {
			doJSONGenerate(client, target, codes)
		})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		runRedirectLoad(concurrency, client, target, codes)
	}()

	wg.Wait()
}

func runLoad(requests, workers int, fn func()) {
	var wg sync.WaitGroup
	tasks := make(chan struct{}, requests)

	for i := 0; i < requests; i++ {
		tasks <- struct{}{}
	}
	close(tasks)

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for range tasks {
				fn()
			}
		}()
	}
	wg.Wait()
}

func runRedirectLoad(workers int, client *resty.Client, target string, codes <-chan string) {
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for code := range codes {
				doRedirect(client, target, code)
			}
		}()
	}

	wg.Wait()
}

func doGenerate(client *resty.Client, target string, out chan<- string) {
	url := fmt.Sprintf("https://p-%d.com", time.Now().UnixNano())

	resp, err := client.R().
		SetHeader("Content-Type", "text/plain").
		SetBody(url).
		Post(target + "/")

	if err != nil {
		return
	}

	if resp.StatusCode() == 201 {
		out <- string(resp.Body())
	}
}

func doJSONGenerate(client *resty.Client, target string, out chan<- string) {
	url := fmt.Sprintf("https://j-%d.com", time.Now().UnixNano())

	req := JSONGenerateURLRequest{URL: url}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Post(target + "/api/shorten")

	if err != nil {
		return
	}

	if resp.StatusCode() == 201 {
		var r JSONGenerateURLResponse
		if json.Unmarshal(resp.Body(), &r) == nil {
			out <- r.Result
		}
	}
}

func doRedirect(client *resty.Client, target, code string) {
	_, _ = client.R().Get(target + "/" + code)
}
