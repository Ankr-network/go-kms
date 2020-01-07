package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Ankr-network/go-kms/approle"
)

type Database interface {
	// get user name and password by role name
	// argument one: user name
	// argument two: password
	Get(roleName string) (string, string, error)
}

type Client struct {
	operatorAddr string
	vaultAddr    string
	c            *http.Client
}

// NewClient create new client by address argument
// kmsAddr should be remote kms server addr which looks like "127.0.0.1:8200"
func NewClient(kmsAddr string) Database {
	var (
		oAddr string
		vAddr string
	)
	if !strings.HasPrefix(kmsAddr, "http") {
		oAddr = fmt.Sprintf("http://%s/ops", kmsAddr)
		vAddr = fmt.Sprintf("http://%s/vault", kmsAddr)
	} else {
		oAddr = fmt.Sprintf("%s/ops", kmsAddr)
		vAddr = fmt.Sprintf("%s/vault", kmsAddr)
	}

	return &Client{
		operatorAddr: oAddr,
		vaultAddr:    vAddr,
		c:            &http.Client{},
	}
}

const (
	xVaultToken = "X-Vault-Token"
)

type DbRsp struct {
	Auth interface{} `json:"auth"`
	Data struct {
		Password string `json:"password"`
		Username string `json:"username"`
	} `json:"data"`
	LeaseDuration int64       `json:"lease_duration"`
	LeaseID       string      `json:"lease_id"`
	Renewable     bool        `json:"renewable"`
	RequestID     string      `json:"request_id"`
	Warnings      interface{} `json:"warnings"`
	WrapInfo      interface{} `json:"wrap_info"`
}

type ErrRsp struct {
	Errors []string `json:"errors"`
}

func (e ErrRsp) Error() string {
	return strings.Join(e.Errors, ",")
}

// Get generate database secret  by role name
// role name get from Ankr organization
func (c *Client) Get(roleName string) (string, string, error) {
	var (
		dbRsp  DbRsp
		errRsp ErrRsp
	)
	// step 1: get token
	token, err := approle.NewClient(c.operatorAddr).Get(roleName)
	if err != nil {
		return "", "", err
	}
	// step 2: get user name and password by token
	dbURL := fmt.Sprintf("%s/v1/database/creds/%s", c.vaultAddr, roleName)
	//fmt.Printf("token: %s url: %s \n", token, dbURL)
	req, err := http.NewRequest("GET", dbURL, nil)
	req.Header.Set(xVaultToken, token)
	rsp, err := c.c.Do(req)
	if err != nil {
		return "", "", err
	}
	defer func() {
		if err = rsp.Body.Close(); err != nil {
			println(err.Error())
		}
	}()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", "", err
	}
	// all the best
	if rsp.StatusCode == http.StatusOK {
		if err = json.Unmarshal(body, &dbRsp); err != nil {
			return "", "", err
		}
		return dbRsp.Data.Username, dbRsp.Data.Password, nil
	}
	// some error happened
	if err = json.Unmarshal(body, &errRsp); err != nil {
		return "", "", err
	}
	return "", "", errRsp
}
