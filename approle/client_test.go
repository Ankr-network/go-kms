package approle

import "testing"

func TestClient_Get(t *testing.T) {
	c := NewClient("http://192.168.39.113:30234/ops")
	token, err := c.Get("ankr-rbac")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("token: ", token)
}
