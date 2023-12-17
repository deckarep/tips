package app

import (
	mapset "github.com/deckarep/golang-set/v2"
	"time"
)

type TailscaleAPICfgCtx struct {
	Timeout time.Duration

	// ApiKey for regular authentication
	ApiKey string

	// OAuthClientID for OAuth based login.
	OAuthClientID string
	// OAuthClientSecret for Oauth based login.
	OAuthClientSecret string

	// ElapsedTime records the time this API call took. It's meant to be mutated during the API call and populated then.
	ElapsedTime time.Duration
}

type TailscaleCLICfgCtx struct {
}

type ConfigCtx struct {
	NoCache      bool
	Filters      map[string]mapset.Set[string]
	Tailnet      string
	TailscaleAPI TailscaleAPICfgCtx
	TailscaleCLI TailscaleCLICfgCtx
}

func NewConfigCtx() *ConfigCtx {
	return &ConfigCtx{}
}
