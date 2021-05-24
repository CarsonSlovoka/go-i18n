package example

import (
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
)

func getTestBundle() *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	return bundle
}

func assertEqual(expected, actual interface{}, errMsg string) error {
	if expected != actual {
		return errors.New(errMsg)
	}
	return nil
}

func MustCheckLegalLang(lang string) {
	// lang: en, es, zh-tw, ...
	_, _, err := language.ParseAcceptLanguage(lang)
	if err != nil {
		log.Fatal(err)
	}
}

type LocalizeTestData struct {
	TemplateData interface{}
	Lang string
	Expected interface{}
}
