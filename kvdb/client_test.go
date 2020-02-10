package kvdb

import "testing"

const (
	kmsAddr = "127.0.0.1:8080"
)

func TestKv_Get(t *testing.T) {
	kvc, err := NewKVer(kmsAddr, "ankr-certmgr")
	if err != nil {
		t.Log(err)
		return
	}

	rsp, err := kvc.Get("/cls-0b705769-d866-4ebc-9c5b-9e7b18e4990c")

	if err != nil {
		t.Log(err)
		return
	}
	t.Logf("res: %#v \n", rsp)
}

func TestKv_Put(t *testing.T) {
	kvc, err := NewKVer(kmsAddr, "ankr-user")
	if err != nil {
		t.Log(err)
		return
	}
	if err = kvc.Put("test", map[string]string{"hello": "world"}); err != nil {
		t.Logf("error: %v\n", err)
		return
	}
	t.Log("work over")
}
