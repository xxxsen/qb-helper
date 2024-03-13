package main

import (
	"flag"
	"github.com/xxxsen/common/logger"
	"go.uber.org/zap"
	"log"
	"qb-helper/cleaner"
	"qb-helper/config"
	"time"
)

var cfg = flag.String("config", "./config.json", "config file")

func main() {
	flag.Parse()
	conf, err := config.Parse(*cfg)
	if err != nil {
		log.Fatalf("Parse config fail, err:%v", err)
	}
	//init log
	logkit := logger.Init(conf.LogConfig.File, conf.LogConfig.Level, int(conf.LogConfig.FileCount), int(conf.LogConfig.FileSize), int(conf.LogConfig.KeepDays), conf.LogConfig.Console)
	svc, err := cleaner.New(
		cleaner.WithQBConfig(conf.QBConfig.Host, conf.QBConfig.Username, conf.QBConfig.Password),
		cleaner.WithInterval(time.Duration(conf.Interval)*time.Second),
		cleaner.WithAddRules(conf.UaList...),
	)
	if err != nil {
		logkit.Fatal("init cleaner failed", zap.Error(err))
	}
	if err := svc.Start(); err != nil {
		logkit.Fatal("run cleaner failed", zap.Error(err))
	}
	select {}
}
