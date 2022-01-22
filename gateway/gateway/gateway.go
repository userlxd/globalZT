package gateway

import (
	"context"
	"globalZT/gateway/proxy"
	"globalZT/tools/config"
	"globalZT/tools/log"
)

var gateway Gateway

type Gateway struct {
	cancle context.CancelFunc
	UUID   uint32
	proxy  proxy.Proxy
}

func init() {
	var uuid uint32 = 123

	// load config
	var c config.Gateway
	c.Load()

	// log tools
	log.Init(c.LOG)

	// gateway
	gateway = Gateway{
		UUID:  uuid,
		proxy: *proxy.NewProxy(uuid),
	}
}

func Run() {

	log.Log.Info("start")
	defer log.Log.Info("quit")

	ctx, cancle := context.WithCancel(context.Background())
	gateway.cancle = cancle

	go gateway.proxy.Run(ctx)

	<-ctx.Done()
}
