package templates

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"html/template"
)

func CreateStringFromTemplate(template *template.Template, data interface{}) string {
	buff := bytes.NewBuffer(nil)
	if err := template.Execute(buff, data); err != nil {
		log.Error(err)
		panic(err)
	}
	return buff.String()
}
