package office

import (
	"context"
	"globalZT/office/proxy"
)

var office Office

type Office struct {
	UUID  string
	proxy proxy.Proxy
}

func init() {
	uuid := "123"
	office = Office{
		UUID: uuid,
		proxy: proxy.Proxy{
			UUID: uuid,
		},
	}
}

func Run() {

	ctx, cancle := context.WithCancel(context.Background())
	office.proxy.Run(cancle)
	<-ctx.Done()
}
