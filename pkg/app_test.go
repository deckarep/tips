package pkg

import "testing"

func TestUserAgent(t *testing.T) {
	const (
		expectedUserAgent = "tips/0.0.1"
	)
	if UserAgent != expectedUserAgent {
		t.Errorf("expected UserAgent to be: %s, got: %s", expectedUserAgent, UserAgent)
	}
}
