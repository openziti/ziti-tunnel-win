package idutil

import (
	"github.com/michaelquigley/pfxlog"
	idcfg "github.com/netfoundry/ziti-sdk-golang/ziti/config"
	"github.com/netfoundry/ziti-tunnel-win/service/cziti"
	"github.com/netfoundry/ziti-tunnel-win/service/ziti-tunnel/dto"
)

var log = pfxlog.Logger()

//Removes the Config from the provided identity and returns a 'cleaned' id
func Clean(id dto.Identity) dto.Identity {
	log.Tracef("cleaning identity: %s", id.Name)
	nid := id
	nid.Config = idcfg.Config{}
	nid.Config.ZtAPI = id.Config.ZtAPI
	nid.Services = make([]*dto.Service, len(id.Services))
	for i, svc := range id.Services{
		nid.Services[i] = svc
	}
	AddMetrics(&nid)
	log.Tracef("Up: %v Down %v", nid.Metrics.Up, nid.Metrics.Down)
	return nid
}

func AddMetrics(id *dto.Identity) {
	up, down, _ := cziti.GetTransferRates(id.NFContext)

	id.Metrics = &dto.Metrics{
		Up:   up,
		Down: down,
	}
}