package auth_test

import (
	"testing"
	"time"

	"github.com/Rohin-Gillgallon/GoJobTracker/internal/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret-key"

func TestGenerateTokenPair(t *testing.T) {
	tokens, err := auth.GenerateTokenPair("user-123", testSecret)

	require.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
	assert.NotEqual(t, tokens.AccessToken, tokens.RefreshToken)
}

func TestValidateToken_Valid(t *testing.T) {
	tokens, err := auth.GenerateTokenPair("user-123", testSecret)
	require.NoError(t, err)

	claims, err := auth.ValidateToken(tokens.AccessToken, testSecret)

	require.NoError(t, err)
	assert.Equal(t, "user-123", claims.UserID)
	assert.True(t, claims.ExpiresAt.After(time.Now()))
}

func TestValidateToken_InvalidSecret(t *testing.T) {
	tokens, err := auth.GenerateTokenPair("user-123", testSecret)
	require.NoError(t, err)

	_, err = auth.ValidateToken(tokens.AccessToken, "wrong-secret")

	assert.Error(t, err)
}

func TestValidateToken_MalformedToken(t *testing.T) {
	_, err := auth.ValidateToken("not.a.valid.token", testSecret)
	assert.Error(t, err)
}

func TestValidateToken_EmptyToken(t *testing.T) {
	_, err := auth.ValidateToken("", testSecret)
	assert.Error(t, err)
}
