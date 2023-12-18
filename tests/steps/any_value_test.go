package steps

import (
	"context"
	"github.com/cucumber/godog"
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

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		//	sdk = utils.GetRPCClient()
		return ctx, nil
	})

	ctx.Step(`^an Any value contains a "([^"]*)" value of "([^"]*)"$`, func(typeName string, value string) error {
		return utils.Pass
	})

	ctx.Step(`^the any value's bytes are "([^"]*)"$`, func(hexBytes string) error {
		return utils.Pass
	})

	ctx.Step(`^that the any value is deployed in a transfer as a named argument$`, func() error {
		return utils.Pass
	})
	ctx.Step(`^the transfer containing the any value is successfully executed$`, func() error {
		return utils.Pass
	})

	ctx.Step(`^the any is read from the deploy$`, func() error {
		return utils.Pass
	})

	ctx.Step(`^the any value's bytes are "([^"]*)"$`, func(hexBytes string) error {
		return utils.Pass
	})

	/////
	ctx.Step(`^that the map of public keys to any types is read from resource "([^"]*)"$`, func(jsonFileName string) error {
		return utils.Pass
	})

	ctx.Step(`^the loaded CLMap will contain (\d+) elements$`, func(count int) error {
		return utils.Pass
	})

	ctx.Step(`^the nested map key type will be "([^"]*)"$`, func(keyType string) error {
		return utils.Pass
	})

	ctx.Step(`^the nested map value type will be "([^"]*)"$`, func(valueType string) error {
		return utils.Pass
	})
	ctx.Step(`^the any value's bytes are "([^"]*)"$`, func(hexBytes string) error {
		return utils.Pass
	})
	ctx.Step(`^the maps bytes will be "([^"]*)"$`, func(hexBytes string) error {
		return utils.Pass
	})

	ctx.Step(`^the nested map keys value will be "([^"]*)"$`, func(keyValue string) error {
		return utils.Pass
	})
	ctx.Step(`the nested map any values bytes length will be (\d+)$`, func(value int) error {
		return utils.Pass
	})

	ctx.Step(`the nested map any values bytes will be "([^"]*)"$`, func(hexBytes string) error {
		return utils.Pass
	})

}
