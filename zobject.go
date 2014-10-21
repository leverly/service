package zc

import "reflect"

const (
	ZC_ATTR_TYPE_ERR = 0
	ZC_ATTR_TYPE_INT = 1
	ZC_ATTR_TYPE_FLOAT = 2
	ZC_ATTR_TYPE_BOOL = 3
	ZC_ATTR_TYPE_STRING = 4
	ZC_ATTR_TYPE_OBJECT = 5
	ZC_ATTR_TYPE_ARRAY = 6
)

type ZObject map[string]interface{}

func (o ZObject) Exists(key string) (bool) {
	return o[key] != nil
}

func GetAttrType(a interface{}) byte {
	switch reflect.TypeOf(a).Kind() {
	case reflect.String:
		return ZC_ATTR_TYPE_STRING
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
	reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return ZC_ATTR_TYPE_INT
	case reflect.Float32, reflect.Float64:
		return ZC_ATTR_TYPE_FLOAT
	case reflect.Bool:
		return ZC_ATTR_TYPE_BOOL
	case reflect.Array:
		return ZC_ATTR_TYPE_ARRAY
	case reflect.Map:
		return ZC_ATTR_TYPE_OBJECT
	}
	return ZC_ATTR_TYPE_ERR
}

func (o ZObject) CheckString(key string) bool {
	v := o[key]
	return v != nil &&
		reflect.TypeOf(o[key]).Kind() == reflect.String
}

func (o ZObject) CheckInt(key string) bool {
	v := o[key]
	return v != nil && GetAttrType(v) == ZC_ATTR_TYPE_INT
}

func (o ZObject) CheckObject(key string) bool {
	v := o[key]
	return v != nil &&
		reflect.TypeOf(o[key]).Kind() == reflect.Map
}

func (o ZObject) CheckExists(keys ...string) (bool) {
	for _, key := range keys {
		if !o.Exists(string(key)) { return false }
	}
	return true
}

func (o ZObject) Put(key string, value interface {}) {
	o[key] = value
}

func (o ZObject) Get(key string) (interface {}) {
	v := o[key]
	if v == nil { return nil }

	kind := reflect.TypeOf(v).Kind()
	switch kind{
	case reflect.String:
		return v.(string)
	case reflect.Map:
		name := reflect.TypeOf(v).Name()
		if name == "ZObject" {
			return o[key].(ZObject)
		}
		return ZObject(v.(map[string]interface{}))
	case reflect.Int:
		return int64(v.(int))
	case reflect.Int64:
		return int64(v.(int64))
	case reflect.Float64:
		return int64(v.(float64))
	case reflect.Int8:
		return int64(v.(int8))
	case reflect.Int16:
		return int64(v.(int16))
	case reflect.Int32:
		return int64(v.(int32))
	case reflect.Uint:
		return int64(v.(uint))
	case reflect.Uint8:
		return int64(v.(uint8))
	case reflect.Uint16:
		return int64(v.(uint16))
	case reflect.Uint32:
		return int64(v.(uint32))
	case reflect.Uint64:
		return int64(v.(uint64))
	}
	return v
}

func (o ZObject) GetKeys() ([]string) {
	var keys = []string{}
	for key, _ := range o {
		keys = append(keys, key)
	}
	return keys
}

func (o ZObject) PutString(key string, value string) {
	o[key] = value
}

func (o ZObject) GetString(key string) (string) {
	v := o[key]
	if v == nil { return "" }
	return v.(string)
}

func (o ZObject) PutInt(key string, value int64) {
	o[key] = value
}

func (o ZObject) GetInt(key string) (int64) {
	v := o.Get(key)
	if v != nil && reflect.TypeOf(v).Kind() == reflect.Int64 {
		return v.(int64)
	}
	return 0
}

func (o ZObject) PutFloat(key string, value float64) {
	o[key] = value
}

func (o ZObject) GetFloat(key string) (float64) {
	if !o.Exists(key) { return 0 }
	return o[key].(float64)
}

func (o ZObject) PutBool(key string, value bool) {
	o[key] = value
}

func (o ZObject) GetBool(key string) (bool) {
	if !o.Exists(key) { return false }
	return o[key].(bool)
}

func (o ZObject) PutObject(key string, value ZObject) {
	o[key] = value
}

func (o ZObject) GetObject(key string) (ZObject) {
	if !o.Exists(key) { return nil }
	name := reflect.TypeOf(o[key]).Name()
	if name == "ZObject" {
		return o[key].(ZObject)
	}
	return ZObject(o[key].(map[string]interface{}))
}

func (o ZObject) AddString(key string, value string) {
	if !o.Exists(key) {
		o[key] = make([]string, 0)
	}
	o[key] = append(o[key].([]string), value)
}

func (o ZObject) AddInt(key string, value int64) {
	if !o.Exists(key) {
		o[key] = make([]int64, 0)
	}
	o[key] = append(o[key].([]int64), value)
}

func (o ZObject) AddFloat(key string, value float64) {
	if !o.Exists(key) {
		o[key] = make([]float64, 0)
	}
	o[key] = append(o[key].([]float64), value)
}

func (o ZObject) AddBool(key string, value bool) {
	if !o.Exists(key) {
		o[key] = make([]bool, 0)
	}
	o[key] = append(o[key].([]bool), value)
}

func (o ZObject) AddObject(key string, value ZObject) {
	if !o.Exists(key) {
		o[key] = make([]ZObject, 0)
	}
	o[key] = append(o[key].([]ZObject), value)
}

func (o ZObject) GetStrings(key string) ([]string) {
	var values = []string{}
	if !o.Exists(key) { return values }
	for _, value := range o[key].([]interface {}) {
		values = append(values, value.(string))
	}
	return values
}

func (o ZObject) GetInts(key string) ([]int64) {
	if !o.Exists(key) { return make([]int64, 0) }
	return o[key].([]int64)
}

func (o ZObject) GetFloats(key string) ([]float64) {
	if !o.Exists(key) { return make([]float64, 0) }
	return o[key].([]float64)
}

func (o ZObject) GetBools(key string) ([]bool) {
	if !o.Exists(key) { return make([]bool, 0) }
	return o[key].([]bool)
}

func (o ZObject) GetObjects(key string) ([]ZObject) {
	if !o.Exists(key) { return make([]ZObject, 0) }
	return o[key].([]ZObject)
}

func NewZObject() (o ZObject) {
	return make(map[string]interface {}, 10)
}
