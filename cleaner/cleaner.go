package cleaner

import (
	"context"
	"fmt"
	"time"

	"github.com/xxxsen/common/logutil"
	"github.com/xxxsen/common/trace"
	"github.com/xxxsen/qbapi"
	"github.com/xxxsen/runner"
	"go.uber.org/zap"
)

const (
	defaultMaxLimit = 1000
)

type Cleaner struct {
	api        *qbapi.QBAPI
	c          *config
	uaRule     *strRuleSet
	regionRule *strRuleSet
	ipRule     *ipRuleSet
	peerIDRule *strRuleSet
}

func New(opts ...Option) (*Cleaner, error) {
	c := &config{}
	for _, opt := range opts {
		opt(c)
	}
	if c.interval == 0 {
		c.interval = 10 * time.Second
	}
	api, err := qbapi.NewAPI(qbapi.WithHost(c.host), qbapi.WithAuth(c.user, c.pwd))
	if err != nil {
		return nil, fmt.Errorf("new api fail, err:%w", err)
	}
	uaRule, err := makeStrRuleSet(c.uaRs)
	if err != nil {
		return nil, fmt.Errorf("make ua rule set failed, err:%w", err)
	}
	regionRule, err := makeStrRuleSet(c.regionRs)
	if err != nil {
		return nil, fmt.Errorf("make region rule set failed, err:%w", err)
	}
	ipRule, err := makeIPRuleSet(c.ipRs)
	if err != nil {
		return nil, fmt.Errorf("make ip rule set failed, err:%w", err)
	}
	peerIDRule, err := makeStrRuleSet(c.peerIDRs)
	if err != nil {
		return nil, fmt.Errorf("make peer rule set failed, err:%w", err)
	}
	return &Cleaner{api: api, c: c, uaRule: uaRule, regionRule: regionRule, ipRule: ipRule, peerIDRule: peerIDRule}, nil
}

func (c *Cleaner) Start() error {
	go c.start()
	return nil
}

func (c *Cleaner) start() {
	var idx uint32 = 0
	for {
		ctx := context.Background()
		idx++
		ctx = trace.WithTraceId(ctx, fmt.Sprintf("%d", idx))
		if err := c.do(ctx); err != nil {
			logutil.GetLogger(ctx).Error("do clean failed", zap.Error(err))
		}
		time.Sleep(c.c.interval)
	}
}

func (c *Cleaner) do(ctx context.Context) error {
	if err := c.ensureLogin(ctx); err != nil {
		return fmt.Errorf("login failed, err:%w", err)
	}
	offset := 0
	limit := defaultMaxLimit
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

func (c *Cleaner) ensureLogin(ctx context.Context) error {
	if err := c.api.Login(ctx); err != nil {
		return err
	}
	return nil
}

func (c *Cleaner) isInBlackList(info *qbapi.TorrentPeerItem) bool {
	if c.uaRule.isMatch(info.Client) {
		return true
	}
	if c.regionRule.isMatch(info.CountryCode) {
		return true
	}
	if c.ipRule.isMatch(info.Ip) {
		return true
	}
	if c.peerIDRule.isMatch(info.PeerIdClient) {
		return true
	}
	return false
}

func (c *Cleaner) doCheckLogic(ctx context.Context, item *qbapi.TorrentListItem) error {
	peers, err := c.readPeerData(ctx, item.Hash)
	if err != nil {
		return err
	}
	banConns := make(map[string]*qbapi.TorrentPeerItem)
	pd := peers.Peers
	for addr, info := range pd {
		//检查ban客户端
		if c.isInBlackList(info) {
			banConns[addr] = info
			continue
		}
	}
	return c.banClients(ctx, banConns)
}

func (c *Cleaner) banClients(ctx context.Context, banMap map[string]*qbapi.TorrentPeerItem) error {
	if len(banMap) == 0 {
		return nil
	}
	var peerList []string
	for addr, item := range banMap {
		peerList = append(peerList, addr)
		logutil.GetLogger(ctx).Info("hit rule, ban it", zap.String("addr", addr),
			zap.String("client", item.Client), zap.String("country", item.Country),
			zap.String("code", item.CountryCode), zap.String("flags", item.Flags),
			zap.Float64("progress", item.Progress), zap.String("peer_id", item.PeerIdClient))
	}
	return c.banClient(ctx, peerList)
}

func (c *Cleaner) banClient(ctx context.Context, peers []string) error {
	req := &qbapi.BanPeersReq{Peers: peers}
	_, err := c.api.BanPeers(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cleaner) readTorrentList(ctx context.Context, offset int, limit int) ([]*qbapi.TorrentListItem, error) {
	filter := "all"
	req := &qbapi.GetTorrentListReq{Offset: &offset, Limit: &limit, Filter: &filter}
	rsp, err := c.api.GetTorrentList(ctx, req)
	if err != nil {
		return nil, err
	}
	return rsp.Items, nil
}

func (c *Cleaner) readPeerData(ctx context.Context, hash string) (*qbapi.TorrentPeerData, error) {
	req := &qbapi.GetTorrentPeerDataReq{Hash: hash, Rid: 0}
	rsp, err := c.api.GetTorrentPeerData(ctx, req)
	if err != nil {
		return nil, err
	}
	return rsp.Data, nil
}
