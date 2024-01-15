package pkg

import (
	"context"
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"

	"github.com/deckarep/tips/pkg/tailscale_cli"

	"github.com/tailscale/tailscale-client-go/tailscale"
)

// WrappedDevice is a type that wraps the core `tailscale.Device` type. It also holds the joined
// `tailscale_cli.DeviceInfo` that may or may not be present when fetched from within the tailnet.
// It also implements the `Indexer` interface, so it may be stored in the DB.
type WrappedDevice struct {
	tailscale.Device
	EnrichedInfo *tailscale_cli.DeviceInfo `json:"enrichedInfo"`
}

// Key returns the device field of how this device gets indexed into the cached db.
// Currently, it just uses the name such: "blade.tail372c.ts.net" which implies devices are stored in alphabetical
// order as ascending via their `name` field.
func (w *WrappedDevice) Key() string {
	return w.Name
}

// EvalColumnField is invoked for each "column" requested per device field. This code was built purposely to be dynamic
// and if it gets more complex it may be worthwhile to break the code up further into discreet functions per field.
// One additional thing I've been considering is the memoize of any redundant "heavy" work but so far there is none here.
func (w *WrappedDevice) EvalColumnField(ctx context.Context, idx int, headerMatchName HeaderMatchName) string {
	//cfg := CtxAsConfig(ctx, CtxKeyConfig)
	enriched := w.EnrichedInfo != nil

	// Safely return an address at index.
	var addrAtIndex = func(idx int) string {
		if idx >= 0 && idx < len(w.Addresses) {
			return w.Addresses[idx]
		}
		return "n/a"
	}

	switch headerMatchName {
	case MatchNameAddress:
		addrs := strings.Join(w.Addresses, ", ")
		return addrs
	case MatchNameIpv4:
		return addrAtIndex(0)
	case MatchNameAuthorized:
		return fmt.Sprintf("%t", w.Authorized)
	case MatchNameBlocksIncomingConnections:
		return fmt.Sprintf("%t", w.BlocksIncomingConnections)
	case MatchNameClientVersion:
		return w.ClientVersion
	case MatchNameExitStatus:
		exitStatus := "no"
		if enriched && w.EnrichedInfo.HasExitNodeOption {
			exitStatus = checkField
		}
		return exitStatus
	case MatchNameFullname:
		return w.Name
	case MatchNameHostname:
		return w.Hostname
	case MatchNameIpv6:
		return addrAtIndex(1)
	case MatchNameLastSeen:
		return fmt.Sprintf("%s", w.LastSeen)
	case MatchNameName, MatchNameMachine:
		name := strings.Split(w.Name, ".")[0]
		return name
	case MatchNameNo:
		no := fmt.Sprintf("%04d", idx)
		return no
	case MatchNameOS:
		return w.OS
	case MatchNameLastSeenAgo:
		lastSeenAgo := humanize.Time(w.LastSeen.Time)
		if enriched && w.EnrichedInfo.Online {
			lastSeenAgo = nowField
		}
		return lastSeenAgo
	case MatchNameTags:
		// Remove all tag: prefixes, and join the tags as a comma delimited string.
		tags := strings.Replace(strings.Join(w.Tags, ", "), "tag:", "", -1)
		return tags
	case MatchNameUser:
		return w.User
	case MatchNameVersion:
		version := fmt.Sprintf("%s - %s", strings.Split(w.ClientVersion, "-")[0], w.OS)
		return version
	default:
		panic(`unknown MatchName column requested (a new MatchName filed was likely introduced but not 
handled here): ` + headerMatchName)
	}
}
