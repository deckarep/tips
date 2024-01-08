package cmd

import (
	"context"
	"errors"
	"strings"

	"github.com/spf13/viper"

	"github.com/deckarep/tips/pkg"
)

func packageCfg(args []string) (*pkg.ConfigCtx, error) {
	cfgCtx := pkg.NewConfigCtx()

	const (
		allFilterCLI = "@" // User uses this on the CLI because * expands in the shell.
		allFilter    = "*" // So bottom line, use this in the codebase.
	)

	// Parse positional args here
	// The 0th arg is the Primary filter, if nothing was specified we consider it to represent: @ for all
	if len(args) > 0 {
		if strings.TrimSpace(args[0]) == allFilterCLI {
			cfgCtx.PrefixFilter = pkg.ParsePrefixFilter("*")
		} else {
			cfgCtx.PrefixFilter = pkg.ParsePrefixFilter(args[0])
		}
	} else {
		cfgCtx.PrefixFilter = pkg.ParsePrefixFilter("*")
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
		return nil, errors.New("the --ips and --json flag must not be used together. Choose one or the other.")
	}

	if strings.TrimSpace(cfgCtx.TailscaleAPI.ApiKey) == "" {
		return nil,
			errors.New("a 'tips_api_key' must be defined either as an environment variable (uppercase), in a config or as a --tips_api_key flag")

	}

	if strings.TrimSpace(cfgCtx.Tailnet) == "" {
		return nil,
			errors.New("at an absolute minimum a tailnet must be specified either in the config file or as flag --tailnet")

	}

	return cfgCtx, nil
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
