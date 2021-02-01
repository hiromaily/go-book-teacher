package tmpl

import (
	"bytes"
	"text/template"
)

func StrTempParser(temp string, params interface{}) (string, error) {
	var writer bytes.Buffer

	tpl := template.Must(template.New("tpl").Parse(temp))

	if err := tpl.Execute(&writer, params); err != nil {
		return "", err
	}

	return writer.String(), nil
}
