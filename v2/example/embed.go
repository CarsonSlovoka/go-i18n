package example

import (
	"embed"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

//go:embed i18n
var i18nDir embed.FS

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

	for _, testData := range []*LocalizeTestData{
		{"Bob", "en", "Hello Bob"},
		{"卡森", "zh-tw", "您好 卡森!"},
	} {
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
