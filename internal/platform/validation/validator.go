package validation

import (
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

// Get returns the singleton validator instance.
func Get() *validator.Validate {
	once.Do(func() {
		validate = validator.New()

		// Register custom JSON tag extractor to map error fields correctly
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			if name == "" {
				return fld.Name
			}
			return name
		})

		// Example Custom Validation: register "agency_code" if needed
		// validate.RegisterValidation("agency_code", customAgencyCodeValidator)
	})
	return validate
}

// Struct validates a struct and returns strongly typed validation errors if they exist.
func Struct(s interface{}) error {
	v := Get()
	return v.Struct(s)
}
