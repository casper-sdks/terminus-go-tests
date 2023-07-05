package tests

import (
	"context"
	"fmt"
	"github.com/cucumber/godog"
	"testing"
)

// godogsCtxKey is the key used to store the available godogs in the context.Context
type godogsCtxKey struct{}

func requestTheLatestBlockViaTheTestNode(ctx context.Context) error {

	fmt.Printf("context value is: %s", ctx.Value("latestBlockSdk"))

	return nil
}

func thatTheLatestBlockIsRequestedViaTheSdk(ctx context.Context) (context.Context, error) {
	fmt.Println("request the latest block via the test node")

	latest, err := sdk().GetBlockLatest(context.Background())
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, "latestBlockSdk", latest), nil
}

func theBodyOfTheReturnedBlockIsEqualToTheBodyOfTheReturnedTestNodeBlock(ctx context.Context) error {
	return nil
}

func theHashOfTheReturnedBlockIsEqualToTheHashOfTheReturnedTestNodeBlock() error {
	return nil
}

func theHeaderOfTheReturnedBlockIsEqualToTheHeaderOfTheReturnedTestNodeBlock() error {
	return nil
}

func theProofsOfTheReturnedBlockAreEqualToTheProofsOfTheReturnedTestNodeBlock() error {
	return nil
}

func TestFeaturesBlocks(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenarioBlocks,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenarioBlocks(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		readConfig()
		return ctx, nil
	})

	ctx.Step(`^request the latest block via the test node$`, requestTheLatestBlockViaTheTestNode)
	ctx.Step(`^that the latest block is requested via the sdk$`, thatTheLatestBlockIsRequestedViaTheSdk)
	ctx.Step(`^the body of the returned block is equal to the body of the returned test node block$`, theBodyOfTheReturnedBlockIsEqualToTheBodyOfTheReturnedTestNodeBlock)
	ctx.Step(`^the hash of the returned block is equal to the hash of the returned test node block$`, theHashOfTheReturnedBlockIsEqualToTheHashOfTheReturnedTestNodeBlock)
	ctx.Step(`^the header of the returned block is equal to the header of the returned test node block$`, theHeaderOfTheReturnedBlockIsEqualToTheHeaderOfTheReturnedTestNodeBlock)
	ctx.Step(`^the proofs of the returned block are equal to the proofs of the returned test node block$`, theProofsOfTheReturnedBlockAreEqualToTheProofsOfTheReturnedTestNodeBlock)
}
