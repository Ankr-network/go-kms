package pki

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Ankr-network/go-kms/approle"
)

type Client struct {
	token   string
	appRole string
	oAddr   string
	vAddr   string
	hc      *http.Client
	encBuf  *bytes.Buffer
	encoder *json.Encoder
}

const headToken = "X-Vault-Token"

func NewPkiClient(kmsAddr, appRole string) (*Client, error) {
	if strings.Contains(appRole, "/") {
		return nil, errors.New("role name can't contain the char /")
	}

	var (
		oAddr string
		vAddr string
	)
	if !strings.HasPrefix(kmsAddr, "http") {
		vAddr = fmt.Sprintf("http://%s/vault/v1/%s", kmsAddr, appRole)
		oAddr = fmt.Sprintf("http://%s/ops", kmsAddr)
	} else {
		vAddr = fmt.Sprintf("%s/vault/v1/%s", kmsAddr, appRole)
		oAddr = fmt.Sprintf("%s/ops", kmsAddr)

	}

	// get token by application role
	ac := approle.NewClient(oAddr)
	token, err := ac.Get(appRole)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buf)

	return &Client{
		token:   token,
		oAddr:   oAddr,
		vAddr:   vAddr,
		appRole: appRole,
		hc:      &http.Client{},
		encBuf:  buf,
		encoder: encoder,
	}, nil
}

type Config struct {
	CommonName string `json:"common_name"`
	// format: "24h"
	Ttl string `json:"ttl"`
}

type Response struct {
	PriKey string
	Pubkey string
	SN     string
}

const (
	notExist = -1
)

type KmsRsp struct {
	Auth interface{} `json:"auth"`
	Data struct {
		CaChain        []string `json:"ca_chain"`
		Certificate    string   `json:"certificate"`
		Expiration     int64    `json:"expiration"`
		IssuingCa      string   `json:"issuing_ca"`
		PrivateKey     string   `json:"private_key"`
		PrivateKeyType string   `json:"private_key_type"`
		SerialNumber   string   `json:"serial_number"`
	} `json:"data"`
	LeaseDuration int64       `json:"lease_duration"`
	LeaseID       string      `json:"lease_id"`
	Renewable     bool        `json:"renewable"`
	RequestID     string      `json:"request_id"`
	Warnings      interface{} `json:"warnings"`
	WrapInfo      interface{} `json:"wrap_info"`
}

type KmsError struct {
	Errors []string `json:"errors"`
}

func (c *Client) Request(cfg *Config) (*Response, error) {
	if cfg.CommonName == "" || cfg.Ttl == "" || strings.LastIndexByte(cfg.Ttl, 'h') == notExist {
		return nil, errors.New("params not valid")
	}

	c.encBuf.Reset()
	if err := c.encoder.Encode(&cfg); err != nil {
		return nil, err
	}

	localRequestAddr := fmt.Sprintf("%s/issue/%s", c.vAddr, c.appRole)
	req, err := http.NewRequest("PUT", localRequestAddr, bytes.NewReader(c.encBuf.Bytes()))
	if err != nil {
		return nil, err
	}
	req.Header.Set(headToken, c.token)
	rsp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != http.StatusOK {
		err := KmsError{}
		if ie := json.Unmarshal(body, &err); ie != nil {
			return nil, ie
		}
		return nil, errors.New(err.Errors[0])
	}

	bodyStruct := &KmsRsp{}
	if err := json.Unmarshal(body, &bodyStruct); err != nil {
		return nil, err
	}

	// get public key
	pub, err := getPublicKey(bodyStruct.Data.PrivateKey)
	if err != nil {
		return nil, err
	}
	return &Response{
		PriKey: bodyStruct.Data.PrivateKey,
		SN:     bodyStruct.Data.SerialNumber,
		Pubkey: pub,
	}, nil
}

func getPublicKey(privateKey string) (string, error) {
	block, _ := pem.Decode([]byte(privateKey))
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	publicKeyDer := x509.MarshalPKCS1PublicKey(&key.PublicKey)
	pubKeyBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   publicKeyDer,
	}
	return string(pem.EncodeToMemory(&pubKeyBlock)), nil
}

type RevokeRequest struct {
	SerialNumber string `json:"serial_number"`
}

func (c *Client) Revoke(serialNumber string) error {
	requstPath := fmt.Sprintf("%s/revoke", c.vAddr)
	c.encBuf.Reset()
	if err := c.encoder.Encode(&RevokeRequest{SerialNumber: serialNumber}); err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", requstPath, bytes.NewReader(c.encBuf.Bytes()))
	if err != nil {
		return err
	}
	req.Header.Set(headToken, c.token)
	rsp, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return err
		}
		re := KmsError{}
		if ie := json.Unmarshal(body, &re); ie != nil {
			return ie
		}
		return errors.New(re.Errors[0])
	}

	return nil
}
