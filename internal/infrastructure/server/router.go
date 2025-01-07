package server

import (
	"fmt"
	"log"
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

	// Create middleware for the authenticator
	mw, err := CreateMiddleware(at)
	if err != nil {
		log.Fatalln("error creating middleware:", err)

	}

	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.CORS()) // CORS don't work if done after e.Use(mw...)
	e.Use(mw...)

	// Create an instance of our handler which satisfies the generated interface
	h := NewHandler()

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
