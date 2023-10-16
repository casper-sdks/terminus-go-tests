package steps

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/key"
	"github.com/make-software/casper-go-sdk/types/keypair"

	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
)

// Test steps for the wasm.feature
func TestFeaturesWasm(t *testing.T) {
	utils.TestFeatures(t, "wasm.feature", InitializeWasmFeature)
}

func InitializeWasmFeature(ctx *godog.ScenarioContext) {
	var sdk casper.RPCClient
	var wasmDeployResult rpc.PutDeployResult
	var wasmBytes []byte
	var faucetKey keypair.PrivateKey
	var stateRootHash string
	var contractHash string

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		sdk = utils.GetRPCClient()
		return ctx, nil
	})

	ctx.Step(`^that a smart contract "([^"]*)" is located in the "([^"]*)" folder$`, func(wasmFileName string, contractsFolder string) error {
		var err error
		wasmPath := fmt.Sprintf("../%s/%s", contractsFolder, wasmFileName)
		wasmBytes, err = os.ReadFile(wasmPath)

		if err == nil {
			err = utils.ExpectEqual(utils.CasperT, "wasmBytes length", len(wasmBytes), 189336)
		}

		return err
	})

	ctx.Step(`^the wasm is loaded as from the file system$`, func() error {
		var err error

		faucetKey, err = casper.NewED25519PrivateKeyFromPEMFile("../../assets/net-1/faucet/secret_key.pem")

		if err != nil {
			return err
		}

		header := types.DefaultHeader()
		header.ChainName = "casper-net-1"
		header.Account = faucetKey.PublicKey()
		header.Timestamp = types.Timestamp(time.Now())

		args := &types.Args{}
		args.AddArgument("token_decimals", *clvalue.NewCLUint8(11)).
			AddArgument("token_name", *clvalue.NewCLString("Acme Token")).
			AddArgument("token_symbol", *clvalue.NewCLString("ACME")).
			AddArgument("token_total_supply", *clvalue.NewCLUInt256(big.NewInt(500000000000)))

		payment := types.StandardPayment(big.NewInt(200000000000))

		session := types.ExecutableDeployItem{
			ModuleBytes: &types.ModuleBytes{
				ModuleBytes: hex.EncodeToString(wasmBytes),
				Args:        args,
			},
		}

		deploy, err := types.MakeDeploy(header, payment, session)

		if err == nil {
			err = deploy.SignDeploy(faucetKey)
		}

		if err == nil {
			wasmDeployResult, err = sdk.PutDeploy(context.Background(), *deploy)
		}

		return err
	})

	ctx.Step(`^the wasm has been successfully deployed$`, func() error {
		deploy, err := utils.WaitForDeploy(wasmDeployResult.DeployHash.String(), 300)

		if err == nil && deploy.ExecutionResults[0].Result.Success == nil {
			err = errors.New("deploy was not successful")
		}

		return err
	})

	ctx.Step(`^the account named keys contain the "([^"]*)" name$`, func(name string) error {
		var stateResult rpc.QueryGlobalStateResult
		latest, err := sdk.GetStateRootHashLatest(context.Background())
		if err != nil {
			return err
		}

		accountHash := "account-hash-" + faucetKey.PublicKey().AccountHash().String()
		stateRootHash = latest.StateRootHash.String()
		path := make([]string, 0)

		stateResult, err = sdk.QueryGlobalStateByStateHash(context.Background(), &stateRootHash, accountHash, path)

		if stateResult.StoredValue.Account == nil {
			err = errors.New("invalid result")
		}

		if err == nil {
			err = utils.ExpectEqual(utils.CasperT,
				"contract-name",
				strings.ToUpper(stateResult.StoredValue.Account.NamedKeys[0].Name),
				strings.ToUpper(name))
		}

		if err == nil {
			if !strings.Contains(stateResult.StoredValue.Account.NamedKeys[0].Key.String(), "hash-") {
				err = errors.New("missing key value")
			} else {
				contractHash = stateResult.StoredValue.Account.NamedKeys[0].Key.String()
			}
		}
		return err
	})

	ctx.Step(`^the contract data "([^"]*)" is a "([^"]*)" with a value of "([^"]*)" and bytes of "([^"]*)"$`,
		func(path string, typeName string, value string, hexBytes string) error {
			var paths []string
			paths = append(paths, path)
			var bytes []byte

			stateItem, err := sdk.QueryGlobalStateByStateHash(context.Background(), &stateRootHash, contractHash, paths)
			if err != nil {
				return err
			}

			clValue, err := stateItem.StoredValue.CLValue.Value()

			if err == nil {
				err = utils.ExpectEqual(utils.CasperT, "type", clValue.Type.Name(), typeName)
			}

			if err == nil {
				err = utils.ExpectEqual(utils.CasperT, "value", clValue.String(), value)
			}

			if err == nil {
				bytes, err = stateItem.StoredValue.CLValue.Bytes()
			}

			if err == nil {
				strBytes := hex.EncodeToString(bytes)
				err = utils.ExpectEqual(utils.CasperT, "bytes", strBytes, hexBytes)
			}
			return err
		},
	)

	ctx.Step(`^the contract entry point is invoked with a transfer amount of "([^"]*)"$`,
		func(transferAmount string) error {
			recipient, err := keypair.GeneratePrivateKey(keypair.ED25519)
			if err != nil {
				return err
			}

			hash, err := casper.NewContractHash(strings.Split(contractHash, "-")[1])

			if err == nil {
				var deploy *types.Deploy

				header := types.DefaultHeader()
				header.ChainName = "casper-net-1"
				header.Account = faucetKey.PublicKey()
				header.Timestamp = types.Timestamp(time.Now())

				txAmt := new(big.Int)
				txAmt.SetString(transferAmount, 10)

				args := &types.Args{}
				accountHashBytes := recipient.PublicKey().AccountHash().Bytes()
				args.AddArgument("recipient", clvalue.NewCLByteArray(accountHashBytes)).
					AddArgument("amount", *clvalue.NewCLUInt256(txAmt))

				session := types.ExecutableDeployItem{
					StoredContractByHash: &types.StoredContractByHash{
						Hash:       hash,
						EntryPoint: "transfer",
						Args:       args,
					},
				}

				payment := types.StandardPayment(big.NewInt(2500000000))

				deploy, err = types.MakeDeploy(header, payment, session)

				if err == nil {
					err = deploy.SignDeploy(faucetKey)
				}

				if err == nil {
					wasmDeployResult, err = sdk.PutDeploy(context.Background(), *deploy)
				}
			}

			return err
		},
	)

	ctx.Step(`^the contract invocation deploy is successful$`, func() error {
		deploy, err := utils.WaitForDeploy(wasmDeployResult.DeployHash.String(), 300)

		if len(deploy.ExecutionResults) == 0 {
			return errors.New("failed to successfully deploy")
		}
		return err
	})

	ctx.Step(`^the the contract is invoked by name "([^"]*)" and a transfer amount of "([^"]*)"$`,
		func(contractName string, transferAmount string) error {
			recipient, err := keypair.GeneratePrivateKey(keypair.ED25519)

			if err == nil {
				var deploy *types.Deploy

				header := types.DefaultHeader()
				header.ChainName = "casper-net-1"
				header.Account = faucetKey.PublicKey()
				header.Timestamp = types.Timestamp(time.Now())

				txAmt := new(big.Int)
				txAmt.SetString(transferAmount, 10)

				args := &types.Args{}
				accountHashBytes := recipient.PublicKey().AccountHash().Bytes()
				args.AddArgument("recipient", clvalue.NewCLByteArray(accountHashBytes)).
					AddArgument("amount", *clvalue.NewCLUInt256(txAmt))

				session := types.ExecutableDeployItem{
					StoredContractByName: &types.StoredContractByName{
						Name:       contractName,
						EntryPoint: "transfer",
						Args:       args,
					},
				}

				payment := types.StandardPayment(big.NewInt(2500000000))

				deploy, err = types.MakeDeploy(header, payment, session)

				if err == nil {
					err = deploy.SignDeploy(faucetKey)
				}

				if err == nil {
					wasmDeployResult, err = sdk.PutDeploy(context.Background(), *deploy)
				}
			}

			return err
		})

	ctx.Step(`^the the contract is invoked by hash and version with a transfer amount of "([^"]*)"$`,
		func(transferAmount string) error {
			recipient, err := keypair.GeneratePrivateKey(keypair.ED25519)

			if err == nil {
				var deploy *types.Deploy
				var version json.Number = "1"
				var hash key.ContractHash
				hash, err = casper.NewContractHash(strings.Split(contractHash, "-")[1])

				if err != nil {
					return err
				}

				header := types.DefaultHeader()
				header.ChainName = "casper-net-1"
				header.Account = faucetKey.PublicKey()
				header.Timestamp = types.Timestamp(time.Now())

				txAmt := new(big.Int)
				txAmt.SetString(transferAmount, 10)

				args := &types.Args{}
				accountHashBytes := recipient.PublicKey().AccountHash().Bytes()
				args.AddArgument("recipient", clvalue.NewCLByteArray(accountHashBytes)).
					AddArgument("amount", *clvalue.NewCLUInt256(txAmt))

				session := types.ExecutableDeployItem{
					StoredVersionedContractByHash: &types.StoredVersionedContractByHash{
						Hash:       hash,
						EntryPoint: "transfer",
						Version:    &version,
						Args:       args,
					},
				}

				payment := types.StandardPayment(big.NewInt(2500000000))

				deploy, err = types.MakeDeploy(header, payment, session)

				if err == nil {
					err = deploy.SignDeploy(faucetKey)
				}

				if err == nil {
					wasmDeployResult, err = sdk.PutDeploy(context.Background(), *deploy)
				}
			}

			return err
		})

	ctx.Step(`^the the contract is invoked by name "([^"]*)" and version with a transfer amount of "([^"]*)"$`,
		func(contractName string, transferAmount string) error {
			recipient, err := keypair.GeneratePrivateKey(keypair.ED25519)

			if err == nil {
				var deploy *types.Deploy
				var version json.Number = "1"

				header := types.DefaultHeader()
				header.ChainName = "casper-net-1"
				header.Account = faucetKey.PublicKey()
				header.Timestamp = types.Timestamp(time.Now())

				txAmt := new(big.Int)
				txAmt.SetString(transferAmount, 10)

				args := &types.Args{}
				accountHashBytes := recipient.PublicKey().AccountHash().Bytes()
				args.AddArgument("recipient", clvalue.NewCLByteArray(accountHashBytes)).
					AddArgument("amount", *clvalue.NewCLUInt256(txAmt))

				session := types.ExecutableDeployItem{
					StoredVersionedContractByName: &types.StoredVersionedContractByName{
						Name:       contractName,
						EntryPoint: "transfer",
						Version:    &version,
						Args:       args,
					},
				}

				payment := types.StandardPayment(big.NewInt(2500000000))

				deploy, err = types.MakeDeploy(header, payment, session)

				if err == nil {
					err = deploy.SignDeploy(faucetKey)
				}

				if err == nil {
					wasmDeployResult, err = sdk.PutDeploy(context.Background(), *deploy)
				}
			}

			return err
		})
}
