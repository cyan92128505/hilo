package jwt

import (
	"crypto/ecdsa"
	"errors"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	jwtPkg "github.com/golang-jwt/jwt"
)

var (
	ErrNoKey = errors.New("not provider jwt private/secret key")
	// ErrTokenExpired variable
	ErrTokenExpired = errors.New("token is expired")
	// ErrParsePrivateKey variable
	ErrParsePrivateKey = errors.New("parse private key error")
)

var (
	now                      = time.Now
	newSigner                = jose.NewSigner
	parseECPrivateKeyFromPEM = jwtPkg.ParseECPrivateKeyFromPEM
	signed                   = jwt.Signed
	parseSigned              = jwt.ParseSigned
)

// IJWT interface
type IJWT interface {
	GenerateToken(claims IJWTClaims) (string, error)
	Validate(raw string) (err error)
	VerifyToken(token string, claims IJWTClaims) (err error)
	RefreshToken(token string, claims IJWTClaims, duration time.Duration) (string, error)
}

func ParseUnverified(raw string, claims IJWTClaims) error {
	tok, errParse := parseSigned(raw)
	if errParse != nil {
		return errParse
	}
	return tok.UnsafeClaimsWithoutVerification(claims)
}

func parseESRaw(signingKey *ecdsa.PrivateKey, raw string, claims IJWTClaims) error {
	tok, errParse := parseSigned(raw)
	if errParse != nil {
		return errParse
	}

	errClaims := tok.Claims(signingKey.Public(), claims)
	if errClaims != nil {
		return errClaims
	}
	return checkExpire(claims)
}

func checkExpire(claims IJWTClaims) error {
	if instance, ok := claims.(IJWTExpire); ok && instance.GetExpiresAfter() != nil && now().UnixNano() > instance.GetExpiresAfter().Time().UnixNano() {
		return ErrTokenExpired
	}
	return nil
}
