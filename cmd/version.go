package cmd

import (
	"fmt"
	"strings"
	"tips/pkg"
	"tips/pkg/tailscale_cli"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: pkg.AppShortName + " empowers you to wrangle your Tailnet",
	Long:  pkg.AppLongName + " empowers you to manage your Tailscale cluster like a pro",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s (%s) v%s\n", pkg.AppLongName, pkg.AppShortName, pkg.AppVersion)
		fmt.Println()

		cliVersion, err := tailscale_cli.TailScaleVersion()
		if err == nil {
			fmt.Println("Tailscale CLI")
			fmt.Println(strings.Repeat("*", 32))
			fmt.Println(cliVersion)
		}
	},
}
