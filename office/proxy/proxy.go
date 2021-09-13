package proxy

import (
	"context"
	"globalZT/pkg/network"
	"globalZT/pkg/tunnel"
	"globalZT/tools/log"
	"io"
	"net"
	"time"
)

type Proxy struct {
	UUID        string
	offcie2Gw   *tunnel.Office
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

	p.tun, err = Open(net.IPv4(4, 4, 4, 4), net.IPv4(0, 0, 0, 0), net.IPv4(255, 255, 255, 0))
	if err != nil {
		cancle()
		return
	}

	p.readBuffer = make(chan []byte, 1024)
	p.writeBuffer = make(chan []byte, 1024)
	apps := make(chan tunnel.OfficeApp, 1024)

	// Read from tun
	go func() {
		if err := p.tun.Read(p.readBuffer); err != nil {
			log.Log.Errorw("[Read Tun]", "msg", err, "obj", "")
			cancle()
			return
		}
	}()

	// Write to tunnel
	go func() {
		for {
			select {
			case data := <-p.readBuffer:
				code, req := network.Parse(data)
				log.Log.Info(code)
				app, new := p.offcie2Gw.GetApp(code)
				if new {
					apps <- *app
				}
				app.ReqChan <- req
			}
		}
	}()

	for {
		select {
		case app := <-apps:
			// Read from tunnel
			go func(a tunnel.OfficeApp) {
				for {
					resp, err := a.Stream.Recv()
					if err != nil {
						if err == io.EOF {
							log.Log.Infof("%s exit by server quit", app.Code)
						} else {
							log.Log.Errorw("[Grpc Recv]", "msg", err, "obj", app.Code)
						}
						a.Stop()
						break
					}
					p.writeBuffer <- resp.Data
				}
			}(app)

			// Write to tunnel
			go func(a tunnel.OfficeApp) {
				defer a.Stream.CloseSend()
				for {
					select {
					case <-a.Done():
						break
					case <-a.Keepalive.C:
						log.Log.Infof("%s exit by client quit", app.Code)
						a.Stop()
					case req := <-a.ReqChan:
						a.Stream.Send(req)
						a.Keepalive.Reset(time.Second * 5)
					}
				}
			}(app)
		}
	}
}

// tun read -> readBuffer
// readBuffer -> oa send
// oa receive -> writeBuffer
// writeBuffer -> tun write
