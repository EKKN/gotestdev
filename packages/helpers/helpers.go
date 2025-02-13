package helpers

import (
	"fmt"
	"reflect"

	"github.com/EKKN/gotestdev/packages/config"
	"github.com/EKKN/gotestdev/packages/models"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func getJSONFieldName(structType reflect.Type, fieldName string) string {
	if field, found := structType.FieldByName(fieldName); found {
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			return jsonTag
		}
	}
	return fieldName
}

func ValidateStruct(validateStruct interface{}) []map[string]string {
	var validationErrors []map[string]string

	structType := reflect.TypeOf(validateStruct)

	err := validate.Struct(validateStruct)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldError := make(map[string]string)

			jsonFieldName := getJSONFieldName(structType, err.StructField())

			switch err.Tag() {
			case "required":
				fieldError[jsonFieldName] = fmt.Sprintf("%s is required", err.Field())
			case "gte":
				fieldError[jsonFieldName] = fmt.Sprintf("%s must be greater than or equal to %s", err.Field(), err.Param())
			case "lte":
				fieldError[jsonFieldName] = fmt.Sprintf("%s must be less than or equal to %s", err.Field(), err.Param())
			case "datetime":
				fieldError[jsonFieldName] = fmt.Sprintf("%s must be a valid date in the format YYYY-MM-DDTHH:MM:SSZ", err.Field())
			default:
				fieldError[jsonFieldName] = fmt.Sprintf("%s failed validation with tag '%s'", err.Field(), err.Tag())
			}

			validationErrors = append(validationErrors, fieldError)
		}
	}
	return validationErrors
}

func HitungKomisi(omzet float64) (float64, float64) {
	var komisi models.KomisiPersen

	config.DB.Where("(min_omzet <= ? AND (max_omzet >= ? OR max_omzet = 0))", omzet, omzet).
		Order("min_omzet DESC").First(&komisi)

	komisiNominal := (komisi.Persentase / 100) * omzet
	return komisi.Persentase, komisiNominal
}

func CalculateInstallment(grandTotal float64, InterestRate float64, month uint) float64 {
	installment := (grandTotal + (grandTotal * InterestRate / 100 * float64(month) / 12)) / float64(month)
	return installment
}
