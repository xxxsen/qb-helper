package cleaner

import (
	"time"
)

type config struct {
	uaRs     []string
	ipRs     []string
	regionRs []string
	peerIDRs []string

	//
	interval time.Duration
	//
	host string
	user string
	pwd  string
}

type Option func(c *config)

func WithQBConfig(host, u, p string) Option {
	return func(c *config) {
		c.host, c.user, c.pwd = host, u, p
	}
}

func WithAddUaRule(rs ...string) Option {
	return func(c *config) {
		c.uaRs = append(c.uaRs, rs...)
	}
}

func WithAddIPRule(rs ...string) Option {
	return func(c *config) {
		c.ipRs = append(c.ipRs, rs...)
	}
}

func WithAddPeerIDRule(rs ...string) Option {
	return func(c *config) {
		c.peerIDRs = append(c.peerIDRs, rs...)
	}
}

func WithAddRegionRule(rs ...string) Option {
	return func(c *config) {
		c.regionRs = append(c.regionRs, rs...)
	}
}

func WithInterval(it time.Duration) Option {
	return func(c *config) {
		c.interval = it
	}
}
