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
	"github.com/charmbracelet/log"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tailscale/tailscale-client-go/tailscale"
	"os"
	"time"
	"tips/app"
	"tips/pkg"
)

var (
	cfgFile       string
	clientTimeout time.Duration
	cliTimeout    time.Duration
	concurrency   int
	filter        string
	sortOrder     string
	tailnet       string
	useCSSHX      bool
	useSSH        bool
)

func init() {
	//cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&sortOrder, "sort", "s", "", "overrides the default/configured sort order --sort 'machine,addresss'")
	rootCmd.PersistentFlags().StringVarP(&tailnet, "tailnet", "t", "", "the tailnet to operate on")
	rootCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "c", 5, "concurrency level when executing requests")
	rootCmd.PersistentFlags().StringVarP(&filter, "filter", "f", "", "if provided, applies filtering logic: --filter 'tag:tunnel'")
	rootCmd.PersistentFlags().BoolVar(&useCSSHX, "csshx", false, "if csshx is installed, opens a multi-window session over all matching hosts")
	rootCmd.PersistentFlags().BoolVar(&useSSH, "ssh", false, "ssh into a matching single host")
	rootCmd.PersistentFlags().DurationVarP(&clientTimeout, "client_timeout", "", time.Second*5, "timeout duration for the Tailscale api")
	// Note: Not sure if this flag is useful.
	rootCmd.PersistentFlags().DurationVarP(&cliTimeout, "cli_timeout", "", time.Second*5, "timeout duration for the Tailscale cli")
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	//rootCmd.PersistentFlags().StringVarP(&projectBase, "projectbase", "b", "", "base project directory eg. github.com/spf13/")
	//rootCmd.PersistentFlags().StringP("author", "a", "Ralph Caraveo <deckarep@gmail.com>", "Author name for copyright attribution")
	//rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "Name of license for the project (can provide `licensetext` in config)")
	//rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
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

		// TODO: make this configurable
		api_key := os.Getenv("tips_api_key")
		// TODO: make this configurable
		tailnet = "deckarep@gmail.com"

		client := func() *tailscale.Client {
			client, err := tailscale.NewClient(
				api_key, // When doing oauth, this field must be blank!!!
				tailnet,
				//tailscale.WithOAuthClientCredentials(oauthClientID, oauthClientSecret, nil),
				tailscale.WithUserAgent(pkg.UserAgent),
			)
			if err != nil {
				log.Fatal("failed to create client with err: ", err)
			}
			return client
		}()

		devList, devEnriched, err := app.DevicesResource(ctx, client)
		if err != nil {
			log.Fatal("problem with resource lookup of devices with err: ", err)
		}

		view, err := app.ProcessDevicesTable(ctx, devList, devEnriched)
		if err != nil {
			log.Fatal("problem processing devices data with err: ", err)
		}

		err = app.RenderTableView(ctx, view, os.Stdout)
		if err != nil {
			log.Fatal("problem rendering table view with err: ", err)
		}
	},
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
