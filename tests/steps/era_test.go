package steps

import (
	"context"
	"fmt"
	"github.com/make-software/casper-go-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"

	"github.com/casper-sdks/terminus-go-tests/tests/utils"
)

// Step Definitions for the eras.feature
func TestFeaturesEra(t *testing.T) {
	utils.TestFeatures(t, "era.feature", InitializeEraFeature)
}

func InitializeEraFeature(ctx *godog.ScenarioContext) {
	var sdk casper.RPCClient
	var eraInfo rpc.ChainGetEraSummaryResult
	var eraInfoNode types.EraSummary

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		sdk = utils.GetRPCClient()
		return ctx, nil
	})

	ctx.Step(`^that the era summary is requested via the sdk$`, func() error {
		var err error
		eraInfo, err = sdk.GetEraSummaryLatest(context.Background())

		return err
	})

	ctx.Step(`^request the era summary via the node$`, func() error {
		summary, err := utils.GetEraSummary(eraInfo.EraSummary.BlockHash.String())

		assert.NotEmpty(utils.CasperT, summary)

		eraInfoNode = summary.EraSummary

		return err
	})

	ctx.Step(`^the block hash of the returned era summary is equal to the block hash of the test node era summary$`,
		func() error {
			return utils.ExpectEqual(utils.CasperT, "blockHash", eraInfo.EraSummary.BlockHash.String(), eraInfoNode.BlockHash.String())
		})

	ctx.Step(`^the era of the returned era summary is equal to the era of the returned test node era summary$`,
		func() error {
			return utils.ExpectEqual(utils.CasperT, "era_id", float64(eraInfo.EraSummary.EraID), float64(eraInfoNode.EraID))
		},
	)

	ctx.Step(`^the merkle proof of the returned era summary is equal to the merkle proof of the returned test node era summary$`,
		func() error {
			//Merkle Proof not returned by the current test node (cctl)
			return utils.Pass
		},
	)

	ctx.Step(`^the state root hash of the returned era summary is equal to the state root hash of the returned test node era summary$`,
		func() error {
			return utils.ExpectEqual(utils.CasperT, "state_root_hash", eraInfo.EraSummary.StateRootHash.String(), eraInfoNode.StateRootHash.String())
		})

	ctx.Step(`^the delegators data of the returned era summary is equal to the delegators data of the returned test node era summary$`, func() error {
		if eraInfo.EraSummary.StoredValue.EraInfo == nil {
			return fmt.Errorf("MissingeraInfo.EraSummary.StoredValue.EraInfo")
		}

		return utils.ExpectEqual(utils.CasperT, "delegators", len(eraInfo.EraSummary.StoredValue.EraInfo.SeigniorageAllocations),
			len(eraInfoNode.StoredValue.EraInfo.SeigniorageAllocations))
	})

	ctx.Step(`^the validators data of the returned era summary is equal to the validators data of the returned test node era summary$`, func() error {
		if eraInfo.EraSummary.StoredValue.EraInfo == nil {
			return fmt.Errorf("MissingeraInfo.EraSummary.StoredValue.EraInfo")
		}

		return utils.ExpectEqual(utils.CasperT, "validators", len(eraInfo.EraSummary.StoredValue.EraInfo.SeigniorageAllocations),
			len(eraInfoNode.StoredValue.EraInfo.SeigniorageAllocations))
	})
}
