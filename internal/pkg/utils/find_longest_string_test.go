package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindLongestStringInSlice(t *testing.T) {
	for _, tc := range []struct {
		name     string
		strings  []string
		expected string
	}{
		{
			name:     "Empty slice",
			strings:  []string{},
			expected: "",
		},
		{
			name:     "One empty string",
			strings:  []string{""},
			expected: "",
		},
		{
			name:     "One string",
			strings:  []string{"qwerty"},
			expected: "qwerty",
		},
		{
			name:     "One longest string",
			strings:  []string{"q", "qwerty", "", "123"},
			expected: "qwerty",
		},
		{
			name:     "Two strings with equal length",
			strings:  []string{"q", "asdfgh", "qwerty", "123"},
			expected: "qwerty",
		},
		{
			name:     "Numbers and non-ASCII symbols",
			strings:  []string{"привет", "!№;%:?*()", "1234567890", "我的名字是艾拉"},
			expected: "1234567890",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, FindLongestStringInSlice(tc.strings))
		})
	}
}
