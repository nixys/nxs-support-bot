package misc

import (
	"bytes"
	ttemplate "text/template"

	"github.com/Masterminds/sprig"
)

// TemplateExec makes message from given template `tpl` and data `d`
func TemplateExec(tpl string, d interface{}) (string, error) {

	var b bytes.Buffer

	// See http://masterminds.github.io/sprig/ for details
	t, err := ttemplate.New("template").Funcs(sprig.TxtFuncMap()).Parse(tpl)
	if err != nil {
		return "", err
	}

	err = t.Execute(&b, d)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}
