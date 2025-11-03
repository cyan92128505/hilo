package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testPrivateKey = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIOChaSphj1MdLSxvU56h9vwmmpqdsQQF2alVwLKTj7dMoAoGCCqGSM49
AwEHoUQDQgAE7gMib5EUeW1An5VkkY4aU3xy+altlU3U0zn3FCO9Ffe/wwNUcUzp
XC9HWu76KhJnPpHczvZZv7Rro+kmqvN5tw==
-----END EC PRIVATE KEY-----`

func TestNewES256JWT(t *testing.T) {
	jwt, err := NewES256JWT(testPrivateKey)
	require.NoError(t, err)
	assert.NotNil(t, jwt)
	assert.NotNil(t, jwt.privateKey)
	assert.NotNil(t, jwt.publicKey)
}

func TestNewES256JWT_InvalidPEM(t *testing.T) {
	_, err := NewES256JWT("invalid pem")
	assert.Error(t, err)
}

func TestGenerateAndVerify_BasicClaims(t *testing.T) {
	jwt, err := NewES256JWT(testPrivateKey)
	require.NoError(t, err)

	claims := NewClaimsBuilder().
		WithSubject("test-subject").
		WithIssuer("test-issuer").
		WithAudience("test-audience").
		WithID("test-id-001").
		ExpiresAfter(1 * time.Hour).
		Build()

	common := NewCommon(claims)
	token, err := jwt.GenerateToken(common)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify token
	verifiedCommon := NewCommon(NewClaimsBuilder().Build())
	err = jwt.VerifyToken(token, verifiedCommon)
	require.NoError(t, err)

	// Assert claims
	assert.Equal(t, "test-subject", verifiedCommon.Subject)
	assert.Equal(t, "test-issuer", verifiedCommon.Issuer)
	assert.Equal(t, "test-id-001", verifiedCommon.ID)
	assert.Contains(t, verifiedCommon.Audience, "test-audience")
}

func TestGenerateAndVerify_WithPermissions(t *testing.T) {
	jwt, err := NewES256JWT(testPrivateKey)
	require.NoError(t, err)

	claims := NewClaimsBuilder().
		WithSubject("user-123").
		ExpiresAfter(1 * time.Hour).
		Build()

	common := NewCommon(claims, WithPermissions("/api/v1/messages", "/api/v1/users"))
	token, err := jwt.GenerateToken(common)
	require.NoError(t, err)

	// Verify with permissions
	verifiedCommon := NewCommon(NewClaimsBuilder().Build())
	err = jwt.VerifyToken(token, verifiedCommon)
	require.NoError(t, err)

	assert.Equal(t, "user-123", verifiedCommon.Subject)
	assert.Len(t, verifiedCommon.Permissions, 2)
	assert.Contains(t, verifiedCommon.Permissions, "/api/v1/messages")
	assert.Contains(t, verifiedCommon.Permissions, "/api/v1/users")
}

func TestVerifyToken_ExpiredToken(t *testing.T) {
	jwt, err := NewES256JWT(testPrivateKey)
	require.NoError(t, err)

	// Create already expired token
	claims := NewClaimsBuilder().
		WithSubject("expired-test").
		ExpiresAfter(-1 * time.Hour). // 1 hour ago
		Build()

	common := NewCommon(claims)
	token, err := jwt.GenerateToken(common)
	require.NoError(t, err)

	// Verify should fail
	verifiedCommon := NewCommon(NewClaimsBuilder().Build())
	err = jwt.VerifyToken(token, verifiedCommon)
	assert.Error(t, err)
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	jwt, err := NewES256JWT(testPrivateKey)
	require.NoError(t, err)

	common := NewCommon(NewClaimsBuilder().Build())
	err = jwt.VerifyToken("invalid.token.string", common)
	assert.Error(t, err)
}

func TestVerifyToken_NoExpiration(t *testing.T) {
	jwt, err := NewES256JWT(testPrivateKey)
	require.NoError(t, err)

	// Token without expiration should still work
	claims := NewClaimsBuilder().
		WithSubject("no-expiry").
		Build()

	common := NewCommon(claims)
	token, err := jwt.GenerateToken(common)
	require.NoError(t, err)

	verifiedCommon := NewCommon(NewClaimsBuilder().Build())
	err = jwt.VerifyToken(token, verifiedCommon)
	require.NoError(t, err)
	assert.Equal(t, "no-expiry", verifiedCommon.Subject)
}

func TestClaimsBuilder_AllMethods(t *testing.T) {
	now := time.Now()
	expiresAt := now.Add(2 * time.Hour)

	claims := NewClaimsBuilder().
		WithSubject("full-test").
		WithIssuer("test-issuer").
		WithAudience("aud1", "aud2").
		WithID("test-id-123").
		ExpiresAt(expiresAt).
		Build()

	assert.Equal(t, "full-test", claims.Subject)
	assert.Equal(t, "test-issuer", claims.Issuer)
	assert.Equal(t, "test-id-123", claims.ID)
	assert.Len(t, claims.Audience, 2)
	assert.Contains(t, claims.Audience, "aud1")
	assert.Contains(t, claims.Audience, "aud2")
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.IssuedAt)
	assert.NotNil(t, claims.NotBefore)
}

func TestCommon_WithSecret(t *testing.T) {
	claims := NewClaimsBuilder().WithSubject("secret-test").Build()
	common := NewCommon(claims, WithSecret("my-secret"))

	assert.Equal(t, "my-secret", common.Secret)
	assert.Equal(t, "secret-test", common.Subject)
}

func TestCommon_GetClaims(t *testing.T) {
	claims := NewClaimsBuilder().WithSubject("get-claims-test").Build()
	common := NewCommon(claims)

	retrievedClaims := common.GetClaims()
	assert.Equal(t, claims, retrievedClaims)
	assert.Equal(t, "get-claims-test", retrievedClaims.Subject)
}

// Benchmark tests
func BenchmarkGenerateToken(b *testing.B) {
	jwt, _ := NewES256JWT(testPrivateKey)
	claims := NewClaimsBuilder().
		WithSubject("benchmark").
		ExpiresAfter(1 * time.Hour).
		Build()
	common := NewCommon(claims)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = jwt.GenerateToken(common)
	}
}

func BenchmarkVerifyToken(b *testing.B) {
	jwt, _ := NewES256JWT(testPrivateKey)
	claims := NewClaimsBuilder().
		WithSubject("benchmark").
		ExpiresAfter(1 * time.Hour).
		Build()
	common := NewCommon(claims)
	token, _ := jwt.GenerateToken(common)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		verifiedCommon := NewCommon(NewClaimsBuilder().Build())
		_ = jwt.VerifyToken(token, verifiedCommon)
	}
}
