package language

import (
	"goapi/language/en"
	"goapi/language/es"
	"goapi/language/ja"
	"goapi/language/ko"
	"goapi/language/ti"
	"goapi/language/ve"
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
	case "en":
		ValidatorCode = &en.ValidatorCode{}
	case "es":
		ValidatorCode = &es.ValidatorCode{}
	case "ja":
		ValidatorCode = &ja.ValidatorCode{}
	case "ko":
		ValidatorCode = &ko.ValidatorCode{}
	case "ti":
		ValidatorCode = &ti.ValidatorCode{}
	case "ve":
		ValidatorCode = &ve.ValidatorCode{}
	case "zh":
		ValidatorCode = &zh.ValidatorCode{}
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
	case "en":
		ErrorCode = &en.ErrorCode{}
	case "es":
		ErrorCode = &es.ErrorCode{}
	case "ja":
		ErrorCode = &ja.ErrorCode{}
	case "ko":
		ErrorCode = &ko.ErrorCode{}
	case "ti":
		ErrorCode = &ti.ErrorCode{}
	case "ve":
		ErrorCode = &ve.ErrorCode{}
	case "zh":
		ErrorCode = &zh.ErrorCode{}
	default:
		ErrorCode = &en.ErrorCode{}
	}
	t := reflect.TypeOf(ErrorCode).Elem()
	field, ok := t.FieldByName(name)
	if !ok {
		return "500", "unknown mistake"
	}
	return field.Tag.Get("code"), field.Tag.Get("msg")
}
