package genesis

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/anmoldh121/falconet/models"
	"github.com/anmoldh121/falconet/proto"
	"google.golang.org/grpc/peer"
)

type Genesis struct {
	conns map[string]string
	send  chan *models.Payload
	exit  chan bool
	wg    *sync.WaitGroup
}

func (g *Genesis) Sender() {
	g.wg.Add(1)
	defer g.wg.Done()

	for {
		select {
		case <-g.exit:
			log.Printf("[Exiting] Sender")
			return
		case p := <-g.send:
			fmt.Println(p)
			return
		}
	}
}

func (g *Genesis) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	_, ok := g.conns[req.Username]
	if !ok {
		c, _ := peer.FromContext(ctx)
		fmt.Println(c.Addr.String())
		g.conns[req.Username] = c.Addr.String()
	}
	return &proto.RegisterResponse{
		Status: 201,
		Message: "SUCCESS",
	},nil
}

func (g *Genesis) GetTarget(ctx context.Context, req *proto.GetTargetRequest) (*proto.GetTargetResponse, error) {
	c, ok := g.conns[req.Username]
	if !ok {
		return &proto.GetTargetResponse{
			Endpoint: "",
			Username: "",
		},nil
	}
	return &proto.GetTargetResponse{
		Endpoint: c,
		Username: req.Username,
	},nil
}
