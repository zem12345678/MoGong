package translate

import (
	"github.com/BurntSushi/toml"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func GinI18nLocalize() gin.HandlerFunc {
	return ginI18n.Localize(
		ginI18n.WithBundle(&ginI18n.BundleCfg{
			RootPath:         "...",
			AcceptLanguage:   []language.Tag{language.Chinese, language.English},
			DefaultLanguage:  language.English,
			FormatBundleFile: "toml",
			UnmarshalFunc:    toml.Unmarshal,
		}),
		ginI18n.WithGetLngHandle(
			func(c *gin.Context, defaultLang string) string {
				lang := c.Request.Header.Get("Accept-Language")
				if lang == "" {
					return defaultLang
				}
				return lang
			},
		),
	)
}

func gettext(param interface{}) string {

	switch param.(type) {
	case string:
		message, err := ginI18n.GetMessage(param)
		if err != nil {
			return param.(string)
		} else {
			return message
		}
	case map[string]string:
		data := param.(map[string]string)
		message, err := ginI18n.GetMessage(&i18n.LocalizeConfig{
			MessageID: data["id"],
			TemplateData: map[string]string{
				data["key"]: data["value"],
			},
		})
		if err != nil {
			return ""
		} else {
			return message
		}
	default:
		message, err := ginI18n.GetMessage(param)
		if err != nil {
			return ""
		} else {
			return message
		}
	}

}
