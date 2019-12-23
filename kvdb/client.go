package kvdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Ankr-network/go-kms/approle"
)

type KVer interface {
	Get(path string) (map[string]interface{}, error)
	Put(path string, value map[string]string) error
}

type kv struct {
	token     string
	appRole   string
	oAddr     string
	vAddr     string
	headToken string
	hc        http.Client
}

func NewKVer(operatorAddr, vaultAddr, appRole string) (KVer, error) {
	if strings.Contains(appRole, "/") {
		return nil, errors.New("role name can't contain the char /")
	}
	if !strings.HasPrefix(vaultAddr, "http") {
		vaultAddr = fmt.Sprintf("http://%s/v1/kv/%s", vaultAddr, appRole)
	} else {
		vaultAddr = fmt.Sprintf("%s/v1/kv/%s", vaultAddr, appRole)
	}

	k := &kv{
		appRole:   appRole,
		oAddr:     operatorAddr,
		vAddr:     vaultAddr,
		headToken: "X-Vault-Token",
		hc:        http.Client{},
	}
	if err := k.init(); err != nil {
		return nil, err
	}
	return k, nil
}

func (k *kv) init() error {
	ac := approle.NewClient(k.oAddr)
	token, err := ac.Get(k.appRole)
	if err != nil {
		return err
	}
	k.token = token
	return nil
}

func (k *kv) Get(path string) (map[string]interface{}, error) {
	if path != "/" {
		path = fmt.Sprintf("%s/%s", k.vAddr, path)
	} else {
		path = k.vAddr
	}
	println("request path: ", path)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set(k.headToken, k.token)
	rsp, err := k.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	var kvRsp map[string]interface{}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &kvRsp); err != nil {
		return nil, err
	}
	if kvRsp["data"] != nil {
		return kvRsp["data"].(map[string]interface{}), nil
	} else {
		return nil, errors.New(fmt.Sprintf("%+v", kvRsp))
	}
}

func (k *kv) Put(path string, value map[string]string) error {
	if path != "/" {
		path = fmt.Sprintf("%s/%s", k.vAddr, path)
	}
	body, err := json.Marshal(&value)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set(k.headToken, k.token)
	rsp, err := k.hc.Do(req)
	if err != nil {
		return err
	}
	if rsp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("status code: %d", rsp.StatusCode)
	}
	return nil
}
