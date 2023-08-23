package steps

import (
	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"testing"
)

type _map struct {
	blockDataNode casper.Block
	blockDataSdk  rpc.ChainGetBlockResult
}

var contextMap _map

var t *testing.T

func TestFeatures(t *testing.T, featureName string, scenarioInitializer func(*godog.ScenarioContext)) {
	suite := godog.TestSuite{
		ScenarioInitializer: scenarioInitializer,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features/" + featureName},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
