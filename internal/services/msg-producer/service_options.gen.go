// Code generated by options-gen. DO NOT EDIT.
package msgproducer

import (
	fmt461e464ebed9 "fmt"

	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	wr KafkaWriter,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.wr = wr

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func WithEncryptKey(opt string) OptOptionsSetter {
	return func(o *Options) {
		o.encryptKey = opt
	}
}

func WithNonceFactory(opt func(size int) ([]byte, error)) OptOptionsSetter {
	return func(o *Options) {
		o.nonceFactory = opt
	}
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("wr", _validate_Options_wr(o)))
	errs.Add(errors461e464ebed9.NewValidationError("encryptKey", _validate_Options_encryptKey(o)))
	return errs.AsError()
}

func _validate_Options_wr(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.wr, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `wr` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_encryptKey(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.encryptKey, "omitempty,hexadecimal"); err != nil {
		return fmt461e464ebed9.Errorf("field `encryptKey` did not pass the test: %w", err)
	}
	return nil
}