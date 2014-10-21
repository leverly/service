package zc

import (
	"time"
	"testing"
	"strconv"
)

var zcloudConfig *ZCloudConfig = new(ZCloudConfig)
var zcloudTestInit bool = false

func initCloud() {
	if (zcloudTestInit) { return }
	zcloudConfig.serviceMap = make(map[string][]string, 10)
	zcloudConfig.serviceMap["echo"] = append(zcloudConfig.serviceMap["echo"],
											 "localhost:12346")
	zcloudConfig.serviceMap["echo"] = append(zcloudConfig.serviceMap["echo"],
											 "localhost:12347")
	SetCloudConfig(zcloudConfig)
}

var firstServiceInit bool = false
func initFirstService() {
	if firstServiceInit { return }
	config := &ZServiceConfig{Port : "12346"}

	MyService.Init("echo", config)
	MyService.Handle("ping", ZServiceHandler(func (req *ZMsg, resp *ZMsg) {
		resp.SetName(ZC_MSG_NAME_ACK)
		resp.SetVersion(uint8(1))
		resp.Put("content", "pong1: " + req.GetString("content"))
		resp.Put("TraceId", req.Get("TraceId"))
	}))
	go MyService.Start()
	firstServiceInit = true
	time.Sleep(time.Second)
}

var secondServiceInit bool = false
func initSecondService() {
	if secondServiceInit { return }
	config := &ZServiceConfig{Port : "12347"}
	var s *ZService = NewZService()
	s.Init("echo", config)
	s.Handle("ping", ZServiceHandler(func (req *ZMsg, resp *ZMsg) {
			resp.SetName(ZC_MSG_NAME_ACK)
			resp.SetVersion(uint8(1))
			resp.Put("content", "pong2: " + req.GetString("content"))
			resp.Put("TraceId", req.Get("TraceId"))
		}))
	go s.Start()
	secondServiceInit = true
	time.Sleep(time.Second)
}

var thirdServiceInit bool = false
func initThirdService() {
	if thirdServiceInit { return }
	config := &ZServiceConfig{Port : "12348"}
	var s *ZService = NewZService()
	s.Init("echo", config)
	s.Handle("ping", ZServiceHandler(func (req *ZMsg, resp *ZMsg) {
			resp.SetName(ZC_MSG_NAME_ACK)
			resp.SetVersion(uint8(1))
			resp.Put("content", "pong3: " + req.GetString("content"))
			resp.Put("TraceId", req.Get("TraceId"))
		}))
	go s.Start()
	thirdServiceInit = true
	time.Sleep(time.Second)
}

func initCloudTest() {
	initCloud()
	initFirstService()
	initSecondService()
	initThirdService()
}

func TestTwoEchoService(t *testing.T) {
	initCloudTest()

	for i := 0; i < 5; i++ {
		req := NewMsg("ping", uint8(1))
		req.Put("content", "hello world from req" + strconv.Itoa(i))
		req.Put("TraceId", strconv.Itoa(i))
		resp, err := Service("echo").Send(req)
		if err != nil {
			t.Errorf("send req%d failed: %s", i, err.Error())
		}

		if resp.GetName() != ZC_MSG_NAME_ACK {
			t.Errorf("get error resp%d: %s", i, resp.GetErr())
		}
	}
}
