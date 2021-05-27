package example

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"io"
	"log"
	"text/template"
)

type I18nTmpl struct {
	bundle *i18n.Bundle
	*i18n.Localizer
	*template.Template
}

type MessageID string
type Context map[string]interface{}

func (i18nTmpl *I18nTmpl) mustLegalLang(lang string) {
	_, _, err := language.ParseAcceptLanguage(lang)
	if err != nil {
		log.Fatal(err)
	}
}

func (i18nTmpl *I18nTmpl) MustCompile(
	lang string,
	expr string,
	TemplateData map[string]interface{},
) {
	i18nTmpl.mustLegalLang(lang)
	i18nTmpl.Localizer = i18n.NewLocalizer(i18nTmpl.bundle, lang)

	i18nFunc := func(messageID MessageID) string {
		// {{ i18n "whatsInThis" . }}
		return i18nTmpl.Localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID:    string(messageID),
			TemplateData: TemplateData, // other = "What's in this {{ .Type }}"
		})
	}

	i18nTmpl.Template = template.Must(template.New("").
		Funcs(template.FuncMap{"i18n": i18nFunc, "T": i18nFunc}). // ``i18n`` or ``T`` are allow
		Parse(expr))
}

type WriterStore struct {
	Data string
}

func (ws *WriterStore) Write(b []byte) (n int, err error) {
	ws.Data += string(b)
	return len(b), nil
}

func (i18nTmpl *I18nTmpl) MustRender(wr io.Writer, ctx Context) {
	if err := i18nTmpl.Template.Execute(wr, ctx); err != nil {
		log.Fatal(err)
		return
	}
}
