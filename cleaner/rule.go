package cleaner

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	defaultFullPrefix   = "full:"
	defaultRegexpPrefix = "regexp:"
)

type uaRuleSet struct {
	regexp []*regexp.Regexp
	full   map[string]struct{}
	prefix []string
}

func (r *uaRuleSet) isMatch(ua string) bool {
	if _, ok := r.full[ua]; ok {
		return true
	}
	for _, prefix := range r.prefix {
		if strings.HasPrefix(ua, prefix) {
			return true
		}
	}
	for _, reg := range r.regexp {
		if reg.MatchString(ua) {
			return true
		}
	}
	return false
}

func makeUserAgentRuleSet(rs []string) (*uaRuleSet, error) {
	set := &uaRuleSet{
		regexp: nil,
		full:   make(map[string]struct{}),
		prefix: nil,
	}
	for _, r := range rs {
		if strings.HasPrefix(r, defaultFullPrefix) {
			set.full[r[len(defaultFullPrefix):]] = struct{}{}
			continue
		}
		if strings.HasPrefix(r, defaultRegexpPrefix) {
			sreg := r[len(defaultRegexpPrefix):]
			expr, err := regexp.Compile(sreg)
			if err != nil {
				return nil, fmt.Errorf("compile regexp failed, err:%w", err)
			}
			set.regexp = append(set.regexp, expr)
		}
		set.prefix = append(set.prefix, r)
	}
	return set, nil
}
