package utils

import (
	"runtime"
	"testing"
)

func TestSelectBinaryPath(t *testing.T) {
	const (
		tsPathBogus  = "/Applications/Something/That/Does/Not/Exist"
		tsPathDarwin = "/Applications/Tailscale.app/Contents/MacOS/Tailscale"
	)

	c := map[string][]string{
		"darwin": {
			tsPathBogus,
			tsPathDarwin,
		},
	}

	result, err := SelectBinaryPath("darwin", c)
	if err != nil {
		t.Errorf("expected nil err, got: %s", err.Error())
	}

	if result != tsPathDarwin {
		t.Errorf("expected binary to be: %s for os: %s", tsPathDarwin, runtime.GOOS)
	}
}
