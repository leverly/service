package zc

import (
	"net/http"
	"errors"
	"strings"
	log"zc-common-go/glog"
	"strconv"
	"net/url"
)

type ZServiceConfig struct {
	Port string
}

type ZServiceHandler func (/*req*/ *ZMsg, /*resp*/ *ZMsg)

type ZService struct {
	name string
	handlers map[string]ZServiceHandler
	port string
}

func NewZService() (*ZService) {
	s := &ZService {name : "", handlers : nil}
	return s
}

func (s *ZService) Init(name string, config *ZServiceConfig) {
	s.name = name
	s.port = config.Port
	s.handlers = make(map[string]ZServiceHandler, 10)
}

func (s *ZService) Start() (error) {
	log.Infof("begin to listen and servie: port=%s", s.port)
	err := http.ListenAndServe(":" + s.port, s)
	if err != nil { log.Infof("listenAndServe failed: %s", err.Error()) }
	return err
}

func (s *ZService) Handle(name string, handler ZServiceHandler) {
	s.handlers[name] = handler
}

func (s *ZService) parseMsgName(uri string) (name string, err error) {
	err = errors.New("wrong uri")
	i := strings.Index(uri, s.name)
	if i < 0 {
		return "", err
	}

	i += len(s.name) + 1
	if i > len(uri) - 1 { return "", err}
	if uri[len(uri)-1] == '/' {
		return uri[i:len(uri)-1], nil
	} else {
		j := strings.Index(uri[i:], "?")
		if j > 0 { return uri[i:i+j], nil }
		return uri[i:], nil
	}
	return "", err
}

func (s *ZService) parseReqMsg(req *http.Request) (reqMsg *ZMsg, err error) {
	msgName, err := s.parseMsgName(req.RequestURI)
	if err != nil { return nil, err }

	reqMsg = NewZMsg()
	reqMsg.SetName(msgName)

	if req.ContentLength > 0 && req.Header.Get("Content-Type") == ZC_MSG_PAYLOAD_OBJECT {
		payload := make([]byte, req.ContentLength)
		l, err := req.Body.Read(payload)
		if l != len(payload) { return nil, errors.New("read payload failed") }
		reqMsg.SetPayload(payload, req.Header.Get("Content-Type"))
		err = reqMsg.decodeObject()
		if err != nil {
			log.Info(string(payload))
			return nil, errors.New("unmarshal payload failed")
		}
		return reqMsg, nil
	} else if len(req.URL.RawQuery) > 0 {
		urlValues, err := url.ParseQuery(req.URL.RawQuery)
		if err != nil { return nil, errors.New("unmarshal raw query params failed") }
		for key, values := range urlValues {
			reqMsg.Put(key, values[0])
		}
		reqMsg.encodeObject()
		log.Info("get req msg from raw query: %", reqMsg)
		return reqMsg, nil
	}

	return nil, errors.New("parse request failed")
}

func (s *ZService) writeResp(resp *ZMsg, w http.ResponseWriter) {
	if resp.hasObjectData() {
		err := resp.encodeObject()
		if err != nil {
			w.Header().Set("Content-Length", "0")
			w.Header().Set(ZC_MSG_NAME_HEADER, ZC_MSG_NAME_ERR)
			return
		}
	}

	payload, payloadFormat := resp.GetPayload()
	if payload == nil {
		w.Header().Set("Content-Length", "0")
	} else {
		w.Header().Set("Content-Length", strconv.Itoa(len(payload)))
	}

	w.Header().Set("Content-Type", payloadFormat)
	w.Header().Set(ZC_MSG_NAME_HEADER, resp.GetName())

	if len(payload) > 0 {
		n, err := w.Write(payload)
		if n != len(payload) || err != nil {
			log.Infof("write payload failed: n=%d, err=%s", n, err.Error())
		}
	}
}

func (s *ZService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	respMsg := NewZMsg()

	if len(req.RequestURI) <= len(s.name) + 1 || req.RequestURI[0] != '/' {
		respMsg.SetErr("wrong service uri: " + req.RequestURI)
		s.writeResp(respMsg, w)
		return
	}

	if req.RequestURI[1:len(s.name)+1] != s.name {
		respMsg.SetErr("wrong service name: " + req.RequestURI)
		s.writeResp(respMsg, w)
		return
	}

	reqMsg, err := s.parseReqMsg(req)
	if err != nil {
		respMsg.SetErr(err.Error())
		s.writeResp(respMsg, w)
		return
	}

	handler := s.handlers[reqMsg.GetName()]
	if handler == nil {
		respMsg.SetName(ZC_MSG_NAME_ERR)
		respMsg.Put("error", "no handler for this msg name: " + reqMsg.GetName())
		s.writeResp(respMsg, w)
		return
	}

	handler(reqMsg, respMsg)
	s.writeResp(respMsg, w)
	return
}

