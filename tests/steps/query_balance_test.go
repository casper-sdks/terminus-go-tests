package steps

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"

	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
)

// Step Definitions for the query_balance.feature
func TestFeaturesQueryGetBalance(t *testing.T) {
	utils.TestFeatures(t, "query_balance.feature", InitializeQueryGetBalance)
}

func InitializeQueryGetBalance(ctx *godog.ScenarioContext) {
	var sdk casper.RPCClient
	var queryGetBalanceResult rpc.QueryBalanceResult
	var queryGetBalanceJson string

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		sdk = utils.GetSdk()
		return ctx, nil
	})

	ctx.Step(`^that a query balance is obtained by main purse public key$`, func() error {
		faucetKey, err := casper.NewED25519PrivateKeyFromPEMFile("../../assets/net-1/faucet/secret_key.pem")

		if err == nil {
			publicKey := faucetKey.PublicKey()
			identifier := rpc.PurseIdentifier{
				MainPurseUnderPublicKey: &publicKey,
			}
			queryGetBalanceResult, err = sdk.QueryBalance(context.Background(), identifier)

			if err == nil {
				queryGetBalanceJson, err = utils.QueryBalance("main_purse_under_public_key", publicKey.ToHex())
			}
		}
		return err
	})

	ctx.Step(`^a valid query_balance_result is returned$`, func() error {
		if len(queryGetBalanceJson) == 0 {
			return errors.New("Missing queryGetBalanceJson")
		}
		return utils.Pass
	})

	ctx.Step(`^the query_balance_result has an API version of "([^"]*)"$`, func(apiVersion string) error {
		return utils.ExpectEqual(utils.CasperT, "api version", queryGetBalanceResult.ApiVersion, apiVersion)
	})

	ctx.Step(`^the query_balance_result has a valid balance$`, func() error {
		balance, err := utils.GetByJsonPath(queryGetBalanceJson, "/result/balance")

		if err == nil {
			actual := queryGetBalanceResult.Balance.Value()
			expected, _ := new(big.Int).SetString(balance, 10)
			return utils.ExpectEqual(utils.CasperT, "balance", actual, expected)
		}

		return err
	})

	ctx.Step(`^that a query balance is obtained by main purse account hash$`, func() error {
		faucetKey, err := casper.NewED25519PrivateKeyFromPEMFile("../../assets/net-1/faucet/secret_key.pem")

		if err == nil {
			accountHash := "account-hash-" + faucetKey.PublicKey().AccountHash().ToHex()
			identifier := rpc.PurseIdentifier{
				MainPurseUnderAccountHash: &accountHash,
			}
			queryGetBalanceResult, err = sdk.QueryBalance(context.Background(), identifier)

			if err == nil {
				queryGetBalanceJson, err = utils.QueryBalance("main_purse_under_account_hash", accountHash)
			}
		}
		return err
	})

	ctx.Step(`^that a query balance is obtained by main purse uref$`, func() error {
		var latest rpc.ChainGetBlockResult
		var accountInfo rpc.StateGetAccountInfo

		faucetKey, err := casper.NewED25519PrivateKeyFromPEMFile("../../assets/net-1/faucet/secret_key.pem")

		if err == nil {
			latest, err = sdk.GetBlockLatest(context.Background())
		}

		if err == nil {
			accountInfo, err = sdk.GetAccountInfoByBlochHash(context.Background(), latest.Block.Hash.String(), faucetKey.PublicKey())
		}

		if err == nil {
			purseUref := accountInfo.Account.MainPurse.String()
			identifier := rpc.PurseIdentifier{
				PurseUref: &purseUref,
			}
			queryGetBalanceResult, err = sdk.QueryBalance(context.Background(), identifier)

			if err == nil {
				queryGetBalanceJson, err = utils.QueryBalance("purse_uref", purseUref)
			}
		}

		return err
	})

	ctx.Step(`^a transfer of (\d+) is made to user-(\d+)'s purse$`, func(amount int, userId int) error {
		return utils.NotImplementError
	})

	ctx.Step(`^that a query balance is obtained by user-(\d+)'s main purse public and latest block identifier$`, func(userId int) error {
		return utils.NotImplementError
	})

	ctx.Step(`^the balance includes the transferred amount$`, func() error {
		return utils.NotImplementError
	})

	ctx.Step(`^that a query balance is obtained by user-1's main purse public key and previous block identifier$`, func() error {
		return utils.NotImplementError
	})

	ctx.Step(`^the balance is the pre transfer amount$`, func() error {
		return utils.NotImplementError
	})

	ctx.Step(`^that a query balance is obtained by user-1's main purse public and latest state root hash identifier$`, func() error {
		return utils.NotImplementError
	})

	ctx.Step(`^that a query balance is obtained by user-1's main purse public key and previous state root hash identifier$`, func() error {
		return utils.NotImplementError
	})
}
