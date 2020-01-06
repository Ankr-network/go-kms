package pki

import "testing"

func TestClient_Request(t *testing.T) {
	cc, err := NewPkiClient("192.168.39.113:30401", "ankr-pki")
	if err != nil {
		t.Error(err)
		return
	}
	if rsp, err := cc.Request(&Config{Ttl: "24h", CommonName: "test.ankr.com"}); err != nil {
		t.Error(err)
		return
	} else {
		t.Logf("output: %+v\n", rsp)
	}
}

func TestClient_Revoke(t *testing.T) {
	cc, err := NewPkiClient("192.168.39.113:30401", "ankr-pki")
	if err != nil {
		t.Error(err)
		return
	}
	if err := cc.Revoke("4a:ad:19:51:7a:47:5e:13:83:26:13:f2:43:25:ea:4a:36:b7:ae:9f"); err != nil {
		t.Error(err)
		return
	}
}
