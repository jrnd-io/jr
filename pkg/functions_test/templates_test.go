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

package functions_test

import (
	"bytes"
	"github.com/ugol/jr/pkg/functions"
	"log"
	"strconv"
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

func TestCondition(t *testing.T) {
	tpl := `{{if eq .Name "Ugo"}}Ugo Landini{{end}}`
	data := struct {
		Name string
	}{"Ugo"}
	if err := runtv(tpl, "Ugo Landini", data); err != nil {
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
	a := `{{set_v "id" "10"}}{{template "sub" .}}`
	s := `{{get_v "id"}}`

	aggregate := template.Must(template.New("aggregate").Funcs(functions.FunctionsMap()).Parse(a))
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

func Test2TemplatesWithCommonId(t *testing.T) {
	userTemplate := `{{set_v "id" (uuid)}}"id":"{{get_v "id"}}`
	orderTemplate := `"id":"{{get_v "id"}}`

	v := template.New("aggregate").Funcs(functions.FunctionsMap())

	user, err := v.New("user").Parse(userTemplate)
	if err != nil {
		log.Fatal(err)
	}
	order, err := v.New("order").Parse(orderTemplate)
	if err != nil {
		log.Fatal(err)
	}

	var expectUser bytes.Buffer

	err = user.Execute(&expectUser, nil)
	if err != nil {
		t.Error(err)
	}

	var expectOrder bytes.Buffer

	err = order.Execute(&expectOrder, nil)
	if err != nil {
		t.Error(err)
	}
	if expectUser.String() != expectOrder.String() {
		t.Errorf("Different IDs, should be equal: '%s' '%s'", expectUser.String(), expectOrder.String())
	}

}

func Test2TemplatesWithValueFromList(t *testing.T) {
	userTemplate := `{{$id:=uuid}}{{add_v_to_list "id_list" $id}}{{$id}}`
	userTemplate2 := `{{$id:=uuid}}{{add_v_to_list "id_list" $id}}{{$id}}`
	orderTemplate := `{{random_v_from_list "id_list"}}`

	v := template.New("aggregate").Funcs(functions.FunctionsMap())

	user, err := v.New("user").Parse(userTemplate)
	if err != nil {
		log.Fatal(err)
	}
	user2, err := v.New("user").Parse(userTemplate2)
	if err != nil {
		log.Fatal(err)
	}
	order, err := v.New("order").Parse(orderTemplate)
	if err != nil {
		log.Fatal(err)
	}

	var expectUser bytes.Buffer

	err = user.Execute(&expectUser, nil)
	if err != nil {
		t.Error(err)
	}

	var expectUser2 bytes.Buffer

	err = user2.Execute(&expectUser2, nil)
	if err != nil {
		t.Error(err)
	}

	var expectOrder bytes.Buffer

	err = order.Execute(&expectOrder, nil)
	if err != nil {
		t.Error(err)
	}
	if expectUser.String() != expectOrder.String() && expectUser2.String() != expectOrder.String() {
		t.Errorf("Different IDs, should be equal to one of these: '%s' '%s' but was '%s'", expectUser.String(), expectUser2.String(), expectOrder.String())
	}

}

func TestTemplatesWithValueFromListAtIndex(t *testing.T) {
	templateOne := `{{add_v_to_list "aList" "a"}}`
	templateTwo := `{{add_v_to_list "aList" "b"}}`
	checkTemplate := `{{get_v_from_list_at_index "aList" 1}}{{get_v_from_list_at_index "aList" 0}}`

	v := template.New("aggregate").Funcs(functions.FunctionsMap())

	tOne, err := v.New("user").Parse(templateOne)
	if err != nil {
		log.Fatal(err)
	}

	var expectedOne bytes.Buffer
	tOne.Execute(&expectedOne, nil)

	tTwo, err := v.New("user").Parse(templateTwo)
	if err != nil {
		log.Fatal(err)
	}

	var expectedTwo bytes.Buffer
	tTwo.Execute(&expectedTwo, nil)

	check, err := v.New("order").Parse(checkTemplate)
	if err != nil {
		log.Fatal(err)
	}

	var expectCheck bytes.Buffer

	err = check.Execute(&expectCheck, nil)
	if err != nil {
		t.Error(err)
	}
	if expectCheck.String() != "ba" {
		t.Errorf("Different Result, should be  '%s' but was '%s'", "ba", expectCheck.String())
	}

}

func TestTemplatesWithValueFromListAtIndex_greater_than_length(t *testing.T) {
	templateOne := `{{add_v_to_list "aList" "a"}}`
	templateTwo := `{{add_v_to_list "aList" "b"}}`
	checkTemplate := `{{get_v_from_list_at_index "aList" 10}}{{get_v_from_list_at_index "aList" 10}}`

	v := template.New("aggregate").Funcs(functions.FunctionsMap())

	tOne, err := v.New("user").Parse(templateOne)
	if err != nil {
		log.Fatal(err)
	}

	var expectedOne bytes.Buffer
	tOne.Execute(&expectedOne, nil)

	tTwo, err := v.New("user").Parse(templateTwo)
	if err != nil {
		log.Fatal(err)
	}

	var expectedTwo bytes.Buffer
	tTwo.Execute(&expectedTwo, nil)

	check, err := v.New("order").Parse(checkTemplate)
	if err != nil {
		log.Fatal(err)
	}

	var expectCheck bytes.Buffer

	err = check.Execute(&expectCheck, nil)
	if err != nil {
		t.Error(err)
	}
	if expectCheck.String() != "" {
		t.Errorf("Different Result, should be  '%s' but was '%s'", "", expectCheck.String())
	}

}

func TestManyTemplates(t *testing.T) {

	v := make([]string, 3)
	v[0] = "{{integer 0 1}}"
	v[1] = "{{integer 1 2}}"
	v[2] = "{{integer 2 3}}"

	tpl := template.New("value").Funcs(functions.FunctionsMap())

	for i := 0; i < len(v); i++ {
		_, err := tpl.New(strconv.Itoa(i)).Parse((v[i]))
		if err != nil {
			log.Fatal(err)
		}
	}

	v1 := tpl.Templates()
	if len(v1) != len(v) {
		t.Errorf("Expected %d templates, got %d", len(v), len(v1))
	}

	var b bytes.Buffer

	for i := 0; i < len(v); i++ {
		_ = tpl.ExecuteTemplate(&b, strconv.Itoa(i), nil)
		result, _ := strconv.Atoi(b.String())
		if i != result {
			t.Errorf("Expected %d, got %d", i, result)
		}
		b.Reset()
	}

}

func TestExtractMeta(t *testing.T) {
	tpl := `01234"_meta:"{............................},56789`
	m, v := functions.ExtractMetaFrom(tpl)
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
	m, v := functions.ExtractMetaFrom(tpl)
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
	m, _ := functions.ExtractMetaFrom(tpl)
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
