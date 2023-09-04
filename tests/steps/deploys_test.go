package steps

import (
	"context"
	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/keypair"
	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
	"github.com/stretchr/testify/assert"
	"math/big"
	"math/rand"
	"testing"
	"time"
)

// The test features implementation for the info_get_peers.feature
func TestFeaturesDeploys(t *testing.T) {
	TestFeatures(t, "deploys.feature", InitializeDeploys)
}

func InitializeDeploys(ctx *godog.ScenarioContext) {

	var sdk casper.RPCClient
	var senderKey keypair.PrivateKey
	var receiverKey keypair.PublicKey
	var transferAmount *big.Int
	var gasPrice int
	var putDeployResult rpc.PutDeployResult

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		sdk = utils.GetSdk()
		return ctx, nil
	})

	ctx.Step(`^that user-(\d+) initiates a transfer to user-(\d+)$`, func(senderId int, receiverId int) error {

		err := utils.Pass

		keyPath := utils.GetUserKeyAssetPath(1, senderId, "secret_key.pem")

		senderKey, err = casper.NewED25519PrivateKeyFromPEMFile(keyPath)

		if err != nil {
			return err
		}

		assert.NotNil(CasperT, senderKey, "senderKey is nil")

		keyPath = utils.GetUserKeyAssetPath(1, receiverId, "secret_key.pem")

		var receiverPrivateKey keypair.PrivateKey
		receiverPrivateKey, err = casper.NewED25519PrivateKeyFromPEMFile(keyPath)

		assert.NotNil(CasperT, receiverPrivateKey, "receiverPrivateKey is nil")

		receiverKey = receiverPrivateKey.PublicKey()

		assert.NotNil(CasperT, receiverKey, "receiverKey is nil")

		return err
	})

	ctx.Step(`^the transfer amount is (\d+)$`, func(amount int64) error {

		transferAmount = big.NewInt(amount)

		assert.NotNil(CasperT, transferAmount, "transferPrice")

		return utils.Pass
	})

	ctx.Step(`^the transfer gas price is (\d+)$`, func(price int) error {

		gasPrice = price

		assert.NotNil(CasperT, gasPrice, "gasPrice")

		return utils.Pass
	})

	ctx.Step(`^the deploy is given a ttl of (\d+)m$`, func(ttl int) error {

		return utils.Pass
	})

	ctx.Step(`^the deploy is put on chain "([^"]*)"$`, func(chainName string) error {

		assert.NotNil(CasperT, chainName, "chainName")

		header := types.DefaultHeader()
		header.ChainName = chainName
		header.Account = receiverKey
		header.Timestamp = types.Timestamp(time.Now())
		payment := types.StandardPayment(big.NewInt(100000000))

		args := &types.Args{}
		args.AddArgument("amount", *clvalue.NewCLUInt512(transferAmount))
		args.AddArgument("target", clvalue.NewCLPublicKey(receiverKey))
		args.AddArgument("id", clvalue.NewCLOption(*clvalue.NewCLUInt64(rand.Uint64())))

		session := types.ExecutableDeployItem{
			Transfer: &types.TransferDeployItem{
				Args: *args,
			},
		}

		deploy, err := types.MakeDeploy(header, payment, session)

		if err != nil {
			return err
		}

		assert.NotNil(CasperT, deploy, "deploy")

		err = deploy.SignDeploy(senderKey)

		if err != nil {
			return err
		}

		result, err := sdk.PutDeploy(context.Background(), *deploy)

		if err != nil {
			return err
		}

		putDeployResult = result

		return utils.Pass
	})

	ctx.Step(`^the deploy response contains a valid deploy hash of length (\d+) and an API version "([^"]*)"$`, func(networkName string) error {

		assert.NotNil(CasperT, putDeployResult, "PutDeployResult")
		return utils.Pass
	})

	ctx.Step(`^wait for a block added event with a timeout of (\d+) seconds$`, func(networkName string) error {
		return utils.NotImplementError
	})

	ctx.Step(`^that a Transfer has been successfully deployed$`, func() error {
		return utils.NotImplementError
	})

	ctx.Step(`^a deploy is requested via the info_get_deploy RCP method$`, func() error {
		return utils.NotImplementError
	})

	ctx.Step(`^the deploy data has an API version of "([^"]*)"$`, func(apiVersion string) error {
		return utils.NotImplementError
	})

	ctx.Step(`^the deploy execution result has "([^"]*)" block hash$`, func(networkName string) error {
		return utils.NotImplementError
	})

	ctx.Step(`^the deploy execution has a cost of (\d+) motes$`, func(networkName string) error {
		return utils.NotImplementError
	})

	ctx.Step(`^the deploy has a payment amount of (\d+)$`, func(networkName string) error {
		return utils.NotImplementError
	})

	ctx.Step(`^the deploy has a valid hash$`, func() error {
		return utils.NotImplementError
	})

	ctx.Step(`^the deploy has a valid timestamp$`, func() error {
		return utils.NotImplementError
	})

	ctx.Step(`^the deploy has a valid body hash$`, func() error {
		return utils.NotImplementError
	})

	ctx.Step(`^the deploy has a session type of "([^"]*)"$`, func(sessionType string) error {
		return utils.NotImplementError
	})

	ctx.Step(`^the deploy is approved by user-(\d+)$`, func(userId int) error {
		return utils.NotImplementError
	})

	ctx.Step(`^the deploy has a gas price of (\d+)$`, func(gasPrice int) error {
		return utils.NotImplementError
	})

	ctx.Step(`^the deploy has a ttl of (\d+)m$`, func(ttl int) error {
		return utils.NotImplementError
	})

	ctx.Step(`^the deploy session has a "([^"]*)" argument value of type "([^"]*)"$`, func(name string, valueTye string) error {
		return utils.NotImplementError
	})

	ctx.Step(`^the deploy session has a "([^"]*)" argument with a numeric value of (\d+)$`, func(name string, value int) error {
		return utils.NotImplementError
	})

	ctx.Step(`^the deploy session has a "([^"]*)" argument with the public key of user-(\d+)$`, func(name string, userId int) error {
		return utils.NotImplementError
	})
}
