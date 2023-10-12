package steps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/make-software/casper-go-sdk/types/keypair"
	"math/big"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/clvalue"

	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
)

/**
 * The test features implementation for the speculative_execution.feature
 */
func TestFeaturesSpeculativeExecution(t *testing.T) {
	utils.TestFeatures(t, "speculative_execution.feature", InitializeSpeculativeExecution)
}

func InitializeSpeculativeExecution(ctx *godog.ScenarioContext) {
	var speculativeExecClient *rpc.SpeculativeClient
	var casperClient casper.RPCClient
	var speculativeExecResult rpc.SpeculativeExecResult
	var speculativeDeploy casper.Deploy

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		speculativeExecClient = utils.GetSpeculativeClient()
		casperClient = utils.GetRPCClient()
		return ctx, nil
	})

	ctx.Step(`that the "faucet" account transfers (\d+) to user-(\d+) account with a gas payment amount of (\d+) using the speculative_exec RPC API`,
		func(transferAmount int64, userId int, paymentAmount int64) error {
			var err error

			if speculativeExecClient == nil {
				return errors.New("unable to create speculative client")
			}

			speculativeDeploy, err = createDeploy()
			if err == nil {
				speculativeExecResult, err = speculativeExecClient.SpeculativeExec(context.Background(), speculativeDeploy, nil)
			}
			return err
		})

	ctx.Step(`^a valid speculative_exec_result will be returned with (\d+) transforms$`, func(transformCount int) error {
		if len(speculativeExecResult.DeployHash.String()) == 0 {
			return errors.New("missing speculativeExecResult")
		}

		return utils.ExpectEqual(utils.CasperT, "transforms", len(speculativeExecResult.ExecutionResult.Success.Effect.Transforms), transformCount)
	})

	ctx.Step(`^the speculative_exec has an api_version of "([^"]*)"`, func(apiVersion string) error {
		return utils.ExpectEqual(utils.CasperT, "api_version", speculativeExecResult.ApiVersion, apiVersion)
	})

	ctx.Step(`^the speculative_exec has a valid block_hash$`, func() error {
		// FIXME BlockHash is missing from the RPC SpeculativeExecResult
		//return utils.ExpectEqual(utils.CasperT, "block_hash", len(speculativeExecResult.BlockHash.Bytes()), 32)
		return utils.NotImplementError
	})

	ctx.Step(`^the execution_results contains a cost of (\d+)$`, func(cost int) error {
		return utils.ExpectEqual(utils.CasperT, "cost", speculativeExecResult.ExecutionResult.Success.Cost, uint64(cost))
	})

	ctx.Step(`^the speculative_exec has a valid execution_result$`, func() error {

		transfer := speculativeExecResult.ExecutionResult.Success.Transfers[0]
		transform, err := getTransform(speculativeExecResult, transfer.ToPrefixedString())
		if err == nil {
			err = utils.ExpectEqual(utils.CasperT, "key", transform.Key.Transfer.ToPrefixedString(), transfer.ToPrefixedString())
		}

		if err == nil {
			err = utils.ExpectEqual(utils.CasperT, "transform.WriteTransfer", transform.Transform.IsWriteTransfer(), true)
		}

		return err
	})

	ctx.Step(`^the speculative_exec execution_result transform wth the transfer key contains the deploy_hash$`, func() error {

		transfer := speculativeExecResult.ExecutionResult.Success.Transfers[0]
		transform, err := getTransform(speculativeExecResult, transfer.ToPrefixedString())
		var writeTransfer *types.WriteTransfer

		if err == nil {
			writeTransfer, err = transform.Transform.ParseAsWriteTransfer()
			err = utils.ExpectEqual(utils.CasperT, "WriteTransfer.deploy_hash", writeTransfer.DeployHash.String(), speculativeDeploy.Hash.String())
		}

		if err == nil {
			actual := writeTransfer.To.String()
			userOneKey, _ := casper.NewED25519PrivateKeyFromPEMFile("../../assets/net-1/user-1/secret_key.pem")
			expected := userOneKey.PublicKey().AccountHash().String()
			err = utils.ExpectEqual(utils.CasperT, "WriteTransfer.To", actual, expected)
		}

		if err == nil {
			actual := writeTransfer.From.String()
			faucetKey, _ := casper.NewED25519PrivateKeyFromPEMFile("../../assets/net-1/faucet/secret_key.pem")
			expected := faucetKey.PublicKey().AccountHash().String()
			err = utils.ExpectEqual(utils.CasperT, "WriteTransfer.from", actual, expected)
		}

		return err
	})

	ctx.Step(`^the speculative_exec execution_result transform with the transfer key has the amount of (\d+)`, func(amount int64) error {

		transfer := speculativeExecResult.ExecutionResult.Success.Transfers[0]
		transform, err := getTransform(speculativeExecResult, transfer.ToPrefixedString())
		var writeTransfer *types.WriteTransfer

		if err == nil {
			writeTransfer, err = transform.Transform.ParseAsWriteTransfer()
		}

		if err == nil {
			// FIXME Should be big.Int
			err = utils.ExpectEqual(utils.CasperT, "WriteTransfer.amount", writeTransfer.Amount, uint64(amount))
		}

		return err
	})

	ctx.Step(`^the speculative_exec execution_result transform with the transfer key has the "([^"]*)" field set to the "([^"]*)" account hash`, func(fieldName string, accountId string) error {

		transfer := speculativeExecResult.ExecutionResult.Success.Transfers[0]
		transform, err := getTransform(speculativeExecResult, transfer.ToPrefixedString())
		var writeTransfer *types.WriteTransfer

		if err == nil {
			writeTransfer, err = transform.Transform.ParseAsWriteTransfer()
		}

		if err == nil {
			var accountHash = getAccountHash(accountId)
			var actual string
			if fieldName == "from" {
				actual = writeTransfer.From.String()
			} else {
				actual = writeTransfer.To.String()
			}

			err = utils.ExpectEqual(utils.CasperT, "WriteTransfer."+fieldName, actual, accountHash)
		}
		return err
	})

	ctx.Step(`^the speculative_exec execution_result transform with the transfer key has the "([^"]*)" field set to the purse uref of the "([^"]*)" account`,
		func(fieldName string, accountId string) error {

			transfer := speculativeExecResult.ExecutionResult.Success.Transfers[0]
			transform, err := getTransform(speculativeExecResult, transfer.ToPrefixedString())
			var writeTransfer *types.WriteTransfer

			if err == nil {
				writeTransfer, err = transform.Transform.ParseAsWriteTransfer()
			}

			if err == nil {
				accountInfo, _ := getAccountInfo(casperClient, accountId)
				var actual string
				if fieldName == "source" {
					actual = writeTransfer.Source.String()
				} else {
					actual = writeTransfer.Target.String()
				}

				expected := accountInfo.Account.MainPurse.String()
				err = utils.ExpectEqual(utils.CasperT, "WriteTransfer."+fieldName, strings.Split(actual, "-")[0], strings.Split(expected, "-")[0])
				err = utils.ExpectEqual(utils.CasperT, "WriteTransfer."+fieldName, strings.Split(actual, "-")[1], strings.Split(expected, "-")[1])
			}
			return err
		})

	ctx.Step(`^the speculative_exec execution_result contains a valid deploy transform$`, func() error {
		transform, err := getTransform(speculativeExecResult, "deploy-"+speculativeDeploy.Hash.String())

		if err == nil {
			t := transform.Transform

			// FIXME Fails should have a method t. isWriteDeployInfo()
			if !t.IsWriteCLValue() {
				return errors.New("should have a method t. isWriteDeployInfo()")
			}
		}

		return nil
	})

	ctx.Step(`^the speculative_exec execution_result contains a valid AddUInt512 transform with a value of (\d+)$`, func(value int64) error {
		return nil
	})
}

