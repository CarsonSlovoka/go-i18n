// Command example runs a sample webserver that uses go-i18n/v2/i18n.
package example

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	textTmp "text/template"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func DemoBasic() {

	var page = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<body>

<h1>{{.Title}}</h1>

{{range .Paragraphs}}<p>{{.}}</p>{{end}}

</body>
</html>
`))

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	// No need to load active.en.toml since we are providing default translations.
	bundle.MustLoadMessageFile("i18n/active.en.toml") // If you ignore this, you must provide ``DefaultMessage``, to let it know the content of `i18n.NewBundle(language.English)` is.
	bundle.MustLoadMessageFile("i18n/active.es.toml")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lang := r.FormValue("lang")
		accept := r.Header.Get("Accept-Language")
		localizer := i18n.NewLocalizer(bundle, lang, accept)

		name := r.FormValue("name")
		if name == "" {
			name = "Bob"
		}

		unreadEmailCount, _ := strconv.ParseInt(r.FormValue("unreadEmailCount"), 10, 64)

		helloPerson := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "HelloPerson",
				Other: "Hello {{.Name}}",
			},
			TemplateData: map[string]string{
				"Name": name,
			},
		})

		myUnreadEmails := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "MyUnreadEmails",
				Description: "The number of unread emails I have",
				One:         "I have {{.PluralCount}} unread email.",
				Other:       "I have {{.PluralCount}} unread emails.",
			},
			PluralCount: unreadEmailCount,
		})

		personUnreadEmails := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "PersonUnreadEmails",
				Description: "The number of unread emails a person has",
				One:         "{{.Name}} has {{.UnreadEmailCount}} unread email.",
				Other:       "{{.Name}} has {{.UnreadEmailCount}} unread emails.",
			},
			PluralCount: unreadEmailCount,
			TemplateData: map[string]interface{}{
				"Name":             name,
				"UnreadEmailCount": unreadEmailCount,
			},
		})

		demoDelim := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:         "HelloPerson",
				LeftDelim:  "<<",
				RightDelim: ">>",
				Other:      "Hello <<.Name>>",
			},
			TemplateData: map[string]string{
				"Name": name,
			},
		})

		demoFuncs := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{ // set the default for your ``i18n.NewBundle(language.English)``
				ID: "UniversalTest", // since, ``active.es.toml`` not belong ``i18n.NewBundle(language.English)`` so this ID must exists on its contents. Otherwise, when you browse "es" will fail.
				Other: `{{largest .Numbers}}
{{sayHi}}
`,
			},
			TemplateData: map[string]interface{}{
				"Numbers": []float64{3, 3.2, 6, 1.2},
			},
			Funcs: textTmp.FuncMap{
				"largest": func(slice []float64) float64 {
					if len(slice) == 0 {
						return 0
					}
					max := slice[0]
					for _, val := range slice[1:] {
						if val > max {
							max = val
						}
					}
					return max
				},
				"sayHi": func() string {
					return "Hello World"
				},
			},
		})

		demoMessageID := localizer.MustLocalize(&i18n.LocalizeConfig{
			// DefaultMessage: // we are not set the ``DefaultMessage``, so we should cancel comment about bundle.MustLoadMessageFile("active.en.toml"), to let it know the DefaultMessage is.
			MessageID: "HelloPerson",
			TemplateData: map[string]interface{}{
				"Name": name,
			},
		})

		err := page.Execute(w, map[string]interface{}{
			"Title": helloPerson,
			"Paragraphs": []string{
				demoDelim,
				myUnreadEmails,
				personUnreadEmails,
				demoFuncs,
				demoMessageID,
			},
		})
		if err != nil {
			panic(err)
		}
	})

	fmt.Println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func demoUseFunc2Render() {
	/*
		This example is pretty like to one below
		> https://gohugo.io/functions/i18n/
	*/
	bundle := getTestBundle()
	i18nTmpl := &I18nTmpl{bundle: bundle}

	ConfigData := map[string]interface{}{ // It simulates you load the data from your config files.
		"User": "Carson",
		"Type": "Interface", // other = "What's in this {{ .Type }}"
	}

	for _, curLang := range []string{"en", "es", "zh-tw"} {
		bundle.MustLoadMessageFile(fmt.Sprintf("i18n/data2/%s.toml", curLang))

		expr := `<h1>{{.Title}}</h1>
{{range .Paragraphs}}<p>{{ i18n .}}</p>{{end}}
{{ i18n "whatsInThis" }}
{{ T "whatsInThis" }}
`
		i18nTmpl.MustCompile(curLang, expr, ConfigData)
		writerStore := &WriterStore{}
		multipleWriter := io.MultiWriter(os.Stdout, writerStore)
		i18nTmpl.MustRender(multipleWriter, Context{
			"Title": "DemoRenderTmpl",
			"Paragraphs": []MessageID{
				"More",
				"readMore",
			}},
		)
		fmt.Println(writerStore.Data)
	}
}

func demoZeroOneTwoFewManyOther() {
	bundle := i18n.NewBundle(language.Arabic)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	localizer := i18n.NewLocalizer(bundle, "ar")

	for _, curCount := range []interface{}{0, 1, 2,
		3, 10, // Few range
		11, 99, // many
		nil, // other
	} {
		localize := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "PersonCats",
				Zero:  "zero 0",
				One:   "1 一",
				Two:   "2 二",
				Few:   "3-10",
				Many:  "11-99",
				Other: "{{.Name}} has {{.Count}} cats.",
			},
			TemplateData: map[string]interface{}{
				"Name":  "Nick",
				"Count": 2,
			},
			PluralCount: curCount,
		})
		fmt.Println(localize)
	}
}

func demoJSON() {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile("i18n/data3-json/en.json") // If you ignore this, you must provide ``DefaultMessage``, to let it know the content of `i18n.NewBundle(language.English)` is.
	bundle.MustLoadMessageFile("i18n/data3-json/zh-tw.json")

	type TmplContext map[string]interface{}
	simpleTest(bundle, "IDS_USER", []*localizeTestData{
		{TemplateData: TmplContext{"Name": "Carson"}, Lang: "en", Expected: "User: Carson"},
		{TmplContext{"Name": "Carson"}, "zh-tw", "使用者: Carson", nil},
	})

	simpleTest(bundle, "IDS_COUNT", []*localizeTestData{
		{TemplateData: nil, Lang: "en", Expected: "1 item", PluralCount: 1},
		{TmplContext{"Count": "unknown"}, "en", "unknown item", nil},
	})
}
