package database

import (
	"fmt"
	"log"
)

const (
	kmsAddr = "192.168.39.113:30401"
	//kmsAddr = "127.0.0.1:6900"
)

func Example_GetSecrets() {
	c := NewClient(kmsAddr)
	username, password, err := c.Get("ankr-rbac")
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("username: %s password: %s \n", username, password)
	// output:
	// username: v-approle-ankr-sms-r-P3OkFNmSZJr password: A1a-wfOtfVAmex12y9Su
}
