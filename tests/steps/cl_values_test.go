package steps

import (
	"context"
	"fmt"
	"testing"

	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/casper"

	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
)

/**
 * The test features implementation for the deploys.feature
 */
func TestClValues(t *testing.T) {
	TestFeatures(t, "cl_values.feature", InitializeClValues)
}

func InitializeClValues(ctx *godog.ScenarioContext) {
	var sdk casper.RPCClient

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		sdk = utils.GetSdk()
		return ctx, nil
	})

	ctx.Step(`^that a CL value of type "([^"]*)" has a value of "([^"]*)"$`, func(typeName string, value string) error {
		err := utils.Pass
		if sdk == nil {
			err = fmt.Errorf("SDK is nil")
		}
		return err
	})

	ctx.Step(`^it's bytes will be "([^"]*)"$`, func(hexBytes string) error {
		return utils.Pass
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
