package conn

import (
	"fmt"

	nebulaClient "github.com/vesoft-inc/nebula-go"
)

const (
	address  = "47.99.243.195"
	port     = 9669
	username = "user"
	password = "password"
)

var (
	Nsid    string
	Pool    *nebulaClient.ConnectionPool
	Session *nebulaClient.Session
)

func InitDb() {

	id, err := Connect(address, port, username, password)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	Nsid = id

	res, err := Execute("use app;")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Print("use space:", res)

}
