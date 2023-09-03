package utils

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func NewHTTPClient() (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %s ", err)
	}
	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}, nil
}

func ResponseToGetRequest(url string, client *http.Client) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %s ", err)
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get response to GET request: %s ", err)
	}
	return response, nil
}

func ResponseToPostRequest(questionPage string, client *http.Client, data url.Values) (*http.Response, error) {
	request, err := http.NewRequest("POST", questionPage, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request for question (%s): %s ", questionPage, err)
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to receive response to POST for question (%s)): %s ", questionPage, err)
	}
	defer response.Body.Close()
	return response, nil
}
