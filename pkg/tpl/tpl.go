// Copyright © 2024 JR team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package tpl

import (
	"bytes"
	"text/template"

	"github.com/rs/zerolog/log"
)

type Tpl struct {
	Context  any
	Template *template.Template
}

func NewTpl(name string, t string, fmap map[string]interface{}, ctx any) (Tpl, error) {

	tp, err := template.New(name).Funcs(fmap).Parse(t)

	tpl := Tpl{
		Context:  ctx,
		Template: tp,
	}
	return tpl, err
}

func (t *Tpl) Execute() string {
	return t.ExecuteWith(t.Context)
}

func (t *Tpl) ExecuteWith(data any) string {
	var buffer bytes.Buffer
	err := t.Template.Execute(&buffer, data)
	if err != nil {
		log.Fatal().Err(err).Msg("Error executing template")
	}
	return buffer.String()
}
