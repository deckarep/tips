/*
Open Source Initiative OSI - The MIT License (MIT):Licensing

The MIT License (MIT)
Copyright Ralph Caraveo (deckarep@gmail.com)

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package utils

import (
	"runtime"
	"testing"
)

func TestSelectBinaryPath(t *testing.T) {
	const (
		pathDarwinBogus     = "/Applications/Something/That/Does/Not/Exist"
		pathDarwinTailscale = "/Applications/Tailscale.app/Contents/MacOS/Tailscale"

		// TODO: On linux, it doesn't work with the full path, it only wants the binary name given.
		// This runs on the CI build server.
		pathLinuxBogus  = "nothingburger"
		pathLinuxPython = "python"
	)

	c := map[string][]string{
		"linux": {
			pathLinuxBogus,
			pathLinuxPython,
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

		if result != "/usr/bin/python" {
			t.Errorf("expected binary to be: %s for os: %s", pathLinuxPython, platform)
		}
	}

}
