package helpers

import (
	"errors"
	"fmt"
	"log"

	"github.com/go-playground/validator"
)

func Validation(model interface{}) (error, int16) {
	var validate = validator.New()
	var count int16
	var errormessage string
	err := validate.Struct(model)
	if err != nil {
		var ErrorMSg validator.ValidationErrors
		errors.As(err, &ErrorMSg)
		for _, validationError := range ErrorMSg {
			errormessage += fmt.Sprintf("Field %s is must required %s \n ", validationError.Field(), validationError.Tag())
			log.Println("Validation Error:", ErrorMSg)
			count++
			if count > 0 {

				return err, count
			}
		}
	}
	log.Println("\n Count is :", count)
	return nil, count
}
