package utils

import (
	"fmt"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/sse"
	"net/http"
)

func GetSdk() casper.RPCClient {

	//goland:noinspection HttpUrlsUsage
	return casper.NewRPCClient(casper.NewRPCHandler(
		fmt.Sprintf("http://%v:%v/rpc", config["host-name"], config["port-rcp"]),
		http.DefaultClient),
	)
}

func GetSse() *sse.Client {
	//goland:noinspection HttpUrlsUsage
	return sse.NewClient(fmt.Sprintf("http://%v:%v/events/main", config["host-name"], config["port-sse"]))
}
