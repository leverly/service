package stub

import (
	"zc"
)

func StartStoreStub(config *zc.ZServiceConfig) {
	storeStub := NewZStoreStub(config)
	storeStub.Start()
}
