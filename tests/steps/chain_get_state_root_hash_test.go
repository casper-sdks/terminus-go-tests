package steps

import (
	"context"
	"testing"

	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"

	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
)

// Step Definitions for the eras.feature
func TestFeaturesChainGetStateRootHash(t *testing.T) {
	utils.TestFeatures(t, "chain_get_state_root_hash.feature", InitializeChainGetStateRootHashFeature)
}

func InitializeChainGetStateRootHashFeature(ctx *godog.ScenarioContext) {
	var sdk casper.RPCClient
	var stateRootHash rpc.ChainGetStateRootHashResult
	var expectedStateRootHash string

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		sdk = utils.GetSdk()
		return ctx, nil
	})

	ctx.Step(`^that the chain_get_state_root_hash RCP method is invoked against nctl$`, func() error {
		var err error
		stateRootHash, err = sdk.GetStateRootHashLatest(context.Background())
		if err == nil {
			expectedStateRootHash, err = utils.GetStateRootHash(1)
		}
		return err
	})

	ctx.Step(`^a valid chain_get_state_root_hash_result is returned$`, func() error {
		return utils.ExpectEqual(utils.CasperT, "state root hash", stateRootHash.StateRootHash.String(), expectedStateRootHash)
	})
}
