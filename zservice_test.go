package zc

import (
	"time"
	"testing"
	"encoding/json"
)

var config *ZServiceConfig = &ZServiceConfig{Port : "12345"}
var s *ZService = NewZService()
var inited bool = false

func initService() {
	if (inited) { return }
	s.Init("echo", config)
	s.Handle("ping", ZServiceHandler(func (req *ZMsg, resp *ZMsg) {
			resp.SetName(ZC_MSG_NAME_ACK)
			resp.SetVersion(uint8(1))
			resp.Put("content", "pong: " + req.GetString("content"))
			resp.Put("TraceId", req.Get("TraceId"))
		}))
	s.Handle("pingjson", ZServiceHandler(func (req *ZMsg, resp *ZMsg) {
			resp.SetName(ZC_MSG_NAME_ACK)
			resp.SetVersion(uint8(1))
			resp.Put("TraceId", req.Get("TraceId"))
			body := make(map[string][]string, 10)
			body["content"] = append(body["content"], "pongjson: " + req.GetString("content"))
			payload, _ := json.Marshal(body)
			resp.SetPayload(payload, ZC_MSG_PAYLOAD_JSON)
		}))
	go s.Start()
	inited = true
	time.Sleep(time.Second)
}

func TestEchoService(t *testing.T) {
	initService()

	c := NewZServiceClient("localhost:12345", "echo")
	req := NewZMsg()
	req.SetName("ping")
	req.Put("content", "hello world!")

	resp, err := c.Send(req)
	if err != nil {
		t.Error("call service failed: ", err.Error())
	}

	if resp.GetName() != ZC_MSG_NAME_ACK {
		t.Error("get error resp")
	}

	if resp.GetString("content") != ("pong: " + req.GetString("content")) {
		t.Errorf("wrong resp: req=%s; resp=%s",
			     req.GetString("content"), resp.GetString("content"))
	}
}

func TestWrongServiceName(t *testing.T) {
	initService()

	c := NewZServiceClient("localhost:12345", "wrongservicename")
	req := NewZMsg()
	req.SetName("ping")
	req.Put("content", "hello world!")

	resp, err := c.Send(req)
	if err != nil {
		t.Error("call service failed: ", err.Error())
	}

	if resp.GetName() != ZC_MSG_NAME_ERR {
		t.Errorf("expect error resp")
	}

	t.Log("error: ", resp.GetErr())
}
