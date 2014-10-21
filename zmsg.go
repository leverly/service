package zc

import "encoding/json"

const (
	ZC_MSG_PAYLOAD_OBJECT = "application/x-zc-object"
	ZC_MSG_PAYLOAD_JSON = "text/json"
	ZC_MSG_PAYLOAD_THRIFT = "application/x-thrift"
	ZC_MSG_PAYLOAD_PROTOBUF = "application/x-protobuf"

	ZC_MSG_NAME_HEADER = "x-zc-msg-name"
	ZC_MSG_NAME_ACK = "x-zc-ack"
	ZC_MSG_NAME_ERR = "x-zc-err"

	ZC_MSG_ERR_VERSION = uint8(2)
)

type ZMsg struct {
	ZObject

	name 			string
	version 		uint8
	payloadFormat 	string
	payload 		[]byte
}

func (m *ZMsg) SetName(name string) {
	m.name = name
}

func (m *ZMsg) GetName() (string) {
	return m.name
}

func (m *ZMsg) SetAck() {
	m.name = ZC_MSG_NAME_ACK
}

func (m *ZMsg) IsAck() (bool) {
	return m.name == ZC_MSG_NAME_ACK
}

func (m *ZMsg) SetErr(err string) {
	m.name = ZC_MSG_NAME_ERR
	m.PutString("error", err)
}

func (m *ZMsg) IsErr() (bool) {
	return m.name == ZC_MSG_NAME_ERR
}

func (m *ZMsg) GetErr() (err string) {
	return m.GetString("error")
}

func (m *ZMsg) SetVersion(v uint8) {
	m.version = v
}

func (m *ZMsg) GetVersion() (uint8) {
	return m.version
}

func (m *ZMsg) SetPayload(payload []byte, format string) {
	m.payload = payload
	m.payloadFormat = format
}

func (m *ZMsg) GetPayload() ([]byte, string) {
	return m.payload, m.payloadFormat
}

func (m *ZMsg) GetPayloadFormat() (string) {
	return m.payloadFormat
}

func (m *ZMsg) hasObjectData() (bool) {
	return len(m.ZObject) > 0
}

func (m *ZMsg) encodeObject() (err error) {
	if len(m.ZObject) <= 0 || len(m.payload) > 0 {
		panic("can not encode object")
	}
	m.payloadFormat = ZC_MSG_PAYLOAD_OBJECT
	m.payload, err = json.Marshal(m.ZObject)
	return err
}

func (m *ZMsg) decodeObject() (err error) {
	if len(m.ZObject) > 0 ||
		len(m.payload) <= 0 ||
		m.payloadFormat != ZC_MSG_PAYLOAD_OBJECT {
		panic("can not decode object")
	}
	return json.Unmarshal(m.payload, &(m.ZObject))// &m.attrs)
}

func NewZMsg() (*ZMsg) {
	m := new(ZMsg)
	m.ZObject = NewZObject()
	return m
}

