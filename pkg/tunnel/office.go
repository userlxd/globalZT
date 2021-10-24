package tunnel

import (
	"context"
	"errors"
	"globalZT/tools/log"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
)

var priorityMap = map[uint]uint{
	0: 1,
	1: 10,
	2: 100,
	3: 1000,
}

type App struct {
	context.Context
	cancle      context.CancelFunc   //
	Keepalive   *time.Timer          //
	Code        string               // Code
	Concurrency uint                 // QOS
	Stream      Office2Gw_DataClient // grpc stream
	ReqChan     chan *OfficeReq      // req chan
	Office      *Office              // office
}

type Office struct {
	sync.RWMutex // lock
	host         string
	UUID         uint32           // device UUID
	active       uint32           // active code
	conn         *grpc.ClientConn // grpc conn
	client       Office2GwClient  // grpc client
	Apps         map[string]*App  // app conn pool
}

func NewOfficeTunnel(UUID uint32) (*Office, error) {

	var err error
	var office = &Office{
		UUID: UUID,
		host: "gw.globalzt.com:31580", // todo 使用resolve解决gw loadblance
		Apps: map[string]*App{},
	}

	conn, err := grpc.Dial(office.host, grpc.WithInsecure())
	if err != nil {
		log.Log.Error("[New Office App Grpc Conn Error]", "msg", err, "obj", office.host)
		return office, err
	}

	office.conn = conn
	office.client = NewOffice2GwClient(conn)

	atomic.CompareAndSwapUint32(&office.active, 0, 1)
	return office, nil
}

func (o *Office) CloseConn() {
	o.conn.Close()
}

func (o *Office) GetApp(code string) (*App, bool) {
	var new = false

	o.RLock()
	app, ok := o.Apps[code]
	o.RUnlock()
	if !ok {
		var err error
		app, err = o.initApp(code)
		if err != nil {
			return app, new
		}
		new = true
		o.Lock()
		o.Apps[code] = app
		o.Unlock()
	}
	return app, new
}

func (o *Office) initApp(code string) (*App, error) {

	if atomic.LoadUint32(&o.active) == 0 {
		return nil, errors.New("out of service")
	}

	var err error
	var oa = &App{
		Code:      code,
		Office:    o,
		Keepalive: time.NewTimer(time.Second * 500),
	}

	oa.Stream, err = o.client.Data(context.Background())
	if err != nil {
		log.Log.Errorw("[New Office App Grpc Stream Error]", "msg", err, "obj", o.host)
		return oa, err
	}

	oa.ReqChan = make(chan *OfficeReq, 1)

	return oa, nil
}

func (oa *App) Stop() {

	oa.Office.Lock()
	delete(oa.Office.Apps, oa.Code)
	oa.Office.Unlock()
	oa.cancle()
}
