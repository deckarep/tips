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
	"context"
	"fmt"
	"os"
	"time"

	"github.com/deckarep/tips/pkg"
	"github.com/spf13/viper"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var (
	//cfgFile       string
	basic         bool
	cacheTimeout  time.Duration
	clientTimeout time.Duration
	cliTimeout    time.Duration
	columns       string
	concurrency   int
	filter        string
	nocache       bool
	nocolor       bool
	slice         string
	sortOrder     string
	tailnet       string
	tipsAPIKey    string
	useCSSHX      bool
	useOauth      bool
	useSSH        bool
	test          bool
	ips           bool
	ips_delimiter string
	jsonn         bool
	page          int
)

// bindRootBoolFlag binds a boolean cobra flag to a viper config flag.
func bindRootBoolFlag(pBool *bool, name, description string, defaultVal bool) {
	rootCmd.PersistentFlags().BoolVar(pBool, name, false, description)
	viper.BindPFlag(name, rootCmd.PersistentFlags().Lookup(name))
}

// bindRootIntFlag binds an int cobra flag to a viper config flag.
func bindRootIntFlag(pInt *int, name, shorthand string, defaultVal int, description string) {
	rootCmd.PersistentFlags().IntVarP(pInt, name, shorthand, defaultVal, description)
	viper.BindPFlag(name, rootCmd.PersistentFlags().Lookup(name))
}

// bindRootStringFlag binds a string cobra flag to a viper config flag.
func bindRootStringFlag(pString *string, name, shorthand, value, description string) {
	rootCmd.PersistentFlags().StringVarP(pString, name, shorthand, value, description)
	viper.BindPFlag(name, rootCmd.PersistentFlags().Lookup(name))
}

// bindRootDurationFlag binds a string cobra flag to a viper config flag.
func bindRootDurationFlag(pDuration *time.Duration, name, shorthand string, value time.Duration, description string) {
	rootCmd.PersistentFlags().DurationVarP(pDuration, name, shorthand, value, description)
	viper.BindPFlag(name, rootCmd.PersistentFlags().Lookup(name))
}

func init() {
	bindRootBoolFlag(&basic, "basic", "when true, renders the table as simple ascii with no color", false)
	bindRootDurationFlag(&cacheTimeout, "cache_timeout", "", time.Minute*5, "timeout duration for local db (db.bolt) cache file")
	bindRootDurationFlag(&clientTimeout, "client_timeout", "", time.Second*5, "timeout duration for the Tailscale api")
	bindRootStringFlag(&columns, "columns", "", "", "columns limits which columns to return")
	bindRootIntFlag(&concurrency, "concurrency", "c", 5, "concurrency level when executing requests")
	bindRootStringFlag(&filter, "filter", "f", "", "if provided, applies filtering logic: --filter 'tag:tunnel'")
	bindRootBoolFlag(&ips, "ips", "when provided returns ips comma-delimited", false)
	bindRootStringFlag(&ips_delimiter, "delimiter", "d", "\n", "delimiter to use when the --ips flag is provided")
	bindRootBoolFlag(&jsonn, "json", "when true returns only json data", false)
	bindRootBoolFlag(&nocache, "nocache", "forces the cache to be expunged", false)
	bindRootBoolFlag(&nocolor, "nocolor", "when --nocolor is provided disables log color highlighting", false)
	bindRootIntFlag(&page, "page", "p", 1, "use with slicing to get the next page of results, paging is 1-based")
	bindRootStringFlag(&slice, "slice", "", "", "slices the results after filtering followed by sorting")
	bindRootStringFlag(&sortOrder, "sort", "s", "",
		"overrides the default/configured sort order --sort 'machine,address:dsc' the default order is always ascending (asc) for each column")
	bindRootStringFlag(&tailnet, "tailnet", "t", "", "the tailnet to operate on (required)")
	bindRootBoolFlag(&test, "test", "when true runs the tool in test mode with mocked data", false)
	bindRootStringFlag(&tipsAPIKey, "tips_api_key", "", "", "tailscale api key for remote requests")
	bindRootBoolFlag(&useCSSHX, "csshx", "if csshx is installed, opens a multi-window session over all matching hosts", false)
	bindRootBoolFlag(&useSSH, "ssh", "ssh into a matching single host", false)
	bindRootBoolFlag(&useOauth, "oauth", "use oauth when flag is provided.", false)
	bindRootDurationFlag(&cliTimeout, "cli_timeout", "", time.Second*5, "timeout duration for the Tailscale cli")

	// Required flags are set here.
	// This doesn't seem compatible with Viper.
	// rootCmd.MarkPersistentFlagRequired("tailnet")
}

var rootCmd = &cobra.Command{
	Use:   "tips",
	Short: "tips: The command-line tool to wrangle your Tailscale tailnet cluster whether large or small.",
	Long: `tips is a robust command-line tool to help you inspect, query, manage and execute commands on your 
				tailnet cluster. Created by @deckarep.
                Complete documentation is available at: github.com/deckarep/tips`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// 0. Package all configuration logic.
		cfgCtx, err := packageCfg(args)
		if err != nil {
			return err
		}

		ctx = context.WithValue(ctx, pkg.CtxKeyConfig, cfgCtx)
		// CONSIDER: should this show all flags?
		ctx = context.WithValue(ctx, pkg.CtxKeyUserQuery, fmt.Sprintf("%s %s", cfgCtx.PrefixFilter.Query(), cfgCtx.RemoteCmd))

		client := pkg.NewClient(ctx)
		if useOauth {
			client = pkg.NewOauthClient(ctx)
		}

		cachedDevRepo := pkg.NewCachedRepo(pkg.NewRemoteDeviceRepo(client))
		var devicesResourceFunc = cachedDevRepo.DevicesResource

		// In test mode, indirect to mocked test data.
		// TODO: refactor this out as it doesn't belong here.
		if cfgCtx.TestMode {
			mockDevRepo := pkg.NewMockedDeviceRepo()
			cachedDevRepo = pkg.NewCachedRepo(mockDevRepo)
			devicesResourceFunc = cachedDevRepo.DevicesResource
		}

		devList, err := devicesResourceFunc(ctx)
		if err != nil {
			return err
		}

		view, err := pkg.ProcessDevicesTable(ctx, devList)
		if err != nil {
			return err
		}

		if cfgCtx.IsRemoteCommand() {
			// It's a remote command, instead of rendering a table execute the remote command over all hosts.
			hosts := getHosts(ctx, view)

			// Do the remote cluster command.
			pkg.ExecuteClusterRemoteCmd(ctx, os.Stdout, hosts, cfgCtx.RemoteCmd)
		} else {
			if cfgCtx.JsonOutput {
				if err = pkg.RenderJson(ctx, view, os.Stdout); err != nil {
					return err
				}
			} else if cfgCtx.IPsOutput {
				if err = pkg.RenderIPs(ctx, view, os.Stdout); err != nil {
					return err
				}
			} else {
				if cfgCtx.Basic {
					if err = pkg.RenderASCIITableView(ctx, view, os.Stdout); err != nil {
						return err
					}
				} else {
					if err = pkg.RenderTableView(ctx, view, os.Stdout); err != nil {
						return err
					}
				}
			}
		}

		return nil
	},
}

//func dumpColors() {
//	for i := 0; i < 256; i++ {
//		s := lipgloss.NewStyle().
//			Bold(true).
//			Foreground(lipgloss.Color(fmt.Sprintf("%d", i)))
//		fmt.Println(s.Render(fmt.Sprintf("Color: %d", i)))
//	}
//}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Print("root command failed", "error", err)
		os.Exit(1)
	}
}
