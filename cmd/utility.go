package cmd

import (
	"context"

	"github.com/deckarep/tips/pkg"
)

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
