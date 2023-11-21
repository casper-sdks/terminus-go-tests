package steps

import (
	"context"
	"github.com/cucumber/godog"
	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
	"testing"
)

// The test features implementation for the nested_tuples.feature
func TestFeaturesOptionValues(t *testing.T) {
	utils.TestFeatures(t, "option_values.feature", InitializeOptionValues)
}

func InitializeOptionValues(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		return ctx, nil
	})

	ctx.Step(`that an Option value has an empty value$`,
		func() error {

			return utils.Pass

		},
	)

	ctx.Step(`the Option value is not present$`,
		func() error {
			return utils.Pass
		},
	)

	ctx.Step(`^the Option value's bytes are "([^"]*)"$`,
		func(strHex string) error {
			return nil // err
		},
	)
	ctx.Step(`^the type of the Option is "([^"]*)" with a value of "([^"]*)"$`,
		func(typeName string, strValue string) error {
			return nil // err
		},
	)

	ctx.Step(`^an Option value contains a "([^"]*)" value of "([^"]*)"$`,
		func(typeName string, strValue string) error {
			return nil
		},
	)

	ctx.Step(`that the Option value is deployed in a transfer as a named argument$`, func() error {
		return nil
	})

	ctx.Step(`the transfer containing the Option value is successfully executed$`, func() error {
		return nil
	})
	ctx.Step(`the Option is read from the deploy$`, func() error {
		return nil
	})
}
