package cleaner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIP(t *testing.T) {
	set, err := makeIPRuleSet([]string{
		"127.0.0.1",
		"192.168.0.0/16",
		"10.0.0.0/8",
	})
	assert.NoError(t, err)
	assert.True(t, set.isMatch("127.0.0.1"))
	assert.True(t, set.isMatch("192.168.50.6"))
	assert.True(t, set.isMatch("10.1.2.3"))
	assert.False(t, set.isMatch("11.1.2.3"))
	assert.False(t, set.isMatch("192.167.50.1"))
	assert.False(t, set.isMatch("127.0.0.2"))

}

func TestStrSet(t *testing.T) {
	set, err := makeStrRuleSet([]string{
		"prefix:abc.com",
		"ccc.com",
		"regexp:^hhhh.*$",
	})
	assert.NoError(t, err)
	assert.True(t, set.isMatch("abc.com.qqq.cc"))
	assert.True(t, set.isMatch("ccc.com"))
	assert.True(t, set.isMatch("hhhhru ok ?"))
	assert.False(t, set.isMatch("aaa.abc.com"))
	assert.False(t, set.isMatch("cccc.com"))
	assert.False(t, set.isMatch("ahhhhqqq"))
}
