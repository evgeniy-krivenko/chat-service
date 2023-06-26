// Code generated by options-gen. DO NOT EDIT.
package server

import (
	fmt461e464ebed9 "fmt"

	"github.com/getkin/kin-openapi/openapi3"
	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	logger *zap.Logger,
	addr string,
	allowOrigins []string,
	v1Swagger *openapi3.T,
	keycloakClient KeycloakClient,
	resource string,
	role string,
	httpErrorHandler echo.HTTPErrorHandler,
	registerHandlers func(*echo.Group),
	shutdown func(),
	wsHandler wsHTTPHandler,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.logger = logger
	o.addr = addr
	o.allowOrigins = allowOrigins
	o.v1Swagger = v1Swagger
	o.keycloakClient = keycloakClient
	o.resource = resource
	o.role = role
	o.httpErrorHandler = httpErrorHandler
	o.registerHandlers = registerHandlers
	o.shutdown = shutdown
	o.wsHandler = wsHandler

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("logger", _validate_Options_logger(o)))
	errs.Add(errors461e464ebed9.NewValidationError("addr", _validate_Options_addr(o)))
	errs.Add(errors461e464ebed9.NewValidationError("allowOrigins", _validate_Options_allowOrigins(o)))
	errs.Add(errors461e464ebed9.NewValidationError("v1Swagger", _validate_Options_v1Swagger(o)))
	errs.Add(errors461e464ebed9.NewValidationError("keycloakClient", _validate_Options_keycloakClient(o)))
	errs.Add(errors461e464ebed9.NewValidationError("resource", _validate_Options_resource(o)))
	errs.Add(errors461e464ebed9.NewValidationError("role", _validate_Options_role(o)))
	errs.Add(errors461e464ebed9.NewValidationError("httpErrorHandler", _validate_Options_httpErrorHandler(o)))
	errs.Add(errors461e464ebed9.NewValidationError("registerHandlers", _validate_Options_registerHandlers(o)))
	errs.Add(errors461e464ebed9.NewValidationError("shutdown", _validate_Options_shutdown(o)))
	errs.Add(errors461e464ebed9.NewValidationError("wsHandler", _validate_Options_wsHandler(o)))
	return errs.AsError()
}

func _validate_Options_logger(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.logger, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `logger` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_addr(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.addr, "required,hostname_port"); err != nil {
		return fmt461e464ebed9.Errorf("field `addr` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_allowOrigins(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.allowOrigins, "min=1"); err != nil {
		return fmt461e464ebed9.Errorf("field `allowOrigins` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_v1Swagger(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.v1Swagger, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `v1Swagger` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_keycloakClient(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.keycloakClient, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `keycloakClient` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_resource(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.resource, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `resource` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_role(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.role, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `role` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_httpErrorHandler(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.httpErrorHandler, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `httpErrorHandler` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_registerHandlers(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.registerHandlers, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `registerHandlers` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_shutdown(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.shutdown, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `shutdown` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_wsHandler(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.wsHandler, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `wsHandler` did not pass the test: %w", err)
	}
	return nil
}
