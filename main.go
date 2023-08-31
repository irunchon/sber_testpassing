package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
)

func main() {
	requestURL := fmt.Sprintf("http://147.78.65.149/start/")
	request, err := http.NewRequest("GET", requestURL, nil)

	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Printf("Failed to create cookie jar: %v\n", err)
		return
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Failed to perform Get request: %v\n", err)
		return
	}
	fmt.Printf("Response:\nstatus code: %v\nheader: %v\ncookies: %v\n",
		response.StatusCode, response.Header, response.Cookies())

	fmt.Printf("*****\n")

	request, err = http.NewRequest("GET", "http://147.78.65.149/question/1", nil)
	if err != nil {
		fmt.Printf("Failed to create request to question 1: %v\n", err)
		return
	}

	response, err = client.Do(request)
	if err != nil {
		fmt.Printf("Failed to perform Get request for quesion: %v\n", err)
		return
	}
	fmt.Printf("Response:\nstatus code: %v\nheader: %v\ncookies: %v\n",
		response.StatusCode, response.Header, response.Cookies())
}
