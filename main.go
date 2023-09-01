package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func main() {
	passingTest("http://147.78.65.149/start/", "http://147.78.65.149/passed")
}

func passingTest(startURLStr, finalURLStr string) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Failed to create cookie jar: %v\n", err)
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}

	getRequest, err := http.NewRequest("GET", startURLStr, nil)
	if err != nil {
		log.Fatalf("Failed to create GET request: %v\n", err)
	}
	response, err := client.Do(getRequest)
	defer response.Body.Close()
	if err != nil {
		log.Fatalf("Failed to get response to GET request: %v\n", err)
	}

	//var location *url.URL
	//var locationError error
	location, locationError := response.Location()
	finalURL, err := url.Parse(finalURLStr)
	if err != nil {
		log.Fatalf("Failed to get URL from string %s: %v\n", finalURLStr, err)
	}
	fmt.Printf("location: %v\nfinalURL: %v\nfinalURL == location? %v\n\n", location, finalURL, finalURL == location) // delete!!!!!!!!!!!!

	for response.StatusCode == 302 {
		getRequest, err = http.NewRequest("GET", fmt.Sprint(location), nil)
		if err != nil {
			log.Fatalf("Failed to create request to %v: %v\n", location, err)
		}
		response, err = client.Do(getRequest)
		if err != nil {
			log.Fatalf("Failed to perform Get request for quesion: %v\n", err)
		}
		if response.StatusCode != 200 {
			log.Fatalf("Wrong status code of response: %v\n", response.StatusCode)
		}

		fmt.Printf("Response to GET:\nstatus: %v\nlocation: %v\n", response.StatusCode, location) // delete!!!!!!!!!!!!
		response = formAndPostAnswer(response.Body, location.String(), client)
		time.Sleep(time.Second)
		location, locationError = response.Location()
		fmt.Printf("Response to POST:\nstatus: %v\nlocation: %v\n\n", response.StatusCode, location) // delete!!!!!!!!!!!!
		if locationError != nil {
			log.Fatalf("Failed to get location in response: %v\n", err)
		}
		if location.String() == finalURLStr {
			fmt.Println("Test successfully passed")
			break
		}
	}
}

func formAndPostAnswer(body io.ReadCloser, questionPage string, client *http.Client) *http.Response {
	keysAndAnswers := parsingHTMLPage(body)
	data := url.Values{}
	for k, v := range keysAndAnswers {
		data.Set(k, v)
	}
	postRequest, err := http.NewRequest("POST", questionPage, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalf("Failed to create POST request for question (%s): %v\n", questionPage, err)
	}
	postRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	postResponse, err := client.Do(postRequest)
	defer postResponse.Body.Close()
	if err != nil {
		log.Fatalf("Failed to receive response to POST for question (%s)): %v\n", questionPage, err)
	}
	return postResponse
}

func parsingHTMLPage(r io.Reader) map[string]string {
	parsedNode, err := html.Parse(r)
	if err != nil {
		log.Fatal(err)
	}

	questionOptions := make(map[string][]string)
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "input" {
			var keyToAdd string
			for _, attr := range n.Attr {
				if attr.Key == "type" {
					continue
				}
				if attr.Key == "name" {
					keyToAdd = attr.Val
					if _, isFound := questionOptions[keyToAdd]; !isFound {
						questionOptions[keyToAdd] = make([]string, 0)
					}
				}
				if attr.Key == "value" {
					questionOptions[keyToAdd] = append(questionOptions[keyToAdd], attr.Val)
				}
			}
		}
		if n.Type == html.ElementNode && (n.Data == "select") {
			var keyToAdd string
			for _, attr := range n.Attr {
				keyToAdd = attr.Val
				if _, found := questionOptions[keyToAdd]; !found {
					questionOptions[keyToAdd] = make([]string, 0)
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Data != "option" {
					continue
				}
				for _, attr := range c.Attr {
					if attr.Val != "" {
						questionOptions[keyToAdd] = append(questionOptions[keyToAdd], attr.Val)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(parsedNode)

	answersByKeys := make(map[string]string)
	for k, v := range questionOptions {
		answersByKeys[k] = formAnswer(v)
	}
	return answersByKeys
}

func formAnswer(options []string) string {
	if len(options) == 0 {
		return "test"
	}
	return findLongestStringInSlice(options)
}

func findLongestStringInSlice(str []string) string {
	if len(str) == 0 {
		return ""
	}
	maxLen := len(str[0])
	answer := str[0]
	for i := 0; i < len(str); i++ {
		if len(str[i]) > maxLen {
			maxLen = len(str[i])
			answer = str[i]
		}
		if len(str[i]) == maxLen && str[i] > answer {
			maxLen = len(str[i])
			answer = str[i]
		}
	}
	return answer
}
