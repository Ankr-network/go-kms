package approle

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type AppRole interface {
	// get token by role name
	// argument 1: token
	Get(roleName string) (string, error)
}

type client struct {
	addr string
	c    *http.Client
}

// NewClient create new client by remote kms server operator address
// address format as "127.0.0.1:8080"
func NewClient(addr string) AppRole {
	if !strings.HasPrefix(addr, "http") {
		addr = fmt.Sprintf("http://%s/role/", addr)
	} else {
		addr = addr + "/role/"
	}
	return &client{addr: addr, c: &http.Client{}}
}

func (c *client) Get(roleName string) (string, error) {
	rsp, err := http.Get(c.addr + roleName)
	if err != nil {
		return "", nil
	}
	defer func() {
		err = rsp.Body.Close()
	}()
	token, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	return string(token), err
}
