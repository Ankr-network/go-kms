package kmsdb

import (
	"fmt"
	"log"
)

const (
	operatorAddr = "192.168.1.93:30386"
	vaultAddr    = "192.168.1.93:30050"
)

func Example_GetSecrets() {
	c := NewClient(operatorAddr, vaultAddr)
	username, password, err := c.Get("ankr-sms-role")
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("username: %s password: %s \n", username, password)
	// output:
	// username: v-approle-ankr-sms-r-P3OkFNmSZJr password: A1a-wfOtfVAmex12y9Su
}
