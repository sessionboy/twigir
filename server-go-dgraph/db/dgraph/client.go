package dgraph

import (
	"os"

	dgo "github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type CancelFunc func()

func NewClient() (*dgo.Dgraph, CancelFunc, error) {
	addr := os.Getenv("DGRAPH_GRPC")
	log.Infof("addr: %s", addr)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
		return nil, nil, err
	}

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	return dg, func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error while closing connection:%v", err)
		}
	}, nil
}
