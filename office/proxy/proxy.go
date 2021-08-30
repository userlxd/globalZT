package proxy

import (
	"context"
	"globalZT/pkg/tunnel"
	"net"
)

type Proxy struct {
	UUID        string
	offcie2Gw   tunnel.Office
	tun         Tun
	readBuffer  chan []byte
	writeBuffer chan []byte
}

func (p *Proxy) Run(cancle context.CancelFunc) {
	var err error
	p.offcie2Gw, err = tunnel.CreateOfficeTunnel(p.UUID)
	if err != nil {
		cancle()
		return
	}

	p.tun, err = Open(net.IPv4(4, 4, 4, 4), net.IPv4(4, 4, 4, 0), net.IPv4(255, 255, 255, 0))
	if err != nil {
		cancle()
		return
	}

	p.readBuffer = make(chan []byte, 1024)
	p.writeBuffer = make(chan []byte, 1024)

	go func() {
		if err := p.tun.Read(p.readBuffer); err != nil {
			cancle()
			return
		}
	}()

	go func() {
		for {
			select {
			case data := <-p.readBuffer:
				println(string(data))
			}
		}
	}()
}
