package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	yml "gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	_, b, _, _        = runtime.Caller(0)
	root              = filepath.Join(filepath.Dir(b), "../..")
	config            map[string]interface{}
	Pass              error = nil
	NotImplementError       = fmt.Errorf("not Implemented.")
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

func Result(err error) error {
	if err != nil {
		return err
	} else {
		return nil
	}
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
