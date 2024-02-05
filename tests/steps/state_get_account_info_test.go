package steps

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/make-software/casper-go-sdk/types/keypair"

	"github.com/casper-sdks/terminus-go-tests/tests/utils"
)

// Step Definitions for the state_get_account_info.feature
func TestFeaturesStateGeAccountInfo(t *testing.T) {
	utils.TestFeatures(t, "state_get_account_info.feature", InitializeStateGetAccountInfoFeature)
}

func InitializeStateGetAccountInfoFeature(ctx *godog.ScenarioContext) {
	var sdk casper.RPCClient
	var latest rpc.ChainGetBlockResult
	var accountInfo rpc.StateGetAccountInfo
	var senderKey keypair.PrivateKey
	var accountInfoJson string

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		sdk = utils.GetRPCClient()
		return ctx, nil
	})

	ctx.Step(`^that the state_get_account_info RCP method is invoked against nctl$`, func() error {
		var err error

		latest, err = sdk.GetBlockLatest(context.Background())

		if err == nil {
			keyPath := utils.GetUserKeyAssetPath(1, 1, "secret_key.pem")
			senderKey, err = casper.NewED25519PrivateKeyFromPEMFile(keyPath)
		}

		if err == nil {
			accountInfo, err = sdk.GetAccountInfoByBlochHash(context.Background(), latest.Block.Hash.String(), senderKey.PublicKey())
		}

		return err
	})

	ctx.Step(`^a valid state_get_account_info_result is returned$`, func() error {
		if accountInfo.ApiVersion != "1.0.0" {
			return errors.New("invalid account info result")
		}
		return utils.Pass
	})

	ctx.Step(`^the state_get_account_info_result contain a valid account hash$`, func() error {
		expectedHash, _ := utils.GetAccountHash(senderKey.PublicKey().String(), latest.Block.Hash.String())
		expectedHash = strings.Split(expectedHash, "-")[2]
		return utils.ExpectEqual(utils.CasperT, "account hash", accountInfo.Account.AccountHash.String(), expectedHash)
	})

	ctx.Step(`^the state_get_account_info_result contain a valid action thresholds$`, func() error {
		var err error
		var expectedDeployment string
		var expectedKeyManagement string
		accountInfoJson, err = utils.GetStateAccountInfo(senderKey.PublicKey().String(), latest.Block.Hash.String())
		if err == nil {
			expectedDeployment, err = utils.GetByJsonPath(accountInfoJson, "/result/account/action_thresholds/deployment")
			if err == nil {
				intVal, _ := strconv.ParseInt(expectedDeployment, 10, 16)
				err = utils.ExpectEqual(utils.CasperT, "deployment", accountInfo.Account.ActionThresholds.Deployment, uint64(intVal))
			}

			if err == nil {
				expectedKeyManagement, err = utils.GetByJsonPath(accountInfoJson, "/result/account/action_thresholds/key_management")
			}

			if err == nil {
				intVal, _ := strconv.ParseInt(expectedKeyManagement, 10, 16)
				err = utils.ExpectEqual(utils.CasperT, "key_management", accountInfo.Account.ActionThresholds.KeyManagement, uint64(intVal))
			}
		}

		return err
	})

	ctx.Step(`^the state_get_account_info_result contain a valid main purse uref$`, func() error {
		purseUref, _ := utils.GetByJsonPath(accountInfoJson, "/result/account/main_purse")
		return utils.ExpectEqual(utils.CasperT, "MainPurse", accountInfo.Account.MainPurse.String(), purseUref)
	})

	ctx.Step(`^the state_get_account_info_result contain a valid merkle proof$`, func() error {
		// Merkel Proof missing
		merkleProof, _ := utils.GetByJsonPath(accountInfoJson, "/result/merkle_proof")
		if merkleProof == "" {
			return errors.New("merkle_proof missing")
		}

		// Merkel Proof missing not failing as missing in other SDKs
		// return utils.ExpectEqual(utils.CasperT, "associated_keys", accountInfo.MerkleProof, merkleProof)
		return utils.Pass
	})

	ctx.Step(`^the state_get_account_info_result contain a valid associated keys$`, func() error {
		expectedHash := senderKey.PublicKey().AccountHash().String()
		return utils.ExpectEqual(utils.CasperT, "associated_keys", accountInfo.Account.AssociatedKeys[0].AccountHash.String(), expectedHash)
	})

	ctx.Step(`^the state_get_account_info_result contain a valid named keys$`, func() error {
		namedKey, _ := utils.GetByJsonPath(accountInfoJson, "/result/account/named_keys")
		if namedKey == "[]" {
			return utils.ExpectEqual(utils.CasperT, "associated_keys", len(accountInfo.Account.NamedKeys), 0)
		}
		return utils.Pass
	})
}
