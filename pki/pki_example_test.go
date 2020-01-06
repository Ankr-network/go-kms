package pki_test

import "github.com/Ankr-network/go-kms/pki"

func ExampleClient_Request() {
	cc, err := pki.NewPkiClient("192.168.39.113:30401", "ankr-pki")
	if err != nil {
		// handle with error
		return
	}
	if rsp, err := cc.Request(&pki.Config{Ttl: "24h", CommonName: "test.ankr.com"}); err != nil {
		// handle with error
		return
	} else {
		_ = rsp
	}
	// output:
}

func ExampleClient_Revoke() {
	cc, err := pki.NewPkiClient("192.168.39.113:30401", "ankr-pki")
	if err != nil {
		// handle with error
		return
	}
	if err := cc.Revoke("4a:ad:19:51:7a:47:5e:13:83:26:13:f2:43:25:ea:4a:36:b7:ae:9f"); err != nil {
		// handle with error
		return
	}
	// output:
}
