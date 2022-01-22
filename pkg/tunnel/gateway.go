package tunnel

import (
	"globalZT/tools/log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type Gateway struct {
	listener net.Listener // net lisnter
	server   *grpc.Server // grpc server
}

func NewGwTunnel() *Gateway {
	return &Gateway{}
}

func (g *Gateway) Run() {
	var err error
	addr := "0.0.0.0:31580"
	g.listener, err = net.Listen("tcp", addr)
	if err != nil {
		log.Log.Errorw("[New Gateway Litener]", "msg", err, "obj", addr)
		return
	}

	g.server = grpc.NewServer()
	RegisterOffice2GwServer(g.server, g)
	g.server.Serve(g.listener)
}

func (g *Gateway) Data(server Office2Gw_DataServer) error {
	for {
		req, err := server.Recv()
		var remoteIP string
		if pr, ok := peer.FromContext(server.Context()); ok {
			if tcpAddr, ok := pr.Addr.(*net.TCPAddr); ok {
				remoteIP = tcpAddr.IP.String()
			} else {
				remoteIP = pr.Addr.String()
			}
		}

		if err != nil {
			log.Log.Errorw("[Gateway Recv]", "msg", err, "obj", remoteIP)
			return err
		}

		log.Log.Infof("uuid:%d data:[%s]", req.GetUUID(), string(req.Data))
	}
}

func (g *Gateway) Control(server Office2Gw_ControlServer) error {
	return nil
}
