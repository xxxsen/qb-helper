package matcher

import (
	"fmt"
	"regexp"
	"strings"
)

func init() {
	regist("regex", NewRegMatcher)
	regist("full", NewFullMatcher)
	regist("", NewFullMatcher)
	regist("prefix", NewPrefixMatcher)
	regist("suffix", NewSuffixMatcher)
}

type Matcher interface {
	IsMatch(wording string) bool
}

type RegMatcher struct {
	key *regexp.Regexp
}

func NewRegMatcher(key string, args interface{}) (Matcher, error) {
	reg := regexp.MustCompile(key)
	return &RegMatcher{key: reg}, nil
}

func Build(name string, key string, args interface{}) (Matcher, error) {
	if len(key) == 0 {
		return nil, fmt.Errorf("invalid key")
	}
	cr, err := get(name)
	if err != nil {
		return nil, err
	}
	return cr(key, args)
}

type Creator func(key string, args interface{}) (Matcher, error)

var mp = make(map[string]Creator)

func get(name string) (Creator, error) {
	if c, ok := mp[name]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("not found")
}

func regist(name string, creator Creator) {
	mp[name] = creator
}

func (m *RegMatcher) IsMatch(wording string) bool {
	return m.key.MatchString(wording)
}

type FullMatcher struct {
	key string
}

func (m *FullMatcher) IsMatch(wording string) bool {
	return strings.EqualFold(wording, m.key)
}

func NewFullMatcher(key string, args interface{}) (Matcher, error) {
	return &FullMatcher{key: key}, nil
}

type PrefixMatcher struct {
	key string
}

func (m *PrefixMatcher) IsMatch(wording string) bool {
	return strings.HasPrefix(strings.ToLower(wording), m.key)
}

func NewPrefixMatcher(key string, args interface{}) (Matcher, error) {
	return &PrefixMatcher{key: strings.ToLower(key)}, nil
}

type SuffixMatcher struct {
	key string
}

func (m *SuffixMatcher) IsMatch(wording string) bool {
	return strings.HasSuffix(strings.ToLower(wording), m.key)
}

func NewSuffixMatcher(key string, args interface{}) (Matcher, error) {
	return &PrefixMatcher{key: strings.ToLower(key)}, nil
}
