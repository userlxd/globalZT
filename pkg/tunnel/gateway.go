package tunnel

import (
	"globalZT/tools/log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

var GW Gateway

type Gateway struct {
	listener net.Listener
	server   *grpc.Server
}

func Run() {
	var err error
	var GW = &Gateway{}

	addr := "0.0.0.0:31580"
	GW.listener, err = net.Listen("tcp", addr)
	if err != nil {
		log.Log.Errorw("[New Gateway Litener]", "msg", err, "obj", addr)
		return
	}

	GW.server = grpc.NewServer()
	RegisterOffice2GwServer(GW.server, GW)
	GW.server.Serve(GW.listener)
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

		log.Log.Infof("uuid:%d", req.GetUUID())
	}
}

func (g *Gateway) Control(server Office2Gw_ControlServer) error {
	return nil
}
