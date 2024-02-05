package steps

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"

	"github.com/make-software/casper-go-sdk/types/clvalue"

	"github.com/casper-sdks/terminus-go-tests/tests/utils"
	"github.com/cucumber/godog"
)

// The test features implementation for the nested_options.feature
func TestFeaturesNestedOptions(t *testing.T) {
	utils.TestFeatures(t, "nested_options.feature", InitializeNestedOptions)
}

func InitializeNestedOptions(ctx *godog.ScenarioContext) {
	var clOption clvalue.CLValue
	var deploy *types.Deploy
	var sdk casper.RPCClient
	var result rpc.PutDeployResult
	var deployResult casper.InfoGetDeployResult

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		sdk = utils.GetRPCClient()
		return ctx, nil
	})

	ctx.Step(`^that a nested Option has an inner type of Option with a type of String and a value of "([^"]*)"$`, func(value string) error {
		clOption = clvalue.NewCLOption(clvalue.NewCLOption(*clvalue.NewCLString(value)))
		if clOption.Option.IsEmpty() {
			return errors.New("Option is empty")
		}
		return utils.Pass
	})

	ctx.Step(`^the inner type is Option with a type of String and a value of  "([^"]*)"$`, func(value string) error {
		return utils.Pass
	})

	ctx.Step(`^the bytes are "([^"]*)"$`, func(hexBytes string) error {
		decoded, err := hex.DecodeString(hexBytes)

		if !bytes.Equal(clOption.Bytes(), decoded) {
			err = fmt.Errorf("%s bytes do not match expected bytes %s", hex.EncodeToString(clOption.Bytes()), hexBytes)
		}

		return err
	})

	ctx.Step(`^that the nested Option is deployed in a transfer$`, func() error {
		var err error

		namedArgs := &types.Args{}
		namedArgs.AddArgument("option", clOption)
		deploy, err = utils.BuildStandardTransferDeploy(*namedArgs)

		result, err = sdk.PutDeploy(context.Background(), *deploy)

		return err
	})

	ctx.Step(`^the transfer containing the nested Option is successfully executed$`, func() error {
		var err error
		deployResult, err = utils.WaitForDeploy(result.DeployHash.String(), 300)
		return err
	})

	ctx.Step(`^the Option is read from the deploy$`, func() error {
		optionArg, err := deployResult.Deploy.Session.Transfer.Args.Find("option")
		if err == nil {
			clOption, err = optionArg.Value()
		}
		return err
	})

	ctx.Step(`^the inner type is Option with a type of String and a value of "([^"]*)"$`, func(value string) error {
		return utils.Pass
	})

	ctx.Step(`^that a nested Option has an inner type of List with a type of U256 and a value of \((\d+), (\d+), (\d+)\)$`, func(val1 int32, val2 int32, val3 int32) error {
		clList := clvalue.NewCLList(cltype.UInt256)
		clList.List.Append(createCLUInt256(val1))
		clList.List.Append(createCLUInt256(val2))
		clList.List.Append(createCLUInt256(val3))
		clOption = clvalue.NewCLOption(clList)
		return utils.Pass
	})

	ctx.Step(`^the list's length is (\d+)$`, func(len int) error {
		return utils.ExpectEqual(utils.CasperT, "length", clOption.Option.Inner.List.Len(), len)
	})

	ctx.Step(`^the list's "([^"]*)" item is a CLValue with U256 value of (\d+)$`, func(nth string, val int32) error {
		index := getNthIndex(nth)
		clValue := clOption.Option.Inner.List.Elements[index]
		return utils.ExpectEqual(utils.CasperT, "value", clValue.UI256.Value(), createCLUInt256(val).UI256.Value())
	})

	ctx.Step(`^that a nested Option has an inner type of Tuple2 with a type of "([^"]*)" values of \("([^"]*)", (\d+)\)$`, func(types string, val1 string, val2 int32) error {
		clTuple2 := clvalue.NewCLTuple2(*clvalue.NewCLString(val1), createCLUInt256(val2))
		clOption = clvalue.NewCLOption(clTuple2)
		return utils.Pass
	})

	ctx.Step(`^the inner type is Tuple2 with a type of "([^"]*)" and a value of \("([^"]*)", (\d+)\)$`, func(types string, val1 string, val2 int32) error {
		err := utils.ExpectEqual(utils.CasperT, "1st tuple", clOption.Option.Value().Tuple2.Inner1.String(), val1)
		if err == nil {
			err = utils.ExpectEqual(utils.CasperT, "2nd tuple", clOption.Option.Value().Tuple2.Inner2.UI256.Value(), createCLUInt256(val2).UI256.Value())
		}
		return err
	})

	ctx.Step(`^that a nested Option has an inner type of Map with a type of "([^"]*)" value of \{"([^"]*)": (\d+)\}$`, func(types string, key string, val int32) error {
		innerClmap := clvalue.NewCLMap(cltype.String, cltype.UInt32)
		err := innerClmap.Map.Append(*clvalue.NewCLString(key), *clvalue.NewCLUInt32(uint32(val)))
		if err == nil {
			clOption = clvalue.NewCLOption(innerClmap)
		}

		return err
	})

	ctx.Step(`^the inner type is Map with a type of "([^"]*)" and a value of \{"([^"]*)": (\d+)\}$`, func(types string, key string, val int32) error {
		foundVal, b := clOption.Option.Value().Map.Find(key)
		if b {
			return utils.ExpectEqual(utils.CasperT, "map", foundVal.UI256.Value(), createCLUInt256(val).UI256.Value())
		} else {
			return errors.New("key not found")
		}
	})

	ctx.Step(`^that a nested Option has an inner type of Any with a value of "([^"]*)"`, func(value string) error {
		decoded, err := hex.DecodeString(value)
		clOption = clvalue.NewCLOption(*clvalue.NewAnyFromBytes(decoded))
		return err
	})

	ctx.Step(`^the inner type is Any with a value of "([^"]*)"$`, func(value string) error {
		decoded, err := hex.DecodeString(value)
		var actual []byte
		if err != nil {
			actual = clOption.Option.Value().Any.Bytes()
		}
		return utils.ExpectEqual(utils.CasperT, "value", actual, decoded)
	})
}

func createCLUInt256(val1 int32) clvalue.CLValue {
	bi := new(big.Int)
	bi.SetUint64(uint64(val1))
	clVal := *clvalue.NewCLUInt256(bi)
	return clVal
}

func getNthIndex(nth string) int {
	return int(nth[0]) - int('1')
}
