package utils

import (
	"encoding/json"
	"fmt"
	"github.com/acarl005/stripansi"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/stretchr/testify/assert"
	yml "gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	root       = filepath.Join(filepath.Dir(b), "../..")
	config     map[string]interface{}
)

func Sdk() casper.RPCClient {

	return casper.NewRPCClient(casper.NewRPCHandler(
		fmt.Sprintf("http://%v:%v/rpc",
			fmt.Sprintf("%v", config["host-name"]),
			fmt.Sprintf("%v", config["port-rcp"])),
		http.DefaultClient))

}

func ReadConfig() {
	f, err := os.ReadFile(root + "/config.yml")
	if err != nil {
		log.Fatal(err)
	}
	err = yml.Unmarshal(f, &config)
	if err != nil {
		log.Fatal(err)
	}
}

func GetNctlLatestBlock() (error, casper.Block) {

	docker := fmt.Sprintf("%v", config["docker-name"])
	block := casper.Block{}

	res, err := exec.Command("/bin/sh", "-c",
		"docker exec  -t "+docker+" /bin/bash -c 'source casper-node/utils/nctl/sh/views/view_chain_block.sh'",
		"| sed -e \"s/\\x1b\\[.\\{1,5\\}m//g\"").Output()
	if err != nil {
		fmt.Println("could not run command: ", err)
	}

	err = json.Unmarshal([]byte(stripansi.Strip(string(res))), &block)
	if err != nil {
		fmt.Println("could not unmarshal: ", err)
	}

	return err, block
}

func AssertExpectedAndActual(a expectedAndActualAssertion, expected, actual interface{}) error {
	var t asserter
	a(&t, expected, actual)
	return t.err
}

type asserter struct {
	err error
}

func (a *asserter) Errorf(format string, args ...interface{}) {
	a.err = fmt.Errorf(format, args...)
}

func Result(err error) error {
	if err != nil {
		return err
	} else {
		return nil
	}
}

type expectedAndActualAssertion func(t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool
