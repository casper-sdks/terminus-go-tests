package steps

import (
	"context"
	"github.com/cucumber/godog"
	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFeaturesBlocks(t *testing.T) {
	TestFeatures(t, "blocks.feature", InitializeBlocksScenario)
}

func InitializeBlocksScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		return ctx, nil
	})

	ctx.Step(`^request the latest block via the test node$`, func() error {
		block, err := utils.GetNctlLatestBlock()
		if err != nil {
			return err
		}

		assert.NotEmpty(CasperT, block)

		contextMap.blockDataNode = block

		return nil
	})

	ctx.Step(`^that the latest block is requested via the sdk$`, func() error {

		block, err := utils.GetSdk().GetBlockLatest(context.Background())
		if err != nil {
			return err
		}

		assert.NotEmpty(CasperT, block)

		contextMap.blockDataSdk = block

		return nil
	})

	ctx.Step(`^the body of the returned block is equal to the body of the returned test node block$`, func() error {
		err := utils.AssertExpectedAndActual(
			assert.Equal, contextMap.blockDataSdk.Block.Body, contextMap.blockDataNode.Body,
		)
		return utils.Result(err)
	})

	ctx.Step(`^the hash of the returned block is equal to the hash of the returned test node block$`, func() error {
		err := utils.AssertExpectedAndActual(
			assert.Equal, contextMap.blockDataSdk.Block.Hash, contextMap.blockDataNode.Hash,
		)
		return utils.Result(err)
	})

	ctx.Step(`^the header of the returned block is equal to the header of the returned test node block$`, func() error {
		err := utils.AssertExpectedAndActual(
			assert.Equal, contextMap.blockDataSdk.Block.Header, contextMap.blockDataNode.Header,
		)
		return utils.Result(err)
	})

	ctx.Step(`^the proofs of the returned block are equal to the proofs of the returned test node block$`, func() error {
		err := utils.AssertExpectedAndActual(
			assert.Equal, contextMap.blockDataSdk.Block.Proofs, contextMap.blockDataNode.Proofs,
		)
		return utils.Result(err)
	})
}
