package app

import (
	"context"
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/tailscale/tailscale-client-go/tailscale"
	"golang.org/x/exp/rand"
	"os"
	"tips/pkg/tailscale_cli"
)

const (
	testDevicesFile = "testmode/devices.json"
)

func DevicesResourceTest(ctx context.Context, client *tailscale.Client) ([]tailscale.Device, map[string]tailscale_cli.DeviceInfo, error) {
	f, err := os.Open(testDevicesFile)
	if err != nil {
		log.Fatal("failed to read file testmode/devices.json with err: ", err)
	}

	defer f.Close()

	var devs []tailscale.Device
	err = json.NewDecoder(f).Decode(&devs)
	if err != nil {
		log.Fatal("failed to Unmarshal file testmode/devices.json with err: ", err)
	}

	// Total shameful hack in the interest of testing.
	// Enrich the data for now with this made up data so that must boxes appear online.
	enrichedDevices := make(map[string]tailscale_cli.DeviceInfo)
	var counter int
	for _, dev := range devs {
		isSelf := counter == 0
		isOnline := true
		if rand.Float32() < 0.01 {
			isOnline = false
		}
		enrichedDevices[dev.NodeKey] = tailscale_cli.DeviceInfo{
			DNSName: "",
			IsSelf:  isSelf,
			Online:  isOnline,
			Tags:    nil,
		}
		counter++
	}

	return devs, enrichedDevices, nil
}
