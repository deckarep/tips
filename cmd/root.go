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
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tailscale/tailscale-client-go/tailscale"
	"os"
	"strings"
	"time"
	"tips/app"
)

var (
	cfgFile       string
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
	useCmd        bool
	useCSSHX      bool
	useOauth      bool
	useSSH        bool
	test          bool
	ips           bool
	jsonn         bool

	foo string
)

func init() {
	//cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&foo, "foo", "", "blah", "foo is a test flag")
	rootCmd.PersistentFlags().StringVarP(&slice, "slice", "", "", "slices the results after filtering followed by sorting")
	rootCmd.PersistentFlags().StringVarP(&sortOrder, "sort", "s", "", "overrides the default/configured sort order --sort 'machine,addresss'")
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
	rootCmd.PersistentFlags().BoolVar(&ips, "ips", false, "when true returns only ips")
	rootCmd.PersistentFlags().BoolVar(&jsonn, "json", false, "when true returns only json data")

	// Note: Not sure if this flag is useful.
	rootCmd.PersistentFlags().DurationVarP(&cliTimeout, "cli_timeout", "", time.Second*5, "timeout duration for the Tailscale cli")
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	//rootCmd.PersistentFlags().StringVarP(&projectBase, "projectbase", "b", "", "base project directory eg. github.com/spf13/")
	//rootCmd.PersistentFlags().StringP("author", "a", "Ralph Caraveo <deckarep@gmail.com>", "Author name for copyright attribution")
	//rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "Name of license for the project (can provide `licensetext` in config)")
	//rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
	viper.BindPFlag("foo", rootCmd.PersistentFlags().Lookup("foo"))
	//viper.SetDefault("foo", "barf")
	//viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	//viper.BindPFlag("projectbase", rootCmd.PersistentFlags().Lookup("projectbase"))
	//viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	//viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	//viper.SetDefault("license", "apache")
}

var rootCmd = &cobra.Command{
	Use:   "tips",
	Short: "tips: The command-line tool to wrangle your Tailscale/tailnet cluster whether large or small.",
	Long: `tips is a robust command-line tool to help you inspect, query, manage and execute commands on your 
				tailnet cluster created by deckarep.
                Complete documentation is available at: github.com/deckarep/tips`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// 0. Package all configuration logic.
		cfgCtx := packageCfg(args)
		ctx = context.WithValue(ctx, app.CtxKeyConfig, cfgCtx)
		ctx = context.WithValue(ctx, app.CtxKeyUserQuery, fmt.Sprintf("%s %s", cfgCtx.PrimaryFilter, cfgCtx.RemoteCmd))

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

		var client *tailscale.Client
		if !useOauth {
			client = app.NewClient(ctx)
		} else {
			client = app.NewOauthClient(ctx)
		}

		var devicesResourceFunc = app.DevicesResource
		if cfgCtx.TestMode {
			// In test mode, indirect to mocked test data.
			devicesResourceFunc = app.DevicesResourceTest
		}

		devList, devEnriched, err := devicesResourceFunc(ctx, client)
		if err != nil {
			log.Fatal("problem with resource lookup of devices with err: ", err)
		}

		view, err := app.ProcessDevicesTable(ctx, devList, devEnriched)
		if err != nil {
			log.Fatal("problem processing devices data with err: ", err)
		}

		if cfgCtx.IsRemoteCommand() {
			// It's a remote command, instead of rendering a table execute the remote command over all hosts.
			var hosts []string
			for _, rows := range view.Rows {
				// TODO: getting back a GeneralTableView in this stage is not ideal, it's too abstract.
				// Column's may change so this is dumb.
				hosts = append(hosts, rows[1])
			}
			// Do the remote cluster command.
			app.ExecuteClusterRemoteCmd(ctx, os.Stdout, hosts, cfgCtx.RemoteCmd)
		} else {
			err = app.RenderTableView(ctx, view, os.Stdout)
			if err != nil {
				log.Fatal("problem rendering table view with err: ", err)
			}
		}
	},
}

func dumpColors() {
	for i := 0; i < 256; i++ {
		s := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(fmt.Sprintf("%d", i)))
		fmt.Println(s.Render(fmt.Sprintf("Color: %d", i)))
	}
}

func packageCfg(args []string) *app.ConfigCtx {
	// Populate context key/values as needed.

	// 0. Validate
	if jsonn == true && ips == true {
		log.Fatal("the --ips and --json flag must not be used together. Choose one or the other.")
	}

	cfgCtx := app.NewConfigCtx()
	cfgCtx.IPsOutput = ips
	cfgCtx.JsonOutput = jsonn
	cfgCtx.NoCache = nocache
	cfgCtx.NoColor = nocolor
	cfgCtx.Slice = app.ParseSlice(slice)
	cfgCtx.Filters = app.ApplyFilter(filter)
	cfgCtx.Columns = app.ParseColumns(columns)
	cfgCtx.Concurrency = concurrency
	cfgCtx.TestMode = test

	// The 0th arg is the Primary filter, if nothing was specified we consider it to represent: @ for all
	if len(args) > 0 {
		cfgCtx.PrimaryFilter = args[0]
	} else {
		cfgCtx.PrimaryFilter = app.PrimaryFilterAll
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

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
