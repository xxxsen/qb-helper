package cron

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/xxxsen/log"
)

type Cron interface {
	Start() error
	Interval() time.Duration
	StartOnce() bool
	Name() string
}

type Creator func(args interface{}) (Cron, error)

var mp = make(map[string]Creator)

func Regist(name string, cr Creator) {
	mp[name] = cr
}

func Get(name string) (Creator, error) {
	if v, ok := mp[name]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("not found")
}

type Boostrap struct {
	cr []Cron
}

func (c *Boostrap) Add(cron Cron) {
	c.cr = append(c.cr, cron)
}

func (c *Boostrap) Start() {
	for _, cr := range c.cr {
		go c.cronjob(cr)
	}
}

func (c *Boostrap) cronjob(cr Cron) {
	for {
		if err := c.startOneJob(cr); err != nil {
			log.Errorf("Start cron but failed, name:%s, err:%v", cr.Name(), err)
		}
		if cr.StartOnce() {
			break
		}
		if cr.Interval() != 0 {
			time.Sleep(cr.Interval())
		}
		log.Debugf("Cron:%s exec finished", cr.Name())
	}
}

func (c *Boostrap) startOneJob(cr Cron) error {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("start cron job panic, name:%s, err:%v, stack:%s", cr.Name(), err, string(debug.Stack()))
		}
	}()
	return cr.Start()
}

func NewBoostrap(crs ...Cron) *Boostrap {
	bst := &Boostrap{}
	for _, cr := range crs {
		bst.Add(cr)
	}
	return bst
}
