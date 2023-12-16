package app

import (
	"github.com/charmbracelet/log"
	"github.com/tailscale/tailscale-client-go/tailscale"
	"tips/pkg"
)

func NewClient(apiKey string, tailnet string) *tailscale.Client {
	client, err := tailscale.NewClient(
		apiKey, // When doing oauth, this field must be blank!!!
		tailnet,
		//tailscale.WithOAuthClientCredentials(oauthClientID, oauthClientSecret, nil),
		tailscale.WithUserAgent(pkg.UserAgent),
	)
	if err != nil {
		log.Fatal("failed to create tailscale api client with err: ", err)
	}
	return client
}

func NewOauthClient(oauthClientID, oauthClientSecret, tailnet string) *tailscale.Client {
	client, err := tailscale.NewClient(
		"", // When doing oauth, this field must be blank!!!
		tailnet,
		tailscale.WithOAuthClientCredentials(oauthClientID, oauthClientSecret, nil),
		tailscale.WithUserAgent(pkg.UserAgent),
	)
	if err != nil {
		log.Fatal("failed to create tailscale api client with err: ", err)
	}
	return client
}
