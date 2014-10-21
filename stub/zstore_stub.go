package stub

import (
	"zc"
	log "zc-common-go/glog"
)

type ZStoreStub struct {
	zc.ZService

	store  map[string]map[string]zc.ZObject
	schema map[string][]string
}

func (s *ZStoreStub) initStub() {
	s.store = make(map[string]map[string]zc.ZObject, 10)
	s.schema = make(map[string][]string, 10)
	s.schema["account"] = []string{"accountid"}
	s.schema["device"] = []string{"deviceid"}
}

func (s *ZStoreStub) checkClass(className string) bool {
	classSchema := s.schema[className]
	return classSchema != nil && len(classSchema) > 0
}

func (s *ZStoreStub) getPrimaryKey(className string, zo zc.ZObject) string {
	var primaryKey string
	for _, key := range s.schema[className] {
		if !zo.CheckString(key) {
			return ""
		}
		primaryKey += ("[" + zo.GetString(key) + "]")
	}
	return primaryKey
}

func (s *ZStoreStub) handleCreate(req *zc.ZMsg, resp *zc.ZMsg) {
	payload, _ := req.GetPayload()
	if !req.CheckString("zc-class") ||
		!req.CheckObject("zc-object") {
		resp.SetErr("invalid create request")
		log.Infof("invalid create request: %s", string(payload))
		return
	}

	className := req.GetString("zc-class")
	if !s.checkClass(className) {
		resp.SetErr("invalid class name")
		log.Infof("invalid create request: %s", string(payload))
		return
	}

	classStore := s.store[className]
	if classStore == nil {
		s.store[className] = make(map[string]zc.ZObject, 10)
		classStore = s.store[className]
	}

	zo := req.GetObject("zc-object")
	primaryKey := s.getPrimaryKey(className, zo)
	classStore[primaryKey] = zo
	resp.SetAck()
	resp.Put("zc-object", zo)
}

func (s *ZStoreStub) handleDelete(req *zc.ZMsg, resp *zc.ZMsg) {
	payload, _ := req.GetPayload()
	if !req.CheckExists("zc-class", "zc-object") {
		resp.SetErr("invalid delete request")
		log.Infof("invalid delete request: %s", string(payload))
		return
	}

	className := req.GetString("zc-class")
	if !s.checkClass(className) {
		resp.SetErr("invalid class name")
		log.Infof("invalid create request: %s", string(payload))
		return
	}

	classStore := s.store[className]
	if classStore == nil || len(classStore) <= 0 {
		resp.SetAck()
		return
	}

	zo := req.GetObject("zc-object")
	primaryKey := s.getPrimaryKey(className, zo)
	delete(classStore, primaryKey)
	resp.SetAck()
	return
}

func (s *ZStoreStub) handleUpdate(req *zc.ZMsg, resp *zc.ZMsg) {
	payload, _ := req.GetPayload()
	if !req.CheckExists("zc-class", "zc-object") {
		resp.SetErr("invalid update request")
		log.Infof("invalid delete request: %s", string(payload))
		return
	}

	className := req.GetString("zc-class")
	if !s.checkClass(className) {
		resp.SetErr("invalid class name")
		log.Infof("invalid update request: %s", string(payload))
		return
	}

	classStore := s.store[className]
	if classStore == nil {
		s.store[className] = make(map[string]zc.ZObject, 10)
		classStore = s.store[className]
	}

	zo := req.GetObject("zc-object")
	primaryKey := s.getPrimaryKey(className, zo)
	oldZO := classStore[primaryKey]
	if oldZO != nil {
		for _, key := range zo.GetKeys() {
			oldZO.Put(key, zo.Get(key))
		}
	}
	resp.SetAck()
	return
}

func (s *ZStoreStub) handlePut(req *zc.ZMsg, resp *zc.ZMsg) {
	payload, _ := req.GetPayload()
	if !req.CheckExists("zc-class", "zc-object") {
		resp.SetErr("invalid put request")
		log.Infof("invalid put request: %s", string(payload))
		return
	}

	className := req.GetString("zc-class")
	if !s.checkClass(className) {
		resp.SetErr("invalid class name")
		log.Infof("invalid put request: %s", string(payload))
		return
	}

	classStore := s.store[className]
	if classStore == nil {
		s.store[className] = make(map[string]zc.ZObject, 10)
		classStore = s.store[className]
	}

	zo := req.GetObject("zc-object")
	primaryKey := s.getPrimaryKey(className, zo)
	classStore[primaryKey] = zo
	resp.SetAck()
	return
}

