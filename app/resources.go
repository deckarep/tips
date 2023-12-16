package app

import (
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/tailscale/tailscale-client-go/tailscale"
	"time"
	"tips/pkg/tailscale_cli"
)

func DevicesResource(ctx context.Context, client *tailscale.Client) ([]tailscale.Device, map[string]tailscale_cli.DeviceInfo, error) {
	// 0. Check cache - return cached results if cache timeout not yet expired.
	if false {
		// TODO: return from cache.
	}

	// 1. Do tailscale api lookup for devices data.
	ctxTimeOut, cancelTimeout := context.WithTimeout(ctx, time.Second*5) //TODO: populate timeout duration from config
	defer cancelTimeout()
	devList, err := client.Devices(ctxTimeOut)
	if err != nil {
		log.Fatal("tailscale api failed during devices lookup: ", err)
	}

	// 2. When available, enrich this data with data from the Tailscale cli, if this is run from a node within the tailnet.
	enrichedDevices, err := tailscale_cli.GetDevicesStatuses()
	if err != nil {
		fmt.Println("failed to get results: ", err)
	}

	return devList, enrichedDevices, nil
}
