package cleaner

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

const (
	defaultPrefixPrefix = "prefix:"
	defaultRegexpPrefix = "regexp:"
)

type strRuleSet struct {
	regexp []*regexp.Regexp
	full   map[string]struct{}
	prefix []string
}

func (r *strRuleSet) isMatch(str string) bool {
	if _, ok := r.full[str]; ok {
		return true
	}
	for _, prefix := range r.prefix {
		if strings.HasPrefix(str, prefix) {
			return true
		}
	}
	for _, reg := range r.regexp {
		if reg.MatchString(str) {
			return true
		}
	}
	return false
}

func makeStrRuleSet(rs []string) (*strRuleSet, error) {
	set := &strRuleSet{
		regexp: nil,
		full:   make(map[string]struct{}),
		prefix: nil,
	}
	for _, r := range rs {
		if strings.HasPrefix(r, defaultPrefixPrefix) {
			set.prefix = append(set.prefix, r[len(defaultPrefixPrefix):])
			continue
		}
		if strings.HasPrefix(r, defaultRegexpPrefix) {
			sreg := r[len(defaultRegexpPrefix):]
			expr, err := regexp.Compile(sreg)
			if err != nil {
				return nil, fmt.Errorf("compile regexp failed, err:%w", err)
			}
			set.regexp = append(set.regexp, expr)
			continue
		}
		set.full[r] = struct{}{}
	}
	return set, nil
}

type ipRuleSet struct {
	ip   map[string]struct{}
	cidr []*net.IPNet
}

func (r *ipRuleSet) isMatch(strip string) bool {
	if _, ok := r.ip[strip]; ok {
		return true
	}
	ip := net.ParseIP(strip)
	if ip == nil {
		return false
	}
	for _, cidr := range r.cidr {
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}

func makeIPRuleSet(rs []string) (*ipRuleSet, error) {
	set := &ipRuleSet{
		ip:   make(map[string]struct{}),
		cidr: nil,
	}
	for _, item := range rs {
		if strings.Contains(item, "/") {
			//TODO: 合并cidr
			_, cidr, err := net.ParseCIDR(item)
			if err != nil {
				return nil, err
			}
			set.cidr = append(set.cidr, cidr)
			continue
		}
		set.ip[item] = struct{}{}
	}
	return set, nil
}
