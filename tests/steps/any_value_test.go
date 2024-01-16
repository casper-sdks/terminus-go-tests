package steps

import (
	"context"
	"encoding/hex"
	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
	"testing"
)

/**
 * The test features implementation for the cl_values.feature
 */
func TestClAnyValue(t *testing.T) {
	utils.TestFeatures(t, "any_value.feature", InitializeClAnyValue)
}

func InitializeClAnyValue(ctx *godog.ScenarioContext) {

	value := clvalue.CLValue{}
	var deploy *types.Deploy
	var rpcClient casper.RPCClient
	var result rpc.PutDeployResult
	var deployResult casper.InfoGetDeployResult

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		rpcClient = utils.GetRPCClient()
		return ctx, nil
	})

	ctx.Step(`^an Any value contains a "([^"]*)" value of "([^"]*)"$`, func(typeName string, hexBytes string) error {
		decoded, err := hex.DecodeString(hexBytes)
		if err == nil {
			value = clvalue.NewCLAny(decoded)
		}
		return err
	})

	ctx.Step(`^the any value's bytes are "([^"]*)"$`, func(hexBytes string) error {
		decoded, err := hex.DecodeString(hexBytes)
		if err == nil {
			err = utils.ExpectEqual(utils.CasperT, "bytes", value.Bytes(), decoded)
		}
		return err
	})

	ctx.Step(`^that the any value is deployed in a transfer as a named argument$`, func() error {
		var err error
		namedArgs := &types.Args{}
		namedArgs.AddArgument("any", value)
		deploy, err = utils.BuildStandardTransferDeploy(*namedArgs)
		result, err = rpcClient.PutDeploy(context.Background(), *deploy)
		return err
	})
	ctx.Step(`^the transfer containing the any value is successfully executed$`, func() error {
		var err error
		deployResult, err = utils.WaitForDeploy(result.DeployHash.String(), 300)
		return err
	})

	ctx.Step(`^the any is read from the deploy$`, func() error {
		arg, err := deployResult.Deploy.Session.Transfer.Args.Find("any")
		if err == nil {
			value, err = arg.Value()
		}
		return err
	})

	ctx.Step(`^the any value's bytes are "([^"]*)"$`, func(hexBytes string) error {
		decoded, err := hex.DecodeString(hexBytes)
		if err == nil {
			err = utils.ExpectEqual(utils.CasperT, "bytes", value.Bytes(), decoded)
		}
		return err
	})

	ctx.Step(`^that the map of public keys to any types is read from resource "([^"]*)"$`, func(jsonFileName string) error {
		return utils.NotImplementError
	})

	ctx.Step(`^the loaded CLMap will contain (\d+) elements$`, func(count int) error {
		return utils.NotImplementError
	})

	ctx.Step(`^the nested map key type will be "([^"]*)"$`, func(keyType string) error {
		return utils.NotImplementError
	})

	ctx.Step(`^the nested map value type will be "([^"]*)"$`, func(valueType string) error {
		return utils.NotImplementError
	})

	ctx.Step(`^the any value's bytes are "([^"]*)"$`, func(hexBytes string) error {
		return utils.NotImplementError
	})

	ctx.Step(`^the maps bytes will be "([^"]*)"$`, func(hexBytes string) error {
		return utils.NotImplementError
	})

	ctx.Step(`^the nested map keys value will be "([^"]*)"$`, func(keyValue string) error {
		return utils.NotImplementError
	})

	ctx.Step(`the nested map any values bytes length will be (\d+)$`, func(value int) error {
		return utils.NotImplementError
	})

	ctx.Step(`the nested map any values bytes will be "([^"]*)"$`, func(hexBytes string) error {
		return utils.NotImplementError
	})
}
