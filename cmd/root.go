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
	"regexp"
	"strings"
	"time"
	"tips/pkg"

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
	useCSSHX      bool
	useOauth      bool
	useSSH        bool
	test          bool
	ips           bool
	ips_delimiter string
	jsonn         bool
	page          int

	foo string
)

func init() {
	//cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().DurationVarP(&cacheTimeout, "cache_timeout", "", time.Minute*5, "timeout duration for local db (db.bolt) cache file")
	rootCmd.PersistentFlags().StringVarP(&foo, "foo", "", "blah", "foo is a test flag")
	rootCmd.PersistentFlags().StringVarP(&slice, "slice", "", "", "slices the results after filtering followed by sorting")
	rootCmd.PersistentFlags().StringVarP(&sortOrder, "sort", "s", "",
		"overrides the default/configured sort order --sort 'machine,address:dsc' the default order is always ascending (asc) for each column")
	rootCmd.PersistentFlags().StringVarP(&tailnet, "tailnet", "t", "", "the tailnet to operate on")
	rootCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "c", 5, "concurrency level when executing requests")
	rootCmd.PersistentFlags().StringVarP(&columns, "columns", "", "", "columns limits which columns to return")
	rootCmd.PersistentFlags().StringVarP(&filter, "filter", "f", "", "if provided, applies filtering logic: --filter 'tag:tunnel'")
	rootCmd.PersistentFlags().BoolVar(&useCSSHX, "csshx", false, "if csshx is installed, opens a multi-window session over all matching hosts")
	rootCmd.PersistentFlags().BoolVar(&useSSH, "ssh", false, "ssh into a matching single host")
	rootCmd.PersistentFlags().BoolVar(&useOauth, "oauth", false, "use oauth when flag is provided.")
	rootCmd.PersistentFlags().DurationVarP(&clientTimeout, "client_timeout", "", time.Second*5, "timeout duration for the Tailscale api")
	rootCmd.PersistentFlags().BoolVarP(&nocache, "nocache", "n", false, "forces the cache to be expunged")
	rootCmd.PersistentFlags().BoolVarP(&nocolor, "nocolor", "", false, "when --nocolor is provided disables log color highlighting")
	rootCmd.PersistentFlags().BoolVar(&test, "test", false, "when true runs the tool in test mode with mocked data")
	rootCmd.PersistentFlags().BoolVarP(&ips, "ips", "", false, "when provided returns ips comma-delimited")
	rootCmd.PersistentFlags().StringVarP(&ips_delimiter, "delimiter", "d", "\n", "delimiter to use when the --ips flag is provided")
	rootCmd.PersistentFlags().BoolVar(&jsonn, "json", false, "when true returns only json data")
	rootCmd.PersistentFlags().BoolVar(&basic, "basic", false, "when true, renders the table as simple ascii with no color")
	rootCmd.PersistentFlags().IntVarP(&page, "page", "p", 1, "use with slicing to get the next page of results, paging is 1-based")

	// Note: Not sure if this flag is useful.
	rootCmd.PersistentFlags().DurationVarP(&cliTimeout, "cli_timeout", "", time.Second*5, "timeout duration for the Tailscale cli")
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	//rootCmd.PersistentFlags().StringVarP(&projectBase, "projectbase", "b", "", "base project directory eg. github.com/spf13/")
	//rootCmd.PersistentFlags().StringP("author", "a", "Ralph Caraveo <deckarep@gmail.com>", "Author name for copyright attribution")
	//rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "Name of license for the project (can provide `licensetext` in config)")
	//rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")

	// TODO: look at bindPflag and how it works.
	//viper.BindPFlag("foo", rootCmd.PersistentFlags().Lookup("foo"))

	//viper.SetDefault("foo", "barf")
	//viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	//viper.BindPFlag("projectbase", rootCmd.PersistentFlags().Lookup("projectbase"))
	//viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	//viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	//viper.SetDefault("license", "apache")
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

		// 0. Package all configuration logic.
		cfgCtx := packageCfg(args)
		ctx = context.WithValue(ctx, pkg.CtxKeyConfig, cfgCtx)
		ctx = context.WithValue(ctx, pkg.CtxKeyUserQuery, fmt.Sprintf("%s %s", cfgCtx.PrimaryFilter, cfgCtx.RemoteCmd))

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

func getHosts(ctx context.Context, view *pkg.GeneralTableView) []string {
	cfg := pkg.CtxAsConfig(ctx, pkg.CtxKeyConfig)
	var hosts []string

	for _, rows := range view.Rows {
		// TODO: getting back a GeneralTableView in this stage is not ideal, it's too abstract.
		// Column's may change so this is dumb.
		if cfg.TestMode {
			hosts = append(hosts, "blade")
		} else {
			hosts = append(hosts, rows[1])
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
	// Populate context key/values as needed.

	// 0. Validate
	if jsonn && ips {
		log.Fatal("the --ips and --json flag must not be used together. Choose one or the other.")
	}

	cfgCtx := pkg.NewConfigCtx()
	cfgCtx.Basic = basic
	cfgCtx.CacheTimeout = cacheTimeout
	cfgCtx.IPsOutput = ips
	cfgCtx.IPsDelimiter = ips_delimiter
	cfgCtx.JsonOutput = jsonn
	cfgCtx.NoCache = nocache
	cfgCtx.NoColor = nocolor
	cfgCtx.Slice = pkg.ParseSlice(slice, page)
	cfgCtx.SortOrder = pkg.ParseSortString(sortOrder)
	cfgCtx.Filters = pkg.ParseFilter(filter)
	cfgCtx.Columns = pkg.ParseColumns(columns)
	cfgCtx.Concurrency = concurrency
	cfgCtx.TestMode = test
	cfgCtx.Page = page

	// The 0th arg is the Primary filter, if nothing was specified we consider it to represent: @ for all
	if len(args) > 0 {
		if strings.TrimSpace(args[0]) == "@" {
			cfgCtx.PrefixFilter = "*"
		} else {
			cfgCtx.PrefixFilter = args[0]
		}
		// The regex filter to be deprecated.
		cfgCtx.PrimaryFilter = regexp.MustCompile(args[0])
	} else {
		cfgCtx.PrefixFilter = "*"
		// This one to be deprecated
		cfgCtx.PrimaryFilter = nil
	}

	// The 1st arg along with the rest - [1:] when provided is a remote command to execute.
	// So we join this up into a single string.
	if len(args) > 1 {
		cfgCtx.RemoteCmd = strings.TrimSpace(strings.Join(args[1:], " "))
	}

	cfgCtx.TailscaleAPI.ApiKey = os.Getenv("tips_api_key")
	cfgCtx.Tailnet = "deckarep@gmail.com"
	cfgCtx.TailscaleAPI.Timeout = time.Second * 5

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
