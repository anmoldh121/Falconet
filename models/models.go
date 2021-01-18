package models

type Payload struct {
	data []byte
}

type Peer struct {
	ID         string
	Username   string
	Endpoint   string
	PublicKey  string
	PrivateKey string
}
