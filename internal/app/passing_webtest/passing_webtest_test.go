package passing_webtest

import (
	"strings"
	"testing"

	"golang.org/x/net/html"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
func TestFindValuesForQuestionOptions(t *testing.T) {
	t.Run("One INPUT[@type=text] field", func(t *testing.T) {
		const htm = `<p>1) <input type="text" name="8bFGjvisdkL5v4V2"></p>`
		node, err := html.Parse(strings.NewReader(htm))
		require.NoError(t, err)

		questionOptions := make(map[string][]string)
		expectedOptions := map[string][]string{
			"8bFGjvisdkL5v4V2": []string{},
		}
		findValuesForQuestionOptions(node, questionOptions)
		assert.True(t, areMapsEqual(t, expectedOptions, questionOptions))
	})
	t.Run("One INPUT[@type=radio] field", func(t *testing.T) {
		const htm = `<p>2) <input type="radio" name="6V7DPsPzmcGv6hKJ" value="g5LGQB6Qb8zR">g5LGQB6Qb8zR &nbsp;&nbsp;&nbsp;
<input type="radio" name="6V7DPsPzmcGv6hKJ" value="sC3F">sC3F &nbsp;&nbsp;&nbsp;
<input type="radio" name="6V7DPsPzmcGv6hKJ" value="KVh8yzk">KVh8yzk &nbsp;&nbsp;&nbsp;</p>`
		node, err := html.Parse(strings.NewReader(htm))
		require.NoError(t, err)

		questionOptions := make(map[string][]string)
		expectedOptions := map[string][]string{
			"6V7DPsPzmcGv6hKJ": []string{
				"g5LGQB6Qb8zR",
				"sC3F",
				"KVh8yzk",
			},
		}
		findValuesForQuestionOptions(node, questionOptions)
		assert.True(t, areMapsEqual(t, expectedOptions, questionOptions))
	})
	t.Run("One SELECT field", func(t *testing.T) {
		const htm = `<p>4) <select name="I7bQvTSIfoCVwN9Y">
<option value=""></option><option value="nQgq">nQgq</option>
<option value="yyWD">yyWD</option>
<option value="xRCCvB">xRCCvB</option>
<option value="nWwu">nWwu</option>
<option value="FrfJFGeBdZf">FrfJFGeBdZf</option></select></p>`
		node, err := html.Parse(strings.NewReader(htm))
		require.NoError(t, err)

		questionOptions := make(map[string][]string)
		expectedOptions := map[string][]string{
			"I7bQvTSIfoCVwN9Y": []string{
				"nQgq",
				"yyWD",
				"xRCCvB",
				"nWwu",
				"FrfJFGeBdZf",
			},
		}
		findValuesForQuestionOptions(node, questionOptions)
		assert.True(t, areMapsEqual(t, expectedOptions, questionOptions))
	})
	t.Run("Mix of the fields", func(t *testing.T) {
		const htm = `<p>1) <input type="text" name="8bFGjvisdkL5v4V2"></p>
<p>2) <input type="radio" name="6V7DPsPzmcGv6hKJ" value="g5LGQB6Qb8zR">g5LGQB6Qb8zR &nbsp;&nbsp;&nbsp;
<input type="radio" name="6V7DPsPzmcGv6hKJ" value="sC3F">sC3F &nbsp;&nbsp;&nbsp;
<input type="radio" name="6V7DPsPzmcGv6hKJ" value="KVh8yzk">KVh8yzk &nbsp;&nbsp;&nbsp;</p>
<p>4) <select name="I7bQvTSIfoCVwN9Y">
<option value=""></option><option value="nQgq">nQgq</option>
<option value="yyWD">yyWD</option>
<option value="xRCCvB">xRCCvB</option>
<option value="nWwu">nWwu</option>
<option value="FrfJFGeBdZf">FrfJFGeBdZf</option></select></p>`
		node, err := html.Parse(strings.NewReader(htm))
		require.NoError(t, err)

		questionOptions := make(map[string][]string)
		expectedOptions := map[string][]string{
			"I7bQvTSIfoCVwN9Y": []string{
				"nQgq",
				"yyWD",
				"xRCCvB",
				"nWwu",
				"FrfJFGeBdZf",
			},
			"6V7DPsPzmcGv6hKJ": []string{
				"g5LGQB6Qb8zR",
				"sC3F",
				"KVh8yzk",
			},
			"8bFGjvisdkL5v4V2": []string{},
		}
		findValuesForQuestionOptions(node, questionOptions)
		assert.True(t, areMapsEqual(t, expectedOptions, questionOptions))
	})
}

func areMapsEqual(t *testing.T, m1, m2 map[string][]string) bool {
	t.Helper()
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		if v2, isFound := m2[k]; !isFound || !assert.ElementsMatch(t, v2, v1) {
			return false
		}
	}
	return true
}

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
