package cleaner

import (
	"time"
)

type config struct {
	rs []string

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

func WithAddRules(rs ...string) Option {
	return func(c *config) {
		c.rs = append(c.rs, rs...)
	}
}

func WithInterval(it time.Duration) Option {
	return func(c *config) {
		c.interval = it
	}
}
