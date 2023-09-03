package passing_webtest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func NewWorker(limiter <-chan time.Time, startURL, finalURL string) (*Worker, error)
func TestNewWorker(t *testing.T) {}

// func (w *Worker) PassingTest() error
func TestPassingTest(t *testing.T) {}

// func formAnswersForSending(body io.ReadCloser) (url.Values, error)
func TestFormAnswersForSending(t *testing.T) {}

// func parsingHTMLPage(r io.Reader) (map[string]string, error)
func TestParsingHTMLPage(t *testing.T) {}

// func findValuesForQuestionOptions(n *html.Node, questionOptions map[string][]string)
func TestFindValuesForQuestionOptions(t *testing.T) {}

func TestFormAnswers(t *testing.T) {
	for _, tc := range []struct {
		name     string
		strings  []string
		expected string
	}{
		{
			name:     "Empty slice",
			strings:  []string{},
			expected: "test",
		},
		{
			name:     "Non-empty slice",
			strings:  []string{"asdfgh", "123", "qwerty"},
			expected: "qwerty",
		},
	} {
		{
			t.Run(tc.name, func(t *testing.T) {
				assert.Equal(t, tc.expected, formAnswers(tc.strings))
			})
		}
	}
}

// func (w *Worker) getPage(url string) (*http.Response, error)
// func (w *Worker) postAnswers(url string, data url.Values) (*http.Response, error)
