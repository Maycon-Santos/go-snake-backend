package validator

type Validator interface {
	Validate(value interface{}) (errType, error)
	Required() Validator
	MinLen(len uint) Validator
	MaxLen(len uint) Validator
	NoContains(noContains []string) Validator
}

func Field(field string) Validator {
	return validatorConfig{field: field}
}

func (v validatorConfig) Validate(value interface{}) (errType, error) {
	for _, fn := range validators {
		if errType, err := fn(v, v.field, value); err != nil {
			return errType, err
		}
	}

	return "", nil
}

func (v validatorConfig) Required() Validator {
	required := true
	v.required = &required
	return v
}

func (v validatorConfig) MinLen(len uint) Validator {
	v.minLen = &len
	return v
}

func (v validatorConfig) MaxLen(len uint) Validator {
	v.maxLen = &len
	return v
}

func (v validatorConfig) NoContains(noContains []string) Validator {
	v.noContains = noContains
	return v
}
