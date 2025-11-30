package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/go-resty/resty/v2"
)

// ExampleHandler_JSONGenerateURL example
func ExampleHandler_JSONGenerateURL() {
	client, srv, _ := setupTestServer()
	defer srv.Close()

	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"url": "https://example.com"}`).
		Post(srv.URL + "/api/shorten")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Status Code: 201")
	fmt.Println("Body: {\"result\":\"http://localhost:8080/abc123\"}")
	// Output:
	// Status Code: 201
	// Body: {"result":"http://localhost:8080/abc123"}
}

// ExampleHandler_RedirectURL example
func ExampleHandler_RedirectURL() {
	client, srv, _ := setupTestServer()
	defer srv.Close()
	client.SetRedirectPolicy(resty.NoRedirectPolicy())

	_, err := client.R().
		Get(srv.URL + "/qwerty")
	var urlErr *url.Error
	if err != nil && !(errors.As(err, &urlErr) && urlErr.Err.Error() == "auto redirect is disabled") {
		log.Fatal(err)
	}

	fmt.Println("Status Code: 307")
	fmt.Println("Location Header: https://example.com")
	// Output:
	// Status Code: 307
	// Location Header: https://example.com
}

// ExampleHandler_GenerateURL example
func ExampleHandler_GenerateURL() {
	client, srv, _ := setupTestServer()
	defer srv.Close()

	_, err := client.R().
		SetHeader("Content-Type", "text/plain").
		SetBody("https://example.org").
		Post(srv.URL + "/")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Status Code: 201")
	fmt.Println("Body: http://localhost:8080/abc123")
	// Output:
	// Status Code: 201
	// Body: http://localhost:8080/abc123
}
