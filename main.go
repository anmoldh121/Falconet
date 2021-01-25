package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/anmoldh121/falconet/server/genesis"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	fmt.Println("Starting Server")
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://anmol:anmol@cluster0.cnhws.mongodb.net/"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	HandleError(err)
	db := client.Database("falconet")
	genesisServer, err := genesis.NewGenesis(tcpAddr, db)
	if err != nil {
		log.Fatal(err)
	}
	genesisServer.Listen()
}

func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
