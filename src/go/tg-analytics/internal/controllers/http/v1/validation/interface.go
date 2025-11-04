package validation

type ValidationInterface interface {
	Validate()
	Normalize()
	ValidateAndNormalize()
}