func createDeploy() (casper.Deploy, error) {
	keyPath := utils.GetUserKeyAssetPath(1, 1, "secret_key.pem")
	receiverPrivateKey, err := casper.NewED25519PrivateKeyFromPEMFile(keyPath)
	if err != nil {
		return types.Deploy{}, err
	}

	var faucetKey keypair.PrivateKey
	faucetKey, err = casper.NewED25519PrivateKeyFromPEMFile("../../assets/net-1/faucet/secret_key.pem")
	if err != nil {
		return types.Deploy{}, err
	}

	header := types.DefaultHeader()
	header.ChainName = "casper-net-1"
	header.Account = faucetKey.PublicKey()
	header.Timestamp = types.Timestamp(time.Now())
	payment := types.StandardPayment(big.NewInt(100000000))
	transferAmount := big.NewInt(2500000000)

	args := &types.Args{}
	args.AddArgument("amount", *clvalue.NewCLUInt512(transferAmount))
	args.AddArgument("target", clvalue.NewCLPublicKey(receiverPrivateKey.PublicKey()))
	args.AddArgument("id", clvalue.NewCLOption(*clvalue.NewCLUInt64(rand.Uint64())))

	session := types.ExecutableDeployItem{
		Transfer: &types.TransferDeployItem{
			Args: *args,
		},
	}

	deploy, err := types.MakeDeploy(header, payment, session)

	if err == nil {
		err = deploy.SignDeploy(faucetKey)
	}

	var deployJson []byte

	deployJson, err = json.Marshal(deploy)
	if err != nil {
		return *deploy, err
	}

	fmt.Println(string(deployJson))

	return *deploy, err
}

func getTransform(speculativeExecResult rpc.SpeculativeExecResult, key string) (types.TransformKey, error) {
	transforms := speculativeExecResult.ExecutionResult.Success.Effect.Transforms
	for _, transform := range transforms {
		if transform.Key.String() == key {
			return transform, nil
		}
	}
	return types.TransformKey{}, fmt.Errorf("unable transform to find with key %s", key)

}

func getPrivateKey(accountId string) (keypair.PrivateKey, error) {
	if "faucet" == accountId {
		return casper.NewED25519PrivateKeyFromPEMFile("../../assets/net-1/faucet/secret_key.pem")
	} else {
		return casper.NewED25519PrivateKeyFromPEMFile("../../assets/net-1/user-1/secret_key.pem")
	}
}

func getAccountHash(accountId string) string {
	key, err := getPrivateKey(accountId)

	if err == nil {
		return key.PublicKey().AccountHash().String()
	} else {
		return ""
	}
}

func getAccountInfo(casperClient casper.RPCClient, accountId string) (rpc.StateGetAccountInfo, error) {
	key, err := getPrivateKey(accountId)

	if err == nil {
		latest, err := casperClient.GetBlockLatest(context.Background())
		if err == nil {
			return casperClient.GetAccountInfoByBlochHash(context.Background(), latest.Block.Hash.String(), key.PublicKey())
		}
	}

	return rpc.StateGetAccountInfo{}, err
}
