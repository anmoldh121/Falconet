package models

import (
	"net"
	"sync"
)

type Payload struct {
	addr *net.UDPAddr
}

type Genesis struct {
	conn  *net.UDPConn
	conns map[string]string
	send  chan *Payload
	exit  chan bool
	wg    *sync.WaitGroup
}

type Peer struct {
	ID         string
	Username   string
	Endpoint   string
	PublicKey  string
	PrivateKey string
}
