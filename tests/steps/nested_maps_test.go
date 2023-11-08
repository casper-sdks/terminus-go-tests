package steps

import (
	"context"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
)

// The test features implementation for the info_get_validator_changes.feature
func TestFeaturesNestedMaps(t *testing.T) {
	utils.TestFeatures(t, "nested_maps.feature", InitializeNestedMaps)
}

func InitializeNestedMaps(ctx *godog.ScenarioContext) {
	//var sdk casper.RPCClient
	//var validatorChanges rpc.InfoGetValidatorChangesResult

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		//sdk = utils.GetRPCClient()
		return ctx, nil
	})

	ctx.Step(`^a map is created \{"([^"]*)": (\d+)\}$`, func(key string, value int64) error {
		return nil
	})

	ctx.Step(`^a nested map is created \{"([^"]*)": \{"([^"]*)": (\d+)}, "([^"]*)": \{"([^"]*)", (\d+)}}$`,
		func(key0 string, key1 string, value1 int, key2 string, key3 string, value3 int) error {
			return nil
		},
	)

	ctx.Step(`^a map is created \{"([^"]*)": (\d+), "([^"]*)": \{"([^"]*)": (\d+), "([^"]*)": \{"([^"]*)": (\d+)}}}}$`,
		func(key0 string,
			value0 int,
			key1 string,
			key2 string,
			value2 int,
			key3 string,
			key4 string,
			value4 int) error {
			return nil
		},
	)

	ctx.Step(`^a nested map is created  \{(\d+): \{(\d+): \{(\d+): "([^"]*)"}, (\d+): \{(\d+): "([^"]*)"}}, (\d+): \{(\d+): \{(\d+): "([^"]*)"}, (\d+): \{(\d+): "([^"]*)"}}}$`,
		func(key1 int,
			key11 int,
			key111 int,
			value111 string,
			key12 int,
			key121 int,
			value121 string,
			key2 int,
			key21 int,
			key211 int,
			value211 string,
			key22 int,
			key221 int,
			value221 string) error {
			return nil
		},
	)

	ctx.Step(`the map's key type is "([^"]*)" and the maps value type is "([^"]*)"$`, func(key string, typeName string) error {
		return utils.Pass
	})

	ctx.Step(`the map's bytes are "([^"]*)"$`, func(strHex string) error {
		return utils.Pass
	})
	ctx.Step(`that the nested map is deployed in a transfer$`, func() error {
		return utils.Pass
	})
	ctx.Step(`the transfer containing the nested map is successfully executed$`, func() error {
		return utils.Pass
	})
	ctx.Step(`the map is read from the deploy$`, func() error {
		return utils.Pass
	})

	ctx.Step(`the map's key is "([^"]*)" and value is "([^"]*)"$`, func(key string, strValue string) error {
		return utils.Pass
	})
	ctx.Step(`the 1st nested map's key is "([^"]*)" and value is "([^"]*)"$`, func(key string, strValue string) error {
		return utils.Pass
	})

}
