package approle

import (
	"fmt"
)

const (
	testAddr = "192.168.1.93:30386"
)

func Example_GetAppToken() {
	c := NewClient(testAddr)
	token, err := c.Get("ankr-sms-role")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("token: %s \n", token)
	// output:
	// token: s.ZIvORIjCInjYZIRoDbV0QgtU
}
