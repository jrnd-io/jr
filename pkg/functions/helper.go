package functions

import (
	"bytes"
	"log"
	"regexp"
	"text/template"
)

func ExecuteTemplate(key *template.Template, value *template.Template, oneline bool) (string, string, error) {

	var kBuffer, vBuffer bytes.Buffer
	var err error

	if err = key.Execute(&kBuffer, JrContext); err != nil {
		log.Println(err)
	}
	k := kBuffer.String()

	if err = value.Execute(&vBuffer, JrContext); err != nil {
		log.Println(err)
	}
	v := vBuffer.String()

	if oneline {
		re := regexp.MustCompile(`\r?\n?`)
		v = re.ReplaceAllString(v, "")
	}

	JrContext.GeneratedObjects++
	JrContext.GeneratedBytes += int64(len(v))

	return k, v, err
}
