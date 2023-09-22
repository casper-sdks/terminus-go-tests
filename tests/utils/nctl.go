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
	"math/big"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// Steps for the state_get_auction_info.feature
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

func GetStateRootHash(nodeId int) (string, error) {
	result, err := nctlExec("view_chain_state_root_hash.sh", fmt.Sprintf("node=%d", nodeId))
	srh := strings.Split(result, "=")[1]
	srh = strings.TrimSpace(srh)
	return srh, err

}

func GetAccountHash(publicKey string, blockHash string) (string, error) {
	jsonStr, _ := GetStateAccountInfo(publicKey, blockHash)
	return GetByJsonPath(jsonStr, "/result/account/account_hash")
}

func StateGetBalance(stateRootHash string, purseUref string) (big.Int, error) {
	var balance = new(big.Int)

	params := fmt.Sprintf("{\"state_root_hash\":\"%s\",\"purse_uref\":\"%s\"}", stateRootHash, purseUref)
	jsonStr, _ := simpleRcp("state_get_balance", params)

	balanceStr, err := GetByJsonPath(jsonStr, "/result/balance_value")

	if err == nil {
		balance.SetString(balanceStr, 10)
	}

	return *balance, err
}

func GetByJsonPath(jsonStr string, path string) (string, error) {
	node, err := GetNodeByJsonPath(jsonStr, path)
	if err == nil {
		return fmt.Sprintf("%v", node.Value()), nil
	} else {
		return "", err
	}
}

func GetNodeByJsonPath(jsonStr string, path string) (*jsonquery.Node, error) {
	doc, err := jsonquery.Parse(strings.NewReader(jsonStr))
	var node *jsonquery.Node
	if err == nil {
		node = jsonquery.FindOne(doc, path)
	}
	return node, err
}

func GetStateAccountInfo(publicKey string, blockHash string) (string, error) {

	params := fmt.Sprintf("{\"public_key\":\"%s\",\"block_identifier\":{\"Hash\":\"%s\"}}", publicKey, blockHash)

	return simpleRcp("state_get_account_info", params)
}

func GetEraSummary(blockHash string) (string, error) {

	params := fmt.Sprintf("[{\"Hash\":\"%s\"}]", blockHash)

	return simpleRcp("chain_get_era_summary", params)
}

func GetAuctionInfoByHash(hash string) (string, error) {
	auctionInfoJson, err := simpleRcp("state_get_auction_info", fmt.Sprintf("[{\"Hash\": \"%s\"}]", hash))
	return auctionInfoJson, err
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
