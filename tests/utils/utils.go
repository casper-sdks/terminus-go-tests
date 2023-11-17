package utils

import (
	"encoding/json"
	"fmt"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/keypair"
	"github.com/stretchr/testify/assert"
	yml "gopkg.in/yaml.v2"
	"log"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

var (
	_, b, _, _        = runtime.Caller(0)
	root              = filepath.Join(filepath.Dir(b), "../..")
	config            map[string]interface{}
	Pass              error = nil
	NotImplementError       = fmt.Errorf("Not Implemented.")
)

func ReadConfig() {
	f, err := os.ReadFile(root + "/config.yml")
	if err != nil {
		log.Fatal(err)
	}
	err = yml.Unmarshal(f, &config)
	if err != nil {
		log.Fatal(err)
	}
}

func AssertExpectedAndActual(a expectedAndActualAssertion, expected, actual interface{}) error {
	var t asserter
	a(&t, expected, actual)
	return t.err
}

type asserter struct {
	err error
}

func (a *asserter) Errorf(format string, args ...interface{}) {
	a.err = fmt.Errorf(format, args...)
}

type expectedAndActualAssertion func(t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool

func GetUserKeyAssetPath(networkId int, userId int, keyFilename string) string {
	return fmt.Sprintf("../../assets/net-%d/user-%d/%s", networkId, userId, keyFilename)
}

func ExpectEqual(t *testing.T, attribute string, actual any, expected any) error {

	if !assert.Equal(t, expected, actual) {
		return CreateExpectError(attribute, actual, expected)
	} else {
		return Pass
	}
}

func CreateExpectError(attribute string, actual any, expected any) error {
	return fmt.Errorf("%s expected %s to be %s", attribute, actual, expected)
}

func GetRpcError(err interface{}) rpc.RpcError {
	rpcError := err.(*rpc.RpcError)
	fmt.Println(rpcError.Code)
	return *rpcError
}

func BuildStandardTransferDeploy(namedArgs types.Args) (*types.Deploy, error) {
	var deploy *types.Deploy

	keyPath := GetUserKeyAssetPath(1, 1, "secret_key.pem")

	senderKey, err := casper.NewED25519PrivateKeyFromPEMFile(keyPath)
	if err != nil {
		return deploy, err
	}

	keyPath = GetUserKeyAssetPath(1, 2, "secret_key.pem")

	var receiverPrivateKey keypair.PrivateKey
	receiverPrivateKey, err = casper.NewED25519PrivateKeyFromPEMFile(keyPath)
	if err != nil {
		return deploy, err
	}

	header := types.DefaultHeader()
	header.ChainName = "casper-net-1"
	header.Account = senderKey.PublicKey()
	header.Timestamp = types.Timestamp(time.Now())
	payment := types.StandardPayment(big.NewInt(100000000))

	args := &types.Args{}
	args.AddArgument("amount", *clvalue.NewCLUInt512(big.NewInt(2500000000)))
	args.AddArgument("target", clvalue.NewCLPublicKey(receiverPrivateKey.PublicKey()))
	args.AddArgument("id", clvalue.NewCLOption(*clvalue.NewCLUInt64(rand.Uint64())))

	var name string
	var val clvalue.CLValue

	for _, arg := range namedArgs {
		name, err = arg.Name()
		if err != nil {
			return deploy, err
		}
		val, err = arg.Value()
		if err != nil {
			return deploy, err
		}
		args.AddArgument(name, val)
	}

	session := types.ExecutableDeployItem{
		Transfer: &types.TransferDeployItem{
			Args: *args,
		},
	}

	deploy, err = types.MakeDeploy(header, payment, session)
	if err != nil {
		return deploy, err
	}

	assert.NotNil(CasperT, deploy, "deploy")

	err = deploy.SignDeploy(senderKey)

	if err != nil {
		return deploy, err
	}

	deployJson, err := json.Marshal(deploy)
	if err != nil {
		return nil, err
	}

	assert.NotNil(CasperT, deployJson)
	fmt.Println(string(deployJson))

	return deploy, err
}
