package validator

import (
	"fmt"
	"github.com/go-playground/locales"
	locEn "github.com/go-playground/locales/en"
	locZh "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	cmap "github.com/orcaman/concurrent-map"
	"goapi/language"
	"reflect"
)

type Type struct {
	Lang string
}

var (
	err        error
	uni        *ut.UniversalTranslator
	importLang locales.Translator
	trans      ut.Translator
	Validate   *validator.Validate
)

func init() {
	Validate = validator.New()
	err = Validate.RegisterValidation("xwMinReg", func(f validator.FieldLevel) bool {
		value := f.Field().String()
		if len(value) < 4 {
			return false
		} else {
			return true
		}
	}) //注册自定义的校验函数 minReg和validate tag值保持一致
	if err != nil {
		fmt.Println("注册自定义验证规则 xwMinReg 失败")
	}
}

func Lang(locale string) *Type {
	lang := new(Type)
	lang.Lang = locale
	switch lang.Lang {
	case "zh":
		importLang = locZh.New()
		uni = ut.New(importLang, importLang)
		trans, _ = uni.GetTranslator(lang.Lang)
		err = zhTranslations.RegisterDefaultTranslations(Validate, trans)
		break
	case "en":
		importLang = locEn.New()
		uni = ut.New(importLang, importLang)
		trans, _ = uni.GetTranslator(lang.Lang)
		err = enTranslations.RegisterDefaultTranslations(Validate, trans)
		break
	default: // 默认中文
		importLang = locZh.New()
		uni = ut.New(importLang, importLang)
		trans, _ = uni.GetTranslator(lang.Lang)
		err = zhTranslations.RegisterDefaultTranslations(Validate, trans)
		break
	}
	if err != nil {
		fmt.Println("validator 翻译出错")
	}
	return lang
}

// Translate 翻译工具
func (l *Type) Translate(err error, s interface{}, lang string) string {
	result := cmap.New()
	t := reflect.TypeOf(s)
	for _, errs := range err.(validator.ValidationErrors) {
		// 使用反射方法获取struct种的json标签作为key --重点2
		var k string
		if field, ok := t.FieldByName(errs.StructField()); ok {
			k = field.Tag.Get("json")
		}
		if k == "" {
			k = errs.StructField()
		}
		diyTag := errs.Tag()
		// 检测自定义 标签语言包
		msg := language.Lang(lang).GetValidatorCode(diyTag)
		if len(msg) > 0 {
			result.Set(k, msg)
		} else {
			result.Set(k, errs.Translate(trans))
		}
	}
	return getFirstMessage(result.Items())
}

func getFirstMessage(messages map[string]interface{}) string {
	for _, val := range messages {
		return val.(string)
	}
	return ""
}
