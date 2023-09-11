package steps

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/keypair"
	"github.com/stretchr/testify/assert"
	"math/big"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
)

/**
 * The test features implementation for the cl_values.feature
 */
func TestClValues(t *testing.T) {
	TestFeatures(t, "cl_values.feature", InitializeClValues)
}

func InitializeClValues(ctx *godog.ScenarioContext) {
	var sdk casper.RPCClient
	testArgs := &types.Args{}
	lastVal := clvalue.CLValue{}
	var clValuesDeploy *types.Deploy
	var clValuesDeployResult rpc.PutDeployResult
	var clValuesInfoGetDeployResult casper.InfoGetDeployResult

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		sdk = utils.GetSdk()
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
		keyPath = utils.GetUserKeyAssetPath(1, 2, "secret_key.pem")

		var receiverPrivateKey keypair.PrivateKey
		receiverPrivateKey, err = casper.NewED25519PrivateKeyFromPEMFile(keyPath)

		header := types.DefaultHeader()
		header.ChainName = "casper-net-1"
		header.Account = senderKey.PublicKey()
		header.Timestamp = types.Timestamp(time.Now())
		payment := types.StandardPayment(big.NewInt(100000000))

		args := &types.Args{}
		args.AddArgument("amount", *clvalue.NewCLUInt512(big.NewInt(2500000000)))
		args.AddArgument("target", clvalue.NewCLPublicKey(receiverPrivateKey.PublicKey()))
		args.AddArgument("id", clvalue.NewCLOption(*clvalue.NewCLUInt64(rand.Uint64())))

		var name string
		var val = clvalue.CLValue{}
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

		assert.NotNil(CasperT, clValuesDeploy, "deploy")

		err = clValuesDeploy.SignDeploy(senderKey)

		if err != nil {
			return err
		}

		deployJson, err := json.Marshal(clValuesDeploy)
		if err != nil {
			return err
		}

		assert.NotNil(CasperT, deployJson)

		fmt.Println(string(deployJson))

		return err
	})

	ctx.Step(`^the deploy is put on chain$`, func() error {

		result, err := sdk.PutDeploy(context.Background(), *clValuesDeploy)
		clValuesDeployResult = result
		return err
	})

	ctx.Step(`^the deploy response contains a valid deploy hash of length (\d+) and an API version "([^"]*)"$`, func(hashLength int, apiVersion string) error {

		err := utils.ExpectEqual(CasperT, "API", clValuesDeployResult.ApiVersion, apiVersion)
		return err
	})

	ctx.Step(`^the deploy has successfully executed$`, func() error {
		_, err := utils.WaitForBlockAdded(clValuesDeployResult.DeployHash.String(), 300)
		return err
	})

	ctx.Step(`^the deploy data has an API version of "([^"]*)"$`, func(apiVersion string) error {
		return utils.NotImplementError
	})

	ctx.Step(`^the deploy is obtained from the node$`, func() error {
		err := utils.Pass
		clValuesInfoGetDeployResult, err = utils.WaitForDeploy(clValuesDeployResult.DeployHash.String(), 300)

		if clValuesInfoGetDeployResult.Deploy.Hash.String() != clValuesDeployResult.DeployHash.String() {
			err = fmt.Errorf("unable to obtain deploy for hash %s", clValuesDeployResult.DeployHash.String())
		}
		return err
	})

	ctx.Step(`^the deploys NamedArgument "([^"]*)" has a value of "([^"]*)" and bytes of "([^"]*)"$`, func(name string, strVal string, hexBytes string) error {

		args := clValuesInfoGetDeployResult.Deploy.Session.Transfer.Args

		arg, err := args.Find(name)
		var expectedValue *clvalue.CLValue
		expectedValue, err = utils.CreateValue(name, strVal)

		if err == nil {
			var value clvalue.CLValue
			value, err = arg.Value()

			err = utils.ExpectEqual(CasperT, "value", value.GetValueByType(), expectedValue.GetValueByType())
			if err == nil {
				err = utils.ExpectEqual(CasperT, "bytes", hex.EncodeToString(value.Bytes()), hexBytes)
			}
		}
		return err
	})

	ctx.Step(`^the deploys NamedArgument Complex value "([^"]*)" has internal types of "([^"]*)" and values of "([^"]*)" and bytes of "([^"]*)"$`,
		func(name string, internalTypes string, values string, hexBytes string) error {
			var value clvalue.CLValue
			args := clValuesInfoGetDeployResult.Deploy.Session.Transfer.Args
			arg, err := args.Find(name)
			value, err = arg.Value()

			if err == nil {
				err = utils.ExpectEqual(CasperT, "bytes", hex.EncodeToString(value.Bytes()), hexBytes)
			}

			var expectedValue *clvalue.CLValue
			expectedValue, err = utils.CreateComplexValue(name, strings.Split(internalTypes, ","), strings.Split(values, ","))
			if err == nil {
				err = utils.ExpectEqual(CasperT, "value", value.GetValueByType().String(), expectedValue.GetValueByType().String())
			}
			return err
		})
}
