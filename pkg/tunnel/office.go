package tunnel

import (
	"context"
	"globalZT/pkg/auth"
	"globalZT/tools/config"
	"globalZT/tools/log"
	"io"
	"time"

	"go.uber.org/atomic"
	"google.golang.org/grpc"
)

type Office2Gw struct {
	Host   string           // server address
	UUID   uint32           // device UUID
	User   string           // user name
	conn   *grpc.ClientConn // grpc conn
	client Office2GwClient  // grpc client
	jwt    *auth.JwtInfo    // jwth auth info

	// Status Signal
	loop   *atomic.Bool       // control recv
	daemon chan struct{}      // exit signal
	Logout chan struct{}      // logout signal
	Cancle context.CancelFunc // cancle when exit

	// Data
	dataClient Office2Gw_DataClient // grpc stream
	DataOut    chan *Out            // in queue
	DataIn     chan *In             // out queue
}

func NewOffice2Gw(UUID uint32, c config.GW) *Office2Gw {

	return &Office2Gw{
		Host:    c.IP + ":" + c.PORT,
		UUID:    UUID,
		loop:    atomic.NewBool(false),
		daemon:  make(chan struct{}, 1),
		Logout:  make(chan struct{}, 1),
		DataIn:  make(chan *In, 100),
		DataOut: make(chan *Out, 100),
	}
}

// close by command
func (o *Office2Gw) Close() {
	// break loop
	o.loop.Store(false)
	// close connection
	o.conn.Close()
	// cancle run
	o.Cancle()
}

// close by loop
// awake daemon to reconnect
func (o *Office2Gw) close() {
	// break loop
	o.loop.Store(false)
	// send exit signal
	o.daemon <- struct{}{}
}

func (o *Office2Gw) connect(ctx context.Context) {
	if conn, err := grpc.Dial(o.Host, grpc.WithInsecure()); err == nil {
		o.conn = conn
	} else {
		log.Log.Error("[NewOfficeTunnel]", "msg", err, "obj", o.Host)
		return
	}

	o.client = NewOffice2GwClient(o.conn)

	// set data client
	if stream, err := o.client.Data(ctx); err == nil {
		o.dataClient = stream
	} else {
		log.Log.Error("[NewDataClient]", "msg", err)
		return
	}
}

func (o *Office2Gw) Run(pctx context.Context) {

	ctx, cacnle := context.WithCancel(pctx)
	o.Cancle = cacnle

	o.daemon <- struct{}{}
	go func() {
		for {
			select {
			case <-o.daemon:
				o.Daemon(ctx)
			}
		}
	}()
	go o.loopSend()
	go o.loopRecv()

	<-ctx.Done()
}

func (o *Office2Gw) Daemon(ctx context.Context) {
	o.connect(ctx)
	o.loop.Store(true)
}

func (o *Office2Gw) loopSend() {
	for {
		if o.loop.Load() {
			select {
			case in := <-o.DataIn:
				if err := o.dataClient.Send(in); err != nil {
					log.Log.Error("[DataClientSend]", "msg", err)
					o.close()
				}
			}
		} else {
			time.Sleep(time.Second)
		}
	}
}

func (o *Office2Gw) loopRecv() {
	for {
		if o.loop.Load() {
			out, err := o.dataClient.Recv()
			if err == io.EOF {
				o.close()
			}
			if err != nil {
				log.Log.Error("[DataClientRecv]", "msg", err)
				o.close()
			}
			o.DataOut <- out
		} else {
			time.Sleep(time.Second)
		}
	}
}
