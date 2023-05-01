package helper

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/validator/v10"
	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

func InitTranslator() ut.Translator {
	var translator ut.Translator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		en := en.New()
		uni := ut.New(en, en)
		translator, _ = uni.GetTranslator("en")
		en_translations.RegisterDefaultTranslations(v, translator)
	}
	return translator
}