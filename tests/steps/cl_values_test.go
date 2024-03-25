package steps

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
	"github.com/make-software/casper-go-sdk/types/keypair"
	"github.com/stretchr/testify/assert"

	"github.com/casper-sdks/terminus-go-tests/tests/utils"
)

/**
 * The test features implementation for the cl_values.feature
 */
func TestClValues(t *testing.T) {
	utils.TestFeatures(t, "cl_values.feature", InitializeClValues)
}

func InitializeClValues(ctx *godog.ScenarioContext) {
	var sdk casper.RPCClient
	testArgs := &types.Args{}
	lastVal := clvalue.CLValue{}
	var clValuesDeploy *types.Deploy
	var clValuesDeployResult rpc.PutDeployResult
	var clValuesInfoGetDeployResult casper.InfoGetDeployResult

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		sdk = utils.GetRPCClient()
		return ctx, nil
	})

	ctx.Step(`^that a CL value of type "([^"]*)" has a value of "([^"]*)"$`, func(typeName string, value string) error {
		clVal, err := utils.CreateValue(typeName, value)
		testArgs.AddArgument(typeName, *clVal)
		lastVal = *clVal
		return err
	})

	ctx.Step(`^it's bytes will be "([^"]*)"$`, func(hexBytes string) error {
		decoded, err := hex.DecodeString(hexBytes)
		if !bytes.Equal(lastVal.Bytes(), decoded) {
			err = fmt.Errorf("%s bytes do not match expected bytes %s", hex.EncodeToString(lastVal.Bytes()), hexBytes)
		}
		return err
	})

	ctx.Step(`^that the CL complex value of type "([^"]*)" with an internal types of "([^"]*)" values of "([^"]*)"$`,
		func(typeName string, internalTypes string, values string) error {
			clVal, err := utils.CreateComplexValue(typeName, strings.Split(internalTypes, ","), strings.Split(values, ","))
			testArgs.AddArgument(typeName, *clVal)
			lastVal = *clVal
			return err
		})

	ctx.Step(`^the values are added as arguments to a deploy$`, func() error {
		keyPath := utils.GetUserKeyAssetPath(1, 1, "secret_key.pem")

		senderKey, err := casper.NewED25519PrivateKeyFromPEMFile(keyPath)
		if err != nil {
			return err
		}

		keyPath = utils.GetUserKeyAssetPath(1, 2, "secret_key.pem")

		var receiverPrivateKey keypair.PrivateKey
		receiverPrivateKey, err = casper.NewED25519PrivateKeyFromPEMFile(keyPath)
		if err != nil {
			return err
		}

		header := types.DefaultHeader()
		header.ChainName = utils.GetChainName()
		header.Account = senderKey.PublicKey()
		header.Timestamp = types.Timestamp(time.Now())
		payment := types.StandardPayment(big.NewInt(100000000))

		args := &types.Args{}
		args.AddArgument("amount", *clvalue.NewCLUInt512(big.NewInt(2500000000)))
		args.AddArgument("target", clvalue.NewCLPublicKey(receiverPrivateKey.PublicKey()))
		args.AddArgument("id", clvalue.NewCLOption(*clvalue.NewCLUInt64(rand.Uint64())))

		var name string
		var val clvalue.CLValue

		for _, arg := range *testArgs {
			name, err = arg.Name()
			if err != nil {
				return err
			}
			val, err = arg.Value()
			if err != nil {
				return err
			}
			args.AddArgument(name, val)
		}

		session := types.ExecutableDeployItem{
			Transfer: &types.TransferDeployItem{
				Args: *args,
			},
		}

		clValuesDeploy, err = types.MakeDeploy(header, payment, session)
		if err != nil {
			return err
		}

		assert.NotNil(utils.CasperT, clValuesDeploy, "deploy")

		err = clValuesDeploy.SignDeploy(senderKey)

		if err != nil {
			return err
		}

		deployJson, err := json.Marshal(clValuesDeploy)
		if err != nil {
			return err
		}

		assert.NotNil(utils.CasperT, deployJson)

		fmt.Println(string(deployJson))

		return err
	})

	ctx.Step(`^the deploy is put on chain$`, func() error {
		result, err := sdk.PutDeploy(context.Background(), *clValuesDeploy)
		clValuesDeployResult = result
		return err
	})

	ctx.Step(`^the deploy response contains a valid deploy hash of length (\d+) and an API version "([^"]*)"$`,
		func(hashLength int, apiVersion string) error {
			err := utils.ExpectEqual(utils.CasperT, "API", clValuesDeployResult.ApiVersion, apiVersion)
			if err == nil {
				err = utils.ExpectEqual(utils.CasperT, "hashLength", len(clValuesDeployResult.DeployHash.String()), hashLength)
			}
			return err
		})

	ctx.Step(`^the deploy has successfully executed$`, func() error {
		_, err := utils.WaitForBlockAdded(clValuesDeployResult.DeployHash.String(), 300)
		return err
	})

	ctx.Step(`^the deploy is obtained from the node$`, func() error {
		var err error

		clValuesInfoGetDeployResult, err = utils.WaitForDeploy(clValuesDeployResult.DeployHash.String(), 300)

		if clValuesInfoGetDeployResult.Deploy.Hash.String() != clValuesDeployResult.DeployHash.String() {
			err = fmt.Errorf("unable to obtain deploy for hash %s", clValuesDeployResult.DeployHash.String())
		}
		return err
	})

	ctx.Step(`^the deploys NamedArgument "([^"]*)" has a value of "([^"]*)" and bytes of "([^"]*)"$`,
		func(name string, strVal string, hexBytes string) error {
			args := clValuesInfoGetDeployResult.Deploy.Session.Transfer.Args

			arg, err := args.Find(name)
			var expectedValue *clvalue.CLValue
			var value clvalue.CLValue

			if err == nil {
				expectedValue, err = utils.CreateValue(name, strVal)
			}

			if err == nil {
				value, err = arg.Value()
			}

			if err == nil {
				err = utils.ExpectEqual(utils.CasperT, "value", value.GetValueByType(), expectedValue.GetValueByType())
			}

			if err == nil {
				err = utils.ExpectEqual(utils.CasperT, "bytes", hex.EncodeToString(value.Bytes()), hexBytes)
			}

			return err
		})

	stepFunc := func(name string, internalTypes string, values string, hexBytes string) error {
		var value clvalue.CLValue
		var expectedValue *clvalue.CLValue

		args := clValuesInfoGetDeployResult.Deploy.Session.Transfer.Args
		arg, err := args.Find(name)
		if err == nil {
			value, err = arg.Value()
		}

		if err == nil {
			err = utils.ExpectEqual(utils.CasperT, "bytes", hex.EncodeToString(value.Bytes()), hexBytes)
		}

		if err == nil {
			expectedValue, err = utils.CreateComplexValue(name, strings.Split(internalTypes, ","), strings.Split(values, ","))
		}

		if err == nil {
			if name == "Map" {
				for i := 0; i < 3; i++ {
					// Map do not maintain order so we need to find the value by key
					key := strconv.FormatInt(int64(i), 10)
					expected, _ := expectedValue.Map.Find(key)
					actual, _ := value.Map.Find(key)
					err = utils.ExpectEqual(utils.CasperT, "value", actual.GetValueByType().String(), expected.GetValueByType().String())
					if err != nil {
						return err
					}
				}
			} else {
				err = utils.ExpectEqual(utils.CasperT, "value", value.GetValueByType().String(), expectedValue.GetValueByType().String())
			}
		}

		return err
	}
	ctx.Step(`^the deploys NamedArgument Complex value "([^"]*)" has internal types of "([^"]*)" and values of "([^"]*)" and bytes of "([^"]*)"$`,
		stepFunc)
}
