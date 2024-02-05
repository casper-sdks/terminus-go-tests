package steps

import (
	"context"
	"encoding/hex"
	"github.com/casper-sdks/terminus-go-tests/tests/utils"
	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
	"testing"
)

// The test features implementation for the option_values.feature
func TestFeaturesOptionValues(t *testing.T) {
	utils.TestFeatures(t, "option_values.feature", InitializeOptionValues)
}

func InitializeOptionValues(ctx *godog.ScenarioContext) {

	var optionValue clvalue.CLValue
	var deploy *types.Deploy
	var sdk casper.RPCClient
	var result rpc.PutDeployResult
	var deployResult casper.InfoGetDeployResult

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		return ctx, nil
	})

	ctx.Step(`that an Option value has an empty value$`, func() error {

		var bytes = []byte{0}

		fromBytes, err := clvalue.NewOptionFromBytes(bytes, cltype.NewOptionType(cltype.Bool))

		optionValue = clvalue.CLValue{
			Option: fromBytes,
		}

		return err
	})

	ctx.Step(`the Option value is not present$`, func() error {
		return utils.ExpectEqual(utils.CasperT, "empty", optionValue.Option.IsEmpty(), true)
	})

	ctx.Step(`^the Option value's bytes are "([^"]*)"$`, func(strHex string) error {
		var actualHex = ""
		if !optionValue.Option.IsEmpty() {
			actualHex = hex.EncodeToString(optionValue.Bytes())
		}
		return utils.ExpectEqual(utils.CasperT, "bytes", actualHex, strHex) // err
	})

	ctx.Step(`^an Option value contains a "([^"]*)" value of "([^"]*)"$`,
		func(typeName string, strValue string) error {
			innerValue, err := utils.CreateValue(typeName, strValue)
			if err == nil {
				optionValue = clvalue.NewCLOption(*innerValue)
			}
			return err
		},
	)

	ctx.Step(`^the type of the Option is "([^"]*)" with a value of "([^"]*)"$`,
		func(typeName string, strValue string) error {
			return nil
		},
	)

	ctx.Step(`that the Option value is deployed in a transfer as a named argument$`, func() error {
		var err error

		namedArgs := &types.Args{}
		namedArgs.AddArgument("OPTION", optionValue)
		deploy, err = utils.BuildStandardTransferDeploy(*namedArgs)

		result, err = sdk.PutDeploy(context.Background(), *deploy)

		return err
	})

	ctx.Step(`the transfer containing the Option value is successfully executed$`, func() error {
		var err error
		deployResult, err = utils.WaitForDeploy(result.DeployHash.String(), 300)
		return err
	})

	ctx.Step(`the Option is read from the deploy$`, func() error {
		arg, err := deployResult.Deploy.Session.Transfer.Args.Find("OPTION")
		if err == nil {
			optionValue, err = arg.Value()
		}
		return err
	})
}
