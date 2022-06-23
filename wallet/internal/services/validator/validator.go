package validator

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

// Service validates structures.
type Service struct {
	validate *validator.Validate
}

func NewValidator() *Service {
	v := validator.New()
	return &Service{validate: v}
}

func (s *Service) Validate(obj interface{}) error {
	value := reflect.ValueOf(obj)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}

	if valueType == reflect.Struct {
		if err := s.validate.Struct(obj); err != nil {
			return errors.Wrap(err, "Error from validator")
		}
	}
	return nil
}
