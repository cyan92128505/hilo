package jwt

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents standard JWT claims
type Claims struct {
	jwt.RegisteredClaims
}

// IJWTClaims interface for custom claim wrappers
type IJWTClaims interface {
	GetClaims() *Claims
}

// IJWT interface for JWT operations
type IJWT interface {
	GenerateToken(claims IJWTClaims) (string, error)
	VerifyToken(token string, claims IJWTClaims) error
}

// ES256JWT implements IJWT using ECDSA P-256
type ES256JWT struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

// NewES256JWT creates JWT handler from PEM-encoded private key
func NewES256JWT(privateKeyPEM string) (*ES256JWT, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &ES256JWT{
		privateKey: key,
		publicKey:  &key.PublicKey,
	}, nil
}

// NewES256JWTFromOptions creates JWT handler from config options
func NewES256JWTFromOptions(config interface{}) (*ES256JWT, error) {
	// Type assertion for config struct with EcdsaPrivateKey field
	type jwtConfig interface {
		GetEcdsaPrivateKey() string
	}

	if cfg, ok := config.(jwtConfig); ok {
		return NewES256JWT(cfg.GetEcdsaPrivateKey())
	}

	// Fallback for direct string field access via reflection
	// This maintains backward compatibility with existing config structures
	return nil, errors.New("invalid config type: must have EcdsaPrivateKey field")
}

func (j *ES256JWT) GenerateToken(claims IJWTClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims.GetClaims())
	return token.SignedString(j.privateKey)
}

func (j *ES256JWT) VerifyToken(tokenString string, claims IJWTClaims) error {
	token, err := jwt.ParseWithClaims(tokenString, claims.GetClaims(), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.publicKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}

// ClaimsBuilder helps construct Claims with fluent API
type ClaimsBuilder struct {
	claims Claims
}

// NewClaimsBuilder creates a new builder with default values
func NewClaimsBuilder() *ClaimsBuilder {
	now := time.Now()
	return &ClaimsBuilder{
		claims: Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
			},
		},
	}
}

func (b *ClaimsBuilder) WithSubject(subject string) *ClaimsBuilder {
	b.claims.Subject = subject
	return b
}

func (b *ClaimsBuilder) WithIssuer(issuer string) *ClaimsBuilder {
	b.claims.Issuer = issuer
	return b
}

func (b *ClaimsBuilder) WithAudience(audience ...string) *ClaimsBuilder {
	b.claims.Audience = audience
	return b
}

func (b *ClaimsBuilder) WithID(id string) *ClaimsBuilder {
	b.claims.ID = id
	return b
}

func (b *ClaimsBuilder) ExpiresAfter(duration time.Duration) *ClaimsBuilder {
	b.claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(duration))
	return b
}

func (b *ClaimsBuilder) ExpiresAt(t time.Time) *ClaimsBuilder {
	b.claims.ExpiresAt = jwt.NewNumericDate(t)
	return b
}

func (b *ClaimsBuilder) Build() *Claims {
	return &b.claims
}

// Common is a wrapper for Claims to implement IJWTClaims
type Common struct {
	*Claims
	Permissions []string `json:"perms,omitempty"`
	Secret      string   `json:"s,omitempty"`
}

// NewCommon creates a Common claim wrapper
func NewCommon(claims *Claims, opts ...CommonOption) *Common {
	c := &Common{Claims: claims}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// GetClaims implements IJWTClaims
func (c *Common) GetClaims() *Claims {
	return c.Claims
}

// CommonOption is a functional option for Common
type CommonOption func(*Common)

// WithPermissions adds permissions to claims
func WithPermissions(perms ...string) CommonOption {
	return func(c *Common) {
		c.Permissions = perms
	}
}

// WithSecret adds secret to claims (legacy support)
func WithSecret(secret string) CommonOption {
	return func(c *Common) {
		c.Secret = secret
	}
}

// NewNumericDate creates jwt.NumericDate from time.Time
func NewNumericDate(t time.Time) *jwt.NumericDate {
	return jwt.NewNumericDate(t)
}
