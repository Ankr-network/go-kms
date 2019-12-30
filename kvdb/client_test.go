package kvdb

import "testing"

const (
	kmsAddr = "127.0.0.1"
)

func TestKv_Get(t *testing.T) {
	kvc, err := NewKVer(kmsAddr, "ankr-user")
	if err != nil {
		t.Log(err)
		return
	}
	rsp, err := kvc.Get("test")
	if err != nil {
		t.Log(err)
		return
	}
	t.Logf("%#v \n", rsp)
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
