package zc

import (
	"reflect"
	"errors"
)

type ZStoreClient struct {
	className 	string
}

type ZStoreCreate struct {
	className string
	object 	ZObject
}

type ZStoreFind struct {
	className string
	find ZObject
}

type ZStoreDelete struct {
	className string
	object ZObject
}

type ZStoreUpdate struct {
	className string
	object 	ZObject
}

type ZStoreReplace struct {
	className string
	object 	ZObject
}

type ZStoreQuery struct {
	className 	string
	query 		ZObject
}

type ZStoreBatch struct {
	ops []interface {}
	batch []ZObject
	client *ZStoreClient
}

func NewZStoreClient(classNames ...string) (*ZStoreClient) {
	if len(classNames) > 1 { panic("can have one class name at most") }
	client := &ZStoreClient{}
	if len(classNames) == 1 { client.className = classNames[0] }
	return client
}

// zo := ZC.NewObject()
// zo.PutString("name", "lihailei")
// err := ZC.Store("account").Create(zo).Execute()
// zo.GetInt("zd-id")
func (c *ZStoreClient) Create(params ...interface {}) (*ZStoreCreate) {
	if len(params) > 1 { panic("invalid params") }
	create := &ZStoreCreate{className : c.className, object : NewObject()}
	if len(params) == 0 { return create }

	create.object = params[0].(ZObject)
	return create
}

func (create *ZStoreCreate) Put(key string, value interface {}) (*ZStoreCreate) {
	create.object.Put(key, value)
	return create
}

func (create *ZStoreCreate) Execute() (error) {
	req := NewMsg("create", uint8(1))
	req.PutString("zc-class", create.className)
	req.PutObject("zc-object", create.object)
	resp, err := Service("zc-store").Send(req)
	if err != nil { return err }
	if resp.IsErr() { return errors.New("get error resp: " + resp.GetErr()) }
	return nil
}

func parseZObjectKVs(zo ZObject, canEmpty bool, kvs []interface {}) {
	if len(kvs) % 2 != 0 || (!canEmpty && len(kvs) == 0) {
		panic("invalid params")
	}
	for i := 0; i < len(kvs) / 2; i++ {
		key, value := kvs[i], kvs[i+1]
		if reflect.TypeOf(key).Kind() != reflect.String ||
				(reflect.TypeOf(value).Kind() != reflect.Int64 &&
					reflect.TypeOf(value).Kind() != reflect.Float64 &&
					reflect.TypeOf(value).Kind() != reflect.String &&
					reflect.TypeOf(value).Kind() != reflect.Bool) {
			panic("wrong param type")
		}
		zo.Put(key.(string), value)
	}
}

// zo, err := ZC.Store("account").Find("zc-id", 1234)
// ZC.Store("account").Find("zc-id", 1234).Select("name").Execute()

func (c *ZStoreClient) Find(kvs ...interface {}) (*ZStoreFind) {
	f := &ZStoreFind{className : c.className, find : NewObject()}
	zo := NewObject()
	parseZObjectKVs(zo, false, kvs)
	f.find.PutObject("zc-object", zo)
	return f
}

func (f *ZStoreFind) setSelectKeys(keys []string) (*ZStoreFind) {
	if len(keys) == 0 { panic("invalid params") }
	for i := 0; i < len(keys); i++ {
		f.find.AddString("zc-select", keys[i])
	}
	return f
}

func (f *ZStoreFind) Select(keys ...string) (*ZStoreFind) {
	return f.setSelectKeys([]string(keys))
}

func (f *ZStoreFind) Execute() (ZObject, error) {
	zo := NewObject()
	req := NewMsg("find", uint8(1))
	req.PutString("zc-class", f.className)
	req.PutObject("zc-find", f.find)
	resp, err := Service("zc-store").Send(req)
	if err != nil { return zo, err }
	if resp.IsErr() { return zo, errors.New("get error resp: " + resp.GetErr()) }
	if !resp.Exists("zc-object") { return zo, errors.New("error resp, no objects returned") }
	zo = resp.GetObject("zc-object")
	return zo, nil
}

// err := ZC.Store("account").Delete("zc-id", 1234).Execute()
func (c *ZStoreClient) Delete(kvs ...interface {}) (*ZStoreDelete) {
	d := &ZStoreDelete{className : c.className, object : NewObject()}
	parseZObjectKVs(d.object, false, kvs)
	return d
}

func (d *ZStoreDelete) Execute() (error) {
	req := NewMsg("delete", uint8(1))
	req.PutString("zc-class", d.className)
	req.PutObject("zc-object", d.object)
	resp, err := Service("zc-store").Send(req)
	if err != nil { return err }
	if resp.IsErr() { return errors.New("error resp: " + resp.GetErr()) }
	return nil
}

// zo := ZC.NewObject()
// zo.PutInt("zc-id", 1234)
// zo.PutString("name", "lihailei")
// err := ZC.Store("account").Update(zo).Execute()

