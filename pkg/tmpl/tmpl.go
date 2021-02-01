package tmpl

import (
	"bytes"
	"text/template"
)

func StrTempParser(temp string, params interface{}) (string, error) {
	var parseResult bytes.Buffer

	tpl := template.Must(template.New("tpl").Parse(temp))

	if err := tpl.Execute(&parseResult, params); err != nil {
		return "", err
	}

	return parseResult.String(), nil
}
