package stub

import (
	"testing"
	"zc"
)

var storeStubConfig = &zc.ZServiceConfig{Port:"12556"}
var storeStub *ZStoreStub = nil
var zcloudConfig = zc.NewCloudConfig()

func startStoreStub() {
	if storeStub != nil { return }

	storeStub = NewZStoreStub(storeStubConfig)
	zcloudConfig.AddServiceAddr("zc-store", "localhost:12556")
	zc.SetCloudConfig(zcloudConfig)
	go storeStub.Start()
}

func TestStoreStub(t *testing.T) {
	startStoreStub()

	// create and find
	err := zc.Store("account").Create().
	                           Put("accountid", "12345").
					           Put("name", "lihailei").
							   Put("age", 1).
					           Execute()
	if err != nil { t.Error(err) }

	zo, err := zc.Store("account").Find("accountid", "12345").Execute()
	if err != nil { t.Error(err) }

	if zo.GetString("accountid") != "12345" ||
	   zo.GetString("name") != "lihailei" {
		t.Error(zo)
	}

	// update and find
	err = zc.Store("account").Update("accountid", "12345").
						      Put("email", "lihailei@zc.com").
							  Put("age", 2).
							  Execute()
	if err != nil { t.Error(err) }

	zo, err = zc.Store("account").Find("accountid", "12345").Execute()
	if err != nil { t.Error(err) }

	if zo.GetString("accountid") != "12345" ||
	   zo.GetString("name") != "lihailei" ||
	   zo.GetString("email") != "lihailei@zc.com" ||
	   zo.GetInt("age") != 2 {
		t.Error(zo)
	}

	// put and find
	err = zc.Store("account").Replace("accountid", "12345").
							  Put("name", "lihailei").
						      Put("age", 3).
						      Execute()
	if err != nil { t.Error(err) }

	zo, err = zc.Store("account").Find("accountid", "12345").Execute()
	if err != nil { t.Error(err) }

	if zo.GetString("accountid") != "12345" ||
		zo.GetString("name") != "lihailei" ||
		zo.Exists("email") ||
		zo.GetInt("age") != 3 {
		t.Error(zo)
	}
}

func TestStoreStubUsingZObject(t *testing.T) {
	startStoreStub()

	// create and find
	account := zc.NewObject()
	account.Put("accountid", "12345")
	account.Put("name", "lihailei")
	account.Put("age", 1)
	err := zc.Store("account").Create(account).Execute()
	if err != nil { t.Error(err) }

	foundAccount, err := zc.Store("account").Find("accountid", "12345").Execute()
	if err != nil { t.Error(err) }
	if foundAccount.GetString("accountid") != "12345" ||
		foundAccount.GetString("name") != "lihailei" ||
		foundAccount.GetInt("age") != 1 {
		t.Error(foundAccount)
	}

	// update and find
	account.Put("email", "lihailei@zc.com")
	account.Put("age", 2)
	err = zc.Store("account").Update(account).Execute()
	if err != nil { t.Error(err) }

	foundAccount, err = zc.Store("account").Find("accountid", "12345").Execute()
	if err != nil { t.Error(err) }

	if foundAccount.GetString("accountid") != "12345" ||
		foundAccount.GetString("name") != "lihailei" ||
		foundAccount.GetString("email") != "lihailei@zc.com" ||
		foundAccount.GetInt("age") != 2 {
		t.Error(foundAccount)
	}

	// put and find
	newAccount := zc.NewObject()
	newAccount.Put("accountid", "12345")
	newAccount.Put("name", "lihailei")
	newAccount.Put("age", 3)
	err = zc.Store("account").Replace(newAccount).Execute()
	if err != nil { t.Error(err) }

	foundAccount, err = zc.Store("account").Find("accountid", "12345").Execute()
	if err != nil { t.Error(err) }

	if foundAccount.GetString("accountid") != "12345" ||
		foundAccount.GetString("name") != "lihailei" ||
		foundAccount.Exists("email") ||
		foundAccount.GetInt("age") != 3 {
		t.Error(foundAccount)
	}
}

func TestStoreStubFindWithSelect(t *testing.T) {
	startStoreStub()

	// create and find all
	account := zc.NewObject()
	account.Put("accountid", "12345")
	account.Put("name", "lihailei")
	account.Put("age", 1)
	account.Put("email", "lihailei@zc.com")

	err := zc.Store("account").Create(account).Execute()
	if err != nil { t.Error(err) }

	foundAccount, err := zc.Store("account").Find("accountid", "12345").Execute()
	if err != nil { t.Error(err) }
	if foundAccount.GetString("accountid") != "12345" ||
		foundAccount.GetString("name") != "lihailei" ||
		foundAccount.GetInt("age") != 1 ||
		foundAccount.GetString("email") != "lihailei@zc.com" {
		t.Error(foundAccount)
	}

	// find with select keys
	foundAccount, err = zc.Store("account").Find("accountid", "12345").
								  		    Select("accountid", "name", "age").
							                Execute()
	if err != nil { t.Error(err) }

	if foundAccount.GetString("accountid") != "12345" ||
		foundAccount.GetString("name") != "lihailei" ||
		foundAccount.Exists("email") ||
		foundAccount.GetInt("age") != 1 {
		t.Error(foundAccount)
	}
}
