/*
Open Source Initiative OSI - The MIT License (MIT):Licensing

The MIT License (MIT)
Copyright (c) 2023 - 2024 Ralph Caraveo (deckarep@gmail.com)

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
	"strings"
	"time"
	"tips/pkg"

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
	// Note: Not sure if this flag is useful.
	bindRootDurationFlag(&cliTimeout, "cli_timeout", "", time.Second*5, "timeout duration for the Tailscale cli")

	// Required flags are set here.
	// This doesn't seem compatible with Viper.
	//rootCmd.MarkPersistentFlagRequired("tailnet")
}

var rootCmd = &cobra.Command{
	Use:   "tips",
	Short: "tips: The command-line tool to wrangle your Tailscale tailnet cluster whether large or small.",
	Long: `tips is a robust command-line tool to help you inspect, query, manage and execute commands on your 
				tailnet cluster. Created by @deckarep.
                Complete documentation is available at: github.com/deckarep/tips`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		//if true {
		//	dumpColors()
		//	//return
		//}

		// 0. Package all configuration logic.
		cfgCtx := packageCfg(args)
		ctx = context.WithValue(ctx, pkg.CtxKeyConfig, cfgCtx)
		// CONSIDER: should this show all flags?
		ctx = context.WithValue(ctx, pkg.CtxKeyUserQuery, fmt.Sprintf("%s %s", cfgCtx.PrefixFilter, cfgCtx.RemoteCmd))

		//if false {
		//	//myCmd := "sudo ls /var/log"
		//	myCmd := "head -n100 /var/log/secure"
		//	//myCmd := "while true; do echo 'hi'; sleep 1; done"
		//	var hosts = []string{
		//		"blade",
		//		"blade",
		//		//"blade",
		//		//"blade",
		//		//"blade",
		//		//"blade",
		//		//"blade",
		//		//"blade",
		//		//"blade",
		//		//"blade",
		//	}
		//
		//	app.ExecuteClusterRemoteCmd(ctx, os.Stdout, hosts, myCmd)
		//}

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
			log.Fatal("problem with resource lookup of devices", "error", err)
		}

		view, err := pkg.ProcessDevicesTable(ctx, devList)
		if err != nil {
			log.Fatal("problem occurred processing devices data", "error", err)
		}

		if cfgCtx.IsRemoteCommand() {
			// It's a remote command, instead of rendering a table execute the remote command over all hosts.
			hosts := getHosts(ctx, view)

			// Do the remote cluster command.
			pkg.ExecuteClusterRemoteCmd(ctx, os.Stdout, hosts, cfgCtx.RemoteCmd)
		} else {
			if cfgCtx.JsonOutput {
				if err = pkg.RenderJson(ctx, view, os.Stdout); err != nil {
					log.Fatal("problem encoding json output", "error", err)
				}
			} else if cfgCtx.IPsOutput {
				if err = pkg.RenderIPs(ctx, view, os.Stdout); err != nil {
					log.Fatal("problem generating ips output", "error", err)
				}
			} else {
				if cfgCtx.Basic {
					if err = pkg.RenderASCIITableView(ctx, view, os.Stdout); err != nil {
						log.Fatal("problem rendering basic table view", "error", err)
					}
				} else {
					if err = pkg.RenderTableView(ctx, view, os.Stdout); err != nil {
						log.Fatal("problem rendering fancy table view", "error", err)
					}
				}
			}
		}
	},
}

func getHosts(ctx context.Context, view *pkg.GeneralTableView) []pkg.RemoteCmdHost {
	cfg := pkg.CtxAsConfig(ctx, pkg.CtxKeyConfig)
	var hosts []pkg.RemoteCmdHost

	for _, rows := range view.Rows {
		// TODO: getting back a GeneralTableView in this stage is not ideal, it's too abstract.
		// Column's may change so this is dumb.
		if cfg.TestMode {
			hosts = append(hosts, pkg.RemoteCmdHost{
				Original: "blade",
				Alias:    rows[1],
			})
		} else {
			hosts = append(hosts, pkg.RemoteCmdHost{
				Original: rows[1],
			})
		}
	}

	return hosts
}

//func dumpColors() {
//	for i := 0; i < 256; i++ {
//		s := lipgloss.NewStyle().
//			Bold(true).
//			Foreground(lipgloss.Color(fmt.Sprintf("%d", i)))
//		fmt.Println(s.Render(fmt.Sprintf("Color: %d", i)))
//	}
//}

func packageCfg(args []string) *pkg.ConfigCtx {
	cfgCtx := pkg.NewConfigCtx()

	const (
		allFilterCLI = "@" // User uses this on the CLI because * expands in the shell.
		allFilter    = "*" // So bottom line, use this in the codebase.
	)

	// Parse positional args here
	// The 0th arg is the Primary filter, if nothing was specified we consider it to represent: @ for all
	if len(args) > 0 {
		if strings.TrimSpace(args[0]) == allFilterCLI {
			cfgCtx.PrefixFilter = allFilter
		} else {
			cfgCtx.PrefixFilter = args[0]
		}
	} else {
		cfgCtx.PrefixFilter = allFilter
	}

	// The 1st arg along with the rest - [1:] when provided is a remote command to execute.
	// So we join this up into a single string.
	if len(args) > 1 {
		cfgCtx.RemoteCmd = strings.TrimSpace(strings.Join(args[1:], " "))
	}

	// Populate flags
	cfgCtx.Basic = viper.GetBool("basic")
	cfgCtx.CacheTimeout = viper.GetDuration("cache_timeout")
	cfgCtx.Columns = pkg.ParseColumns(viper.GetString("columns"))
	cfgCtx.Concurrency = viper.GetInt("concurrency")
	cfgCtx.Filters = pkg.ParseFilter(viper.GetString("filter"))
	cfgCtx.IPsOutput = viper.GetBool("ips")
	cfgCtx.IPsDelimiter = viper.GetString("delimiter")
	cfgCtx.JsonOutput = viper.GetBool("json")
	cfgCtx.NoCache = viper.GetBool("nocache")
	cfgCtx.NoColor = viper.GetBool("nocolor")
	cfgCtx.Page = viper.GetInt("page")
	cfgCtx.Slice = pkg.ParseSlice(viper.GetString("slice"), viper.GetInt("page"))
	cfgCtx.SortOrder = pkg.ParseSortString(viper.GetString("sort"))
	cfgCtx.Tailnet = viper.GetString("tailnet")
	cfgCtx.TailscaleAPI.ApiKey = viper.GetString("tips_api_key")
	cfgCtx.TailscaleAPI.Timeout = viper.GetDuration("client_timeout")
	cfgCtx.TestMode = viper.GetBool("test")

	// Validate flags
	if cfgCtx.JsonOutput && cfgCtx.IPsOutput {
		log.Fatal("the --ips and --json flag must not be used together. Choose one or the other.")
	}

	if strings.TrimSpace(cfgCtx.TailscaleAPI.ApiKey) == "" {
		log.Fatal("a 'tips_api_key' must be defined either as an environment variable (uppercase), in a config or as a --tips_api_key flag")
	}

	if strings.TrimSpace(cfgCtx.Tailnet) == "" {
		log.Fatal("at an absolute minimum a tailnet must be specified either in the config file or as flag --tailnet")
	}

	return cfgCtx
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//func initConfig() {
//	// Don't forget to read config either from cfgFile or from home directory!
//	if cfgFile != "" {
//		// Use config file from the flag.
//		viper.SetConfigFile(cfgFile)
//	} else {
//		// Find home directory.
//		home, err := homedir.Dir()
//		if err != nil {
//			fmt.Println(err)
//			os.Exit(1)
//		}
//
//		// Search config in home directory with name ".cobra" (without extension).
//		viper.AddConfigPath(home)
//		viper.SetConfigName(".cobra")
//	}
//
//	if err := viper.ReadInConfig(); err != nil {
//		fmt.Println("Can't read config:", err)
//		os.Exit(1)
//	}
//}
