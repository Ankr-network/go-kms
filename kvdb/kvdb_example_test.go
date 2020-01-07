package kvdb_test

import "github.com/Ankr-network/go-kms/kvdb"

func ExampleKV_Get() {
	kvc, err := kvdb.NewKVer("127.0.0.1", "ankr-user")
	if err != nil {
		// handle error
		return
	}

	rsp, err := kvc.Get("/")
	if err != nil {
		// handle error
		return
	}
	// handle response here
	_ = rsp
	// output:
}

func ExampleKV_Put() {
	kvc, err := kvdb.NewKVer("127.0.0.1", "ankr-user")
	if err != nil {
		// handle error
		return
	}
	if err = kvc.Put("test", map[string]string{"hello": "world"}); err != nil {
		// handle error
		return
	}
	// output:
}
