package language

import (
	"goapi/language/en"
	"goapi/language/zh"
	"reflect"
)

type Type struct {
	Lang string
}

func Lang(zh string) *Type {
	lang := new(Type)
	lang.Lang = zh
	return lang
}

func (l *Type) GetValidatorCode(name string) string {
	var ValidatorCode interface{}
	switch l.Lang {
	case "zh":
		ValidatorCode = &zh.ValidatorCode{}
	case "en":
		ValidatorCode = &en.ValidatorCode{}
	default:
		ValidatorCode = &en.ValidatorCode{}
	}
	t := reflect.TypeOf(ValidatorCode).Elem()
	field, ok := t.FieldByName(name)
	if !ok {
		return ""
	}
	return field.Tag.Get("msg")
}

func (l *Type) GetErrorCode(name string) (string, string) {
	var ErrorCode interface{}
	switch l.Lang {
	case "zh":
		ErrorCode = &zh.ErrorCode{}
	case "en":
		ErrorCode = &en.ErrorCode{}
	default:
		ErrorCode = &en.ErrorCode{}
	}
	t := reflect.TypeOf(ErrorCode).Elem()
	field, ok := t.FieldByName(name)
	if !ok {
		return "500", "未知错误"
	}
	return field.Tag.Get("code"), field.Tag.Get("msg")
}
