package localization

import (
	"html"
	"html/template"

	"github.com/Masterminds/sprig"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Lang struct {
	l          *i18n.Localizer
	botButtons map[Button]string
}

func (l *Lang) BotButton(button Button) string {
	return l.botButtons[button]
}

func (l *Lang) MessageCreate(message Message, in any) (string, error) {
	return l.l.Localize(
		&i18n.LocalizeConfig{
			MessageID: message.String(),
			Funcs: func() template.FuncMap {

				// Add standard functions
				fm := sprig.HtmlFuncMap()

				// Add additional functions
				fm["escapeHTML"] = html.EscapeString
				return fm
			}(),
			TemplateData: in,
		},
	)
}
