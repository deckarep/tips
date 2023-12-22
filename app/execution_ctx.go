package app

import (
	"context"

	"github.com/charmbracelet/log"
)

type contextKey string

func (c contextKey) String() string {
	return "mypackage context key " + string(c)
}

var (
	// CtxKeyConfig holds all config settings that were resolved from the environment/config file/cli flags
	CtxKeyConfig    = contextKey("configuration")
	CtxKeyUserQuery = contextKey("user-query")
)

func CtxAsString(ctx context.Context, key contextKey) string {
	if val, ok := ctx.Value(key).(string); ok {
		return val
	}
	log.Fatalf("failed to get context value as a string with key: %s", key)
	return ""
}

func CtxAsBool(ctx context.Context, key contextKey) bool {
	if val, ok := ctx.Value(key).(bool); ok {
		return val
	}
	log.Fatalf("failed to get context value as a bool with key: %s", key)
	return false
}

func CtxAsInt(ctx context.Context, key contextKey) int {
	if val, ok := ctx.Value(key).(int); ok {
		return val
	}
	log.Fatalf("failed to get context value as an int with key: %s", key)
	return 0
}

func CtxAsConfig(ctx context.Context, key contextKey) *ConfigCtx {
	if val, ok := ctx.Value(key).(*ConfigCtx); ok {
		return val
	}
	log.Fatalf("failed to get context value as an int with key: %s", key)
	return &ConfigCtx{}
}
