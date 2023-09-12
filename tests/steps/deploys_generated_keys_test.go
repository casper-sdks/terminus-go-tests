package steps

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/keypair"
	"github.com/stretchr/testify/assert"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/casper"

	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
)

var keysDeployResult rpc.PutDeployResult

const ED25519 = "Ed25519"
const SECP256K1 = "Secp256k1"

// The test features implementation for the deploys_generated_keys.feature
func TestFeaturesGeneratedKeys(t *testing.T) {
	TestFeatures(t, "deploys_generated_keys.feature", InitializeGeneratedKeys)
}

func InitializeGeneratedKeys(ctx *godog.ScenarioContext) {

	var sdk casper.RPCClient
	var senderKey keypair.PrivateKey
	var faucetKey keypair.PrivateKey
	var receiverKey keypair.PrivateKey

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		sdk = utils.GetSdk()
		return ctx, nil
	})

	ctx.Step(`^that a "([^"]*)" sender key is generated$`, func(keyAlgo string) error {

		var err error

		senderKey, err = generateKey(keyAlgo)

		if err == nil && senderKey.PublicKey().Bytes() == nil {
			err = fmt.Errorf("missing sender public key")
		}

		if err == nil {
			msg := []byte("this is the sender")
			sign, err := senderKey.Sign(msg)
			if err == nil {
				err = senderKey.PublicKey().VerifySignature(msg, sign)
			}
		}

		return err
	})

	ctx.Step(`^fund the account from the faucet user with a transfer amount of (\d+) and a payment amount of (\d+)$`, func(transfer int64, payment int64) error {

		var err error

		faucetKey, err = casper.NewED25519PrivateKeyFromPEMFile("../../assets/net-1/faucet/secret_key.pem")

		if err == nil {
			err = doDeploy(sdk, faucetKey, senderKey.PublicKey(), transfer, payment)
		}

		return err
	})

	ctx.Step(`^wait for a block added event with a timeout of (\d+) seconds$`, func(timeoutSeconds int) error {

		_, err := utils.WaitForBlockAdded(keysDeployResult.DeployHash.String(), timeoutSeconds)
		return err
	})

	ctx.Step(`^that a "([^"]*)" receiver key is generated$`, func(keyAlgo string) error {

		var err error

		receiverKey, err = generateKey(keyAlgo)

		if err == nil && receiverKey.PublicKey().Bytes() == nil {
			err = fmt.Errorf("missing receiver public key")
		}

		return err
	})

	ctx.Step(`^transfer to the receiver account the transfer amount of (\d+) and the payment amount of (\d+)$`, func(transfer int64, payment int64) error {
		return doDeploy(sdk, senderKey, receiverKey.PublicKey(), transfer, payment)
	})

	ctx.Step(`the transfer approvals signer contains the "([^"]*)" algo$`, func(keyAlg string) error {

		deploy, err := sdk.GetDeploy(context.Background(), keysDeployResult.DeployHash.String())

		approval := deploy.Deploy.Approvals[0]
		tagByte := approval.Signer.Bytes()[0]

		if (ED25519 == keyAlg && tagByte != 1) || (SECP256K1 == keyAlg && tagByte != 2) {
			err = fmt.Errorf("invalid key algorithm %s for tag byte %b", keyAlg, tagByte)
		}

		return err
	})
}

func doDeploy(sdk casper.RPCClient,

	faucet keypair.PrivateKey,
	receiverKey keypair.PublicKey,
	transfer int64,
	payment int64) error {

	var deployJson []byte
	header := types.DefaultHeader()
	header.ChainName = "casper-net-1"
	header.Account = faucet.PublicKey()
	header.Timestamp = types.Timestamp(time.Now())
	stdPayment := types.StandardPayment(big.NewInt(payment))

	args := &types.Args{}
	args.AddArgument("amount", *clvalue.NewCLUInt512(big.NewInt(transfer)))
	args.AddArgument("target", clvalue.NewCLPublicKey(receiverKey))
	args.AddArgument("id", clvalue.NewCLOption(*clvalue.NewCLUInt64(rand.Uint64())))

	session := types.ExecutableDeployItem{
		Transfer: &types.TransferDeployItem{
			Args: *args,
		},
	}

	deploy, err := types.MakeDeploy(header, stdPayment, session)

	if err == nil {
		assert.NotNil(CasperT, deploy, "deploy")
		err = deploy.SignDeploy(faucet)
	}

	if err == nil {
		deployJson, err = json.Marshal(deploy)
		assert.NotNil(CasperT, deployJson)
		fmt.Println(string(deployJson))
	}

	if err == nil {
		keysDeployResult, err = sdk.PutDeploy(context.Background(), *deploy)

		if keysDeployResult.DeployHash.Bytes() == nil {
			err = fmt.Errorf("missing deploy hash")
		}
	}

	return err
}

func generateKey(keyAlgo string) (keypair.PrivateKey, error) {

	if keyAlgo == ED25519 {
		return keypair.GeneratePrivateKey(keypair.ED25519)
	} else if keyAlgo == SECP256K1 {
		return keypair.GeneratePrivateKey(keypair.SECP256K1)
	} else {
		return keypair.PrivateKey{}, fmt.Errorf("unsupported keyAlgo %s", keyAlgo)
	}
}
