package steps

import (
	"context"
	"fmt"
	"testing"

	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/stretchr/testify/assert"

	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
)

// The test features implementation for the info_get_validator_changes.feature
func TestFeaturesInfoGetValidatorChanges(t *testing.T) {
	utils.TestFeatures(t, "info_get_validator_changes.feature", InitializeInfoGetValidatorChanges)
}

func InitializeInfoGetValidatorChanges(ctx *godog.ScenarioContext) {
	var sdk casper.RPCClient
	var validatorChanges rpc.InfoGetValidatorChangesResult

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		sdk = utils.GetRPCClient()
		return ctx, nil
	})

	ctx.Step(`^that the info_get_validator_changes method is invoked against a node$`, func() error {
		var err error

		validatorChanges, err = sdk.GetValidatorChangesInfo(context.Background())

		return err
	})

	ctx.Step(`^a valid info_get_validator_changes_result is returned$`, func() error {
		assert.NotNil(utils.CasperT, validatorChanges, "validatorChanges is nil")

		assert.NotNil(utils.CasperT, validatorChanges.Changes)

		// TODO need to investigate how to ensure changes exist as we have none on nctl startup

		return utils.Pass
	})

	ctx.Step(`the info_get_validator_changes_result contains a valid API version$`, func() error {
		apiVersion := "1.0.0"

		if apiVersion != validatorChanges.APIVersion {
			return fmt.Errorf("expected %s ApiVersion to be %s", validatorChanges.APIVersion, apiVersion)
		}
		return utils.Pass
	})
}
