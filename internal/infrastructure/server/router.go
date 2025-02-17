package server

import (
	"fmt"
	"os"
	"sample-go-server/api"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"

	echomiddleware "github.com/labstack/echo/v4/middleware"
	middleware "github.com/oapi-codegen/echo-middleware"
)

func NewRouter() *echo.Echo {
	e := echo.New()

	const secret string = "secret"

	// Create a authenticator
	at := NewAuthenticator(secret)

	e.Use(echomiddleware.Logger())

	// Create middleware for the authenticator
	mw, err := CreateMiddleware(at)
	if err != nil {
		e.Logger.Fatal("error creating middleware:", err)
		os.Exit(1)
	}

	e.Use(echomiddleware.CORS()) // CORS don't work if done after e.Use(mw...)
	e.Use(mw...)

	// Create an instance of our handler which satisfies the generated interface
	h, err := NewHandler()
	if err != nil {
		e.Logger.Fatal("error creating handler:", err)
		os.Exit(1)
	}

	// We now register our petStore above as the handler for the interface
	api.RegisterHandlers(e, h)

	return e
}

func CreateMiddleware(v JWSValidator) ([]echo.MiddlewareFunc, error) {
	swagger, err := api.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("loading spec: %w", err)
	}

	// Clear out the servers array in the swagger spec
	swagger.Servers = nil

	validator := middleware.OapiRequestValidatorWithOptions(swagger,
		&middleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: v.Authenticate,
			},
		})

	return []echo.MiddlewareFunc{validator}, nil
}
