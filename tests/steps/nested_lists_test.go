package steps

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/casper-sdks/terminus-go-tests/tests/utils"
	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"math/big"
	"testing"
)

// The test features implementation for the nested_lists.feature
func TestFeaturesNestedLists(t *testing.T) {
	utils.TestFeatures(t, "nested_lists.feature", InitializedNestedLists)
}

func InitializedNestedLists(ctx *godog.ScenarioContext) {

	var clList clvalue.CLValue
	var deploy *types.Deploy
	var sdk casper.RPCClient
	var result rpc.PutDeployResult
	var deployResult casper.InfoGetDeployResult

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		return ctx, nil
	})

	ctx.Step(`a list is created with U(\d+) values of \((\d+), (\d+), (\d+)\)$`,
		func(numLen int, val1 int, val2 int, val3 int) error {

			valUns1 := createUnsignedValue(numLen, val1)
			clList = clvalue.NewCLList(valUns1.Type)
			clList.List.Append(valUns1)
			clList.List.Append(createUnsignedValue(numLen, val2))
			clList.List.Append(createUnsignedValue(numLen, val3))
			return utils.Pass
		},
	)

	ctx.Step(`a nested list is created with U(\d+) values of \(\((\d+), (\d+), (\d+)\),\((\d+), (\d+), (\d+)\)\)$`,
		func(numLen int, val1 int, val2 int, val3 int, val4 int, val5 int, val6 int) error {
			/*valUns1 := createUnsignedValue(numLen, val1)
			clInnerList1 := clvalue.NewCLList(valUns1.Type)
			clInnerList1.List.Append(clInnerList1)
			clInnerList1.List.Append(createUnsignedValue(numLen, val2))
			clInnerList1.List.Append(createUnsignedValue(numLen, val3))

			clInnerList2 := clvalue.NewCLList(valUns1.Type)
			clInnerList2.List.Append(createUnsignedValue(numLen, val4))
			clInnerList2.List.Append(createUnsignedValue(numLen, val5))
			clInnerList2.List.Append(createUnsignedValue(numLen, val6))

			clList = clvalue.NewCLList(clInnerList1.Type)
			clList.List.Append(clInnerList1)
			clList.List.Append(clInnerList2)
			return utils.Pass*/
			return errors.New("nested lists produce \"fatal error: stack overflow\"")
		},
	)

	ctx.Step(`a list is created with "([^"]*)" values of \("([^"]*)", "([^"]*)", "([^"]*)"\)$`,
		func(dataType string, val1 string, val2 string, val3 string) error {
			clVal := createValue(dataType, val1)
			clList = clvalue.NewCLList(clVal.Type)
			clList.List.Append(clVal)
			clList.List.Append(createValue(dataType, val2))
			clList.List.Append(createValue(dataType, val3))

			return utils.Pass

		},
	)

	ctx.Step(`a list is created with I(\d+) values of \((\d+), (\d+), (\d+)\)$`,
		func(numLen int, val1 int, val2 int, val3 int) error {
			valSign1 := createSignedValue(numLen, val1)
			clList = clvalue.NewCLList(valSign1.Type)
			clList.List.Append(valSign1)
			clList.List.Append(createSignedValue(numLen, val2))
			clList.List.Append(createSignedValue(numLen, val3))
			return utils.Pass
		},
	)

	ctx.Step(`the list's bytes are "([^"]*)"$`, func(strHex string) error {
		return utils.ExpectEqual(utils.CasperT, "bytes", hex.EncodeToString(clList.Bytes()), strHex)
	})

	ctx.Step(`the list's length is (\d+)$`, func(len int) error {
		return utils.ExpectEqual(utils.CasperT, "length", clList.List.Len(), len)
	})

	ctx.Step(`the list's "([^"]*)" item is a CLValue with "([^"]*)" value of "([^"]*)"$`,
		func(nth string, valueType string, strValue string) error {
			clVal := getListElement(clList, nth)

			err := utils.ExpectEqual(utils.CasperT, "type", clVal.Type.Name(), valueType)

			if err == nil {
				err = utils.ExpectEqual(utils.CasperT, "value", clVal.String(), strValue)
			}
			return err
		},
	)

	ctx.Step(`the list's "([^"]*)" item is a CLValue with U(\d+) value of (\d+)$`,
		func(nth string, numLen int, value int) error {

			clVal := getListElement(clList, nth)
			err := utils.ExpectEqual(utils.CasperT, "type", clVal.Type.Name(), fmt.Sprintf("U%d", numLen))
			if err == nil {
				err = utils.ExpectEqual(utils.CasperT, "value", clVal.String(), fmt.Sprintf("%d", value))
			}
			return err
		},
	)

	ctx.Step(`the list's "([^"]*)" item is a CLValue with I(\d+) value of (\d+)$`,
		func(nth string, numLen int, value int) error {
			clVal := getListElement(clList, nth)
			err := utils.ExpectEqual(utils.CasperT, "type", clVal.Type.Name(), fmt.Sprintf("I%d", numLen))
			if err == nil {
				err = utils.ExpectEqual(utils.CasperT, "value", clVal.String(), fmt.Sprintf("%d", value))
			}
			return err
		},
	)

	ctx.Step(`the "([^"]*)" nested list's "([^"]*)" item is a CLValue with U(\d+) value of (\d+)$`,
		func(nth string, nestedNth string, numLen int, value string) error {
			return errors.New("no methods exposed to obtain list elements")
		},
	)

	ctx.Step(`that the list is deployed in a transfer$`, func() error {
		var err error

		namedArgs := &types.Args{}
		namedArgs.AddArgument("list", clList)
		deploy, err = utils.BuildStandardTransferDeploy(*namedArgs)

		// Fails here raised issue https://github.com/make-software/casper-go-sdk/issues/70
		result, err = sdk.PutDeploy(context.Background(), *deploy)

		return err
	})

	ctx.Step(`the transfer containing the list is successfully executed$`, func() error {
		var err error
		deployResult, err = utils.WaitForDeploy(result.DeployHash.String(), 300)
		return err
	})

	ctx.Step(`the list is read from the deploy$`, func() error {
		mapVal, err := deployResult.Deploy.Session.Transfer.Args.Find("list")
		if err == nil {
			clList, err = mapVal.Value()
		}
		return err
	})
}

func getListElement(list clvalue.CLValue, nth string) clvalue.CLValue {
	index := int(nth[0]) - int('1')
	return list.List.Elements[index]
}

func createValue(typeName string, strValue string) clvalue.CLValue {
	value, _ := utils.CreateValue(typeName, strValue)
	return *value
}

func createUnsignedValue(numLen int, value int) clvalue.CLValue {
	if numLen == 8 {
		return *clvalue.NewCLUint8(uint8(value))
	} else if numLen == 32 {
		return *clvalue.NewCLUInt32(uint32(value))
	} else if numLen == 64 {
		return *clvalue.NewCLUInt64(uint64(value))
	} else if numLen == 128 {
		return *clvalue.NewCLUInt128(big.NewInt(int64(value)))
	} else if numLen == 256 {
		return *clvalue.NewCLUInt256(big.NewInt(int64(value)))
	} else if numLen == 512 {
		return *clvalue.NewCLUInt512(big.NewInt(int64(value)))
	} else {
		return clvalue.CLValue{}
	}
}

func createSignedValue(numLen int, value int) clvalue.CLValue {
	if numLen == 32 {
		return clvalue.NewCLInt32(int32(value))
	} else if numLen == 64 {
		return *clvalue.NewCLInt64(int64(value))
	} else {
		return clvalue.CLValue{}
	}
}
