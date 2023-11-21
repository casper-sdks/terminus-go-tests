package steps

import (
	"context"
	"github.com/cucumber/godog"
	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
	"testing"
)

// The test features implementation for the nested_lists.feature
func TestFeaturesNestedLists(t *testing.T) {
	utils.TestFeatures(t, "nested_lists.feature", InitializedNestedLists)
}

func InitializedNestedLists(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		return ctx, nil
	})

	ctx.Step(`a list is created with "([^"]*)" values of \("([^"]*)", "([^"]*)", "([^"]*)"\)$`,
		func(dataType string, val1 string, val2 string, val3 string) error {
			return nil
		})

	ctx.Step(`a list is created with U(\d+) values of \((\d+), (\d+), (\d+)\)$`,
		func(numLen int, val1 int, val2 int, val3 int) error {
			return nil
		},
	)

	ctx.Step(`a nested list is created with U(\d+) values of \(\((\d+), (\d+), (\d+)\),\((\d+), (\d+), (\d+)\)\)$`,
		func(numLen int, val1 int, val2 int, val3 int, val4 int, val5 int, val6 int) error {
			return nil
		},
	)

	ctx.Step(`a list is created with I(\d+) values of \((\d+), (\d+), (\d+)\)$`,
		func(numLen int, val1 int, val2 int, val3 int) error {
			return nil
		},
	)

	ctx.Step(`the list's bytes are "([^"]*)"$`, func() error {
		return nil
	})

	ctx.Step(`the list's length is (\d+)$`, func(len int) error {
		return nil
	})

	ctx.Step(`the list's "([^"]*)" item is a CLValue with "([^"]*)" value of "([^"]*)"$`,
		func(nth string, valueType string, strValue string) error {
			return nil
		},
	)

	ctx.Step(`the list's "([^"]*)" item is a CLValue with U(\d+) value of (\d+)$`,
		func(nth string, numLen int, value string) error {
			return nil
		},
	)

	ctx.Step(`the list's "([^"]*)" item is a CLValue with I(\d+) value of (\d+)$`,
		func(nth string, numLen int, value string) error {
			return nil
		},
	)

	ctx.Step(`the "([^"]*)" nested list's "([^"]*)" item is a CLValue with U(\d+) value of (\d+)$`,
		func(nth string, nestedNth string, numLen int, value string) error {
			return nil
		},
	)

	ctx.Step(`that the list is deployed in a transfer$`, func() error {
		return nil
	})

	ctx.Step(`the transfer containing the list is successfully executed$`, func() error {
		return nil
	})

	ctx.Step(`the list is read from the deploy$`, func() error {
		return nil
	})
}
