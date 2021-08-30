package main

import (
	"globalZT/pkg/tunnel"
	"net"

	"google.golang.org/grpc"
)

type Gateway struct {
	UUID string
}

func Run() error {
	listener, err := net.Listen("tcp", "0.0.0.1:31580")
	if err != nil {
		return err
	}

	gw := Gateway{}
	server := grpc.NewServer()
	tunnel.RegisterOffice2GwServer(server, &gw)

	return server.Serve(listener)
}

func (g *Gateway) Data(server tunnel.Office2Gw_DataServer) error {

	return nil
}

func (g *Gateway) Control(server tunnel.Office2Gw_ControlServer) error {
	return nil
}
