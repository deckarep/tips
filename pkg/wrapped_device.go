package pkg

import (
	"tips/pkg/tailscale_cli"

	"github.com/tailscale/tailscale-client-go/tailscale"
)

// WrappedDevice is a type that wraps the core `tailscale.Device` type. It also holds the joined
// `tailscale_cli.DeviceInfo` that may or may not be returned when fetched from within the tailnet.
// It also implements the `Indexer` interface, so it may be stored in the DB.
type WrappedDevice struct {
	tailscale.Device
	EnrichedInfo *tailscale_cli.DeviceInfo `json:"enrichedInfo"`
}

func (w *WrappedDevice) Key() string {
	return w.Name
}
