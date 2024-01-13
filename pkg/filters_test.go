package pkg

import (
	"context"
	"testing"

	"github.com/deckarep/tips/pkg/tailscale_cli"

	"github.com/tailscale/tailscale-client-go/tailscale"

	"github.com/stretchr/testify/assert"
)

func TestParseFilter(t *testing.T) {
	// Note: much more robust testing of the parser occurs in the filtercomp package.
	// It makes this testing somewhat redundant.

	// A filter expression with imbalanced parenthesis should return an error.
	_, err := ParseFilter("(hello, world")
	assert.Error(t, err, "imbalanced parenthesis detected")

	// An empty filter should return no error with a nil ast.
	ast, err := ParseFilter("")
	assert.NoError(t, err)
	assert.Nil(t, ast)

	// An empty filter should return no error with a nil ast.
	_, err = ParseFilter("how are you doing?")
	assert.Error(t, err, "parser did not run to completion, tokens were not fully consumed")
}

func TestExecuteFilters(t *testing.T) {
	ctx := context.Background()
	cfg := NewConfigCtx()

	ast, err := ParseFilter("user@gmail.com")
	assert.NoError(t, err)

	cfg.Filters = ast

	ctx = context.WithValue(ctx, CtxKeyConfig, cfg)

	devs := []*WrappedDevice{
		{Device: tailscale.Device{
			User:          "user@gmail.com",
			Tags:          []string{"foo", "bar", "baz"},
			Addresses:     []string{"127.0.0.1"},
			OS:            "rasbarbarian",
			ClientVersion: "1.23.45-deadbeef",
		}},
		{Device: tailscale.Device{
			User:          "user@gmail.com",
			Tags:          []string{"foo", "biz", "bang"},
			Addresses:     []string{"127.0.0.2"},
			OS:            "loonix",
			ClientVersion: "1.23.46-deadbeef",
		}},
		{Device: tailscale.Device{
			User:          "user@gmail.com",
			Tags:          []string{"poo", "par", "paz"},
			Addresses:     []string{"127.0.0.3"},
			OS:            "windoze",
			ClientVersion: "1.23.2-deadbeef",
		}},
		{Device: tailscale.Device{
			User:          "user@gmail.com",
			Tags:          []string{"voo", "foo", "var", "vaz"},
			Addresses:     []string{"127.0.0.4"},
			OS:            "bigmacos",
			ClientVersion: "1.23.1-deadbeef",
		}, EnrichedInfo: &tailscale_cli.DeviceInfo{
			HasExitNodeOption: true,
		}},
	}

	// Apply Users filter.
	filteredResults := executeFilters(ctx, devs)
	assert.NotNil(t, filteredResults)
	assert.Equal(t, len(filteredResults), 4)

	// Apply Tags filter.
	ast, err = ParseFilter("foo")
	assert.NoError(t, err)
	cfg.Filters = ast

	filteredResults = executeFilters(ctx, devs)
	assert.NotNil(t, filteredResults)
	assert.Equal(t, len(filteredResults), 3)

	// Apply Tags filter.
	ast, err = ParseFilter("127.0.0.4")
	assert.NoError(t, err)
	cfg.Filters = ast

	filteredResults = executeFilters(ctx, devs)
	assert.NotNil(t, filteredResults)
	assert.Equal(t, len(filteredResults), 1)

	// Apply OS filter.
	ast, err = ParseFilter("(windoze | loonix)")
	assert.NoError(t, err)
	cfg.Filters = ast

	filteredResults = executeFilters(ctx, devs)
	assert.NotNil(t, filteredResults)
	assert.Equal(t, len(filteredResults), 2)

	// Apply Version filter.
	ast, err = ParseFilter("1.23.4*")
	assert.NoError(t, err)
	cfg.Filters = ast

	filteredResults = executeFilters(ctx, devs)
	assert.NotNil(t, filteredResults)
	assert.Equal(t, len(filteredResults), 2)

	// Exit node filter.
	ast, err = ParseFilter("+exit")
	assert.NoError(t, err)
	cfg.Filters = ast

	filteredResults = executeFilters(ctx, devs)
	assert.NotNil(t, filteredResults)
	assert.Equal(t, len(filteredResults), 1)
}
