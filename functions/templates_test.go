//Copyright © 2022 Ugo Landini <ugo.landini@gmail.com>
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in
//all copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
//THE SOFTWARE.

package functions

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
		{"Roma", "00100"},
		{"Milano", "20100"},
		{"Napoli", "80100"},
	}

	tpl := `{{range $i, $e := .}}{{$i}},{{$e.ZIP}} {{end}}`
	if err := runtv(tpl, "0,00100 1,20100 2,80100 ", cities); err != nil {
		t.Error(err)
	}
}

func TestNested(t *testing.T) {
	tpl := template.Must(template.New("inside").Parse(`{{.Name}}`))
	_, _ = tpl.New("outside").Parse(`{{template "inside"}}`)

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
	_, err := tpl.New("header").Parse(`<HEADER/>`)
	if err != nil {
		return
	}
	_, err = tpl.New("footer").Parse(`<FOOTER/>`)
	if err != nil {
		return
	}

	data := struct {
		Name string
	}{"Ugo"}
	var b bytes.Buffer

	err = tpl.Execute(&b, data)
	if err != nil {
		t.Error(err)
	}

	expect := "<HEADER/> Ugo <FOOTER/>"
	if expect != b.String() {
		t.Errorf("Expected '%s', got '%s'", expect, b.String())
	}

}

func TestNestedPassingContext(t *testing.T) {
	a := `{{template "sub" .}}`
	s := `{{.C}}`

	aggregate := template.Must(template.New("aggregate").Parse(a))
	_, err := aggregate.New("sub").Parse(s)

	if err != nil {
		t.Error(err)
	}

	var b bytes.Buffer

	//*
	data := struct {
		C string
	}{"10"}
	//*/

	err = aggregate.Execute(&b, data)
	if err != nil {
		t.Error(err)
	}

	expect := "10"
	if expect != b.String() {
		t.Errorf("Expected '%s', got '%s'", expect, b.String())
	}
}

func TestRelationship(t *testing.T) {
	a := `{{set "id" "10"}}{{template "sub" .}}`
	s := `{{get "id"}}`

	aggregate := template.Must(template.New("aggregate").Funcs(FunctionsMap()).Parse(a))
	sub, err := aggregate.New("sub").Parse(s)

	if err != nil {
		t.Error(err)
	}

	var b bytes.Buffer

	err = aggregate.Execute(&b, nil)
	if err != nil {
		t.Error(err)
	}

	expect := "10"
	if expect != b.String() {
		t.Errorf("Expected '%s', got '%s'", expect, b.String())
	}

	b.Reset()

	err = sub.Execute(&b, nil)
	if err != nil {
		t.Error(err)
	}
	if expect != b.String() {
		t.Errorf("Expected '%s', got '%s'", expect, b.String())
	}

	templates := aggregate.Templates()

	if len(templates) != 2 {
		t.Errorf("Expected 2 templates, got %d", len(templates))
	}
}
func TestExtractMeta(t *testing.T) {
	tpl := `01234"_meta:"{............................},56789`
	m, v := ExtractMetaFrom(tpl)
	expect := `0123456789`
	if expect != v {
		t.Errorf("Expected '%s', got '%s'", expect, v)
	}
	expect = `{............................}`
	if expect != m {
		t.Errorf("Expected '%s', got '%s'", expect, m)
	}
}

func TestExtractEmptyMeta(t *testing.T) {
	tpl := `0123456789`
	m, v := ExtractMetaFrom(tpl)
	expect := `0123456789`
	if expect != v {
		t.Errorf("Expected '%s', got '%s'", expect, v)
	}
	expect = ``
	if expect != m {
		t.Errorf("Expected '%s', got '%s'", expect, m)
	}
}

func TestExtractUserMeta(t *testing.T) {
	tpl := `{
  "_meta":{
                      "topic": "users",
                      "key": "id",
                      "relationships": [
                          {
                              "topic": "purchases",
                              "parent_field": "id",
                              "child_field": "user_id",
                              "records_per": 4
                          }
                      ]
                  },
  "guid": "{{uuid}}",
  "isActive": {{bool}},
  "balance": "{{amount 100 10000 "€"}}",
  "picture": "http://placehold.it/32x32",
  "age": {{integer 20 60}},
  "eyeColor": "{{randoms "blue|brown|green"}}",
  "name": "{{name}} {{surname}}",
  "gender": "{{gender}}",
  "company": "{{company}}",
  "email": "{{email}}",
  "about": "{{lorem 20}}",
  "address": "{{city}}, {{street}} {{building 2}}, {{zip}}",
  "phone_number": "{{land_prefix}} {{regex "[0-9]{7}"}}",
  "latitude": {{latitude}},
  "longitude": {{longitude}}
}`
	m, _ := ExtractMetaFrom(tpl)
	expect := `{
                      "topic": "users",
                      "key": "id",
                      "relationships": [
                          {
                              "topic": "purchases",
                              "parent_field": "id",
                              "child_field": "user_id",
                              "records_per": 4
                          }
                      ]
                  }`
	if expect != m {
		t.Errorf("Expected '%s', got '%s'", expect, m)
	}
}
