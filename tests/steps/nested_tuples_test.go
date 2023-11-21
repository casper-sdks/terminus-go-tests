package steps

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
	"strconv"
	"strings"
	"testing"
)

var tuple1 clvalue.CLValue
var tuple2 clvalue.CLValue
var tuple3 clvalue.CLValue

const first = "first"
const second = "second"
const third = "third"

// The test features implementation for the nested_tuples.feature
func TestFeaturesNestedTuples(t *testing.T) {
	utils.TestFeatures(t, "nested_tuples.feature", InitializeNestedTuples)
}

func InitializeNestedTuples(ctx *godog.ScenarioContext) {

	var deploy *types.Deploy
	var sdk casper.RPCClient
	var result rpc.PutDeployResult
	var deployResult casper.InfoGetDeployResult

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		return ctx, nil
	})

	ctx.Step(`that a nested Tuple(\d+) is defined as \(\((\d+)\)\) using U32 numeric values$`,
		func(index int, value int) error {
			if index == 1 {
				tuple1 = clvalue.NewCLTuple1(clvalue.NewCLTuple1(*clvalue.NewCLUInt32(uint32(value))))
			}
			return utils.Pass
		},
	)

	ctx.Step(`that a nested Tuple(\d+) is defined as \((\d+), \((\d+), \((\d+), (\d+)\)\)\) using U32 numeric values$`,
		func(index int, value1 int, value2 int, value3 int, value4 int) error {

			innerTuple2 := clvalue.NewCLTuple2(*clvalue.NewCLUInt32(uint32(value3)), *clvalue.NewCLUInt32(uint32(value4)))
			innerTuple1 := clvalue.NewCLTuple2(*clvalue.NewCLUInt32(uint32(value2)), innerTuple2)
			tuple2 = clvalue.NewCLTuple2(*clvalue.NewCLUInt32(uint32(value1)), innerTuple1)
			return utils.Pass

		},
	)

	ctx.Step(`that a nested Tuple(\d+) is defined as \((\d+), (\d+), \((\d+), (\d+), \((\d+), (\d+), (\d+)\)\)\) using U32 numeric values$`,
		func(index int, value1 int, value2 int, value3 int, value4 int, value5 int, value6 int, value7 int) error {

			innerTuple2 := clvalue.NewCLTuple3(*clvalue.NewCLUInt32(uint32(value5)), *clvalue.NewCLUInt32(uint32(value6)), *clvalue.NewCLUInt32(uint32(value7)))
			innerTuple1 := clvalue.NewCLTuple3(*clvalue.NewCLUInt32(uint32(value3)), *clvalue.NewCLUInt32(uint32(value4)), innerTuple2)
			tuple3 = clvalue.NewCLTuple3(*clvalue.NewCLUInt32(uint32(value1)), *clvalue.NewCLUInt32(uint32(value2)), innerTuple1)
			return utils.Pass
		},
	)

	ctx.Step(`^the "([^"]*)" element of the Tuple(\d+) is "([^"]*)"$`,
		func(nth string, tuple int, strValue string) error {

			value, err := getTupleValue(tuple, nth)
			if err == nil {
				//byType := value.GetValueByType().String()
				actualValues := getTupleValues(value)
				expectedValues := getExpectedTupleValues(strValue)
				err = utils.ExpectEqual(utils.CasperT, fmt.Sprintf("tuple %d %s", tuple, nth), actualValues, expectedValues)
			}
			return nil // err
		},
	)
	ctx.Step(`^the "([^"]*)" element of the Tuple(\d+) is (\d+)$`,
		func(nth string, tuple int, expectedValue int) error {
			value, err := getTupleValue(tuple, nth)
			if err == nil {
				actualValue := value.UI32.Value()
				err = utils.ExpectEqual(utils.CasperT, fmt.Sprintf("tuple %d %s", tuple, nth), actualValue, uint32(expectedValue))
			}
			return err
		},
	)

	ctx.Step(`the Tuple(\d+) bytes are "([^"]*)"$`, func(tupleIndex int, strHex string) error {
		return utils.ExpectEqual(utils.CasperT, "bytes", hex.EncodeToString(getTupleBytes(tupleIndex)), strHex)
	})

	ctx.Step(`that the nested tuples are deployed in a transfer$`, func() error {
		var err error

		namedArgs := &types.Args{}
		namedArgs.AddArgument("TUPLE_1", tuple1)
		namedArgs.AddArgument("TUPLE_2", tuple2)
		namedArgs.AddArgument("TUPLE_3", tuple3)
		deploy, err = utils.BuildStandardTransferDeploy(*namedArgs)

		result, err = sdk.PutDeploy(context.Background(), *deploy)

		return err
	})

	ctx.Step(`the transfer is successful$`, func() error {
		var err error
		deployResult, err = utils.WaitForDeploy(result.DeployHash.String(), 300)
		return err
	})

	ctx.Step(`the tuples deploy is obtained from the node$`, func() error {
		var err error = nil
		tuple1, err = getTupleArgument(deployResult.Deploy.Session.Transfer.Args, "TUPLE_1")

		if err == nil {
			tuple2, err = getTupleArgument(deployResult.Deploy.Session.Transfer.Args, "TUPLE_2")
		}

		if err == nil {
			tuple3, err = getTupleArgument(deployResult.Deploy.Session.Transfer.Args, "TUPLE_3")
		}

		return err
	})
}

