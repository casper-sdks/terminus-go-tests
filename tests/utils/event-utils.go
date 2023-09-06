package utils

import (
	"context"
	"fmt"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/sse"
	"strings"
	"time"
)

func WaitForDeploy(deployHash string, timeoutSeconds int) (casper.InfoGetDeployResult, error) {

	sdk := GetSdk()

	var timeout = int64(timeoutSeconds*1000) + time.Now().UnixMilli()
	var deploy = casper.InfoGetDeployResult{}
	var err error

	for len(deploy.ExecutionResults) == 0 {

		deploy, err = sdk.GetDeploy(context.Background(), deployHash)

		if len(deploy.ExecutionResults) == 0 && time.Now().UnixMilli() > timeout {
			deploy = casper.InfoGetDeployResult{}
			return deploy, fmt.Errorf("timed-out waiting for deploy hash %s", deployHash)
		}

		if err != nil {
			return deploy, err
		}
	}

	return deploy, err
}

func WaitForBlockAdded(deployHash string, timeoutSeconds int) (sse.BlockAddedEvent, error) {

	var blockAddedEvent sse.BlockAddedEvent
	err := Pass
	sseClient := GetSse()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	sseClient.RegisterHandler(sse.BlockAddedEventType, func(ctx context.Context, event sse.RawEvent) error {
		if event.EventType == sse.BlockAddedEventType {
			blockAddedEvent, err = event.ParseAsBlockAddedEvent()

			if len(blockAddedEvent.BlockAdded.Block.Body.TransferHashes) > 0 && blockAddedEvent.BlockAdded.Block.Body.TransferHashes[0].String() == deployHash {
				// Cancel so we stop listening
				cancel()
			} else {
				// Clear the event
				blockAddedEvent = sse.BlockAddedEvent{}
			}

		}
		return Pass
	})

	err = sseClient.Start(ctx, 0)

	if strings.Contains(fmt.Sprint(err), "context canceled") {
		// We invoked the cancel so clear it
		err = nil
	}

	return blockAddedEvent, err
}
