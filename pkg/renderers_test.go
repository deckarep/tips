package pkg

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRenderRemoteSummary(t *testing.T) {
	ctx := context.Background()

	var b bytes.Buffer
	err := RenderRemoteSummary(ctx, &b, 2, 0, time.Millisecond*333)
	assert.NoError(t, err, "RenderRemoteSummary should have returned no error")

	assert.Equal(t, b.String(), "Finished: successes: 2, failures: 0, elapsed (secs): 0.33\n")

	b.Reset()
	err = RenderRemoteSummary(ctx, &b, 0, 3, time.Millisecond*777)
	assert.NoError(t, err, "RenderRemoteSummary should have returned no error")

	assert.Equal(t, b.String(), "Finished: successes: 0, failures: 3, elapsed (secs): 0.78\n")
}
