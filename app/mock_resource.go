package app

import (
	"context"

	"os"
	"time"
	"tips/pkg/tailscale_cli"

	"github.com/charmbracelet/log"
	jsoniter "github.com/json-iterator/go"
	"github.com/tailscale/tailscale-client-go/tailscale"
	"golang.org/x/exp/rand"
)

const (
	testDevicesFile = "testmode/devices.json"
)

func DevicesResourceTest(ctx context.Context, client *tailscale.Client) ([]tailscale.Device, map[string]tailscale_cli.DeviceInfo, error) {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)
	startTime := time.Now()
	defer func() {
		cfg.TailscaleAPI.ElapsedTime = time.Since(startTime)
	}()

	f, err := os.Open(testDevicesFile)
	if err != nil {
		log.Fatal("failed to read file", "error", err)
	}
	defer f.Close()

	var devs []tailscale.Device

	// Using jsoniter for now...stdlib json was a bit slower for large blobs.
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.NewDecoder(f).Decode(&devs)
	if err != nil {
		log.Fatal("failed to Unmarshal file", "error", err)
	}

	// Total shameful hack in the interest of testing.
	// Enrich the data for now with this made up data so that must boxes appear online.
	enrichedDevices := make(map[string]tailscale_cli.DeviceInfo, len(devs))
	var counter int
	for _, dev := range devs {
		isSelf := counter == 0
		isOnline := true
		offersExitNode := false
		if rand.Float32() < 0.01 {
			isOnline = false
		}
		if rand.Float32() < 0.05 {
			offersExitNode = true
		}
		enrichedDevices[dev.NodeKey] = tailscale_cli.DeviceInfo{
			DNSName:           "", // currently unused
			HasExitNodeOption: offersExitNode,
			IsSelf:            isSelf,
			Online:            isOnline,
			Tags:              nil, // currently unused
		}
		counter++
	}

	return devs, enrichedDevices, nil
}
