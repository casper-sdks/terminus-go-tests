package steps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/make-software/casper-go-sdk/sse"
	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/keypair"
	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
	"github.com/stretchr/testify/assert"
	"log"
	"math/big"
	"math/rand"
	"strings"
	"testing"
	"time"
)

// The test features implementation for the query_global_state.feature
func TestFeaturesQueryGlobalState(t *testing.T) {
	utils.TestFeatures(t, "query_global_state.feature", InitializeQueryGlobalState)
}

func InitializeQueryGlobalState(ctx *godog.ScenarioContext) {

	var sdk casper.RPCClient
	var deployResult rpc.PutDeployResult
	var lastBlockAdded sse.BlockAddedEvent
	var globalState rpc.QueryGlobalStateResult
	var stateRootHash rpc.ChainGetStateRootHashResult
	var expectedErr rpc.RpcError

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		sdk = utils.GetSdk()
		return ctx, nil
	})

	ctx.Step(`^that a valid block hash is known$`, func() error {

		err := utils.Pass

		deployResult, err = createTransfer(sdk)

		if err == nil {
			lastBlockAdded, err = utils.WaitForBlockAdded(deployResult.DeployHash.String(), 300)
		}

		return err
	})

	ctx.Step(`^the query_global_state RCP method is invoked with the block hash as the query identifier$`, func() error {
		var err error
		key := "deploy-" + deployResult.DeployHash.String()
		globalState, err = sdk.QueryGlobalStateByBlockHash(context.Background(), lastBlockAdded.BlockAdded.BlockHash, key, nil)
		return err
	})

	ctx.Step(`a valid query_global_state_result is returned$`, func() error {

		err := utils.ExpectEqual(utils.CasperT, "ApiVersion", globalState.ApiVersion, "1.0.0")

		if err == nil && len(globalState.MerkleProof) == 0 {
			err = fmt.Errorf("invalid merkle proof")
		}

		if err == nil && len(globalState.MerkleProof) == 0 {
			err = fmt.Errorf("invalid merkle proof")
		}

		if err == nil && globalState.BlockHeader.Timestamp.ToTime().UnixMilli() < 0 {
			err = fmt.Errorf("invalid timestamp")
		}

		if err == nil && globalState.BlockHeader.EraID < 1 {
			err = fmt.Errorf("invalid EraId")
		}

		if err == nil && len(globalState.BlockHeader.AccumulatedSeed.Bytes()) < 1 {
			err = fmt.Errorf("invalid AccumulatedSeed")
		}

		if err == nil && len(globalState.BlockHeader.BodyHash.String()) < 32 {
			err = fmt.Errorf("invalid BodyHash")
		}

		if err == nil && len(globalState.BlockHeader.ParentHash.String()) < 32 {
			err = fmt.Errorf("invalid ParentHash")
		}

		return err
	})

	ctx.Step(`the query_global_state_result contains a valid deploy info stored value$`, func() error {

		if globalState.StoredValue.DeployInfo == nil {
			return fmt.Errorf("missing value in global state")
		}
		return utils.Pass
	})

	ctx.Step(`the query_global_state_result's stored value from is the user-1 account hash$`, func() error {
		if globalState.StoredValue.DeployInfo == nil {
			return fmt.Errorf("missing value in global state")
		}

		return utils.ExpectEqual(utils.CasperT, "DeployHash", globalState.StoredValue.DeployInfo.DeployHash.String(), deployResult.DeployHash.String())
	})

	ctx.Step(`the query_global_state_result's stored value contains a gas price of (\d+)$`, func(gasPrice int64) error {
		return utils.ExpectEqual(utils.CasperT, "Gas", globalState.StoredValue.DeployInfo.Gas, uint64(gasPrice))
	})

	ctx.Step(`the query_global_state_result stored value contains the transfer hash$`, func() error {
		if !strings.Contains(globalState.StoredValue.DeployInfo.Transfers[0].ToPrefixedString(), "transfer-") {
			return fmt.Errorf("invalid transfer")
		} else {
			return utils.Pass
		}
	})

	ctx.Step(`the query_global_state_result stored value contains the transfer source uref$`, func() error {
		if !strings.Contains(globalState.StoredValue.DeployInfo.Source.ToPrefixedString(), "uref-") {
			return fmt.Errorf("invalid transfer")
		} else {
			return utils.Pass
		}
	})

	ctx.Step(`that the state root hash is known$`, func() error {
		var err error
		stateRootHash, err = sdk.GetStateRootHashLatest(context.Background())
		return err
	})

	ctx.Step(`the query_global_state RCP method is invoked with the state root hash as the query identifier and an invalid key$`, func() error {
		key := "uref-d0343bb766946f9f850a67765aae267044fa79a6cd50235ffff248a37534"
		srh := stateRootHash.StateRootHash.String()
		_, err := sdk.QueryGlobalStateByStateHash(context.Background(), &srh, key, nil)

		if err == nil {
			return fmt.Errorf("Should have error ")
		}

		expectedErr = getRpcError(errors.Unwrap(err))

		return utils.Pass
	})

	ctx.Step(`an error code of -(\d+) is returned$`, func(errorCode int) error {
		if expectedErr.Code != -1*errorCode && expectedErr.Code != -32001 {
			return fmt.Errorf("invalid error code %d", expectedErr.Code)
		}

		return utils.Pass
	})

	ctx.Step(`an error message of "([^"]*)" is returned$`, func(errorMessage string) error {
		if expectedErr.Data != errorMessage && expectedErr.Message != "No such block" {
			return fmt.Errorf("invalid error message %s", expectedErr.Data)
		}
		return utils.Pass
	})

	ctx.Step(`the query_global_state_result stored value contains the transfer source uref$`, func() error {
		return utils.Pass
	})

	ctx.Step(`the query_global_state RCP method is invoked with an invalid block hash as the query identifier$`, func() error {

		blockHash := "06e04c9e3b8b084d169b4908ac68a797374d325cbe919e01d290d6b7f5c720d0"

		key := "deploy-" + deployResult.DeployHash.String()

		_, err := sdk.QueryGlobalStateByBlockHash(context.Background(), blockHash, key, nil)

		if err == nil {
			return fmt.Errorf("Should have error ")
		}

		expectedErr = getRpcError(errors.Unwrap(err))

		return utils.Pass
	})

}

