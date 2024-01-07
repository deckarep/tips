package cmd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/deckarep/tips/pkg"
)

func TestGetHosts(t *testing.T) {
	ctx := context.Background()
	cfgCtx := pkg.NewConfigCtx()
	ctx = context.WithValue(ctx, pkg.CtxKeyConfig, cfgCtx)

	// Regular user mode.
	tv := &pkg.GeneralTableView{
		ContextView: pkg.ContextView{},
		TailnetView: pkg.TailnetView{},
		SelfView:    pkg.SelfView{},
		Headers:     []string{},
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
