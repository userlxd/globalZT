package proxy

import (
	"context"
	"globalZT/pkg/tunnel"
	"globalZT/pkg/tuntap"
)

type Proxy struct {
	cancle context.CancelFunc

	UUID uint32
	gw2O *tunnel.Gateway
	tun  tuntap.Tun
}

func NewProxy(uuid uint32) *Proxy {
	return &Proxy{
		UUID: uuid,
	}
}

func (p *Proxy) Run(pctx context.Context) {
	ctx, cancle := context.WithCancel(pctx)
	p.cancle = cancle

	p.gw2O = tunnel.NewGwTunnel()
	go p.gw2O.Run()

	<-ctx.Done()
}
