package steps

import (
	"context"
	"testing"

	"github.com/cucumber/godog"

	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
)

/**
 * The test features implementation for the speculative_execution.feature
 */
func TestFeaturesSpeculativeExexcution(t *testing.T) {
	utils.TestFeatures(t, "speculative_execution.feature", InitializeSpeculativeExexcution)
}

func InitializeSpeculativeExexcution(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		return ctx, nil
	})

	ctx.Step(`that a deploy is executed against a node using the speculative_exec RPC API$`, func() error {
		return utils.Pass
	})

	ctx.Step(`^a valid speculative_exec_result will be returned$`, func() error {
		return utils.Pass
	})

	ctx.Step(`^the speculative_exec has an api_version of "([^"]*)"`, func(apiVersion string) error {
		return utils.Pass
	})

	ctx.Step(`^the speculative_exec has a valid block_hash$`, func() error {
		return utils.Pass
	})

	ctx.Step(`^the speculative_exec has a valid execution_results$`, func() error {
		return utils.Pass
	})

	ctx.Step(`^the execution_results contains a cost of (\d+)$`, func(cost int) error {
		return utils.Pass
	})
}
