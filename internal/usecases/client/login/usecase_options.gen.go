// Code generated by options-gen. DO NOT EDIT.
package login

import (
	fmt461e464ebed9 "fmt"

	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	authClient authClient,
	usrGetter userGetter,
	profilesRepo profilesRepository,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.authClient = authClient
	o.usrGetter = usrGetter
	o.profilesRepo = profilesRepo

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("authClient", _validate_Options_authClient(o)))
	errs.Add(errors461e464ebed9.NewValidationError("usrGetter", _validate_Options_usrGetter(o)))
	errs.Add(errors461e464ebed9.NewValidationError("profilesRepo", _validate_Options_profilesRepo(o)))
	return errs.AsError()
}

func _validate_Options_authClient(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.authClient, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `authClient` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_usrGetter(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.usrGetter, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `usrGetter` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_profilesRepo(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.profilesRepo, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `profilesRepo` did not pass the test: %w", err)
	}
	return nil
}
