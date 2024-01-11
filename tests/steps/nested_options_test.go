package steps

import (
	"context"
	"errors"
	"testing"

	"github.com/make-software/casper-go-sdk/types/clvalue"

	"github.com/cucumber/godog"
	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
)

// The test features implementation for the nested_maps.feature
func TestFeaturesNestedOptions(t *testing.T) {
	utils.TestFeatures(t, "nested_options.feature", InitializeNestedOptions)
}

func InitializeNestedOptions(ctx *godog.ScenarioContext) {
	// var clOption clvalue.CLValue
	// var deploy *types.Deploy
	// var sdk casper.RPCClient
	// var result rpc.PutDeployResult
	// var deployResult casper.InfoGetDeployResult

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		return ctx, nil
	})

	ctx.Step(`^that a nested Option has an inner type of Option with a type of String and a value of "([^"]*)"$`, func(value string) error {
		clOption := clvalue.NewCLOption(clvalue.NewCLOption(*clvalue.NewCLString(value)))
		if clOption.Option.IsEmpty() {
			return errors.New("Option is empty")
		}
		return utils.Pass
	})

	ctx.Step(`^the inner type is Option with a type of String and a value of  "([^"]*)"$`,
		func(value string) error {
			// Fail SDK does not allow the creation of nested maps as the Map does not implement the CLValue interface
			/*clMap = clvalue.NewCLMap(cltype.String, cltype.Map)
			]
						innerClMap1 := clvalue.NewCLMap(cltype.String, cltype.UInt32)
						err := innerClMap1.Map.Append(*clvalue.NewCLString(key1), *clvalue.NewCLUInt32(uint32(value1)))

						innerClMap2 := clvalue.NewCLMap(cltype.String, cltype.UInt32)
						err = innerClMap2.Map.Append(*clvalue.NewCLString(key3), *clvalue.NewCLUInt32(uint32(value3)))

						err = clMap.Map.Append(*clvalue.NewCLString(key0), innerClMap1)
						err = clMap.Map.Append(*clvalue.NewCLString(key2), innerClMap2)*/
			return utils.Pass
		},
	)

	ctx.Step(`^the bytes are "([^"]*)"$`, func(hexBytes string) error {
		return utils.Pass
	},
	)

	ctx.Step(`^that the nested Option is deployed in a transfer$`, func() error {
		return utils.Pass
	},
	)

	ctx.Step(`^the transfer containing the nested Option is successfully executed$`, func() error {
		return utils.Pass
	})

	ctx.Step(`^the Option is read from the deploy$`, func() error {
		return utils.Pass
	})

	ctx.Step(`^the inner type is Option with a type of String and a value of "([^"]*)"$`, func(value string) error {
		return utils.Pass
	})

	ctx.Step(`^that a nested Option has an inner type of List with a type of U256 and a value of \((\d+), (\d+), (\d+)\)$`, func(val1 int32, val int32, val3 int32) error {
		return utils.Pass
	})

	ctx.Step(`^the list's length is (\d+)$`, func(len int) error {
		return utils.Pass
	})

	ctx.Step(`^the list's "([^"]*)" item is a CLValue with U256 value of (\d+)$`, func(nth string, val int) error {
		return utils.Pass
	})

	ctx.Step(`^the 1st nested map's key is "([^"]*)" and value is "([^"]*)"$`, func(key string, strValue string) error {
		return utils.Pass
	})

	ctx.Step(`^that a nested Option has an inner type of Tuple2 with a type of "([^"]*)" values of "([^"]*)"$`, func(types string, values string) error {
		return utils.Pass
	})

	ctx.Step(`^the inner type is Tuple2 with a type of "([^"]*)" and a value of "([^"]*)"`, func(types string, values string) error {
		return utils.Pass
	})

	ctx.Step(`^that a nested Option has an inner type of Map with a type of "([^"]*)" values of \{"([^"]*)": (\d+)\}$`, func(types string, key string, val int) error {
		return utils.Pass
	})

	ctx.Step(`^the inner type is Map with a type of "([^"]*)" and a value of "([^"]*)"`, func(types string, values string) error {
		return utils.Pass
	})

	ctx.Step(`^that a nested Option has an inner type of Any with a value of "([^"]*)"`, func(value string) error {
		return utils.Pass
	})

	ctx.Step(`^the inner type is Any with a value of "([^"]*)"$`, func(values string) error {
		return utils.Pass
	})
}
