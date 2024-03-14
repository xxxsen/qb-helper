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
		"2001:db8::/66",
		"2001:db8:0:0:8000::5",
	})
	assert.NoError(t, err)
	assert.True(t, set.isMatch("127.0.0.1"))
	assert.True(t, set.isMatch("192.168.50.6"))
	assert.True(t, set.isMatch("10.1.2.3"))
	assert.True(t, set.isMatch("2001:db8::1"))
	assert.True(t, set.isMatch("2001:db8:0:0:8000::5"))
	assert.False(t, set.isMatch("11.1.2.3"))
	assert.False(t, set.isMatch("192.167.50.1"))
	assert.False(t, set.isMatch("127.0.0.2"))
	assert.False(t, set.isMatch("2001:db8:0:0:8000::4"))
}

func TestInvalidIP(t *testing.T) {
	_, err := makeIPRuleSet([]string{
		"2001:db8.3.4::/66",
	})
	assert.Error(t, err)
	_, err = makeIPRuleSet([]string{
		"1.2.3.4.5",
	})
	assert.Error(t, err)
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
