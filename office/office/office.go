package office

import (
	"context"
	"globalZT/office/proxy"
	"globalZT/tools/config"
	"globalZT/tools/log"
)

var office Office

// 1.web todo
// 2.proxy now
type Office struct {
	cancle context.CancelFunc
	UUID   uint32
	proxy  *proxy.Proxy
}

func init() {
	var uuid uint32 = 123

	// load config
	var c config.Office
	c.Load()

	// log tools
	log.Init(c.LOG)

	// office
	office = Office{
		UUID:  uuid,
		proxy: proxy.NewProxy(c, uuid),
	}
}

func Run() {

	log.Log.Info("start")
	defer log.Log.Info("quit")

	ctx, cancle := context.WithCancel(context.Background())
	office.cancle = cancle

	office.proxy.Run(ctx)

	<-ctx.Done()
}
