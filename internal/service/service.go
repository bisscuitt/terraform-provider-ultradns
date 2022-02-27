package service

import (
	"github.com/ultradns/ultradns-go-sdk/pkg/client"
	"github.com/ultradns/ultradns-go-sdk/pkg/probe"
	"github.com/ultradns/ultradns-go-sdk/pkg/record"
	"github.com/ultradns/ultradns-go-sdk/pkg/zone"
)

type Service struct {
	ZoneService   *zone.Service
	RecordService *record.Service
	ProbeService  *probe.Service
}

func NewService(c *client.Client) (*Service, error) {
	service := &Service{}

	var err error

	if service.ZoneService, err = zone.Get(c); err != nil {
		return nil, err
	}

	if service.RecordService, err = record.Get(c); err != nil {
		return nil, err
	}

	if service.ProbeService, err = probe.Get(c); err != nil {
		return nil, err
	}

	return service, nil
}
