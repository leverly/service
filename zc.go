package zc

import (
	"math/rand"
	"time"
)

type ZCloudConfig struct {
	serviceMap map[string][]string
}

func (config *ZCloudConfig) ParseFrom(file string) (error) {
	if len(file) > 0 {
		return nil
	}
	return nil
}

func NewCloudConfig() (*ZCloudConfig) {
	return &ZCloudConfig{serviceMap : make(map[string][]string, 0)}
}

func (config *ZCloudConfig) AddServiceAddr(serviceName string, addr string) {
	addrs := config.serviceMap[serviceName]
	if addrs == nil {
		config.serviceMap[serviceName] = []string{addr}
	} else {
		config.serviceMap[serviceName] = append(addrs, addr)
	}
}

var cloudConfig *ZCloudConfig
var MyService ZService

func NewMsg(name string, version uint8) (*ZMsg) {
	m := NewZMsg()
	m.SetName(name)
	m.SetVersion(uint8(version))
	return m
}

func NewObject() ZObject {
	return NewZObject()
}

func Service(name string) (client *ZServiceClient) {
	serviceAddrs := cloudConfig.serviceMap[name]
	if len(serviceAddrs) <= 0 { panic("call wrong service: " + name) }
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(serviceAddrs))
	return &ZServiceClient{cloudConfig.serviceMap[name][i], name}
}

func Store(classNames ...string) (client *ZStoreClient) {
	return NewZStoreClient(classNames...)
	return nil
}

func Init() (error) {
	if cloudConfig == nil {
		config := NewCloudConfig()
		err := config.ParseFrom("zc.ini")
		if err != nil { return err }
		SetCloudConfig(config)
	}
	return nil
}

func SetCloudConfig(config *ZCloudConfig) {
	cloudConfig = config
}
