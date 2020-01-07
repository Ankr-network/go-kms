package pki

type Handler interface {
	// request private key and public key
	Request(cfg *Config) (*Response, error)
	// revoke certificates by serial number
	Revoke(serialNumber string) error
}
