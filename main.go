package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	requestURL := fmt.Sprintf("http://147.78.65.149/start/")
	request, err := http.NewRequest("GET", requestURL, nil)

	if err != nil {
		log.Fatalf("Failed to create request: %v\n", err)
	}

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
	response, err := client.Do(request)
	defer response.Body.Close()
	if err != nil {
		log.Fatalf("Failed to perform Get request: %v\n", err)
	}
	fmt.Printf("Response:\nstatus code: %v\nheader: %v\ncookies: %v\n",
		response.StatusCode, response.Header, response.Cookies())

	fmt.Printf("*****\n")

	// ------------ QUESTION 1 ------------ //
	request, err = http.NewRequest("GET", "http://147.78.65.149/question/1", nil)
	if err != nil {
		log.Fatalf("Failed to create request to question 1: %v\n", err)
	}

	response, err = client.Do(request)
	if err != nil {
		log.Fatalf("Failed to perform Get request for quesion: %v\n", err)
	}
	fmt.Printf("Response:\nstatus code: %v\nheader: %v\ncookies: %v\n",
		response.StatusCode, response.Header, response.Body)

	// ------------ POST ANSWERS ------------ //
	keysAndAnswers := parsingHTMLPage(response.Body)
	data := url.Values{}
	for k, v := range keysAndAnswers {
		data.Set(k, v)
	}
	encodedData := data.Encode()
	fmt.Printf("\nencodedData: %v\n", encodedData)
	postRequest, err := http.NewRequest("POST", "http://147.78.65.149/question/1", strings.NewReader(encodedData))
	if err != nil {
		log.Fatalf("Failed to create POST request for question 1: %v\n", err)
	}
	postRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	postResponse, err := client.Do(postRequest)
	defer postResponse.Body.Close()
	if err != nil {
		log.Fatalf("Failed to receive response to POST for question 1: %v\n", err)
	}
	fmt.Printf("\nResponse to POST:\n%v\n", postResponse.Status)
	body, err := ioutil.ReadAll(postResponse.Body)
	if err != nil {
		log.Fatalf("Failed to read body: %v\n", err)

	}
	fmt.Printf("body:\n%s\n", body)
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
	//------------ DELETE THIS: ------------ //
	fmt.Println("\nquestionOptions:")
	for k, v := range questionOptions {
		fmt.Printf("key: %s\tvalues: %v\n", k, v)
	}
	fmt.Println("question - qnswer:")
	for k, v := range questionOptions {
		fmt.Printf("%s: %s\n", k, formAnswer(v)) /////////////////////
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
