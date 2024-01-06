package utils

import (
	"runtime"
	"testing"
)

func TestSelectBinaryPath(t *testing.T) {
	const (
		pathDarwinBogus     = "/Applications/Something/That/Does/Not/Exist"
		pathDarwinTailscale = "/Applications/Tailscale.app/Contents/MacOS/Tailscale"
		pathLinuxBogus      = "/usr/bin/nothingburger"
		pathLinuxPython3    = "/usr/bin/python"
	)

	c := map[string][]string{
		"linux": {
			pathLinuxBogus,
			pathLinuxPython3,
		},
		"darwin": {
			pathDarwinBogus,
			pathDarwinTailscale,
		},
	}

	platform := runtime.GOOS
	switch platform {
	case "darwin":
		result, err := SelectBinaryPath(platform, c)
		if err != nil {
			t.Errorf("expected nil err, got: %s", err.Error())
		}

		if result != pathDarwinTailscale {
			t.Errorf("expected binary to be: %s for os: %s", pathDarwinTailscale, platform)
		}
	case "linux":
		result, err := SelectBinaryPath(platform, c)
		if err != nil {
			t.Errorf("expected nil err, got: %s", err.Error())
		}

		if result != pathDarwinTailscale {
			t.Errorf("expected binary to be: %s for os: %s", pathLinuxPython3, platform)
		}
	}

}
