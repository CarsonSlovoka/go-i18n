package example

import (
	"embed"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"log"
	"path"
)

//go:embed i18n/active.en.toml i18n/active.zh-tw.toml
var i18nES embed.FS

func demoSpecifyFile() error {
	CHFilePath := "i18n/active.zh-tw.toml"
	bytesCH, err := i18nES.ReadFile(CHFilePath)
	if err != nil {
		panic(err)
	}

	bundle := getTestBundle()
	bundle.MustParseMessageFileBytes(bytesCH, CHFilePath)
	// No need to load active.en.toml since we are providing default translations.
	// bundle.MustParseMessageFileBytes(bytesCH, ENFilePath)

	for _, testData := range []*localizeTestData{
		{"Bob", "en", "Hello Bob", nil},
		{"卡森", "zh-tw", "您好 卡森!", nil},
	} {
		mustCheckLegalLang(testData.Lang) // check only
		localizer := i18n.NewLocalizer(bundle, testData.Lang)
		name := testData.TemplateData.(string)
		resultStr := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "HelloPerson",
				Other: "Hello {{.Name}}",
			},
			TemplateData: map[string]string{
				"Name": name,
			},
		})
		errMsg := fmt.Sprintf("%s != %s", testData.Expected, resultStr)
		if err := assertEqual(testData.Expected, resultStr, errMsg); err != nil {
			return err
		}
	}
	return nil
}

//go:embed i18n
var i18nDirFS embed.FS

func demoDir() {
	bundle := getTestBundle()
	myLangDirPath := "i18n"
	dirEntrySlice, err := i18nDirFS.ReadDir(myLangDirPath)
	if err != nil {
		log.Fatal(err)
	}

	// load all language file, similar as ``bundle.MustLoadMessageFile("i18n/xxx.toml")``
	for _, dirEntry := range dirEntrySlice {
		if dirEntry.IsDir() {
			continue
		}
		langFilePath := path.Join(myLangDirPath, dirEntry.Name())
		bytesLang, err := i18nDirFS.ReadFile(langFilePath)
		if err != nil {
			log.Fatal(err)
		}
		bundle.MustParseMessageFileBytes(bytesLang, langFilePath)
	}

	// Start Test
	type TmplContext map[string]interface{}
	simpleTest(bundle, "HelloPerson", []*localizeTestData{
		{TemplateData: TmplContext{"Name": "Bob"}, Lang: "en", Expected: "Hello Bob", PluralCount: nil},
		{TmplContext{"Name": "Bar"}, "es", "Hola Bar", nil},
		{TmplContext{"Name": "卡森"}, "zh-tw", "您好 卡森!", nil},
	})
}
