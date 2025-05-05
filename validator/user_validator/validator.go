package uservalidator

import (
	"fmt"
	"regexp"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"task1/service/userservice"
)

type Validator struct {
}

func New() Validator {

	return Validator{}
}
func (v Validator) ValidateRegisterRequest(req userservice.SignUpRequest) (error, map[string]string) {

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Email, validation.Required, is.Email),
		validation.Field(&req.Password, validation.Required, validation.Match(regexp.MustCompile("^[A-Za-z0-9]{8,}$"))),
	); err != nil {

		fieldErrors := make(map[string]string)

		errV, ok := err.(validation.Errors)
		if ok {
			for key, value := range errV {
				if value != nil {

					fieldErrors[key] = value.Error()

				}
			}
		}
		return fmt.Errorf("Invalid Input"), fieldErrors
	}

	return nil, nil

}
