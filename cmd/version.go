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

package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/deckarep/tips/pkg/tailscale_cli"

	"github.com/deckarep/tips/pkg"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: pkg.AppShortName + " empowers you to wrangle your Tailnet",
	Long:  pkg.AppLongName + " empowers you to manage your Tailscale cluster like a pro",
	RunE: func(cmd *cobra.Command, args []string) error {
		return printVersion(os.Stdout, tailscale_cli.GetVersion)
	},
}

type versionGetter func() (string, error)

func printVersion(w io.Writer, getVersion versionGetter) error {
	if _, err := fmt.Fprintf(w, "%s (%s) %s\n", pkg.AppLongName, pkg.AppShortName, pkg.AppVersion); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}

	cliVersion, err := getVersion()
	if err == nil {
		if _, err := fmt.Fprintln(w, "Tailscale CLI"); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, strings.Repeat("*", 32)); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, cliVersion); err != nil {
			return err
		}
	}
	return nil
}
