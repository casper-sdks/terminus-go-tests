package steps

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/stretchr/testify/assert"

	"github.com/casper-sdks/terminus-go-tests/tests/utils"
)

// The test features implementation for the info_get_peers.feature
func TestFeaturesInfoGetStatus(t *testing.T) {
	utils.TestFeatures(t, "info_get_status.feature", InitializeInfoGetStatus)
}

func InitializeInfoGetStatus(ctx *godog.ScenarioContext) {
	var sdk casper.RPCClient
	var infoGetStatusResult rpc.InfoGetStatusResult
	nodeGetStatus := casper.InfoGetStatusResult{}

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		sdk = utils.GetRPCClient()
		return ctx, nil
	})

	ctx.Step(`^that the info_get_status is invoked against nctl$`, func() error {
		var err error

		infoGetStatusResult, err = sdk.GetStatus(context.Background())

		if err == nil {
			nodeGetStatus, err = utils.GetNodeStatus(1)
		}

		return err
	})

	ctx.Step(`^an info_get_status_result is returned`, func() error {
		assert.NotNil(utils.CasperT, infoGetStatusResult, "infoGetStatusResult is nil")
		return utils.Pass
	})

	ctx.Step(`^the info_get_status_result api_version is "([^"]*)"$`, func(apiVersion string) error {
		if apiVersion != infoGetStatusResult.APIVersion {
			return fmt.Errorf("expected %s ApiVersion to be %s", infoGetStatusResult.APIVersion, apiVersion)
		}
		return utils.Pass
	})

	ctx.Step(`^the info_get_status_result chainspec_name is "([^"]*)"$`, func(_ string) error {
		if utils.GetConfigChainName() != infoGetStatusResult.ChainSpecName {
			return fmt.Errorf("expected %s ChainSpecName to be %s", infoGetStatusResult.ChainSpecName, utils.GetConfigChainName())
		}

		return utils.Pass
	})

	ctx.Step(`^the info_get_status_result has a valid last_added_block_info$`, func() error {
		err := utils.ExpectEqual(utils.CasperT, "hash", infoGetStatusResult.LastAddedBlockInfo.Hash, nodeGetStatus.LastAddedBlockInfo.Hash)

		if err == nil {
			err = utils.ExpectEqual(utils.CasperT, "Timestamp", infoGetStatusResult.LastAddedBlockInfo.Timestamp, nodeGetStatus.LastAddedBlockInfo.Timestamp)
		}

		if err == nil {
			err = utils.ExpectEqual(utils.CasperT, "EraID", infoGetStatusResult.LastAddedBlockInfo.EraID, nodeGetStatus.LastAddedBlockInfo.EraID)
		}

		if err == nil {
			err = utils.ExpectEqual(utils.CasperT, "Height", infoGetStatusResult.LastAddedBlockInfo.Height, nodeGetStatus.LastAddedBlockInfo.Height)
		}

		if err == nil {
			err = utils.ExpectEqual(utils.CasperT, "StateRootHash", infoGetStatusResult.LastAddedBlockInfo.StateRootHash, nodeGetStatus.LastAddedBlockInfo.StateRootHash)
		}

		if err == nil {
			err = utils.ExpectEqual(utils.CasperT, "Creator", infoGetStatusResult.LastAddedBlockInfo.Creator, nodeGetStatus.LastAddedBlockInfo.Creator)
		}

		return err
	})

	ctx.Step(`^the info_get_status_result has a valid our_public_signing_key$`, func() error {
		return utils.ExpectEqual(utils.CasperT, "our_public_signing_key", infoGetStatusResult.OutPublicSigningKey, nodeGetStatus.OutPublicSigningKey)
	})

	ctx.Step(`^the info_get_status_result has a valid starting_state_root_hash$`, func() error {
		return utils.ExpectEqual(utils.CasperT, "StartingStateRootHash", infoGetStatusResult.StartingStateRootHash, nodeGetStatus.StartingStateRootHash)
	})

	ctx.Step(`^the info_get_status_result has a valid build_version$`, func() error {
		return utils.ExpectEqual(utils.CasperT, "BuildVersion", infoGetStatusResult.BuildVersion, nodeGetStatus.BuildVersion)
	})

	ctx.Step(`^the info_get_status_result has a valid round_length$`, func() error {
		return utils.ExpectEqual(utils.CasperT, "RoundLength", infoGetStatusResult.RoundLength, nodeGetStatus.RoundLength)
	})

	ctx.Step(`^the info_get_status_result has a valid uptime$`, func() error {
		if !strings.HasSuffix(infoGetStatusResult.Uptime, "ms") {
			return fmt.Errorf("missing ms from uptime %s", infoGetStatusResult.Uptime)
		}

		return utils.Pass
	})

	ctx.Step(`^the info_get_status_result has a valid peers$`, func() error {
		return utils.ExpectEqual(utils.CasperT, "Peers", infoGetStatusResult.Peers, nodeGetStatus.Peers)
	})
}
