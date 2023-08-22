package steps

import (
	"context"
	"github.com/cucumber/godog"
	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
	"testing"
)

func TestFeaturesWasm(t *testing.T) {
	TestFeatures(t, "wasm.feature", InitializeWasmFeature)
}

func InitializeWasmFeature(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		return ctx, nil
	})

	ctx.Step(`^that a smart contract "([^"]*)" is located in the "([^"]*)" folder$`, func(wasmFileName string, contractsFolder string) error {
		return nil
	})

	ctx.Step(`^the wasm is loaded as from the file system$`, func() error {
		return nil
	})

	ctx.Step(`^the wasm has been successfully deployed$`, func() error {
		return nil
	})

	ctx.Step(`^the account named keys contain the "([^"]*)" name$`,
		func(name string) error {
			return nil
		},
	)

	ctx.Step(`^the contract data "([^"]*)" is a "([^"]*)" with a value of "([^"]*)" and bytes of "([^"]*)"$`,
		func(path string, typeName string, value string, hexBytes string) error {
			return nil
		},
	)

	ctx.Step(`^the contract entry point is invoked with a transfer amount of "([^"]*)"$`,
		func(transferAmount string) error {
			return nil
		},
	)

	ctx.Step(`^the contract invocation deploy is successful$`, func() error {
		return nil
	})

}
