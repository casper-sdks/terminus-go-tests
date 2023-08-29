package utils

import (
	"encoding/json"
	"fmt"
	"github.com/acarl005/stripansi"
	"github.com/make-software/casper-go-sdk/casper"
	"log"
	"os/exec"
	"strings"
)

func GetNctlLatestBlock() (casper.Block, error) {

	block := casper.Block{}

	res, err := nctlExec("view_chain_block.sh", "")

	err = json.Unmarshal([]byte(res), &block)

	if err != nil {
		log.Fatal(err)
	}

	return block, err
}

func GetNodeStatus(nodeId int) (casper.InfoGetStatusResult, error) {

	res, err := nctlExec("view_node_status.sh", fmt.Sprintf("node=%d", nodeId))

	index := strings.Index(res, "{")
	jsonStr := res[index:]

	infoGetStatusResult := casper.InfoGetStatusResult{}

	if err := json.Unmarshal([]byte(jsonStr), &infoGetStatusResult); err != nil {
		log.Fatal(err)
	}

	return infoGetStatusResult, err
}

func nctlExec(command string, params string) (string, error) {

	docker := fmt.Sprintf("%v", config["docker-name"])
	cmd := fmt.Sprintf("docker exec  -t %s /bin/bash -c 'source casper-node/utils/nctl/sh/views/%s %s'", docker, command, params)

	strRes := ""

	res, err := exec.Command("/bin/sh", "-c", cmd).Output()

	if err != nil {
		log.Printf("Could not run command: %s", cmd)
		log.Fatal(err)
	} else {
		// Strip out ANSI control characters from response
		strRes = stripansi.Strip(string(res))
	}

	return strRes, err
}
