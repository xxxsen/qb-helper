package main

import (
	"flag"
	"qb-helper/client"
	"qb-helper/config"
	"qb-helper/cron"

	"github.com/xxxsen/log"
)

var cfg = flag.String("cfg", "./config.json", "config file")

func main() {
	flag.Parse()
	conf, err := config.Parse(*cfg)
	if err != nil {
		log.Fatalf("Parse config fail, err:%v", err)
	}
	//init log
	log.Init(conf.Log.File, log.StringToLevel(conf.Log.Level), conf.Log.Rotate, conf.Log.Size, conf.Log.KeepDay, conf.Log.Console)

	//init qbapi
	if err := client.Init(conf.Auth.Username, conf.Auth.Password, conf.Auth.Host); err != nil {
		log.Fatalf("Init qbclient failed, user:%s, pwd:%s, host:%s, err:%v", conf.Auth.Username, conf.Auth.Password, conf.Auth.Host, err)
	}

	//init cron
	cg := cron.NewBoostrap()
	for _, cr := range conf.CronList {
		if !cr.Enable {
			log.Infof("Cron:%s, disabled, skip", cr.Name)
			continue
		}
		log.Debugf("Init cron:%s with args:%v", cr.Name, cr.Args)
		creator, err := cron.Get(cr.Name)
		if err != nil {
			log.Fatalf("Get cron:%s failed, err:%v", cr.Name, err)
		}
		cri, err := creator(cr.Args)
		if err != nil {
			log.Fatalf("Init cron:%s failed, args:%v, err:%v")
		}
		log.Infof("Cron:%s init succ", cr.Name)
		cg.Add(cri)
	}
	log.Infof("All cron start finished")
	cg.Start()
	select {}
}
