package localization

import (
	"fmt"
	"io/ioutil"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type Bundle struct {
	b  *i18n.Bundle
	dl Lang
}

func Init(path string) (Bundle, error) {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return Bundle{}, fmt.Errorf("localization init: %w", err)
	}

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	for _, f := range files {
		if _, err := bundle.LoadMessageFile(path + "/" + f.Name()); err != nil {
			return Bundle{}, fmt.Errorf("localization init: %w", err)
		}
	}

	// Init default language
	dl, err := langSwitch(bundle, "")
	if err != nil {
		return Bundle{}, fmt.Errorf("localization init: %w", err)
	}

	return Bundle{
		b:  bundle,
		dl: dl,
	}, nil
}

func (b *Bundle) LangSwitch(tag string) (Lang, error) {
	return langSwitch(b.b, tag)
}

func (b *Bundle) DefaultLanguage() Lang {
	return b.dl
}

func langSwitch(b *i18n.Bundle, tag string) (Lang, error) {

	l := i18n.NewLocalizer(b, tag)

	bb := make(map[Button]string)

	for _, btn := range buttons {
		button, err := l.Localize(&i18n.LocalizeConfig{
			MessageID: btn.String(),
		})
		if err != nil {
			return Lang{}, fmt.Errorf("lang switch: %w", err)
		}

		bb[btn] = button
	}

	return Lang{
		l:          l,
		botButtons: bb,
	}, nil
}
