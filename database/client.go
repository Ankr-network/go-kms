package database

import (
	"net/http"

	"github.com/Ankr-network/go-kms/approle"
)

type Database interface {
	// get user name and password by role name
	// argument one: user name
	// argument two: password
	Get(roleName string) (string, string, error)
}

type client struct {
	operatorAddr string
	vaultAddr    string
	c            *http.Client
}

// NewClient create new client by address argument
// vaultAddr should be remote vault server addr which looks like "127.0.0.1:8200"
// operatorAddr should be remote operator server addr which looks like "127.0.0.1:8080"
func NewClient(operatorAddr, vaultAddr string) Database {
	return &client{
		operatorAddr: operatorAddr,
		vaultAddr:    vaultAddr,
		c:            &http.Client{},
	}
}

func (c *client) Get(roleName string) (string, string, error) {
	var (
		username string
		password string
	)
	// step 1: get token
	token, err := approle.NewClient(c.operatorAddr).Get(roleName)
	return username, password, nil

}
