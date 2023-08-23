package utils

import (
	"fmt"
	"github.com/make-software/casper-go-sdk/casper"
	"net/http"
)

func GetSdk() casper.RPCClient {

	return casper.NewRPCClient(casper.NewRPCHandler(
		fmt.Sprintf("http://%v:%v/rpc", config["host-name"], config["port-rcp"]),
		http.DefaultClient),
	)

}
