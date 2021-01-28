package genesis

import (
	// "bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/anmoldh121/falconet/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	Purpose int    `json:"purpose,omitempty"`
	PeerId  string `json:"peerId, omitempty" `
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
	fmt.Println(conn.RemoteAddr())
	remoteAddr := conn.RemoteAddr()
	fmt.Println(remoteAddr.String())
	decoder := json.NewDecoder(conn)
	var message Message
	err := decoder.Decode(&message)
	fmt.Println(message, err)
	conn.Close()
	if message.Purpose == 1 {
		err := g.SavePeer(remoteAddr, message)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func UnmarshalMessage(req []byte) Message {
	var message Message
	err := json.Unmarshal(req, &message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(message)
	return message
}

func (g *Genesis) SavePeer(addr net.Addr, message Message) error {
	collection := g.db.Collection("peers")
	id, _ := primitive.ObjectIDFromHex(message.PeerId)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"endpoint", addr.String()}}}}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	updated, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	fmt.Println(updated.UpsertedID)
	return nil
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
