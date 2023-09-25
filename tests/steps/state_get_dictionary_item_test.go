package steps

import (
	"context"
	"testing"

	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"

	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
)

// The test features implementation for the state_get_dictionary_item.feature
func TestFeaturesStateGetDictionaryItem(t *testing.T) {
	utils.TestFeatures(t, "state_get_dictionary_item.feature", InitializeStateGetDictionaryItem)
}

func InitializeStateGetDictionaryItem(ctx *godog.ScenarioContext) {
	var sdk casper.RPCClient
	var mainPurse string
	var dictionaryItem rpc.StateGetDictionaryResult
	var faucetAccountHash string

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		sdk = utils.GetSdk()
		return ctx, nil
	})

	ctx.Step(`^that the state_get_dictionary_item RCP method is invoked$`, func() error {
		var latestBlock rpc.ChainGetBlockResult
		var accountInfo rpc.StateGetAccountInfo
		var stateRootHash rpc.ChainGetStateRootHashResult

		faucetKey, err := casper.NewED25519PrivateKeyFromPEMFile("../../assets/net-1/faucet/secret_key.pem")

		if err == nil {
			faucetAccountHash = faucetKey.PublicKey().AccountHash().String()
			stateRootHash, err = sdk.GetStateRootHashLatest(context.Background())
		}

		if err == nil {
			latestBlock, err = sdk.GetBlockLatest(context.Background())
		}

		if err == nil {
			accountInfo, err = sdk.GetAccountInfoByBlochHash(context.Background(), latestBlock.Block.Hash.String(), faucetKey.PublicKey())
			mainPurse = accountInfo.Account.MainPurse.String()
		}
		if err == nil {
			// TODO find an example of how to use URef to obtain dictionary item

			srh := stateRootHash.StateRootHash.String()
			dictionaryItem, err = sdk.GetDictionaryItem(context.Background(), &srh, mainPurse, "wtfDoIPutHereCantSeeAnyWorkingExampleInAnyCodeBase")
		}

		return err
	})

	ctx.Step(`^a valid state_get_dictionary_item_result is returned$`, func() error {
		var err error
		accountHash := dictionaryItem.StoredValue.Account.AccountHash.String()

		err = utils.ExpectEqual(utils.CasperT, "accountHash", accountHash, faucetAccountHash)

		if err == nil {
			err = utils.ExpectEqual(utils.CasperT, "mainPurse", dictionaryItem.StoredValue.Account.MainPurse.String(), mainPurse)
		}

		return err
	})
}
