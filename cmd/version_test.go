package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintVersion(t *testing.T) {
	var b bytes.Buffer

	// When we can get the version.
	var getVersion = func() (string, error) {
		return "whatever 2.0.1", nil
	}

	err := printVersion(&b, getVersion)
	assert.NoError(t, err)

	assert.Equal(t, b.String(),
		"Tailscale IPs (tips) 0.0.1\n\nTailscale CLI\n********************************\nwhatever 2.0.1\n")

	b.Reset()

	// When we can't get the version.
	getVersion = func() (string, error) {
		return "", errors.New("couldn't get version cause some crazy reason!")
	}

	err = printVersion(&b, getVersion)
	assert.NoError(t, err)

	assert.Equal(t, b.String(),
		"Tailscale IPs (tips) 0.0.1\n\n")
}
