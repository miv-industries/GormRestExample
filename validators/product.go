package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/miv-industries/GormRestExample/models"
)

func ValidateProduct(product models.Product) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()

	err := validate.Struct(product)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
