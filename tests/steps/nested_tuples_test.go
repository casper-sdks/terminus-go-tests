package steps

import (
	"context"
	"github.com/cucumber/godog"
	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
	"testing"
)

// The test features implementation for the nested_tuples.feature
func TestFeaturesNestedTuples(t *testing.T) {
	utils.TestFeatures(t, "nested_tuples.feature", InitializeNestedTuples)
}

func InitializeNestedTuples(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		return ctx, nil
	})

	ctx.Step(`that a nested Tuple(\d+) is defined as \(\((\d+)\)\) using U32 numeric values$`,
		func(index int, value int) error {
			return nil
		},
	)

	ctx.Step(`that a nested Tuple(\d+) is defined as \((\d+), \((\d+), \((\d+), (\d+)\)\)\) using U32 numeric values$`,
		func(index int, value1 int, value2 int, value3 int) error {
			return nil
		},
	)
	ctx.Step(`that a nested Tuple(\d+) is defined as \((\d+), (\d+), \((\d+), (\d+), \((\d+), (\d+), (\d+)\)\)\) using U32 numeric values$`,
		func(index int, value1 int, value2 int, value3 int, value4 int, value5 int, value6 int, value7 int) error {
			return nil
		},
	)

	ctx.Step(`^the "([^"]*)" element of the Tuple(\d+) is "([^"]*)"$`,
		func(nth string, index int, strValue string) error {
			return nil
		},
	)
	ctx.Step(`^the "([^"]*)" element of the Tuple(\d+) is (\d+)$`,
		func(nth string, index int, value int) error {
			return nil
		},
	)

	ctx.Step(`the Tuple(\d+) bytes are "([^"]*)"$`, func(tupleIndex int, strHex string) error {
		//return utils.ExpectEqual(utils.CasperT, "bytes", hex.EncodeToString(clMap.Map.Bytes()), strHex)
		return nil
	})

	ctx.Step(`that the nested tuples are deployed in a transfer$`, func() error {
		return nil
	})

	ctx.Step(`the transfer is successful$`, func() error {
		return nil
	})

	ctx.Step(`the tuples deploy is obtained from the node$`, func() error {
		return nil
	})
}
