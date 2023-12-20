package app

import (
	"context"
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/tailscale/tailscale-client-go/tailscale"
	"os"
	"tips/pkg/tailscale_cli"
)

func DevicesResourceTest(ctx context.Context, client *tailscale.Client) ([]tailscale.Device, map[string]tailscale_cli.DeviceInfo, error) {
	b, err := os.ReadFile("testmode/devices.json")
	if err != nil {
		log.Fatal("failed to read file testmode/devices.json with err: ", err)
	}

	var devs []tailscale.Device
	err = json.Unmarshal(b, &devs)
	if err != nil {
		log.Fatal("failed to Unmarshal file testmode/devices.json with err: ", err)
	}

	return devs, nil, nil
}
