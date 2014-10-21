package zc

import (
	"testing"
	log"zc-common-go/glog"
	"encoding/json"
)

type MyStruct struct {
	value int
}

func TestZObject(t *testing.T) {
	o := NewZObject()
	o.PutString("name", "lihailei")
	o.PutInt("age", 1234)
	o.PutBool("vip", false)
	o.AddInt("myson", 1)
	o.AddInt("myson", 2)

	if o.GetString("name") != "lihailei" {
		t.Error("wrong name")
	}

	if o.GetBool("vip") != false {
		t.Error("wrong vip flag")
	}

	if len(o.GetBools("vips")) != 0 {
		t.Error("wrong vip flags")
	}

	if len(o.GetStrings("names")) != 0 {
		t.Error("wrong names")
	}

	o2 := NewZObject()
	o2.PutString("name", "lihailei2")
	o2.PutInt("age", 233)
	o2.PutBool("vip", true)
	o.AddObject("sons", o2)

	o3 := NewZObject()
	o3.PutString("name", "lihailei3")
	o3.PutInt("age", 23)
	o3.PutBool("vip", false)
	o.AddObject("sons", o3)

	opayload, _ := json.Marshal(o)
	o2payload, _ := json.Marshal(o2)
	o3payload, _ := json.Marshal(o3)

	log.Info("o: ", string(opayload))
	log.Info("o2: ", string(o2payload))
	log.Info("o3: ", string(o3payload))

	o4 := NewZObject()
	o5 := o4
	o4["testkey"] = "testvalue"
	log.Info(o4, o5)

	var gt = getMap()
	log.Infof("%p", gt)
	log.Info(len(gt))
	gt["lihailei"] = 2
	log.Info(gt)
}

func TestZObjectFormat(t *testing.T) {

	// create and find all
	account := ZObject{
		"accountid" : "12345",
		"name" : "lihailei",
		"age" : 1,
		"email" : "lihailei@zc.com",
		"addr" : ZObject{
			"province" : "hebei",
			"city" : "shijiazhuang",
		},
	}

	if account.GetString("accountid") != "12345" ||
	   account.GetString("name") != "lihailei" ||
	   account.GetInt("age") != 1 ||
	   account.GetString("email") != "lihailei@zc.com" ||
	   account.GetObject("addr").GetString("province") != "hebei" ||
	   account.GetObject("addr").GetString("city") != "shijiazhuang" {
		t.Error(account)
	}
}

func getMap() map[string]int {
	var test map[string]int = make(map[string]int)
	test2 := test
	log.Infof("%p", test)
	log.Infof("%p", test2)
	return test2
}
