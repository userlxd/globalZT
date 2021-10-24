package proxy

import (
	"context"
	"globalZT/pkg/network"
	"globalZT/pkg/tunnel"
	"globalZT/pkg/tuntap"
	"globalZT/tools/log"
	"io"
	"net"
	"time"
)

type Proxy struct {
	UUID        uint32
	Cancle      context.CancelFunc
	o2Gw        *tunnel.Office
	tun         tuntap.Tun
	readBuffer  chan []byte
	writeBuffer chan []byte
	apps        chan *tunnel.App
}

func NewProxy(uuid uint32) *Proxy {
	return &Proxy{
		UUID:        uuid,
		readBuffer:  make(chan []byte, 1024),
		writeBuffer: make(chan []byte, 1024),
		apps:        make(chan *tunnel.App, 1024),
	}
}

func (p *Proxy) Run() {
	var err error

	p.o2Gw, err = tunnel.NewOfficeTunnel(p.UUID)
	if err != nil {
		p.Cancle()
		return
	}
	defer p.o2Gw.CloseConn()

	p.tun, err = tuntap.Open(net.IPv4(4, 4, 4, 4), net.IPv4(0, 0, 0, 0), net.IPv4(255, 255, 255, 0))
	if err != nil {
		p.Cancle()
		return
	}

	// Read from tun
	go p.TunRead()

	// Locate to app
	go p.AppLocate()

	for {
		select {
		case app := <-p.apps:
			// Read from tunnel
			go p.TunnelRead(app)
			// Write to tunnel
			go p.TunnelWrite(app)
		}
	}
}

// tun read -> readBuffer
func (p *Proxy) TunRead() {
	if err := p.tun.Read(p.readBuffer); err != nil {
		log.Log.Errorw("[Read Tun]", "msg", err, "obj", "")
		p.Cancle()
		return
	}
}

// writeBuffer -> tun write
func (p *Proxy) TunWrite() {

}

// oa receive -> writeBuffer
func (p *Proxy) TunnelRead(app *tunnel.App) {
	for {
		resp, err := app.Stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Log.Info(app.Code, "[EOF]")
			} else {
				log.Log.Errorw("[Grpc Recv]", "msg", err, "obj", app.Code)
			}
			app.Stop()
			break
		}
		p.writeBuffer <- resp.Data
	}
}

// readBuffer -> oa send
func (p *Proxy) TunnelWrite(app *tunnel.App) {
	for {
		select {
		case <-app.Keepalive.C:
			log.Log.Info(app.Code, "[timeout]")
			app.Stop()
		case req := <-app.ReqChan:
			log.Log.Info(app.Code, "[write chan]")
			req.UUID = p.UUID
			app.Stream.Send(req)
			app.Keepalive.Reset(time.Second * 500)
		}
	}
}

func (p *Proxy) AppLocate() {
	for {
		select {
		case data := <-p.readBuffer:
			code, req := network.Parse(data)
			app, new := p.o2Gw.GetApp(code)
			if new {
				log.Log.Info(code, "[new app]")
				p.apps <- app
			}
			app.ReqChan <- req
		}
	}
}
