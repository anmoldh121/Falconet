package main

import (
	"github.com/anmoldh121/falconet/proto"
	"github.com/anmoldh121/falconet/server/genesis"

	"google.golang.org/grpc"
	"net"
	"log"
)

func main() {
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()
	proto.RegisterGenesisServer(server, &genesis.Genesis{})
	server.Serve(lis)
}
