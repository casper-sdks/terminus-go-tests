package steps

import (
	"context"
	"encoding/hex"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	cltype "github.com/make-software/casper-go-sdk/types/clvalue/cltype"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
)

// The test features implementation for the nested_maps.feature
func TestFeaturesNestedMaps(t *testing.T) {
	utils.TestFeatures(t, "nested_maps.feature", InitializeNestedMaps)
}

func InitializeNestedMaps(ctx *godog.ScenarioContext) {
	var clMap clvalue.CLValue

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		//sdk = utils.GetRPCClient()
		return ctx, nil
	})

	ctx.Step(`^a map is created \{"([^"]*)": (\d+)\}$`, func(key string, value int64) error {
		clMap = clvalue.NewCLMap(cltype.String, cltype.UInt32)
		return clMap.Map.Append(*clvalue.NewCLString(key), *clvalue.NewCLUInt32(uint32(value)))
	})

	ctx.Step(`^a nested map is created \{"([^"]*)": \{"([^"]*)": (\d+)}, "([^"]*)": \{"([^"]*)", (\d+)}}$`,
		func(key0 string, key1 string, value1 int, key2 string, key3 string, value3 int) error {

			// Fail SDK does not allow the creation of nested maps as the Map does not implement the CLValue interface
			/*clMap = clvalue.NewCLMap(cltype.String, cltype.Map)

			innerClMap1 := clvalue.NewCLMap(cltype.String, cltype.UInt32)
			err := innerClMap1.Map.Append(*clvalue.NewCLString(key1), *clvalue.NewCLUInt32(uint32(value1)))

			innerClMap2 := clvalue.NewCLMap(cltype.String, cltype.UInt32)
			err = innerClMap2.Map.Append(*clvalue.NewCLString(key3), *clvalue.NewCLUInt32(uint32(value3)))

			err = clMap.Map.Append(*clvalue.NewCLString(key0), innerClMap1)
			err = clMap.Map.Append(*clvalue.NewCLString(key2), innerClMap2)*/
			return utils.NotImplementError
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
			return utils.NotImplementError
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
			return utils.NotImplementError
		},
	)

	ctx.Step(`the map's key type is "([^"]*)" and the maps value type is "([^"]*)"$`, func(keyType string, valueType string) error {

		err := utils.ExpectEqual(utils.CasperT, "keyType", clMap.Map.Type.Key.Name(), keyType)
		if err != nil {
			err = utils.ExpectEqual(utils.CasperT, "valueType", clMap.Map.Type.Val.Name(), valueType)
		}
		return err
	})

	ctx.Step(`the map's bytes are "([^"]*)"$`, func(strHex string) error {
		return utils.ExpectEqual(utils.CasperT, "bytes", hex.EncodeToString(clMap.Map.Bytes()), strHex)
	})

	ctx.Step(`that the nested map is deployed in a transfer$`, func() error {
		return utils.NotImplementError
	})

	ctx.Step(`the transfer containing the nested map is successfully executed$`, func() error {
		return utils.NotImplementError
	})

	ctx.Step(`the map is read from the deploy$`, func() error {
		return utils.NotImplementError
	})

	ctx.Step(`the map's key is "([^"]*)" and value is "([^"]*)"$`, func(key string, strValue string) error {
		return utils.NotImplementError
	})

	ctx.Step(`the 1st nested map's key is "([^"]*)" and value is "([^"]*)"$`, func(key string, strValue string) error {
		return utils.NotImplementError
	})
}
