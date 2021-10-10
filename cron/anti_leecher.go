package cron

import (
	"context"
	"fmt"
	"qb-helper/client"
	"qb-helper/matcher"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/xxxsen/log"
	"github.com/xxxsen/qbapi"
	"github.com/xxxsen/runner"
)

const (
	antiLeecherName = "anti-leecher"
)

type MatchConfig struct {
	Mode string `json:"mode"`
	Key  string `json:"key"`
}

type AntiLeecherConfig struct {
	BanClient []*MatchConfig `mapstructure:"ban_client"`
}

func init() {
	Regist(antiLeecherName, func(args interface{}) (Cron, error) {
		cfg := &AntiLeecherConfig{}
		if err := mapstructure.Decode(args, cfg); err != nil {
			return nil, err
		}
		return NewAntiLeecher(cfg)
	})
}

type AntiLeecher struct {
	BaseCron
	matchers []matcher.Matcher
}

func NewAntiLeecher(cfg *AntiLeecherConfig) (Cron, error) {
	var lst []matcher.Matcher
	for _, item := range cfg.BanClient {
		mt, err := matcher.Build(item.Mode, item.Key, nil)
		if err != nil {
			return nil, err
		}
		lst = append(lst, mt)
	}
	return &AntiLeecher{matchers: lst}, nil
}

func (c *AntiLeecher) Start() error {
	ctx := context.Background()
	offset := 0
	limit := 1000
	//qb无法正确处理 offset, 先一次性拉回来
	run := runner.New(10)
	for {
		lst, err := c.readTorrentList(ctx, offset, limit)
		if err != nil {
			return err
		}
		if len(lst) == 0 {
			break
		}
		for _, item := range lst {
			item := item
			run.Add(fmt.Sprintf("check_%s", item.Hash), func(ctx context.Context) error {
				return c.doCheckLogic(ctx, item)
			})
		}
		offset += limit
		break
	}
	return run.Run(ctx)
}

func (c *AntiLeecher) isBanedClient(ctx context.Context, info *qbapi.TorrentPeerItem) bool {
	client := info.Client
	for _, mt := range c.matchers {
		if mt.IsMatch(client) {
			return true
		}
	}
	return false
}

func (c *AntiLeecher) isLeechClient(ctx context.Context, info *qbapi.TorrentPeerItem) bool {
	return false
}

func (c *AntiLeecher) doCheckLogic(ctx context.Context, item *qbapi.TorrentListItem) error {
	peers, err := c.readPeerData(ctx, item.Hash)
	if err != nil {
		return err
	}
	banConns := make(map[string]*qbapi.TorrentPeerItem)
	pd := peers.Peers
	for addr, info := range pd {
		//检查ban客户端
		if c.isBanedClient(ctx, info) {
			banConns[addr] = info
			continue
		}
		//检查吸血规则
		if c.isLeechClient(ctx, info) {
			banConns[addr] = info
			continue
		}
	}
	return c.banClients(ctx, banConns)
}

func (c *AntiLeecher) banClients(ctx context.Context, banMap map[string]*qbapi.TorrentPeerItem) error {
	if len(banMap) == 0 {
		return nil
	}
	var peerList []string
	for addr, item := range banMap {
		peerList = append(peerList, addr)
		log.Infof("Addr:%s hit rule, ban it, client:%s, contry:%s, code:%s, flag:%s, progress:%f", addr, item.Client, item.Country, item.CountryCode, item.Flags, item.Progress)
	}
	return c.banClient(ctx, peerList)
}

func (c *AntiLeecher) banClient(ctx context.Context, peers []string) error {
	req := &qbapi.BanPeersReq{Peers: peers}
	_, err := client.Instance().BanPeers(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *AntiLeecher) readTorrentList(ctx context.Context, offset int, limit int) ([]*qbapi.TorrentListItem, error) {
	filter := "all"
	req := &qbapi.GetTorrentListReq{Offset: &offset, Limit: &limit, Filter: &filter}
	rsp, err := client.Instance().GetTorrentList(ctx, req)
	if err != nil {
		return nil, err
	}
	return rsp.Items, nil
}

func (c *AntiLeecher) readPeerData(ctx context.Context, hash string) (*qbapi.TorrentPeerData, error) {
	req := &qbapi.GetTorrentPeerDataReq{Hash: hash, Rid: 0}
	rsp, err := client.Instance().GetTorrentPeerData(ctx, req)
	if err != nil {
		return nil, err
	}
	return rsp.Data, nil
}

func (c *AntiLeecher) Interval() time.Duration {
	return 20 * time.Second
}

func (c *AntiLeecher) StartOnce() bool {
	return false
}

func (c *AntiLeecher) Name() string {
	return antiLeecherName
}
