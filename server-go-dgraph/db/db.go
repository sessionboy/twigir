package db

import (
	"server/db/dgraph"

	dgo "github.com/dgraph-io/dgo/v210"
)

var (
	Dgraph  *dgo.Dgraph
	CloseDb dgraph.CancelFunc
)

func InitDb() {
	dg, close, err := dgraph.NewClient()
	if err != nil {
		println("db error:", err)
		return
	}
	Dgraph = dg
	CloseDb = close
}
