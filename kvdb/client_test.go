package kvdb

import "testing"

const (
	operaterAddr = "http://192.168.1.92:30386"
	vaultAddr    = "http://192.168.1.92:30050"
)

func TestKv_Get(t *testing.T) {
	kvc, err := NewKVer(operaterAddr, vaultAddr, "ankr-user")
	if err != nil {
		t.Log(err)
		return
	}
	rsp, err := kvc.Get("world")
	if err != nil {
		t.Log(err)
		return
	}
	t.Logf("%#v \n", rsp)
}

func TestKv_Put(t *testing.T) {
	kvc, err := NewKVer(operaterAddr, vaultAddr, "ankr-sms")
	if err != nil {
		t.Log(err)
		return
	}
	if err = kvc.Put("hello", map[string]string{"hello": "world"}); err != nil {
		t.Logf("error: %v\n", err)
		return
	}
	t.Log("work over")
}
