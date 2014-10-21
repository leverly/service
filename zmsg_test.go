package zc

import (
	"testing"
	log"zc-common-go/glog"
)


func TestMarshalParams(t *testing.T) {
	m := NewZMsg()
	m.SetName("hello")
	m.SetVersion(uint8(1))
	m.Put("text", "hello world")
	m.Put("from", "lihailei")
	m.Put("to", "zhidong")
	m.encodeObject()
	payload, payloadFormat := m.GetPayload()
	if payloadFormat != ZC_MSG_PAYLOAD_OBJECT {
		t.Log("wrong palyload format")
		t.FailNow()
	}

	s := NewZMsg()
	s.SetName(m.GetName())
	s.SetVersion(m.GetVersion())
	s.SetPayload(payload, payloadFormat)
	s.decodeObject()
	spayload, spayloadFormat := s.GetPayload()
	if spayloadFormat != ZC_MSG_PAYLOAD_OBJECT {
		t.Error("wrong palyload format")
	}

	if len(payload) != len(spayload) || len(m.ZObject) != len(s.ZObject) {
		t.Errorf("wrong marshaled payload: %s => %s\n",
				 string(payload), string(spayload))
	}

	for key, value := range m.ZObject {
		if value != s.Get(key) {
			t.Errorf("wrong marshaled payload: %s => %s\n",
					 string(payload), string(spayload))
		}
	}

	zo := NewZObject()
	zo.Put("name", "lihailei")
	create := NewZMsg()
	create.Put("zc-class", "account")
	create.Put("zc-object", zo)
	create.encodeObject()
	createPayload, _ := create.GetPayload()
	log.Info(string(createPayload))
}
