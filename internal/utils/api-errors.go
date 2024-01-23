package utils

import (
	"github.com/go-playground/validator/v10"
)

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "url":
		return "Invalid url"
	case "email":
		return "Invalid email"
	case "max":
		return "Max length is " + fe.Param()
	}

	return fe.Error() // default error
}

func ConvertValidationErrorsToObject(err error) map[string]interface{} {
	validationErrors := err.(validator.ValidationErrors)

	var errorsObject = make(map[string]interface{})

	for _, err := range validationErrors {
		field := err.Field()

		if errorsObject[field] != nil {
			errorsObject[field] = append(errorsObject[field].([]string), msgForTag(err))
		} else {
			errorsObject[field] = []string{msgForTag(err)}
		}
	}

	return errorsObject
}
