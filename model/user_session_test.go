package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserSessionStatus(t *testing.T) {
	assert.Equal(t, UserSessionStatus(1), UserSessionStatusActive)
	assert.Equal(t, UserSessionStatus(2), UserSessionStatusDeleted)

}
