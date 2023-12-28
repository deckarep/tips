/*
Open Source Initiative OSI - The MIT License (MIT):Licensing

The MIT License (MIT)
Copyright (c) 2023 - 2024 Ralph Caraveo (deckarep@gmail.com)

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

	"github.com/tailscale/tailscale-client-go/tailscale"
)

func NewClient(ctx context.Context) *tailscale.Client {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)

	client, err := tailscale.NewClient(
		cfg.TailscaleAPI.ApiKey,
		cfg.Tailnet,
		tailscale.WithUserAgent(UserAgent),
	)
	if err != nil {
		log.Fatal("failed to create tailscale api client", "error", err)
	}
	return client
}

func NewOauthClient(ctx context.Context) *tailscale.Client {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)
	client, err := tailscale.NewClient(
		// When doing oauth, this field must be blank.
		"",
		cfg.Tailnet,
		tailscale.WithOAuthClientCredentials(cfg.TailscaleAPI.OAuthClientID, cfg.TailscaleAPI.OAuthClientSecret, nil),
		tailscale.WithUserAgent(UserAgent),
	)
	if err != nil {
		log.Fatal("failed to create tailscale api oauth client", "error", err)
	}
	return client
}
