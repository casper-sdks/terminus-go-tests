package steps

import (
	"context"
	"fmt"
	"github.com/antchfx/jsonquery"
	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
	"strings"
	"testing"
)

func TestFeaturesEra(t *testing.T) {
	utils.TestFeatures(t, "era.feature", InitializeEraFeature)
}

func InitializeEraFeature(ctx *godog.ScenarioContext) {
	var sdk casper.RPCClient
	var eraInfo rpc.ChainGetEraInfoResult
	var latest rpc.ChainGetBlockResult
	var doc *jsonquery.Node

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		sdk = utils.GetSdk()
		return ctx, nil
	})

	ctx.Step(`^that the era summary is requested via the sdk$`, func() error {
		var err error
		latest, err = sdk.GetBlockLatest(context.Background())

		if err == nil {
			eraInfo, err = sdk.GetEraInfoByBlockHash(context.Background(), latest.Block.Hash.String())
		}

		return err
	})

	ctx.Step(`^request the era summary via the node$`, func() error {

		summary, err := utils.GetEraSummary(latest.Block.Hash.String())

		if len(summary) == 0 {
			err = fmt.Errorf("unable to obtain summary")
		}

		if err == nil {
			doc, err = jsonquery.Parse(strings.NewReader(summary))
		}

		return err
	})

	ctx.Step(`^the block hash of the returned era summary is equal to the block hash of the test node era summary$`,
		func() error {
			blockHash := jsonquery.FindOne(doc, "result/era_summary/block_hash").Value()
			return utils.ExpectEqual(utils.CasperT, "block_hash", eraInfo.EraSummary.BlockHash.String(), blockHash)
		})

	ctx.Step(`^the era of the returned era summary is equal to the era of the returned test node era summary$`,
		func() error {
			var eraId = jsonquery.FindOne(doc, "result/era_summary/era_id").Value()
			return utils.ExpectEqual(utils.CasperT, "era_id", float64(eraInfo.EraSummary.EraID), eraId)
		},
	)

	ctx.Step(`^the merkle proof of the returned era summary is equal to the merkle proof of the returned test node era summary$`,
		func() error {
			merkleProof := jsonquery.FindOne(doc, "result/era_summary/merkle_proof").Value()
			return utils.ExpectEqual(utils.CasperT, "merkle_proof", eraInfo.EraSummary.EraID, merkleProof)
		},
	)

	ctx.Step(`^the state root hash of the returned era summary is equal to the state root hash of the returned test node era summary$`,
		func() error {
			stateRootHash := jsonquery.FindOne(doc, "result/era_summary/state_root_hash").Value()
			return utils.ExpectEqual(utils.CasperT, "state_root_hash", eraInfo.EraSummary.StateRootHash.String(), stateRootHash)
		})

	ctx.Step(`^the delegators data of the returned era summary is equal to the delegators data of the returned test node era summary$`, func() error {

		if eraInfo.EraSummary.StoredValue.EraInfo == nil {
			return fmt.Errorf("MissingeraInfo.EraSummary.StoredValue.EraInfo")
		}

		var delegators = jsonquery.FindOne(doc, "result/era_summary/stored_value/EraInfo/seigniorage_allocations").ChildNodes()
		return utils.ExpectEqual(utils.CasperT, "delegators", len(eraInfo.EraSummary.StoredValue.EraInfo.SeigniorageAllocations), len(delegators))
	})

	ctx.Step(`^the validators data of the returned era summary is equal to the validators data of the returned test node era summary$`, func() error {

		if eraInfo.EraSummary.StoredValue.EraInfo == nil {
			return fmt.Errorf("MissingeraInfo.EraSummary.StoredValue.EraInfo")
		}

		var validators = jsonquery.FindOne(doc, "result/era_summary/stored_value/EraInfo/seigniorage_allocations").ChildNodes()
		return utils.ExpectEqual(utils.CasperT, "validators", len(eraInfo.EraSummary.StoredValue.EraInfo.SeigniorageAllocations), len(validators))
	})
}
