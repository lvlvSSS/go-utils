package validation

import (
	"github.com/go-playground/validator/v10"
	"testing"
)

func TestTagName(t *testing.T) {
	u := user{
		FirstName: "123",
		LastName:  "234",
		Age:       11,
		Email:     "hhs@sss",
	}
	validate := validator.New()
	validate.SetTagName("v")
	validate.RegisterTagNameFunc(TagNameFunc)
	if err := validate.Struct(u); err != nil {
		t.Logf("%v", err)
	}
}

func TestUserStructLevelValidationFunc(t *testing.T) {
	u := user{
		FirstName: "",
		LastName:  "234",
		Age:       11,
		Email:     "hhs@sss.com",
	}
	validate := validator.New()
	validate.SetTagName("v")
	validate.RegisterStructValidation(UserStructLevelValidationFunc, u)
	if err := validate.Struct(u); err != nil {
		t.Logf("%v", err)
		if errs, ok := err.(validator.ValidationErrors); ok {
			t.Logf(errs[0].Tag(), errs[0].Param())
		}
	}
}

func TestIsNickValidationFunc(t *testing.T) {
	u := customUser{
		Name:   "123",
		Gender: -1,
	}
	validate := validator.New()
	validate.SetTagName("v")
	validate.RegisterValidation("is-nick", IsNickValidationFunc)
	validate.RegisterValidation("gender", GenderValidationFunc)
	if err := validate.Struct(u); err != nil {
		t.Logf("%v", err)
	}
}
