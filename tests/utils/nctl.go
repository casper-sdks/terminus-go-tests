package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/acarl005/stripansi"
	"github.com/antchfx/jsonquery"
	"github.com/make-software/casper-go-sdk/casper"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
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

func GetAccountHash(publicKey string, blockHash string) (string, error) {
	jsonStr, _ := GetStateAccountInfo(publicKey, blockHash)
	return GetByJsonPath(jsonStr, "/result/account/account_hash")
}

func GetByJsonPath(jsonStr string, path string) (string, error) {
	doc, err := jsonquery.Parse(strings.NewReader(jsonStr))
	value := jsonquery.FindOne(doc, path).Value()
	return fmt.Sprintf("%v", value), err
}

func GetStateAccountInfo(publicKey string, blockHash string) (string, error) {

	params := fmt.Sprintf("{\"public_key\":\"%s\",\"block_identifier\":{\"Hash\":\"%s\"}}", publicKey, blockHash)

	return simpleRcp("state_get_account_info", params)
}

func GetEraSummary(blockHash string) (string, error) {

	params := fmt.Sprintf("[{\"Hash\":\"%s\"}]}", blockHash)

	return simpleRcp("chain_get_era_summary", params)
}

func simpleRcp(method string, params string) (string, error) {
	var nctlJson string
	id := time.Now().UnixMilli()
	payload := fmt.Sprintf(`{"id": %d, "jsonrpc":"2.0","method":"%s","params":%s}`, id, method, params)
	bufferString := bytes.NewBufferString(payload)

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%v:%v/rpc", config["host-name"], config["port-rcp"]), bufferString)
	request.Header.Add("Content-Type", "application/json")

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	var response *http.Response

	response, err = client.Do(request)

	if err == nil {
		if response == nil && response.StatusCode != http.StatusOK {
			err = fmt.Errorf("invalid response %d", response.StatusCode)
		} else {
			bodyBytes, err := io.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			nctlJson = string(bodyBytes)
		}
	}
	return nctlJson, err
}

func nctlExec(command string, params string) (string, error) {

	docker := fmt.Sprintf("%v", config["docker-name"])
	cmd := fmt.Sprintf("docker exec  -t %s /bin/bash -c 'source /home/casper/casper-node/utils/nctl/sh/views/%s %s'", docker, command, params)

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
