package tunnel

import (
	"context"
	"errors"
	"globalZT/tools/log"
	"sync/atomic"

	"google.golang.org/grpc"
)

var priorityMap = map[uint]uint{
	0: 1,
	1: 10,
	2: 100,
	3: 1000,
}

type OfficeApp struct {
	App         string
	Concurrency uint
	stream      Office2Gw_DataClient
	ReqChan     chan *OfficeReq
	RespChan    chan *OfficeResp
	Office      *Office
}

type Office struct {
	UUID   string
	active uint32
	conn   *grpc.ClientConn
	client Office2GwClient
}

func CreateOfficeTunnel(UUID string) (Office, error) {

	var err error
	var office = Office{UUID: UUID}
	var host = "gw.globalzt.com:31580" // todo 使用resolve解决gw loadblance

	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Log.Error("[New Office App Grpc Conn Error]", "msg", err, "obj", host)
		return office, err
	}
	defer conn.Close()

	office.client = NewOffice2GwClient(conn)

	atomic.CompareAndSwapUint32(&office.active, 0, 1)
	return office, nil
}

func (o *Office) CreateOfficeApp(ctx context.Context) (*OfficeApp, error) {

	if atomic.LoadUint32(&o.active) == 0 {
		return nil, errors.New("out of service")
	}

	var err error
	var oa = &OfficeApp{
		Office: o,
	}

	oa.stream, err = o.client.Data(ctx)
	if err != nil {
		log.Log.Errorw("[New Office App Grpc Stream Error]", "msg", err, "obj", o.conn.Target())
		return oa, err
	}

	return oa, nil
}

func (oa *OfficeApp) SetApp(app string, priority uint) {
	oa.App = app
	oa.ReqChan = make(chan *OfficeReq, priorityMap[priority])
	oa.RespChan = make(chan *OfficeResp, priorityMap[priority])
}

func (oa *OfficeApp) DataRun(cancle context.CancelFunc) {

	go func() {
		for {
			resp, err := oa.stream.Recv()
			if err != nil {
				log.Log.Errorw("[Recv Office App Resp Error]", "msg", err, "obj", oa.App)
				cancle()
				return
			}
			oa.RespChan <- resp
		}
	}()

	for {
		select {
		case req := <-oa.ReqChan:
			req.UUID = oa.Office.UUID
			if err := oa.stream.Send(req); err != nil {
				log.Log.Errorw("[Send Office App Req Error]", "msg", err, "obj", oa.App)
				cancle()
				return
			}
		}
	}
}
