package cron

import (
	"fmt"
	"time"
)

type BaseCron struct {
}

func (c *BaseCron) Start() error {
	return fmt.Errorf("impl it plz")
}

func (c *BaseCron) Interval() time.Duration {
	return 0
}

func (c *BaseCron) StartOnce() bool {
	return false
}

func (c *BaseCron) Name() string {
	return "base"
}
