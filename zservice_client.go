package zc

import (
	"net/http"
	"bytes"
	"errors"
	log"zc-common-go/glog"
)

type ZServiceClient struct {
	serviceAddr string
	serviceName string
}

func NewZServiceClient(serviceAddr string, serviceName string) (*ZServiceClient) {
	return &ZServiceClient{serviceAddr : serviceAddr, serviceName : serviceName}
}

func (c *ZServiceClient) Send(req *ZMsg) (resp *ZMsg, err error) {
	if req.hasObjectData() {
		err = req.encodeObject()
		if err != nil { return nil, err }
	}

	url := "http://" + c.serviceAddr + "/" + c.serviceName +"/" + req.GetName()
	payload, payloadFormat := req.GetPayload()
	log.Infof("client send req: serive=%s; name=%s; payload=%s",
		c.serviceName, req.GetName(), string(payload))

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil { return nil, errors.New("failed to create http req") }
	httpReq.Header.Set("Content-Type", payloadFormat)

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil { return nil, err }

	resp = NewZMsg()
	resp.SetName(httpResp.Header.Get(ZC_MSG_NAME_HEADER))

	if httpResp.ContentLength <= 0 { return resp, nil }
	if httpResp.Header.Get("Content-Type") != ZC_MSG_PAYLOAD_OBJECT {
		return nil, errors.New("wrong resp format, only support text/json")
	}

	respPayload := make([]byte, httpResp.ContentLength)
	l, err := httpResp.Body.Read(respPayload)
	if err != nil { return nil, err }
	if l != len(respPayload) { return nil, errors.New("failed to read http resp body") }

	log.Infof("client get resp: serive=%s; name=%s; payload=%s",
		       c.serviceName, resp.GetName(), string(respPayload))
	resp.SetPayload(respPayload, httpResp.Header.Get("Content-Type"))
	err = resp.decodeObject()
	if err != nil { return nil, err }
	return resp, err
}
