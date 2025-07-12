package validator

type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) SetFieldError(name, reason string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, ok := v.FieldErrors[name]; !ok {
		v.FieldErrors[name] = reason
	}
}

func (v *Validator) Check(ok bool, name, reason string) {
	if !ok {
		v.SetFieldError(name, reason)
	}
}

func (v *Validator) IsValid() bool {
	return len(v.FieldErrors) == 0
}
