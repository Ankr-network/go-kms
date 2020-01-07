package database_test

import (
	"fmt"
	"log"

	"github.com/Ankr-network/go-kms/database"
)

func ExampleClient_Get() {
	c := database.NewClient("127.0.0.1")
	username, password, err := c.Get("ankr-user")
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("username: %s password: %s \n", username, password)
	// output:
	// username: v-approle-ankr-sms-r-P3OkFNmSZJr password: A1a-wfOtfVAmex12y9Su
}
