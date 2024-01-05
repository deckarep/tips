/*
Open Source Initiative OSI - The MIT License (MIT):Licensing

The MIT License (MIT)
Copyright Ralph Caraveo (deckarep@gmail.com)

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package pkg

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
