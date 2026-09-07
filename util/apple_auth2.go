package util

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/royanqodri/Login-Gateway-API/config"
)

const (
	appleKeysURL = "https://appleid.apple.com/auth/keys"
	appleIssuer  = "https://appleid.apple.com"
)

type AppleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

type appleJWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type appleJWKSet struct {
	Keys []appleJWK `json:"keys"`
}

func fetchApplePublicKey(kid string) (*rsa.PublicKey, error) {
	resp, err := http.Get(appleKeysURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch apple public keys: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("apple keys endpoint returned status %d", resp.StatusCode)
	}

	var jwkSet appleJWKSet
	if err := json.NewDecoder(resp.Body).Decode(&jwkSet); err != nil {
		return nil, fmt.Errorf("failed to decode apple jwks: %w", err)
	}

	for _, key := range jwkSet.Keys {
		if key.Kid == kid {
			return buildRSAPublicKey(key.N, key.E)
		}
	}

	return nil, fmt.Errorf("no matching apple public key found for kid: %s", kid)
}

func buildRSAPublicKey(nStr, eStr string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(nStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %w", err)
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(eStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %w", err)
	}

	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	return &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}, nil
}

func VerifyAppleIDToken(identityToken string) (*AppleUserInfo, error) {
	if identityToken == "" {
		return nil, fmt.Errorf("identity_token is required")
	}

	clientID := config.Get().Apple.ClientId
	if clientID == "" {
		return nil, fmt.Errorf("APPLE_CLIENT_ID is not configured")
	}

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(identityToken, claims, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok || kid == "" {
			return nil, fmt.Errorf("kid not found in token header")
		}

		return fetchApplePublicKey(kid)
	})
	if err != nil {
		return nil, fmt.Errorf("invalid apple identity token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("apple identity token is not valid")
	}

	iss, _ := claims["iss"].(string)
	if iss != appleIssuer {
		return nil, fmt.Errorf("invalid token issuer: %s", iss)
	}

	aud, _ := claims["aud"].(string)
	if aud != clientID {
		return nil, fmt.Errorf("invalid token audience: %s", aud)
	}

	// Validasi expired
	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().After(time.Unix(int64(exp), 0)) {
		return nil, fmt.Errorf("apple identity token is expired")
	}

	sub, _ := claims["sub"].(string)
	if sub == "" {
		return nil, fmt.Errorf("sub not found in apple identity token")
	}

	email, _ := claims["email"].(string)
	if email == "" {
		return nil, fmt.Errorf("email not found in apple identity token")
	}

	emailVerified := parseAppleBool(claims["email_verified"])
	if !emailVerified {
		return nil, fmt.Errorf("email is not verified by apple")
	}

	return &AppleUserInfo{
		Sub:           sub,
		Email:         email,
		EmailVerified: emailVerified,
	}, nil
}

func parseAppleBool(v any) bool {
	switch val := v.(type) {
	case bool:
		return val
	case string:
		return val == "true"
	default:
		return false
	}
}
