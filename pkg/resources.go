package pkg

import (
	"context"
	"time"
	"tips/pkg/tailscale_cli"

	"github.com/charmbracelet/log"
	"github.com/tailscale/tailscale-client-go/tailscale"
)

func DevicesResource(ctx context.Context, client *tailscale.Client) ([]tailscale.Device, map[string]tailscale_cli.DeviceInfo, error) {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)
	startTime := time.Now()
	defer func() {
		cfg.TailscaleAPI.ElapsedTime = time.Since(startTime)
	}()

	// 0. Check cache - return cached results if cache timeout not yet expired.
	if cfg.NoCache {
		log.Warn("--nocache is not yet support, but this should force a refresh of data in cache")
	}

	// 1. Do tailscale api lookup for devices data.
	ctxTimeOut, cancelTimeout := context.WithTimeout(ctx, cfg.TailscaleAPI.Timeout)
	defer cancelTimeout()
	devList, err := client.Devices(ctxTimeOut)
	if err != nil {
		log.Fatal("tailscale api failed during devices lookup: ", err)
	}

	// 2. When available, enrich this data with data from the Tailscale cli, if this is run from a node within the tailnet.
	enrichedDevices, err := tailscale_cli.GetDevicesState()
	if err != nil {
		log.Debug("unable to get enriched data from tailscale cli", "error", err)
	}

	return devList, enrichedDevices, nil
}
