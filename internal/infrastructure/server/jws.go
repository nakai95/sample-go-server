package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/golang-jwt/jwt/v5"
)

type JWSValidator interface {
	Authenticate(c context.Context, input *openapi3filter.AuthenticationInput) error
}

type jwsAuthenticator struct {
	SecretKey string
}

var (
	ErrNoAuthHeader            = errors.New("authorization header is missing")
	ErrInvalidAuthHeader       = errors.New("authorization header is malformed")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrParsingJWS              = errors.New("parsing JWS failed")
	ErrClaimsInvalid           = errors.New("provided claims do not match expected scopes")
)

func GetJWSFromRequest(req *http.Request) (string, error) {
	authHdr := req.Header.Get("Authorization")
	// Check for the Authorization header.
	if authHdr == "" {
		return "", ErrNoAuthHeader
	}
	// We expect a header value of the form "Bearer <token>", with 1 space after
	// Bearer, per spec.
	prefix := "Bearer "
	if !strings.HasPrefix(authHdr, prefix) {
		return "", ErrInvalidAuthHeader
	}
	return strings.TrimPrefix(authHdr, prefix), nil
}

func NewAuthenticator(secretKey string) JWSValidator {
	return &jwsAuthenticator{
		SecretKey: secretKey,
	}
}

func (j *jwsAuthenticator) Authenticate(c context.Context, input *openapi3filter.AuthenticationInput) error {
	// Our security scheme is named BearerAuth, ensure this is the case
	if input.SecuritySchemeName != "BearerAuth" {
		return fmt.Errorf("security scheme %s != 'BearerAuth'", input.SecuritySchemeName)
	}

	// Now, we need to get the JWS from the request, to match the request expectations
	// against request contents.
	jws, err := GetJWSFromRequest(input.RequestValidationInput.Request)
	if err != nil {
		return fmt.Errorf("getting jws: %w", err)
	}

	// if the JWS is valid, we have a JWT, which will contain a bunch of claims.
	token, err := j.ValidateJWS(jws)
	if err != nil {
		return fmt.Errorf("validating JWS: %w", err)
	}

	// Check that the claims are valid
	err = j.CheckClaims(token)
	if err != nil {
		return fmt.Errorf("checking claims: %w", err)
	}

	return nil
}

func (j *jwsAuthenticator) ValidateJWS(jws string) (*jwt.Token, error) {
	token, err := jwt.Parse(jws, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return nil, ErrParsingJWS
	}
	return token, nil
}

func (j *jwsAuthenticator) CheckClaims(token *jwt.Token) error {
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	} else {
		return ErrClaimsInvalid
	}
}
