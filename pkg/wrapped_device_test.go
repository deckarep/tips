package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tailscale/tailscale-client-go/tailscale"
)

func TestWrappedDevice_Key(t *testing.T) {
	w := &WrappedDevice{
		Device: tailscale.Device{
			Name:          "pleebus.serv",
			User:          "user@gmail.com",
			Tags:          []string{"foo", "bar", "baz"},
			Addresses:     []string{"127.0.0.1"},
			OS:            "rasbarbarian",
			ClientVersion: "1.23.45-deadbeef",
		}}

	// Key should just return whatever is inside the wrapped device Name field.
	assert.Equal(t, w.Key(), "pleebus.serv")
}
