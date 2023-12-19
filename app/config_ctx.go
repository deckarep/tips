package app

import (
	"github.com/charmbracelet/log"
	mapset "github.com/deckarep/golang-set/v2"
	"strconv"
	"strings"
	"time"
)

const (
	// PrimaryFilterAll basically means glob: *, but since this expands in the terminal we use @
	PrimaryFilterAll = "@"
)

type SliceCfg struct {
	From *int
	To   *int
}

func (s *SliceCfg) IsDefined() bool {
	return s != nil && (s.From != nil || s.To != nil)
}

func ParseSlice(s string) *SliceCfg {
	if strings.TrimSpace(s) == "" {
		return nil
	}

	sli := &SliceCfg{
		From: nil,
		To:   nil,
	}

	const (
		invalidSliceMsg = "the --slice flag or slice configuration is invalid; format is: [from:to] where from and to are both optional"
	)

	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "[") {
		log.Fatal(invalidSliceMsg)
	}

	if !strings.HasSuffix(s, "]") {
		log.Fatal(invalidSliceMsg)
	}

	colonIdx := strings.Index(s, ":")
	if colonIdx == -1 {
		log.Fatal(invalidSliceMsg)
	}

	fromNumStr := s[1:colonIdx]
	toNumStr := s[colonIdx+1 : len(s)-1]

	if n, err := strconv.Atoi(strings.TrimSpace(fromNumStr)); err == nil {
		sli.From = &n
	}

	if n, err := strconv.Atoi(strings.TrimSpace(toNumStr)); err == nil {
		sli.To = &n
	}

	return sli
}

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
	Columns       mapset.Set[string]
	Concurrency   int
	Filters       map[string]mapset.Set[string]
	NoCache       bool
	NoColor       bool
	PrimaryFilter string
	RemoteCmd     string
	Slice         *SliceCfg
	Tailnet       string
	TailscaleAPI  TailscaleAPICfgCtx
	TailscaleCLI  TailscaleCLICfgCtx
}

func NewConfigCtx() *ConfigCtx {
	return &ConfigCtx{}
}

func ParseColumns(s string) mapset.Set[string] {
	if len(strings.TrimSpace(s)) == 0 {
		return nil
	}

	m := mapset.NewSet[string]()

	parts := strings.Split(s, ",")
	for _, p := range parts {
		m.Add(strings.ToLower(strings.TrimSpace(p)))
	}

	return m
}
