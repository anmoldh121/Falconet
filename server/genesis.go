package genesis

import (
	"fmt"
	"log"
	"sync"

	"github.com/anmoldh121/falconet/models"
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
