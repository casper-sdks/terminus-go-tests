package steps

import (
	"context"
	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
	"testing"
)

// The test features implementation for the info_get_peers.feature
func TestFeaturesReadDeploy(t *testing.T) {
	utils.TestFeatures(t, "read_deploy.feature", InitializeReadDeploy)
}

func InitializeReadDeploy(ctx *godog.ScenarioContext) {

	var sdk casper.RPCClient

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		sdk = utils.GetSdk()
		if sdk == nil {

		}
		return ctx, nil
	})

	ctx.Step(`^that the "transfer.json" JSON deploy is loaded$`, func() error {

		err := utils.Pass

		return err
	})

	ctx.Step(`^a valid transfer deploy is created$`, func() error {

		return utils.Pass
	})

	ctx.Step(`^the deploy hash is "([^"]*)"$`, func(apiVersion string) error {

		return utils.Pass
	})

	ctx.Step(`^the deploy hash is "([^"]*)"$`, func(chainSpecName string) error {

		return utils.Pass
	})

	ctx.Step(`^the account is "([^"]*)"$`, func() error {

		return nil
	})

	ctx.Step(`^the timestamp is "([^"]*)"$`, func() error {
		return nil
	})

	ctx.Step(`^the ttl is (\d+)m$`, func() error {
		return nil
	})

	ctx.Step(`^the gas price is (\d+)$`, func() error {
		return nil
	})

	ctx.Step(`^the body_hash is "([^"]*)"$`, func() error {
		return nil
	})

	ctx.Step(`^the chain name is "([^"]*)"$`, func() error {

		return utils.Pass
	})

	ctx.Step(`^dependency (\d+) is "([^"]*)"$`, func() error {
		return nil
	})
	ctx.Step(`^the payment amount is (\d+)$`, func() error {
		return nil
	})

	ctx.Step(`^the session is a transfer$`, func() error {
		return nil
	})

	ctx.Step(`^the session "([^"]*)" is (\d+)$`, func() error {
		return nil
	})

	ctx.Step(`^the session "([^"]*)" is "([^"]*)"$`, func() error {
		return nil
	})

	ctx.Step(`^the session "([^"]*)" type is "([^"]*)"$`, func() error {
		return nil
	})

	ctx.Step(`^the session "([^"]*)" bytes is "([^"]*)"$`, func() error {
		return nil
	})

	ctx.Step(`^the session "([^"]*)" parsed is "([^"]*)"$`, func() error {
		return nil
	})

	ctx.Step(`^the deploy has (\d+) approval$`, func() error {
		return nil
	})

	ctx.Step(`^the approval signer is "([^"]*)"$`, func() error {
		return nil
	})

	ctx.Step(`^the approval signature is "([^"]*)"$`, func() error {
		return nil
	})
}
