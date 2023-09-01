package passing_webtest

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

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

func postData(questionPage string, client *http.Client, data url.Values) (*http.Response, error) {
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
