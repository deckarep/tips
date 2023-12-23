package pkg

import (
	"context"

	"github.com/charmbracelet/log"

	"github.com/tailscale/tailscale-client-go/tailscale"
)

func NewClient(ctx context.Context) *tailscale.Client {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)

	client, err := tailscale.NewClient(
		cfg.TailscaleAPI.ApiKey, // When doing oauth, this field must be blank!!!
		cfg.Tailnet,
		//tailscale.WithOAuthClientCredentials(oauthClientID, oauthClientSecret, nil),
		tailscale.WithUserAgent(UserAgent),
	)
	if err != nil {
		log.Fatal("failed to create tailscale api client with err: ", err)
	}
	return client
}

func NewOauthClient(ctx context.Context) *tailscale.Client {
	cfg := CtxAsConfig(ctx, CtxKeyConfig)
	client, err := tailscale.NewClient(
		"", // When doing oauth, this field must be blank!!!
		cfg.Tailnet,
		tailscale.WithOAuthClientCredentials(cfg.TailscaleAPI.OAuthClientID, cfg.TailscaleAPI.OAuthClientSecret, nil),
		tailscale.WithUserAgent(UserAgent),
	)
	if err != nil {
		log.Fatal("failed to create tailscale api client with err: ", err)
	}
	return client
}