func (s *ZStoreStub) handleFind(req *zc.ZMsg, resp *zc.ZMsg) {
	payload, _ := req.GetPayload()
	if !req.CheckExists("zc-class", "zc-find") {
		resp.SetErr("invalid find request")
		log.Infof("invalid find request: %s", string(payload))
		return
	}

	className := req.GetString("zc-class")
	if !s.checkClass(className) {
		resp.SetErr("invalid class name")
		log.Infof("invalid find request: %s", string(payload))
		return
	}

	classStore := s.store[className]
	if classStore == nil {
		resp.SetAck()
		return
	}

	find := req.GetObject("zc-find")
	zo := find.GetObject("zc-object")
	selectKeys := find.GetStrings("zc-select")
	primaryKey := s.getPrimaryKey(className, zo)
	dataZO := classStore[primaryKey]
	if dataZO == nil {
		resp.SetAck()
		return
	}

	if len(selectKeys) <= 0 {
		resp.SetAck()
		resp.Put("zc-object", dataZO)
		return
	}

	selectZO := zc.NewObject()
	for _, key := range selectKeys {
		selectZO.Put(key, dataZO.Get(key))
	}
	resp.SetAck()
	resp.Put("zc-object", selectZO)
	return
}

func (s *ZStoreStub) handleQuery(req *zc.ZMsg, resp *zc.ZMsg) {
	payload, _ := req.GetPayload()
	if !req.CheckExists("zc-class", "zc-query") {
		resp.SetErr("invalid query request")
		log.Infof("invalid query request: %s", string(payload))
		return
	}

	className := req.GetString("zc-class")
	if !s.checkClass(className) {
		resp.SetErr("invalid class name")
		log.Infof("invalid query request: %s", string(payload))
		return
	}

	classStore := s.store[className]
	if classStore == nil {
		resp.SetAck()
		return
	}

	query := req.GetObject("zc-query")
	if !query.CheckExists("zc-eq") {
		resp.SetErr("no eq condition to find objects")
		return
	}

	eq := query.GetObject("zc-eq")
	selectKeys := query.GetStrings("zc-select")
	primaryKey := s.getPrimaryKey(className, eq)
	dataZO := classStore[primaryKey]
	if dataZO == nil {
		resp.SetAck()
		return
	}

	selectZO := zc.NewObject()
	if len(selectKeys) <= 0 {
		selectZO = dataZO
	} else {
		for _, key := range selectKeys {
			selectZO.Put(key, dataZO.Get(key))
		}
	}
	resp.SetAck()
	resp.AddObject("zc-objects", selectZO)
	return
}

func NewZStoreStub(config *zc.ZServiceConfig) *ZStoreStub {
	s := &ZStoreStub{}
	s.initStub()

	s.Init("zc-store", config)
	s.Handle("create", zc.ZServiceHandler(func(req *zc.ZMsg, resp *zc.ZMsg) {
		s.handleCreate(req, resp)
	}))
	s.Handle("delete", zc.ZServiceHandler(func(req *zc.ZMsg, resp *zc.ZMsg) {
		s.handleDelete(req, resp)
	}))
	s.Handle("update", zc.ZServiceHandler(func(req *zc.ZMsg, resp *zc.ZMsg) {
		s.handleUpdate(req, resp)
	}))
	s.Handle("put", zc.ZServiceHandler(func(req *zc.ZMsg, resp *zc.ZMsg) {
		s.handlePut(req, resp)
	}))
	s.Handle("find", zc.ZServiceHandler(func(req *zc.ZMsg, resp *zc.ZMsg) {
		s.handleFind(req, resp)
	}))
	s.Handle("query", zc.ZServiceHandler(func(req *zc.ZMsg, resp *zc.ZMsg) {
		s.handleQuery(req, resp)
	}))
	return s
}
