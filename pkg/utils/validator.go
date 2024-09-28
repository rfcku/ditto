package utils

type Validator interface {
	Validate() error
}

func ValidateObject(v Validator) error {
	return v.Validate()
}
