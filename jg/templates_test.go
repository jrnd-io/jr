package jg

import (
	"bytes"
	"testing"
	"text/template"
)

func TestSimpleContext(t *testing.T) {
	tpl := `Hello, {{.Name}}!`
	data := struct {
		Name string
	}{"Ugo Landini"}

	if err := runtv(tpl, "Hello, Ugo Landini!", data); err != nil {
		t.Error(err)
	}
}

func TestRange(t *testing.T) {

	type City struct {
		Name string
		ZIP  string
	}

	cities := []City{
		City{"Roma", "00100"},
		City{"Milano", "20100"},
		City{"Napoli", "80100"},
	}

	tpl := `{{range $i, $e := .}}{{$i}},{{$e.ZIP}} {{end}}`
	if err := runtv(tpl, "0,00100 1,20100 2,80100 ", cities); err != nil {
		t.Error(err)
	}
}

func TestNested(t *testing.T) {
	tpl := template.Must(template.New("inside").Parse(`{{.Name}}`))
	tpl.New("outside").Parse(`{{template "inside"}}`)

	data := struct {
		Name string
	}{"Ugo"}
	var b bytes.Buffer

	err := tpl.Execute(&b, data)
	if err != nil {
		t.Error(err)
	}

	expect := "Ugo"
	if expect != b.String() {
		t.Errorf("Expected '%s', got '%s'", expect, b.String())
	}

}

func TestComplexNested(t *testing.T) {
	tpl := template.Must(template.New("test").Parse(`{{template "header"}} {{.Name}} {{template "footer"}}`))
	tpl.New("header").Parse(`<HEADER/>`)
	tpl.New("footer").Parse(`<FOOTER/>`)

	data := struct {
		Name string
	}{"Ugo"}
	var b bytes.Buffer

	err := tpl.Execute(&b, data)
	if err != nil {
		t.Error(err)
	}

	expect := "<HEADER/> Ugo <FOOTER/>"
	if expect != b.String() {
		t.Errorf("Expected '%s', got '%s'", expect, b.String())
	}

}
