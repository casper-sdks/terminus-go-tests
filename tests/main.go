package tests

import (
	"fmt"
	"github.com/make-software/casper-go-sdk/casper"
	yml "gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	root       = filepath.Join(filepath.Dir(b), "..")
	config     map[string]interface{}
)

func main() {
}

func sdk() casper.RPCClient {

	host := fmt.Sprintf("%v", config["host-name"])
	port := fmt.Sprintf("%v", config["port-rcp"])

	var handler = casper.NewRPCHandler("http://"+host+":"+port+"/rpc", http.DefaultClient)
	return casper.NewRPCClient(handler)

}

func readConfig() {
	f, err := os.ReadFile(root + "/config.yml")
	if err != nil {
		log.Fatal(err)
	}
	err = yml.Unmarshal(f, &config)
	if err != nil {
		log.Fatal(err)
	}
}
