package genesis

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/anmoldh121/falconet/models"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Genesis struct {
	conn  *net.TCPListener
	conns map[string]string
	send  chan *models.Payload
	exit  chan bool
	db    *mongo.Database
	wg    *sync.WaitGroup
}

type Message struct {
	Purpose int
	PeerId  string
}

type Response struct {
	addr string
}

func (g *Genesis) Receiver() {
	for {
		con, err := g.conn.AcceptTCP()
		if err != nil {
			continue
		}
		go g.handleConnection(con)
	}
}

func (g *Genesis) handleConnection(conn *net.TCPConn) {
	var buffer [4096]byte
	fmt.Println(conn.RemoteAddr())
	remoteAddr := conn.RemoteAddr()
	n, err := conn.Read(buffer[0:])
	if err != nil {
		log.Fatal(err)
	}
	mess := UnmarshalMessage(buffer[:n])
	if mess.Purpose == 1 {
		g.SavePeer(remoteAddr, mess)
	}
}

func UnmarshalMessage(req []byte) Message {
	var message Message
	err := json.Unmarshal(req, &message)
	if err != nil {
		log.Fatal(err)
	}
	return message
}

func (g *Genesis) SavePeer(addr net.Addr, message Message) {
	// collection := g.db.Collection("peers")
	// filter := bson.D{{"_id", message.PeerId}}
	// update := bson.D{{"endpoint"}}
	fmt.Println(message)
}

func (g *Genesis) Listen() {
	g.Receiver()
}

func NewGenesis(addr *net.TCPAddr, db *mongo.Database) (*Genesis, error) {
	c, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Genesis{
		conn:  c,
		conns: make(map[string]string),
		send:  make(chan *models.Payload),
		exit:  make(chan bool),
		db:    db,
		wg:    &sync.WaitGroup{},
	}, nil
}