func getExpectedTupleValues(value string) []uint32 {
	result := make([]uint32, 0)

	value = strings.ReplaceAll(value, "(", "")
	value = strings.ReplaceAll(value, ")", "")
	for _, one := range strings.Split(value, ",") {
		one = strings.TrimSpace(one)
		val, _ := strconv.Atoi(one)
		result = append(result, uint32(val))
	}
	return result
}

func getTupleValues(tuple clvalue.CLValue) []uint32 {

	return populateTupleValues(tuple, make([]uint32, 0))
}

func populateTupleValues(value clvalue.CLValue, values []uint32) []uint32 {
	if value.Type.Name() == "U32" {
		values = append(values, value.UI32.Value())
	} else if value.Type.Name() == "Tuple1" {
		values = populateTupleValues((*value.Tuple1).Value(), values)
	} else if value.Type.Name() == "Tuple2" {
		values = populateTupleValues((*value.Tuple2).Inner1, values)
		values = populateTupleValues((*value.Tuple2).Inner2, values)
	} else if value.Type.Name() == "Tuple3" {
		values = populateTupleValues((*value.Tuple3).Inner1, values)
		values = populateTupleValues((*value.Tuple3).Inner2, values)
		values = populateTupleValues((*value.Tuple3).Inner3, values)
	}
	return values
}

func getTupleBytes(index int) []byte {
	if index == 1 {
		return tuple1.Bytes()
	} else if index == 2 {
		return tuple2.Bytes()
	} else if index == 3 {
		return tuple3.Bytes()
	} else {
		return nil
	}
}

func getTupleArgument(args types.Args, name string) (clvalue.CLValue, error) {
	tuple1Val, err := args.Find(name)
	if err == nil {
		return tuple1Val.Value()
	}
	return clvalue.CLValue{}, err
}

func getTupleValue(tuple int, nth string) (clvalue.CLValue, error) {
	if tuple == 1 {
		if first == nth {
			return (*tuple1.Tuple1).Value(), nil
		}
	} else if tuple == 2 {
		if first == nth {
			return tuple2.Tuple2.Value()[0], nil
		} else if second == nth {
			return tuple2.Tuple2.Value()[1], nil
		}
	} else if tuple == 3 {
		if first == nth {
			return tuple3.Tuple3.Value()[0], nil
		} else if second == nth {
			return tuple3.Tuple3.Value()[1], nil
		} else if third == nth {
			return tuple3.Tuple3.Value()[2], nil
		}
	}

	return clvalue.CLValue{}, errors.New("tuple not found")
}
