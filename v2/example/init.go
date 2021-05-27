package example

import (
	"errors"
	"fmt"
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

type localizeTestData struct {
	TemplateData interface{}
	Lang string
	Expected interface{}
	PluralCount interface{}
}

func mustAssertEqual(expected, actual interface{}, errMsg string) {
	if err := assertEqual(expected, actual, errMsg); err != nil {
		log.Fatal(err)
	}
}

func mustCheckLegalLang(lang string) {
	// lang: en, es, zh-tw, ...
	_, _, err := language.ParseAcceptLanguage(lang)
	if err != nil {
		log.Fatal(err)
	}
}

func simpleTest(bundle *i18n.Bundle, messageID string, testDataSlice []*localizeTestData) {
	for _, testData := range testDataSlice {

		mustCheckLegalLang(testData.Lang) // check only
		localizer := i18n.NewLocalizer(bundle, testData.Lang)
		resultStr := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: messageID,
			TemplateData: testData.TemplateData,
			PluralCount: testData.PluralCount,
		})
		errMsg := fmt.Sprintf("%s != %s", testData.Expected, resultStr)
		mustAssertEqual(testData.Expected, resultStr, errMsg)
	}
}
