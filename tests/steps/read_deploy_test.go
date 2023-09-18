package steps

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"
)

// The test features implementation for the read_deploy.feature
func TestFeaturesReadDeploy(t *testing.T) {
	utils.TestFeatures(t, "read_deploy.feature", InitializeReadDeploy)
}

func InitializeReadDeploy(ctx *godog.ScenarioContext) {

	var deploy types.Deploy

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		return ctx, nil
	})

	ctx.Step(`^that the "transfer.json" JSON deploy is loaded$`, func() error {

		deployBytes, err := os.ReadFile("../json/transfer.json")

		if err == nil {
			err = json.Unmarshal(deployBytes, &deploy)
		}

		return err
	})

	ctx.Step(`^a valid transfer deploy is created$`, func() error {
		if &deploy == nil {
			return errors.New("deploy not created")
		}
		return utils.Pass
	})

	ctx.Step(`^the deploy hash is "([^"]*)"$`, func(deployHash string) error {
		return utils.ExpectEqual(utils.CasperT, "deployHash", deploy.Hash.String(), deployHash)
	})

	ctx.Step(`^the account is "([^"]*)"$`, func(account string) error {
		return utils.ExpectEqual(utils.CasperT, "account", deploy.Header.Account.String(), account)
	})

	ctx.Step(`^the timestamp is "([^"]*)"$`, func(timestamp string) error {
		timestamp = strings.ReplaceAll(timestamp, ".104", "")
		actual := fmt.Sprintf(deploy.Header.Timestamp.ToTime().Format(time.RFC3339))
		return utils.ExpectEqual(utils.CasperT, "timestamp", actual, timestamp)
	})

	ctx.Step(`^the ttl is (\d+)m$`, func(ttl int64) error {

		var expected = types.Duration(ttl * time.Minute.Nanoseconds())
		return utils.ExpectEqual(utils.CasperT, "ttl", deploy.Header.TTL, expected)
	})

	ctx.Step(`^the gas price is (\d+)$`, func(gasPrice int64) error {
		return utils.ExpectEqual(utils.CasperT, "gas price", deploy.Header.GasPrice, uint64(gasPrice))
	})

	ctx.Step(`^the body_hash is "([^"]*)"$`, func(bodyHash string) error {
		return utils.ExpectEqual(utils.CasperT, "bodyHash", deploy.Header.BodyHash.String(), bodyHash)
	})

	ctx.Step(`^the chain name is "([^"]*)"$`, func(chainName string) error {
		return utils.ExpectEqual(utils.CasperT, "bodyHash", deploy.Header.ChainName, chainName)
	})

	ctx.Step(`^dependency (\d+) is "([^"]*)"$`, func(index int, hex string) error {
		return utils.ExpectEqual(utils.CasperT, "dependency", deploy.Header.Dependencies[index].String(), hex)
	})
	ctx.Step(`^the payment amount is (\d+)$`, func(payment int64) error {
		actual, err := deploy.Payment.ModuleBytes.Args.Find("amount")
		if err == nil {
			value, _ := actual.Value()
			err = utils.ExpectEqual(utils.CasperT, "dependency", value.GetValueByType(), clvalue.NewCLUInt512(big.NewInt(payment)).UI512)
		}
		return err
	})

	ctx.Step(`^the session is a transfer$`, func() error {
		if deploy.Session.Transfer == nil {
			return errors.New("expected transfer")
		}
		return utils.Pass
	})

	ctx.Step(`^the session "([^"]*)" is (\d+)$`, func(parameterName string, amount int64) error {

		parameter, err := deploy.Session.Transfer.Args.Find(parameterName)

		if err == nil {
			value, _ := parameter.Value()
			err = utils.ExpectEqual(utils.CasperT, parameterName, value.GetValueByType(), clvalue.NewCLUInt512(big.NewInt(amount)).UI512)
		}
		return err
	})

	ctx.Step(`^the session "([^"]*)" is "([^"]*)"$`, func(parameterName string, strValue string) error {
		parameter, err := deploy.Session.Transfer.Args.Find(parameterName)

		if err == nil {
			value, _ := parameter.Value()
			err = utils.ExpectEqual(utils.CasperT, parameterName, value.String(), strValue)
		}
		return err
	})

	ctx.Step(`^the session "([^"]*)" type is "([^"]*)"$`, func(parameterName string, typeName string) error {
		parameter, err := deploy.Session.Transfer.Args.Find(parameterName)

		if err == nil {
			value, _ := parameter.Value()
			err = utils.ExpectEqual(utils.CasperT, parameterName, value.Type.Name(), typeName)
		}
		return err
	})

	ctx.Step(`^the session "([^"]*)" bytes is "([^"]*)"$`, func(parameterName string, hexBytes string) error {

		parameter, err := deploy.Session.Transfer.Args.Find(parameterName)

		if err == nil {
			value, _ := parameter.Value()
			actual := hex.EncodeToString(value.Bytes())
			err = utils.ExpectEqual(utils.CasperT, parameterName, actual, hexBytes)
		}
		return err
	})

	ctx.Step(`^the session "([^"]*)" parsed is "([^"]*)"$`, func(parameterName string, parsed string) error {
		parameter, err := deploy.Session.Transfer.Args.Find(parameterName)

		if err == nil {
			value, _ := parameter.Value()
			actual := value.GetValueByType().String()
			err = utils.ExpectEqual(utils.CasperT, parameterName, actual, parsed)
		}
		return err
	})

	ctx.Step(`^the deploy has (\d+) approval$`, func(approvalSize int) error {
		return utils.ExpectEqual(utils.CasperT, "approvalSize", len(deploy.Approvals), approvalSize)
	})

	ctx.Step(`^the approval signer is "([^"]*)"$`, func(signer string) error {
		actual := deploy.Approvals[0].Signer.String()
		return utils.ExpectEqual(utils.CasperT, "signer", actual, signer)
	})

	ctx.Step(`^the approval signature is "([^"]*)"$`, func(signature string) error {
		actual := deploy.Approvals[0].Signature.String()
		return utils.ExpectEqual(utils.CasperT, "signer", actual, signature)
	})
}
