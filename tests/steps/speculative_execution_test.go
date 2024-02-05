package steps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/make-software/casper-go-sdk/types/keypair"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/clvalue"

	"github.com/casper-sdks/terminus-go-tests/tests/utils"
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
				speculativeExecResult, err = speculativeExecClient.SpeculativeExec(
					context.Background(),
					speculativeDeploy,
					nil)
			}
			return err
		})

	ctx.Step(`^a valid speculative_exec_result will be returned with (\d+) transforms$`, func(transformCount int) error {
		return utils.ExpectEqual(utils.CasperT, "transforms",
			len(speculativeExecResult.ExecutionResult.Success.Effect.Transforms),
			transformCount)
	})

	ctx.Step(`^the speculative_exec has an api_version of "([^"]*)"`, func(apiVersion string) error {
		return utils.ExpectEqual(utils.CasperT, "api_version", speculativeExecResult.ApiVersion, apiVersion)
	})

	ctx.Step(`^the speculative_exec has a valid block_hash$`, func() error {
		return utils.ExpectEqual(utils.CasperT, "block_hash", len(speculativeExecResult.BlockHash.Bytes()), 32)
	})

	ctx.Step(`^the execution_results contains a cost of (\d+)$`, func(cost int) error {
		return utils.ExpectEqual(utils.CasperT, "cost", speculativeExecResult.ExecutionResult.Success.Cost, uint64(cost))
	})

	ctx.Step(`^the speculative_exec has a valid execution_result$`, func() error {

		transfer := speculativeExecResult.ExecutionResult.Success.Transfers[0]
		transform, err := getTransform(speculativeExecResult, transfer.ToPrefixedString())
		if err == nil {
			err = utils.ExpectEqual(utils.CasperT,
				"key",
				transform.Key.Transfer.ToPrefixedString(),
				transfer.ToPrefixedString())
		}

		if err == nil {
			err = utils.ExpectEqual(utils.CasperT,
				"transform.WriteTransfer",
				transform.Transform.IsWriteTransfer(),
				true)
		}
		return err
	})

	ctx.Step(`^the speculative_exec execution_result transform wth the transfer key contains the deploy_hash$`,
		func() error {

			transfer := speculativeExecResult.ExecutionResult.Success.Transfers[0]
			transform, err := getTransform(speculativeExecResult, transfer.ToPrefixedString())
			var writeTransfer *types.WriteTransfer

			if err == nil {
				writeTransfer, err = transform.Transform.ParseAsWriteTransfer()
				err = utils.ExpectEqual(utils.CasperT,
					"WriteTransfer.deploy_hash",
					writeTransfer.DeployHash.String(),
					speculativeDeploy.Hash.String())
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

	ctx.Step(`^the speculative_exec execution_result transform with the transfer key has the amount of (\d+)`,
		func(amount int64) error {

			transfer := speculativeExecResult.ExecutionResult.Success.Transfers[0]
			transform, err := getTransform(speculativeExecResult, transfer.ToPrefixedString())
			var writeTransfer *types.WriteTransfer

			if err == nil {
				writeTransfer, err = transform.Transform.ParseAsWriteTransfer()
			}

			if err == nil {
				err = utils.ExpectEqual(utils.CasperT, "WriteTransfer.amount", writeTransfer.Amount.String(), strconv.FormatInt(amount, 10))
			}

			return err
		})

	ctx.Step(`^the speculative_exec execution_result transform with the transfer key has the "([^"]*)" field set to the "([^"]*)" account hash`,
		func(fieldName string, accountId string) error {

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
				err = utils.ExpectEqual(utils.CasperT,
					"WriteTransfer."+fieldName,
					strings.Split(actual, "-")[0],
					strings.Split(expected, "-")[0])
				err = utils.ExpectEqual(utils.CasperT,
					"WriteTransfer."+fieldName,
					strings.Split(actual, "-")[1],
					strings.Split(expected, "-")[1])
			}
			return err
		})

	ctx.Step(`the speculative_exec execution_result transform with the deploy key has the deploy_hash of the transfer's hash$`,
		func() error {
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

	ctx.Step(`the speculative_exec execution_result transform with a deploy key has a gas field of (\d+)$`,
		func(value int64) error {
			return utils.NotImplementError
		})

	ctx.Step(`the speculative_exec execution_result transform with a deploy key has (\d+) transfer with a valid transfer hash$`,
		func(transfers int) error {
			return utils.NotImplementError
		})

	ctx.Step(`the speculative_exec execution_result transform with a deploy key has as from field of the "([^"]*)" account hash$`,
		func(faucet string) error {
			return utils.NotImplementError
		})

	ctx.Step(`the speculative_exec execution_result transform with a deploy key has as source field of the "([^"]*)" account purse uref$`,
		func(faucet string) error {
			return utils.NotImplementError
		})

	ctx.Step(`the speculative_exec execution_result contains at least (\d+) valid balance transforms$`,
		func(min int) error {
			transforms, err := getFaucetBalanceTransforms(casperClient, speculativeExecResult.ExecutionResult.Success.Effect.Transforms)
			if err == nil {
				err = utils.ExpectEqual(utils.CasperT, "balance transforms", len(transforms), min)
			}
			return err
		})

	ctx.Step(`the speculative_exec execution_result (\d+)st balance transform is an Identity transform$`,
		func(first int) error {
			transforms, err := getFaucetBalanceTransforms(casperClient, speculativeExecResult.ExecutionResult.Success.Effect.Transforms)
			if err == nil {
				transform := transforms[first-1]
				err = utils.ExpectEqual(utils.CasperT, "balance transform identity", string(transform.Transform), "\"Identity\"")
			}
			return err
		})

	ctx.Step(`the speculative_exec execution_result last balance transform is an Identity transform is as WriteCLValue of type "([^"]*)"$`,
		func(typeName string) error {
			transforms, err := getFaucetBalanceTransforms(casperClient, speculativeExecResult.ExecutionResult.Success.Effect.Transforms)
			if err == nil {
				transform := transforms[len(transforms)-1]
				err = utils.ExpectEqual(utils.CasperT, "IsWriteCLValue", transform.Transform.IsWriteCLValue(), true)
				clValue, err := transform.Transform.ParseAsWriteCLValue()
				if err == nil {
					value, err := clValue.Value()
					if err == nil {
						err = utils.ExpectEqual(utils.CasperT, "clValue", value.Type.Name(), typeName)
					}

					if err == nil && value.UI512.Value().Int64() < 9999 {
						err = fmt.Errorf("clValue value %d is less than 9999", value.UI512.Value().Int64())
					}
				}
			}
			return err
		})

	ctx.Step(`the speculative_exec execution_result contains a valid AddUInt512 transform with a value of (\d+)$`,
		func(val int64) error {

			lastEntry := speculativeExecResult.ExecutionResult.Success.Effect.Transforms[len(speculativeExecResult.ExecutionResult.Success.Effect.Transforms)-1]
			err := utils.ExpectEqual(utils.CasperT, "balance transform identity", string(lastEntry.Transform), "{\"AddUInt512\":\"100000000\"}")

			if err == nil {
				//addInt := lastEntry.Transform.ParseAsAddUInt512();
				err = errors.New("not implemented .ParseAsAddUInt512()")
			}

			return err
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

func getFaucetBalanceTransforms(casperClient casper.RPCClient, transforms []types.TransformKey) ([]types.TransformKey, error) {
	balanceTransforms := make([]types.TransformKey, 0)

	info, err := getAccountInfo(casperClient, "faucet")

	if err == nil {
		key := "balance-" + strings.Split(info.Account.MainPurse.String(), "-")[1]

		for _, transform := range transforms {
			if transform.Key.String() == key {
				balanceTransforms = append(balanceTransforms, transform)
			}
		}
	}
	return balanceTransforms, err
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
