package utils

import (
	"encoding/hex"
	"fmt"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/key"
	"math/big"
	"strconv"
)

func CreateValue(typeName string, strValue string) (*clvalue.CLValue, error) {
	var clVal = clvalue.CLValue{}
	var err error = nil

	switch typeName {

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
		bytes, err = hex.DecodeString(strValue)
		uRef, err = key.NewURefFromBytes(bytes)
		if err == nil {
			clVal = clvalue.NewCLUref(uRef)
		}

	default:
		err = fmt.Errorf("invalid type %s or value %s", typeName, strValue)
	}
	return &clVal, err
}
