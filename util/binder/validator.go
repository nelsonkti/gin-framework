package binder

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh_Hans_CN"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
	"sync"
)

type Validator struct {
	once     sync.Once
	validate *validator.Validate
	trans    ut.Translator
}

type SliceValidationError []error

// Error concatenates all error elements in SliceValidationError into a single string separated by \n.
func (err SliceValidationError) Error() string {
	n := len(err)
	switch n {
	case 0:
		return ""
	default:
		var b strings.Builder
		if err[0] != nil {
			fmt.Fprintf(&b, "[%d]: %s", 0, err[0].Error())
		}
		if n > 1 {
			for i := 1; i < n; i++ {
				if err[i] != nil {
					b.WriteString("\n")
					fmt.Fprintf(&b, "[%d]: %s", i, err[i].Error())
				}
			}
		}
		return b.String()
	}
}

var _ binding.StructValidator = (*Validator)(nil)

// ValidateStruct receives any kind of type, but only performed struct or pointer to struct type.
func (v *Validator) ValidateStruct(obj any) error {
	if obj == nil {
		return nil
	}

	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr:
		return v.ValidateStruct(value.Elem().Interface())
	case reflect.Struct:
		return v.validateStruct(obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(SliceValidationError, 0)
		for i := 0; i < count; i++ {
			if err := v.ValidateStruct(value.Index(i).Interface()); err != nil {
				validateRet = append(validateRet, err)
			}
		}
		if len(validateRet) == 0 {
			return nil
		}
		return validateRet
	default:
		return nil
	}
}

// validateStruct receives struct type
func (v *Validator) validateStruct(obj any) error {
	v.lazyinit()
	err := v.validate.Struct(obj)
	if err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			// 将验证错误翻译成中文
			return errors.New(v.translateErrors(errs))
		}
	}
	return err
}

// translateErrors 将验证错误翻译成中文并返回字符串
func (v *Validator) translateErrors(errs validator.ValidationErrors) string {
	for _, err := range errs {
		return err.Translate(v.trans)
	}
	return ""
}

// Engine returns the underlying validator engine which powers the default
// Validator instance. This is useful if you want to register custom validations
// or struct level validations. See validator GoDoc for more info -
// https://pkg.go.dev/github.com/go-playground/validator/v10
func (v *Validator) Engine() any {
	v.lazyinit()
	return v.validate
}

func (v *Validator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")

		// 初始化中文翻译器
		zh := zh_Hans_CN.New()
		uni := ut.New(zh)
		v.trans, _ = uni.GetTranslator("zh_Hans_CN")

		// 注册中文翻译器
		zhTrans.RegisterDefaultTranslations(v.validate, v.trans)

		// 将验证字段名映射为中文名
		v.validate.RegisterTagNameFunc(func(field reflect.StructField) string {
			label := field.Tag.Get("label")
			if label == "" {
				label = field.Name
			}
			return label
		})
	})
}