func getRpcError(err interface{}) rpc.RpcError {
	rpcError := err.(*rpc.RpcError)
	fmt.Println(rpcError.Code)
	return *rpcError
}

func createTransfer(sdk casper.RPCClient) (rpc.PutDeployResult, error) {

	var err error = nil
	var amount = big.NewInt(2500000000)
	var deployJson []byte
	var senderKey keypair.PrivateKey
	var receiverKey keypair.PrivateKey
	var deploy *types.Deploy

	keyPath := utils.GetUserKeyAssetPath(1, 1, "secret_key.pem")
	senderKey, err = casper.NewED25519PrivateKeyFromPEMFile(keyPath)

	if err != nil {
		log.Fatal(err)
	}

	keyPath = utils.GetUserKeyAssetPath(1, 2, "secret_key.pem")
	receiverKey, err = casper.NewED25519PrivateKeyFromPEMFile(keyPath)
	if err != nil {

		log.Fatal(err)
	}

	header := types.DefaultHeader()
	header.ChainName = "casper-net-1"
	header.Account = senderKey.PublicKey()
	header.Timestamp = types.Timestamp(time.Now())
	stdPayment := types.StandardPayment(big.NewInt(100000000))

	args := &types.Args{}
	args.AddArgument("amount", *clvalue.NewCLUInt512(amount))
	args.AddArgument("target", clvalue.NewCLPublicKey(receiverKey.PublicKey()))
	args.AddArgument("id", clvalue.NewCLOption(*clvalue.NewCLUInt64(rand.Uint64())))

	session := types.ExecutableDeployItem{
		Transfer: &types.TransferDeployItem{
			Args: *args,
		},
	}

	deploy, err = types.MakeDeploy(header, stdPayment, session)

	if err == nil {
		assert.NotNil(utils.CasperT, deploy, "deploy")
		err = deploy.SignDeploy(senderKey)
	}

	if err == nil {
		deployJson, err = json.Marshal(deploy)
		assert.NotNil(utils.CasperT, deployJson)
		fmt.Println(string(deployJson))
	}

	return sdk.PutDeploy(context.Background(), *deploy)
}
