package cmd

import (
	"context"
	"testing"

	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"

	"github.com/deckarep/tips/pkg"
)

func TestPackageCfg(t *testing.T) {
	// Test error cases first.
	args := []string{"@"}

	_, err := packageCfg(args)
	assert.Error(t, err, "it has an api_key_error")

	// Test non-error case.
	args = []string{"@", "echo 'hello world'", "&& sleep 0.5", "&& ps aux | grep foo"}

	viper.Set("tips_api_key", "foo")
	viper.Set("tailnet", "bar")
	cfg, err := packageCfg(args)

	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.True(t, cfg.PrefixFilter.IsAll())
	assert.Equal(t, cfg.RemoteCmd, "echo 'hello world' && sleep 0.5 && ps aux | grep foo")
}

func TestGetHosts(t *testing.T) {
	ctx := context.Background()
	cfgCtx := pkg.NewConfigCtx()
	ctx = context.WithValue(ctx, pkg.CtxKeyConfig, cfgCtx)

	// Regular user mode.
	tv := &pkg.GeneralTableView{
		ContextView: pkg.ContextView{},
		TailnetView: pkg.TailnetView{},
		Self:        &pkg.SelfView{},
		Headers:     []pkg.Header{},
		Rows: [][]string{
			{"a", "a1", "a2", "a3", "a4", "a5"},
			{"b", "b1", "b2", "b3", "b4", "b5"},
		},
	}

	hostList := getHosts(ctx, tv)

	assert.NotNil(t, hostList)
	assert.Equal(t, hostList, []pkg.RemoteCmdHost{
		{Original: "a1", Alias: ""},
		{Original: "b1", Alias: ""},
	})

	// Test mode: uses a pseudo test-server named blade to simulate certain test scenarios.
	cfgCtx.TestMode = true
	hostList = getHosts(ctx, tv)

	assert.NotNil(t, hostList)
	assert.Equal(t, hostList, []pkg.RemoteCmdHost{
		{Original: "blade", Alias: "a1"},
		{Original: "blade", Alias: "b1"},
	})
}
