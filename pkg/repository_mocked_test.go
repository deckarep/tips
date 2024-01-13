package pkg

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockedDeviceRepo_DevicesResource(t *testing.T) {
	mocked := NewMockedDeviceRepoWithPath("../testmode/devices.json")
	ctx := context.Background()
	cfgCtx := NewConfigCtx()
	ctx = context.WithValue(ctx, CtxKeyConfig, cfgCtx)

	results, err := mocked.DevicesResource(ctx)
	assert.NoError(t, err)

	assert.Equal(t, len(results), 3000)
}
