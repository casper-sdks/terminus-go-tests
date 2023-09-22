package steps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/make-software/casper-go-sdk/sse"
	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
	"github.com/make-software/casper-go-sdk/types/keypair"
	"github.com/stretchr/testify/assert"

	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
)

/**
 * The test features implementation for the deploys.feature
 */
func TestFeaturesDeploys(t *testing.T) {
	utils.TestFeatures(t, "deploys.feature", InitializeDeploys)
}

var (
	putDeployResult     rpc.PutDeployResult
	infoGetDeployResult rpc.InfoGetDeployResult
	blockHash           string
	putDeploy           casper.Deploy
	senderKey           keypair.PrivateKey
	receiverPrivateKey  keypair.PrivateKey
)

func InitializeDeploys(ctx *godog.ScenarioContext) {
	var sdk casper.RPCClient
	var receiverKey keypair.PublicKey
	var transferAmount *big.Int
	var gasPrice int

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

		assert.NotNil(utils.CasperT, senderKey, "senderKey is nil")

		keyPath = utils.GetUserKeyAssetPath(1, receiverId, "secret_key.pem")

		receiverPrivateKey, err = casper.NewED25519PrivateKeyFromPEMFile(keyPath)

		assert.NotNil(utils.CasperT, receiverPrivateKey, "receiverPrivateKey is nil")

		receiverKey = receiverPrivateKey.PublicKey()

		assert.NotNil(utils.CasperT, receiverKey, "receiverKey is nil")

		return err
	})

	ctx.Step(`^the transfer amount is (\d+)$`, func(amount int64) error {
		transferAmount = big.NewInt(amount)

		assert.NotNil(utils.CasperT, transferAmount, "transferPrice")

		return utils.Pass
	})

	ctx.Step(`^the transfer gas price is (\d+)$`, func(price int) error {
		gasPrice = price

		assert.NotNil(utils.CasperT, gasPrice, "gasPrice")

		return utils.Pass
	})

	ctx.Step(`^the deploy is given a ttl of (\d+)m$`, func(ttl int) error {
		return utils.Pass
	})

	ctx.Step(`^the deploy is put on chain "([^"]*)"$`, func(chainName string) error {
		assert.NotNil(utils.CasperT, chainName, "chainName")

		header := types.DefaultHeader()
		header.ChainName = chainName
		header.Account = senderKey.PublicKey()
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

		assert.NotNil(utils.CasperT, deploy, "deploy")

		err = deploy.SignDeploy(senderKey)

		if err != nil {
			return err
		}

		deployJson, err := json.Marshal(deploy)
		if err != nil {
			return err
		}

		assert.NotNil(utils.CasperT, deployJson)

		fmt.Println(string(deployJson))

		result, err := sdk.PutDeploy(context.Background(), *deploy)
		if err != nil {
			return err
		}

		putDeployResult = result
		putDeploy = *deploy

		return utils.Pass
	})

	ctx.Step(`^the deploy response contains a valid deploy hash of length (\d+) and an API version "([^"]*)"$`, func(hashLength int, apiVersion string) error {
		err := utils.Pass
		assert.NotNil(utils.CasperT, putDeployResult, "PutDeployResult")

		err = utils.ExpectEqual(utils.CasperT, "putDeployResult.DeployHash", len(putDeployResult.DeployHash.String()), hashLength)

		if err == nil {
			err = utils.ExpectEqual(utils.CasperT, "putDeployResult.ApiVersion", putDeployResult.ApiVersion, apiVersion)
		}

		return err
	})

	ctx.Step(`^wait for a block added event with a timeout of (\d+) seconds$`, func(timeoutSeconds int) error {
		err := utils.Pass
		var blockAddedEvent sse.BlockAddedEvent
		blockAddedEvent, err = utils.WaitForBlockAdded(putDeployResult.DeployHash.String(), timeoutSeconds)

		if err == nil {
			blockHash = blockAddedEvent.BlockAdded.BlockHash
			infoGetDeployResult, err = utils.WaitForDeploy(putDeployResult.DeployHash.String(), timeoutSeconds)
		}

		return err
	})

	ctx.Step(`^that a Transfer has been successfully deployed$`, func() error {
		err := utils.Pass

		if infoGetDeployResult.Deploy.Hash.String() != putDeployResult.DeployHash.String() {
			err = fmt.Errorf("deploy does not match hash %s", putDeployResult.DeployHash.String())
		}

		if len(infoGetDeployResult.ExecutionResults) == 0 || infoGetDeployResult.ExecutionResults[0].Result.Success == nil {
			err = fmt.Errorf("deploy %s was not succesfuly deployed", putDeployResult.DeployHash.String())
		}

		return err
	})

	ctx.Step(`^a deploy is requested via the info_get_deploy RCP method$`, func() error {
		return utils.Pass
	})

	ctx.Step(`^the deploy data has an API version of "([^"]*)"$`, func(apiVersion string) error {
		return utils.ExpectEqual(utils.CasperT, "apiVersion", infoGetDeployResult.ApiVersion, apiVersion)
	})

	ctx.Step(`^the deploy execution result has "([^"]*)" block hash$`, func(blockName string) error {
		return utils.ExpectEqual(utils.CasperT, "blockHash", infoGetDeployResult.ExecutionResults[0].BlockHash.String(), blockHash)
	})

	ctx.Step(`^the deploy execution has a cost of (\d+) motes$`, func(cost int64) error {
		return utils.ExpectEqual(utils.CasperT, "cost", infoGetDeployResult.ExecutionResults[0].Result.Success.Cost, uint64(cost))
	})

	ctx.Step(`^the deploy has a payment amount of (\d+)$`, func(payment int64) error {
		amount, err := infoGetDeployResult.Deploy.Payment.ModuleBytes.Args.Find("amount")
		if err != nil {
			return err
		}

		// ERROR the SDK only provides the named argument bytes it has not deserialized named arguments name or value fields
		value, err := amount.Value()
		if err != nil {
			return err
		}

		err = utils.ExpectEqual(utils.CasperT, "value type", value.Type.GetTypeID(), cltype.UInt512.GetTypeID())

		if err != nil {
			return err
		}

		return utils.ExpectEqual(utils.CasperT, "value", value.GetValueByType(), clvalue.NewCLUInt512(big.NewInt(payment)).GetValueByType())
	})

	ctx.Step(`^the deploy has a valid hash$`, func() error {
		return utils.ExpectEqual(utils.CasperT, "hash", infoGetDeployResult.Deploy.Hash.String(), putDeployResult.DeployHash.String())
	})

	ctx.Step(`^the deploy has a valid timestamp$`, func() error {
		if len(infoGetDeployResult.Deploy.Header.Timestamp.ToTime().String()) < 26 {
			return errors.New("invalid timestamp")
		}
		return utils.Pass
	})
	ctx.Step(`^the deploy has a valid body hash$`, func() error {
		return utils.ExpectEqual(utils.CasperT, "body hash",
			infoGetDeployResult.Deploy.Header.BodyHash.String(),
			putDeploy.Header.BodyHash.String())
	})

	ctx.Step(`^the deploy has a session type of "([^"]*)"$`, func(sessionType string) error {

		if infoGetDeployResult.Deploy.Session.Transfer == nil {
			return errors.New("missing transfer")
		}

		return utils.Pass
	})

	ctx.Step(`^the deploy is approved by user-(\d+)$`, func(userId int) error {

		approval := infoGetDeployResult.Deploy.Approvals[0]
		return utils.ExpectEqual(utils.CasperT, "approval", approval.Signer.String(), senderKey.PublicKey().String())
	})

	ctx.Step(`^the deploy has a gas price of (\d+)$`, func(gasPrice int) error {

		return utils.ExpectEqual(utils.CasperT, "gas price", infoGetDeployResult.Deploy.Header.GasPrice, uint64(gasPrice))
	})

	ctx.Step(`^the deploy has a ttl of (\d+)m$`, func(ttl int64) error {
		var expected = types.Duration(ttl * time.Minute.Nanoseconds())
		return utils.ExpectEqual(utils.CasperT, "gas price", infoGetDeployResult.Deploy.Header.TTL, expected)
	})

	ctx.Step(`^the deploy session has a "([^"]*)" argument value of type "([^"]*)"$`, func(name string, valueType string) error {
		actual, err := infoGetDeployResult.Deploy.Session.Transfer.Args.Find(name)
		if err == nil {
			value, _ := actual.Value()
			err = utils.ExpectEqual(utils.CasperT, "dependency", value.Type.Name(), valueType)
		}
		return err
	})

	ctx.Step(`^the deploy session has a "([^"]*)" argument with a numeric value of (\d+)$`, func(name string, expectedValue int) error {
		parameter, err := infoGetDeployResult.Deploy.Session.Transfer.Args.Find(name)

		if err == nil {
			value, _ := parameter.Value()
			err = utils.ExpectEqual(utils.CasperT, name, value.GetValueByType().String(), fmt.Sprintf("%d", expectedValue))
		}
		return err
	})

	ctx.Step(`^the deploy session has a "([^"]*)" argument with the public key of user-(\d+)$`, func(name string, userId int) error {
		parameter, err := infoGetDeployResult.Deploy.Session.Transfer.Args.Find(name)

		if err == nil {
			value, _ := parameter.Value()
			var pubKey = *value.PublicKey
			err = utils.ExpectEqual(utils.CasperT, name, pubKey.String(), receiverPrivateKey.PublicKey().String())
		}
		return err
	})
}
