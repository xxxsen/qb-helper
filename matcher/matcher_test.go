package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegMatcher(t *testing.T) {
	type pair struct {
		wording string
		match   bool
	}
	key := "7\\.\\d+\\.\\d+\\.\\d+"
	matcher, err := Build("regex", key, nil)
	assert.NoError(t, err)
	lst := []pair{
		{"7.5.6.10", true},
		{"7.10.50.336", true},
		{"8.3.45.6", false},
	}
	for _, item := range lst {
		assert.Equal(t, matcher.IsMatch(item.wording), item.match)
	}
}
