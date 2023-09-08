package steps

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
)

/**
 * The test features implementation for the deploys.feature
 */
func TestClValues(t *testing.T) {
	TestFeatures(t, "cl_values.feature", InitializeClValues)
}

func InitializeClValues(ctx *godog.ScenarioContext) {
	//var sdk casper.RPCClient
	args := &types.Args{}
	lastVal := clvalue.CLValue{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		//sdk = utils.GetSdk()
		return ctx, nil
	})

	ctx.Step(`^that a CL value of type "([^"]*)" has a value of "([^"]*)"$`, func(typeName string, value string) error {

		clVal, err := utils.CreateValue(typeName, value)
		args.AddArgument(typeName, *clVal)
		lastVal = *clVal
		return err
	})

	ctx.Step(`^it's bytes will be "([^"]*)"$`, func(hexBytes string) error {
		decoded, err := hex.DecodeString(hexBytes)
		if !bytes.Equal(lastVal.Bytes(), decoded) {
			err = fmt.Errorf("bytes do not match expected bytes %s", hexBytes)
		}
		return err
	})

	ctx.Step(`^that the CL complex value of type "([^"]*)" with an internal types of "([^"]*)" values of "([^"]*)"$`, func(typeName string, internalTypes string, values string) error {
		return utils.Pass
	})

	ctx.Step(`^the values are added as arguments to a deploy$`, func() error {
		return utils.Pass
	})

	ctx.Step(`^the deploy is put on chain$`, func() error {
		return utils.Pass
	})

	ctx.Step(`^the deploy response contains a valid deploy hash of length (\d+) and an API version "([^"]*)"$`, func(hashLength int, apiVersion string) error {
		return utils.Pass
	})

	ctx.Step(`^the deploy response contains a valid deploy hash of length 64 and an API version "([^"]*)"$`, func() error {
		return utils.Pass
	})

	ctx.Step(`^the deploy has successfully executed$`, func() error {
		return utils.Pass
	})

	ctx.Step(`^the deploy data has an API version of "([^"]*)"$`, func(apiVersion string) error {
		return utils.Pass
	})

	ctx.Step(`^the deploy is obtained from the node$`, func() error {
		return utils.Pass
	})

	ctx.Step(`^the deploys NamedArgument "([^"]*)" has a value of "([^"]*)" and bytes of "([^"]*)"$`, func(cost int64) error {
		return utils.Pass
	})

	ctx.Step(`^the deploys NamedArgument Complex value "([^"]*)" has internal types of "([^"]*)" and values of "([^"]*)" and bytes of "([^"]*)"$`, func(payment int64) error {
		return utils.Pass
	})
}
