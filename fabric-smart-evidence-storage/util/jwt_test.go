package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJWT(t *testing.T) {
	username := "testuser"
	role := "admin"
	token, err := GenerateJWT(username, role)
	assert.NoError(t, err)
	usernameResult, roleResult, err := ParseJWT(token)
	assert.NoError(t, err)
	assert.Equal(t, username, usernameResult)
	assert.Equal(t, role, roleResult)
}
