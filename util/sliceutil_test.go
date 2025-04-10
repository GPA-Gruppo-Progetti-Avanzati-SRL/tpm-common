package util_test

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	SizeOfBatch = 10000
)

type BatchOfEvents struct {
	Events   [SizeOfBatch]string
	Capacity int
	Size     int
}

func TestSlice(t *testing.T) {

	// sliceOfSize := make([]string, 0, 10)
	events := BatchOfEvents{Capacity: SizeOfBatch}
	err := underLyingFunction(&events, 50)
	require.NoError(t, err)

	t.Log(events.Size)
}

func underLyingFunction(batchOfEvents *BatchOfEvents, maxItems int) error {

	_, err := util.NewErrorRandomizer("1/M")
	if err != nil {
		return err
	}

	if batchOfEvents.Capacity-batchOfEvents.Size < maxItems {
		maxItems = batchOfEvents.Capacity - batchOfEvents.Size
	}

	for i := batchOfEvents.Size; i < maxItems; i++ {
		batchOfEvents.Events[i] = fmt.Sprintf("hello-%d", i)
		batchOfEvents.Size++
	}

	return nil
}