// err := ZC.Store("account").Update("zc-id", 1234)
// 				  .PutString("name", "lihailei2")
// 				  .PutInt("age", 24)
// 				  .Execute()
func (c *ZStoreClient) Update(params ...interface {}) (*ZStoreUpdate) {
	update := &ZStoreUpdate{className : c.className, object : NewObject()}
	if len(params) == 0 { return update }

	if len(params) == 1 {
		update.object = params[0].(ZObject)
		return update
	}

	parseZObjectKVs(update.object, true, params)
	return update
}

func (update *ZStoreUpdate) Put(key string, value interface {}) (*ZStoreUpdate) {
	update.object.Put(key, value)
	return update
}

func (update *ZStoreUpdate) Execute() (error) {
	req := NewMsg("update", uint8(1))
	req.PutString("zc-class", update.className)
	req.PutObject("zc-object", update.object)
	resp, err := Service("zc-store").Send(req)
	if err != nil { return err }
	if resp.IsErr() { return errors.New("error resp: " + resp.GetErr()) }
	return nil
}

// zo := ZC.NewObject()
// zo.PutInt("zc-id", 1234)
// zo.PutString("name", "lihailei")
// err := ZC.Store("account").Replace(zo).Execute()

// err := ZC.Store("account").Replace("zc-id", 1234)
// 				  .PutString("name", "lihailei2")
// 				  .PutInt("age", 24)
// 				  .Execute()
func (c *ZStoreClient) Replace(params ...interface {}) (*ZStoreReplace) {
	replace := &ZStoreReplace{className : c.className, object : NewObject()}
	if len(params) == 0 { return replace }

	if len(params) == 1 {
		replace.object = params[0].(ZObject)
		return replace
	}

	parseZObjectKVs(replace.object, true, params)
	return replace
}

func (replace *ZStoreReplace) Put(key string, value interface {}) (*ZStoreReplace) {
	replace.object.Put(key, value)
	return replace
}

func (replace *ZStoreReplace) Execute() (error) {
	req := NewMsg("put", uint8(1))
	req.PutString("zc-class", replace.className)
	req.PutObject("zc-object", replace.object)
	resp, err := Service("zc-store").Send(req)
	if err != nil { return err }
	if resp.IsErr() { return errors.New("error resp: " + resp.GetErr()) }
	return nil
}

// zos, err := ZC.Store("DeviceEvent").Query()
// 								      .WhereEQ("homeid", 1234, "deviceid", 345)
//									  .Select("timestamp", "status")
// 									  .Execute()
func (c *ZStoreClient) Query() (*ZStoreQuery) {
	return &ZStoreQuery{className : c.className, query : NewObject()}
}

func (q *ZStoreQuery) WhereEQ(kvs ...interface {}) (*ZStoreQuery) {
	zo := NewObject()
	parseZObjectKVs(zo, false, kvs)
	q.query.PutObject("zc-eq", zo)
	return q
}

func (q *ZStoreQuery) Select(keys ...string) (*ZStoreQuery) {
	if len(keys) % 2 != 0 || len(keys) == 0 { panic("invalid params") }
	for i := 0; i < len(keys); i++ {
		q.query.AddString("zc-select", keys[i])
	}
	return q
}

func (q *ZStoreQuery) Execute() ([]ZObject, error) {
	zos := make([]ZObject, 0)
	req := NewMsg("query", uint8(1))
	req.PutString("zc-class", q.className)
	req.PutObject("zc-query", q.query)
	resp, err := Service("zc-store").Send(req)
	if err != nil { return zos, err }
	if resp.IsErr() { return zos, errors.New("error resp: " + resp.GetErr()) }
	if !resp.Exists("zc-objects") { return zos, errors.New("error resp, no objects returned") }
	zos = resp.GetObjects("zc-objects")
	return zos, nil
}

type ZStoreBatchClass struct {
	className string
	batch *ZStoreBatch
}

func (c *ZStoreBatchClass) Find(kvs ...interface {}) (*ZStoreFind)  {
	f := c.batch.client.Find(kvs)
	c.batch.ops = append(c.batch.ops, f)
	return f
}

func (c *ZStoreBatchClass) Delete(params ...interface {}) (*ZStoreDelete)  {
	delete := c.batch.client.Delete(params)
	c.batch.ops = append(c.batch.ops, delete)
	return delete
}

func (c *ZStoreBatchClass) Update(params ...interface {}) (*ZStoreUpdate) {
	update := c.batch.client.Update(params)
	c.batch.ops = append(c.batch.ops, update)
	return update
}

func (c *ZStoreBatchClass) Replace(params ...interface {}) (*ZStoreReplace)  {
	replace := c.batch.client.Replace(params)
	c.batch.ops = append(c.batch.ops, replace)
	return replace
}

func (c *ZStoreBatchClass) Query() (*ZStoreQuery) {
	q := c.batch.client.Query()
	c.batch.ops = append(c.batch.ops, q)
	return q
}

func (c *ZStoreClient) Batch() (*ZStoreBatch) {
	if len(c.className) > 0 { panic("this client can not do batch") }
	return &ZStoreBatch{ops : make([]interface {}, 0),
		                batch : make([]ZObject, 0),
		                client : c}
}

func (b *ZStoreBatch) Class(className string) (*ZStoreBatchClass) {
	return &ZStoreBatchClass{className : className, batch : b}
}

func (b *ZStoreBatch) Execute() (error) {
	return nil
}


