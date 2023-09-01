package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

func main() {
	startURL := os.Getenv("START_PAGE")
	finalURL := os.Getenv("FINAL_PAGE")
	qtyOfThreads, err := strconv.Atoi(os.Getenv("QTY_OF_THREADS"))
	if err != nil {
		log.Fatalf("Failed to parse quantity of threads: %s\n", os.Getenv("QTY_OF_THREADS"))
	}

	wg := sync.WaitGroup{}
	wg.Add(qtyOfThreads)

	for i := 0; i < qtyOfThreads; i++ {
		go func(n int) {
			result := passingTest(startURL, finalURL)
			log.Printf("Process #%d: ", n)
			if result == nil {
				log.Println("Test successfully passed")
			} else {
				log.Println(result)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func passingTest(startURL, finalURL string) error {
	client, err := newHTTPClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %s ", err)
	}

	response, err := responseToHTTPGetRequest(startURL, client)
	if err != nil {
		return fmt.Errorf("failed to get response for start page: %s ", err)
	}
	defer response.Body.Close()

	location, locationError := response.Location()
	if err != nil {
		return fmt.Errorf("Failed to get URL from string %s: %s ", finalURL, err)
	}

	for locationError == nil && response.StatusCode == 302 && location.String() != finalURL {
		//fmt.Printf("%v\n", location) // TODO: delete!!!!!!!!!!!!!!!!!!!!!!!
		response, err = responseToHTTPGetRequest(location.String(), client)
		if err != nil {
			return fmt.Errorf("failed to get response for page with question: %s ", err)
		}
		if response.StatusCode != 200 {
			return fmt.Errorf("Wrong status code of response: %d ", response.StatusCode)
		}

		data, dataError := formAnswersForSending(response.Body)
		if dataError != nil {
			return fmt.Errorf("failed to form data with answers: %s ", dataError)
		}
		response, err = postAnswer(location.String(), client, data)
		if err != nil {
			return fmt.Errorf("failed to post answers: %s ", err)
		}
		// For not to exceed 3 requests per second:
		time.Sleep(time.Second)
		location, locationError = response.Location()
	}
	if response.StatusCode != 302 {
		return fmt.Errorf("wrong status code of response: %d ", response.StatusCode)
	}
	if locationError != nil {
		return fmt.Errorf("failed to get location in response: %s ", err)
	}
	return nil
}

func newHTTPClient() (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create cookie jar: %s ", err)
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}
	return client, nil
}

func responseToHTTPGetRequest(url string, client *http.Client) (*http.Response, error) {
	getRequest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %s ", err)
	}
	response, err := client.Do(getRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to get response to GET request: %s ", err)
	}
	return response, nil
}

func formAnswersForSending(body io.ReadCloser) (url.Values, error) {
	keysAndAnswers, err := parsingHTMLPage(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML page: %s ", err)
	}
	data := url.Values{}
	for k, v := range keysAndAnswers {
		data.Set(k, v)
	}
	return data, nil
}

func postAnswer(questionPage string, client *http.Client, data url.Values) (*http.Response, error) {
	postRequest, err := http.NewRequest("POST", questionPage, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request for question (%s): %s ", questionPage, err)
	}
	postRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	postResponse, err := client.Do(postRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to receive response to POST for question (%s)): %s ", questionPage, err)
	}
	defer postResponse.Body.Close()
	return postResponse, nil
}

func parsingHTMLPage(r io.Reader) (map[string]string, error) {
	parsedNode, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	questionOptions := make(map[string][]string)
	findValuesForQuestionOptions(parsedNode, questionOptions)

	answersByKeys := make(map[string]string)
	for k, v := range questionOptions {
		answersByKeys[k] = formAnswers(v)
	}
	return answersByKeys, nil
}

func findValuesForQuestionOptions(n *html.Node, questionOptions map[string][]string) {
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
	if n.Type == html.ElementNode && n.Data == "select" {
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
		findValuesForQuestionOptions(c, questionOptions)
	}
}

func formAnswers(options []string) string {
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
