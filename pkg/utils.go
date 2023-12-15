package pkg

import (
	"log"
	"tips/pkg/tailscale_cli"
)

func SelfDevice(devices map[string]tailscale_cli.DeviceInfo) tailscale_cli.DeviceInfo {
	for _, d := range devices {
		if d.IsSelf {
			return d
		}
	}

	log.Fatal("unexpected condition: if this func was invoked it must always return the `self` device")
	return tailscale_cli.DeviceInfo{}
}
