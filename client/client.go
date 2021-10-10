package client

import (
	"context"

	"github.com/xxxsen/qbapi"
)

var client *qbapi.QBAPI

func Init(user string, pwd string, host string) error {
	cli, err := qbapi.NewAPI(qbapi.WithAuth(user, pwd), qbapi.WithHost(host))
	if err != nil {
		return err
	}
	err = cli.Login(context.Background())
	if err != nil {
		return err
	}
	client = cli
	return nil
}

func Instance() *qbapi.QBAPI {
	return client
}
