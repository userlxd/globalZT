package office

import (
	"context"
	"globalZT/office/proxy"
)

var office Office

// 1.web todo
// 2.proxy now
type Office struct {
	UUID  uint32
	proxy *proxy.Proxy
}

func init() {
	var uuid uint32 = 123
	office = Office{
		UUID:  uuid,
		proxy: proxy.NewProxy(uuid),
	}
}

func Run() {

	ctx, cancle := context.WithCancel(context.Background())

	office.proxy.Cancle = cancle
	office.proxy.Run()

	<-ctx.Done()
}
