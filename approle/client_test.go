package approle

import (
	"testing"
)

const (
	testAddr = "192.168.1.93:30386"
)

func TestClientGet(t *testing.T) {
	c := NewClient(testAddr)
	token, err := c.Get("ankr-sms-role")
	if err != nil {
		t.Log(err)
		return
	}
	t.Logf("token: %s \n", token)
}
