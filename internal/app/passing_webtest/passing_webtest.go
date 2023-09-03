package passing_webtest

import (
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/irunchon/sber_testpassing/internal/pkg/utils"

	"golang.org/x/net/html"
)

type Worker struct {
	limiter <-chan time.Time
}

func NewWorker(limiter <-chan time.Time) *Worker {
	return &Worker{limiter: limiter}
}

func (w *Worker) PassingTest(startURL, finalURL string) error {
	client, err := utils.NewHTTPClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %s ", err)
	}

	<-w.limiter // Rate limiter for 3 requests per second
	response, err := utils.ResponseToHTTPGetRequest(startURL, client)
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
		<-w.limiter // Rate limiter for 3 requests per second
		response, err = utils.ResponseToHTTPGetRequest(location.String(), client)
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

		<-w.limiter // Rate limiter for 3 requests per second
		response, err = utils.PostData(location.String(), client, data)
		if err != nil {
			return fmt.Errorf("failed to post answers: %s ", err)
		}

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
	return utils.FindLongestStringInSlice(options)
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
