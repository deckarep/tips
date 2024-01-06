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
	"errors"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/log"
)

func SelectBinaryPath(platform string, candidates map[string][]string) (string, error) {
	osSelected := runtime.GOOS

	//r, err := exec.LookPath("python")
	//log.Warn("first LookPath: ", "val", r, "err", err)
	//
	//r, err = exec.LookPath("python3")
	//log.Warn("second LookPath: ", "val", r, "err", err)

	if paths, exists := candidates[platform]; exists {
		for _, p := range paths {
			if confirmedPath, err := exec.LookPath(p); err == nil {
				return confirmedPath, nil
			}
		}
		return "", errors.New("no binary exists for this os: " + osSelected)
	}

	log.Fatal("binary is not setup for this os", "os", osSelected)
	return "", nil
}
