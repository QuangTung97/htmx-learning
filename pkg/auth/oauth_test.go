package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildState(t *testing.T) {
	state := buildOAuthState("some-token")
	assert.Equal(t, "csrf=some-token", state)

	assert.Equal(t, "some-token", getCSRFTokenFromState(state))
}
