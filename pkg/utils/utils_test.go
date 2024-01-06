package utils

import (
	"runtime"
	"testing"
)

func TestSelectBinaryPath(t *testing.T) {
	// TODO: this will fail on other platforms.
	const (
		tsPathDarwin = "/Applications/Tailscale.app/Contents/MacOS/Tailscale"
	)

	c := map[string][]string{
		"darwin": {
			tsPathDarwin,
		},
	}

	result, err := SelectBinaryPath(c)
	if err != nil {
		t.Errorf("expected nil err, got: %s", err.Error())
	}

	if result != tsPathDarwin {
		t.Errorf("expected binary to be: %s for os: %s", tsPathDarwin, runtime.GOOS)
	}
}
