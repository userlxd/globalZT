package proxy

import (
	"context"
	"globalZT/pkg/network"
	"globalZT/pkg/tunnel"
	"globalZT/pkg/tuntap"
	"globalZT/tools/config"
	"net"

	"go.uber.org/atomic"
)

type Proxy struct {
	cancle context.CancelFunc // proxy cancle
	UUID   uint32             // device UUID
	Active *atomic.Bool       // connect status

	// grpc tunnel
	o2Gw *tunnel.Office2Gw // grpc tunnel

	// Tun
	tun         tuntap.Tun  // tun tap interface
	tunReadBuf  chan []byte // tun read buffer
	tunWriteBuf chan []byte // tun write buffer
}

func NewProxy(c config.Office, uuid uint32) *Proxy {
	var o = &Proxy{
		UUID:        uuid,
		tunReadBuf:  make(chan []byte, 1024),
		tunWriteBuf: make(chan []byte, 1024),
		o2Gw:        tunnel.NewOffice2Gw(uuid, c.GW),
	}

	return o
}

func (p *Proxy) Run(pctx context.Context) {

	ctx, cancle := context.WithCancel(pctx)
	p.cancle = cancle

	if intf, err := tuntap.Open(
		net.IPv4(4, 4, 4, 4),
		net.IPv4(0, 0, 0, 0),
		net.IPv4(255, 255, 255, 0),
	); err != nil {
		return
	} else {
		p.tun = intf
	}

	go p.o2Gw.Run(ctx)
	// Read from tun
	go p.tun.Read(p.tunReadBuf)
	go p.tun.Write(p.tunWriteBuf)
	go p.exchangeIn()
	go p.exchangeOut()

	<-ctx.Done()
}

func (p *Proxy) exchangeIn() {
	for {
		select {
		case in := <-p.tunReadBuf:
			p.o2Gw.DataIn <- network.Parse(in)
		}
	}
}

func (p *Proxy) exchangeOut() {
	for {
		select {
		case out := <-p.o2Gw.DataOut:
			p.tunWriteBuf <- out.Data
		}
	}
}
