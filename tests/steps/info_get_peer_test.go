package steps

import (
	"context"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/stormeye2000/cspr-sdk-standard-tests-go/tests/utils"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

// The test features implementation for the info_get_peers.feature
func TestFeaturesInfoGetPeer(t *testing.T) {
	TestFeatures(t, "info_get_peer.feature", InitializeInfoGetPeers)
}

func InitializeInfoGetPeers(ctx *godog.ScenarioContext) {

	var sdk casper.RPCClient
	var infoGetPeersResult rpc.InfoGetPeerResult

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		utils.ReadConfig()
		sdk = utils.GetSdk()
		return ctx, nil
	})

	ctx.Step(`^that the info_get_peers RPC method is invoked against a node$`, func() error {

		result, err := sdk.GetPeers(context.Background())

		infoGetPeersResult = result

		return err
	})

	ctx.Step(`^the node returns an info_get_peers_result`, func() error {

		assert.NotNil(CasperT, infoGetPeersResult, "infoGetPeersResult is nil")

		return utils.Pass
	})

	ctx.Step(`^the info_get_peers_result has an API version of "([^"]*)"$`, func(apiVersion string) error {

		if apiVersion != infoGetPeersResult.ApiVersion {
			return fmt.Errorf("expected %s ApiVersion to be %s", infoGetPeersResult.ApiVersion, apiVersion)
		}

		return utils.Pass
	})

	ctx.Step(`^the info_get_peers_result contains (\d+) peers$`, func(peerCount int) error {
		if peerCount != len(infoGetPeersResult.Peers) {
			return fmt.Errorf("expected infoGetPeersResult.Peers length to be %d", peerCount)
		}
		return utils.Pass
	})

	ctx.Step(`^the info_get_peers_result contains a valid peer with a port number of (\d+)$`, func(portNumber int) error {

		var found = false

		for i := 0; i < len(infoGetPeersResult.Peers); i++ {

			port := strings.Split(infoGetPeersResult.Peers[i].Address, ":")[1]

			actualPort, err := strconv.Atoi(port)

			if err != nil {
				return err
			}

			if actualPort == portNumber {
				found = true
			}
		}

		if !found {
			return fmt.Errorf("expected port number %d was not ", portNumber)
		}

		return utils.Pass
	})
}
