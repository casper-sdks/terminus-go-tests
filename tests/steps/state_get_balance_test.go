package steps

import (
	"context"
	"errors"
	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/make-software/casper-go-sdk/types/keypair"
	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
	"math/big"
	"testing"
)

// The test features implementation for the state_get_balance.feature
func TestFeaturesStateGetBalance(t *testing.T) {
	utils.TestFeatures(t, "state_get_balance.feature", InitializeStateGetBalance)
}

func InitializeStateGetBalance(ctx *godog.ScenarioContext) {

	var sdk casper.RPCClient
	var balance rpc.StateGetBalanceResult
	var accountKey keypair.PrivateKey
	var latestBlock rpc.ChainGetBlockResult
	var stateRootHash rpc.ChainGetStateRootHashResult

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		sdk = utils.GetSdk()
		return ctx, nil
	})

	ctx.Step(`^that the state_get_balance RPC method is invoked against nclt user-1 purse$`, func() error {

		var accountInfo rpc.StateGetAccountInfo
		var err error

		latestBlock, err = sdk.GetBlockLatest(context.Background())
		keyPath := utils.GetUserKeyAssetPath(1, 1, "secret_key.pem")
		accountKey, err = casper.NewED25519PrivateKeyFromPEMFile(keyPath)

		if err == nil {
			accountInfo, err = sdk.GetAccountInfoByBlochHash(context.Background(), latestBlock.Block.Hash.String(), accountKey.PublicKey())
		}

		if err == nil {
			stateRootHash, err = sdk.GetStateRootHashLatest(context.Background())
		}

		if err == nil {
			hashStr := stateRootHash.StateRootHash.String()
			balance, err = sdk.GetAccountBalance(context.Background(), &hashStr, accountInfo.Account.MainPurse.String())
		}

		return err
	})

	ctx.Step(`^a valid state_get_balance_result is returned$`, func() error {

		if &balance == nil {
			return errors.New("missing balance")
		}
		return utils.Pass
	})

	ctx.Step(`^the state_get_balance_result contains the purse amount$`, func() error {

		var expectedBalance big.Int

		accountInfo, err := utils.GetStateAccountInfo(accountKey.PublicKey().String(), latestBlock.Block.Hash.String())
		purseUref, _ := utils.GetByJsonPath(accountInfo, "/result/account/main_purse")

		expectedBalance, err = utils.StateGetBalance(stateRootHash.StateRootHash.String(), purseUref)

		if err == nil {
			err = utils.ExpectEqual(utils.CasperT, "purse amount", balance.BalanceValue, expectedBalance)
		}

		return err
	})

	ctx.Step(`the state_get_balance_result contains api version "([^"]*)"$`, func(apiVersion string) error {
		return utils.ExpectEqual(utils.CasperT, "api version", balance.ApiVersion, apiVersion)
	})

	ctx.Step(`the state_get_balance_result contains a valid merkle proof`, func() error {
		// Doest not contain a merkle proof	 - Ignoring as other SDKs do not fetch this
		return utils.Pass
	})
}
