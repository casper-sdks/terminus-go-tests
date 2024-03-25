package utils

import (
	"encoding/hex"
	"fmt"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
	"github.com/make-software/casper-go-sdk/types/key"
	"math/big"
	"strconv"
)

func CreateValue(typeName string, strValue string) (*clvalue.CLValue, error) {
	var clVal = clvalue.CLValue{}
	var err error = nil

	switch typeName {

	case "Any":
		var bytes []byte
		bytes, err = hex.DecodeString(strValue)
		clVal = clvalue.NewCLAny(bytes)

	case "Bool":
		var b bool
		b, err = strconv.ParseBool(strValue)
		clVal = clvalue.NewCLBool(b)

	case "U8":
		var u64 uint64
		u64, err = strconv.ParseUint(strValue, 10, 8)
		clVal = *clvalue.NewCLUint8(uint8(u64))

	case "U32":
		var u64 uint64
		u64, err = strconv.ParseUint(strValue, 10, 32)
		clVal = *clvalue.NewCLUInt32(uint32(u64))

	case "U64":
		var u64 uint64
		u64, err = strconv.ParseUint(strValue, 10, 64)
		clVal = *clvalue.NewCLUInt64(u64)

	case "U128":
		var bi = new(big.Int)
		bi.SetString(strValue, 10)
		clVal = *clvalue.NewCLUInt128(bi)

	case "U256":
		var bi = new(big.Int)
		bi.SetString(strValue, 10)
		clVal = *clvalue.NewCLUInt256(bi)

	case "U512":
		var bi = new(big.Int)
		bi.SetString(strValue, 10)
		clVal = *clvalue.NewCLUInt512(bi)

	case "I32":
		var i64 int64
		i64, err = strconv.ParseInt(strValue, 10, 32)
		clVal = clvalue.NewCLInt32(int32(i64))

	case "I64":
		var i64 int64
		i64, err = strconv.ParseInt(strValue, 10, 64)
		clVal = *clvalue.NewCLInt64(i64)

	case "String":
		clVal = *clvalue.NewCLString(strValue)

	case "ByteArray":
		var ba []byte
		ba, err = hex.DecodeString(strValue)
		clVal = clvalue.NewCLByteArray(ba)

	case "Key":
		var ba []byte
		ba, err = hex.DecodeString(strValue)
		if err == nil {
			var clKey = casper.Key{}
			err = clKey.Scan(ba)
			clVal = clvalue.NewCLKey(clKey)
		}

	case "PublicKey":
		var ba []byte
		ba, err = hex.DecodeString(strValue)
		if err == nil {
			var clKey = casper.PublicKey{}
			err = clKey.Scan(ba)
			clVal = clvalue.NewCLPublicKey(clKey)
		}

	case "URef":
		uRef := key.URef{}
		var bytes []byte
		bytes, err = hex.DecodeString(strValue + "07")
		uRef, err = key.NewURefFromBytes(bytes)
		if err == nil {
			clVal = clvalue.NewCLUref(uRef)
		}

	default:
		err = fmt.Errorf("invalid type %s for value %s", typeName, strValue)
	}
	return &clVal, err
}

func CreateComplexValue(typeName string, innerTypes []string, strValues []string) (*clvalue.CLValue, error) {
	var clVal = clvalue.CLValue{}
	var innerValues []clvalue.CLValue

	var err error = nil
	innerValues, err = createInnerValues(strValues, innerTypes)
	if err != nil {
		return nil, err
	}

	switch typeName {

	case "Option":
		clVal = clvalue.NewCLOption(innerValues[0])

	case "Tuple1":
		clVal = clvalue.NewCLTuple1(innerValues[0])

	case "Tuple2":
		clVal = clvalue.NewCLTuple2(innerValues[0], innerValues[1])

	case "Tuple3":
		clVal = clvalue.NewCLTuple3(innerValues[0], innerValues[1], innerValues[2])

	case "List":
		clList := clvalue.NewCLList(innerValues[0].Type)

		for _, innerValue := range innerValues {
			if innerTypes[0] == "String" {
				clList.List.Append(innerValue)
			} else {
				clList.List.Append(clvalue.NewCLByteArray(innerValue.Bytes()))
			}
		}
		clVal = clList

	case "Map":
		clMap := clvalue.NewCLMap(cltype.String, cltype.String)
		var mapKey *clvalue.CLValue
		for i, innerValue := range innerValues {
			formatInt := strconv.FormatInt(int64(i), 10)
			mapKey, err = CreateValue("String", formatInt)
			if err != nil {
				break
			}
			err = clMap.Map.Append(*mapKey, innerValue)
			if err != nil {
				break
			}
		}
		clVal = clMap

	default:
		err = fmt.Errorf("invalid type %s with innerTypes %s and values %s", typeName, innerTypes, strValues)
	}

	return &clVal, err
}

func createInnerValues(strValues []string, innerTypes []string) ([]clvalue.CLValue, error) {

	var innerValues []clvalue.CLValue

	for i, innerValue := range strValues {
		innerType := innerTypes[i]
		value, err := CreateValue(innerType, innerValue)
		if err != nil {
			return nil, err
		}
		innerValues = append(innerValues, *value)
	}
	return innerValues, nil
}
