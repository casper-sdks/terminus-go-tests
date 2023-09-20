package steps

import (
	"context"
	"errors"
	"fmt"
	"github.com/antchfx/jsonquery"
	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
	"math/big"
	"strconv"
	"testing"
)

// Step Definitions for the state_get_account_info.feature
func TestFeaturesStateGetAuctionInfo(t *testing.T) {
	utils.TestFeatures(t, "state_get_auction_info.feature", InitializeStateAuctionInfoFeature)
}

func InitializeStateAuctionInfoFeature(ctx *godog.ScenarioContext) {
	var sdk casper.RPCClient
	var auctionInfo rpc.StateGetAuctionInfoResult
	var jsonAuctionInfo string
	var rpcErr rpc.RpcError

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		sdk = utils.GetSdk()
		return ctx, nil
	})

	ctx.Step(`^that the state_get_auction_info RPC method is invoked by hash block identifier$`, func() error {

		latest, err := sdk.GetBlockLatest(context.Background())

		jsonAuctionInfo, err = utils.GetAuctionInfoByHash(latest.Block.Header.ParentHash.String())

		auctionInfo, err = sdk.GetAuctionInfoByHash(context.Background(), latest.Block.Header.ParentHash.String())

		return err
	})

	ctx.Step(`^that the state_get_auction_info RPC method is invoked by height block identifier$`, func() error {

		latest, err := sdk.GetBlockLatest(context.Background())

		jsonAuctionInfo, err = utils.GetAuctionInfoByHash(latest.Block.Header.ParentHash.String())

		auctionInfo, err = sdk.GetAuctionInfoByHeight(context.Background(), latest.Block.Header.Height)

		return err
	})

	ctx.Step(`^that the state_get_auction_info RPC method is invoked by an invalid block hash identifier$`, func() error {

		_, err := sdk.GetAuctionInfoByHash(context.Background(), "9608b4b7029a18ae35373eab879f523850a1b1fd43a3e6da774826a343af4ad2")

		if err != nil {
			rpcErr = utils.GetRpcError(err)
			return utils.Pass
		} else {
			return errors.New("should have given rpc error")
		}
	})

	ctx.Step(`^a valid state_get_auction_info_result is returned$`, func() error {

		if len(auctionInfo.AuctionState.StateRootHash) == 0 {
			return errors.New("invalid auction info")
		}
		return utils.Pass

	})

	ctx.Step(`^the state_get_auction_info_result has and api version of "([^"]*)"$`, func(apiVersion string) error {
		return utils.ExpectEqual(utils.CasperT, "apiVersion", auctionInfo.Version, apiVersion)
	})

	ctx.Step(`^the state_get_auction_info_result action_state has a valid state root hash$`, func() error {

		expectedStateRootHash, err := utils.GetByJsonPath(jsonAuctionInfo, "/result/auction_state/state_root_hash")

		if err == nil {
			err = utils.ExpectEqual(
				utils.CasperT,
				"state_root_hash",
				auctionInfo.AuctionState.StateRootHash,
				expectedStateRootHash)
		}
		return err
	})

	ctx.Step(`^the state_get_auction_info_result action_state has a valid height$`, func() error {

		var height int64

		expectedHeight, err := utils.GetByJsonPath(jsonAuctionInfo, "/result/auction_state/block_height")

		if err == nil {
			height, err = strconv.ParseInt(expectedHeight, 10, 64)
		}

		if err == nil {
			err = utils.ExpectEqual(utils.CasperT,
				"height",
				// There is the possibility that the height may have incremented so account for that too
				auctionInfo.AuctionState.BlockHeight == uint64(height) || auctionInfo.AuctionState.BlockHeight == uint64(height+1),
				true)
		}

		return err
	})

	ctx.Step(`^the state_get_auction_info_result action_state has valid bids$`, func() error {

		bidsNode, err := utils.GetNodeByJsonPath(jsonAuctionInfo, "/result/auction_state/bids")

		err = utils.ExpectEqual(utils.CasperT,
			"bids length",
			len(auctionInfo.AuctionState.Bids),
			len(bidsNode.ChildNodes()))

		if err == nil {
			var publicKey *jsonquery.Node
			publicKey, err = jsonquery.Query(bidsNode.FirstChild, "/public_key")

			err = utils.ExpectEqual(utils.CasperT,
				"public_key",
				auctionInfo.AuctionState.Bids[0].PublicKey.String(),
				publicKey.Value())
		}

		if err == nil {
			var bondingPurse *jsonquery.Node
			bondingPurse, err = jsonquery.Query(bidsNode.FirstChild, "/bid/bonding_purse")

			err = utils.ExpectEqual(
				utils.CasperT,
				"public_key",
				auctionInfo.AuctionState.Bids[0].Bid.BondingPurse.String(),
				bondingPurse.Value())
		}

		if err == nil {
			var delegationRate *jsonquery.Node
			delegationRate, err = jsonquery.Query(bidsNode.FirstChild, "/bid/delegation_rate")

			err = utils.ExpectEqual(
				utils.CasperT,
				"public_key",
				auctionInfo.AuctionState.Bids[0].Bid.DelegationRate,
				float32(delegationRate.Value().(float64)))
		}

		if err == nil {
			var inactive *jsonquery.Node
			inactive, err = jsonquery.Query(bidsNode.FirstChild, "/bid/inactive")

			err = utils.ExpectEqual(
				utils.CasperT,
				"public_key",
				auctionInfo.AuctionState.Bids[0].Bid.Inactive,
				inactive.Value())
		}

		if err == nil {
			var stakedAmount *jsonquery.Node
			stakedAmount, err = jsonquery.Query(bidsNode.FirstChild, "/bid/staked_amount")
			val := big.Int{}
			val.SetString(fmt.Sprintf("%v", stakedAmount.Value()), 10)

			err = utils.ExpectEqual(
				utils.CasperT,
				"public_key",
				auctionInfo.AuctionState.Bids[0].Bid.StakedAmount,
				val.Uint64())
		}

		if err == nil {
			var delegators *jsonquery.Node
			var delegatee *jsonquery.Node
			delegators, err = jsonquery.Query(bidsNode.FirstChild, "/bid/delegators")
			delegatee, err = jsonquery.Query(delegators.FirstChild, "/delegatee")

			err = utils.ExpectEqual(
				utils.CasperT,
				"public_key",
				auctionInfo.AuctionState.Bids[0].Bid.Delegators[0].Delegatee.String(),
				delegatee.Value())
		}

		if err == nil {
			var delegators *jsonquery.Node
			var publicKey *jsonquery.Node
			delegators, err = jsonquery.Query(bidsNode.FirstChild, "/bid/delegators")
			publicKey, err = jsonquery.Query(delegators.FirstChild, "/public_key")

			err = utils.ExpectEqual(
				utils.CasperT,
				"public_key",
				auctionInfo.AuctionState.Bids[0].Bid.Delegators[0].PublicKey.String(),
				publicKey.Value())
		}

		if err == nil {
			var delegators *jsonquery.Node
			var stakedAmount *jsonquery.Node
			delegators, err = jsonquery.Query(bidsNode.FirstChild, "/bid/delegators")
			stakedAmount, err = jsonquery.Query(delegators.FirstChild, "/staked_amount")
			val := big.Int{}
			val.SetString(fmt.Sprintf("%v", stakedAmount.Value()), 10)

			err = utils.ExpectEqual(
				utils.CasperT,
				"public_key",
				auctionInfo.AuctionState.Bids[0].Bid.Delegators[0].StakedAmount,
				val.Uint64())
		}

		return err
	})

	ctx.Step(`^the state_get_auction_info_result action_state has valid era validators$`, func() error {

		validatorsNode, err := utils.GetNodeByJsonPath(jsonAuctionInfo, "/result/auction_state/era_validators")

		err = utils.ExpectEqual(
			utils.CasperT,
			"bids length",
			len(auctionInfo.AuctionState.EraValidators),
			len(validatorsNode.ChildNodes()))

		if err == nil {
			eraId := jsonquery.FindOne(validatorsNode.FirstChild, "/era_id")

			err = utils.ExpectEqual(
				utils.CasperT,
				"eraId",
				auctionInfo.AuctionState.EraValidators[0].EraID,
				uint32(eraId.Value().(float64)))
		}

		validatorWeights := jsonquery.FindOne(validatorsNode.FirstChild, "/validator_weights")

		if err == nil {
			publicKey := jsonquery.FindOne(validatorWeights.FirstChild, "/public_key")

			err = utils.ExpectEqual(
				utils.CasperT,
				"eraId",
				auctionInfo.AuctionState.EraValidators[0].ValidatorWeights[0].Validator.String(),
				publicKey.Value())
		}

		if err == nil {
			weight := jsonquery.FindOne(validatorWeights.FirstChild, "/weight")
			val := big.Int{}
			val.SetString(fmt.Sprintf("%v", weight.Value()), 10)

			err = utils.ExpectEqual(
				utils.CasperT,
				"eraId",
				auctionInfo.AuctionState.EraValidators[0].ValidatorWeights[0].Weight.Value().String(),
				val.String())
		}

		return err
	})

	ctx.Step(`^an error code of -(\d+) is returned$`, func(errorCode int) error {
		return utils.ExpectEqual(utils.CasperT, "error code", rpcErr.Code, -1*errorCode)
	})

	ctx.Step(`^an error message of "([^"]*)" is returned$`, func(errorMessage string) error {
		return utils.ExpectEqual(utils.CasperT, "error code", rpcErr.Message, errorMessage)
	})
}
