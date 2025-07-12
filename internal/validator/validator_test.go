package validator

import "testing"

func TestValidatorLifecycle(t *testing.T) {
	var validator Validator
	if !validator.IsValid() {
		t.Fatal("expected new validator to be valid")
	}

	validator.Check(true, "foo", "bar")
	if !validator.IsValid() {
		t.Fatal("expected Check(true, ..., ...) to not add an error")
	}

	validator.Check(false, "foo", "bar")
	if validator.IsValid() {
		t.Fatal("expected Check(false, ..., ...) to add an error")
	}
}
