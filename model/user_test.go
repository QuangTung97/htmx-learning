package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserStatus(t *testing.T) {
	assert.Equal(t, UserStatus(1), UserStatusActive)
}
