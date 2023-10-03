package utils

import (
	"fmt"
	"net/http"

	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/make-software/casper-go-sdk/sse"
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

func GetSpeculativeClient() *rpc.SpeculativeClient {
	//goland:noinspection HttpUrlsUsage
	return rpc.NewSpeculativeClient(casper.NewRPCHandler(
		fmt.Sprintf("http://%v:%v/rpc", config["host-name"], config["port-spd"]), http.DefaultClient),
	)
}
