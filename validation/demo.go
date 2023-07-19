package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

type user struct {
	FirstName string
	LastName  string
	Age       int    `v:"gte=0,lte=100"`
	Email     string `v:"email" fld:"e-mail,othervalue"`
}

// 自定义字段类型名称
func TagNameFunc(field reflect.StructField) string {
	tagValue := field.Tag.Get("fld")
	attrs := strings.Split(tagValue, ",")
	if len(attrs) > 1 {
		name := attrs[0]
		if name == "-" {
			return ""
		}
		return name
	}
	return tagValue
}

// 结构体补充校验
func UserStructLevelValidationFunc(sl validator.StructLevel) {
	u := sl.Current().Interface().(user)
	if len(u.FirstName) == 0 {
		sl.ReportError(u.FirstName, "first-name", "FirstName", "first name[%s] is empty", u.FirstName)
	}
}

type customUser struct {
	Name   string `v:"required,is-nick"` //is-nick为自定义tag
	Gender int    `v:"gender"`           // gender为自定义tag
}

// 自定义tag的校验规则
func IsNickValidationFunc(fl validator.FieldLevel) bool {
	return fl.Field().String() == "nick"
}

// 自定义tag的校验规则
func GenderValidationFunc(fl validator.FieldLevel) bool {
	return fl.Field().Int() <= 1 && fl.Field().Int() >= 0
}

type Data struct {
	Name    string
	Email   string
	Details *Details
}

type Details struct {
	FamilyMembers *FamilyMembers
	Salary        string
}

type FamilyMembers struct {
	FatherName string
	MotherName string
}

type Data2 struct {
	Name string
	Age  uint32
}

var validate = validator.New()

func validateStruct() {
	data := Data2{
		Name: "leo",
		Age:  1000,
	}

	rules := map[string]string{
		"Name": "min=4,max=6",
		"Age":  "min=4,max=6",
	}
	// 当引用第三方框架，我们无法修改第三方框架声明的struct的约束时，可以用该方式进行约束校验。
	validate.RegisterStructValidationMapRules(rules, Data2{})

	err := validate.Struct(data)
	fmt.Println(err)
	fmt.Println()
}

func validateStructNested() {
	data := Data{
		Name:  "11sdfddd111",
		Email: "zytel3301@mail.com",
		Details: &Details{
			Salary: "1000",
		},
	}

	rules1 := map[string]string{
		"Name":    "min=4,max=6",
		"Email":   "required,email",
		"Details": "required",
	}

	rules2 := map[string]string{
		"Salary":        "number",
		"FamilyMembers": "required",
	}

	rules3 := map[string]string{
		"FatherName": "required,min=4,max=32",
		"MotherName": "required,min=4,max=32",
	}

	validate.RegisterStructValidationMapRules(rules1, Data{})
	validate.RegisterStructValidationMapRules(rules2, Details{})
	validate.RegisterStructValidationMapRules(rules3, FamilyMembers{})
	err := validate.Struct(data)

	fmt.Println(err)
}
