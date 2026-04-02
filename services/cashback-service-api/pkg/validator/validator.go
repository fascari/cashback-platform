package validator

import (
	"errors"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/locales/en"
	unitrans "github.com/go-playground/universal-translator"
	govalidator "github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
)

var (
	once sync.Once
	v    *govalidator.Validate
	uni  *unitrans.UniversalTranslator
)

func init() {
	once.Do(func() {
		eng := en.New()
		uni = unitrans.New(eng, eng)
		v = govalidator.New()

		trans, _ := uni.GetTranslator("en")
		_ = entranslations.RegisterDefaultTranslations(v, trans)

		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	})
}

func Validate(s any) error {
	err := v.Struct(s)
	if err == nil {
		return nil
	}

	var validationErrors govalidator.ValidationErrors
	ok := errors.As(err, &validationErrors)
	if !ok {
		return err
	}

	trans, _ := uni.GetTranslator("en")
	errs := Errors{}

	for _, e := range validationErrors {
		errs[e.Field()] = e.Translate(trans)
	}

	return errs
}
