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
	"fmt"
	"testing"
	"text/template"
)

func TestSubstr(t *testing.T) {
	tpl := `{{"fooo" | substr 0 3 }}`
	if err := runt(tpl, "foo"); err != nil {
		t.Error(err)
	}
}

func TestSplit(t *testing.T) {
	tpl := `{{split "a|b" "|"}}`
	if err := runt(tpl, "[a b]"); err != nil {
		t.Error(err)
	}
}

func TestTitle(t *testing.T) {
	tpl := `{{"foo" | title}}`
	if err := runt(tpl, "Foo"); err != nil {
		t.Error(err)
	}
}

func TestMax(t *testing.T) {
	tpl := `{{max 1 4}}`
	if err := runt(tpl, "4"); err != nil {
		t.Error(err)
	}
}

func TestMin(t *testing.T) {
	tpl := `{{min 1 4}}`
	if err := runt(tpl, "1"); err != nil {
		t.Error(err)
	}
}

func TestUSState(t *testing.T) {

	hawaii := "{{capital_at 10}} {{state_at 10}} {{state_short_at 10}} {{zip_at 10}}"
	massachussets := `{{capital_at 20}} {{state_at 20}} {{state_short_at 20}} {{zip_at 20}}`
	newyork := `{{capital_at 31}} {{state_at 31}} {{state_short_at 31}} {{zip_at 31}}`
	texas := `{{capital_at 42}} {{state_at 42}} {{state_short_at 42}} {{zip_at 42}}`
	virginia := `{{capital_at 45}} {{state_at 45}} {{state_short_at 45}} {{zip_at 45}}`
	wyoming := `{{capital_at 49}} {{state_at 49}} {{state_short_at 49}} {{zip_at 49}}`
	if err := runt(hawaii, "Honolulu Hawaii HI 96813"); err != nil {
		t.Error(err)
	}
	if err := runt(massachussets, "Boston Massachusetts MA 02201"); err != nil {
		t.Error(err)
	}
	if err := runt(newyork, "Albany New York NY 12207"); err != nil {
		t.Error(err)
	}
	if err := runt(texas, "Austin Texas TX 78701"); err != nil {
		t.Error(err)
	}
	if err := runt(virginia, "Richmond Virginia VA 23219"); err != nil {
		t.Error(err)
	}
	if err := runt(wyoming, "Cheyenne Wyoming WY 82001"); err != nil {
		t.Error(err)
	}
}

func TestCache(t *testing.T) {

	v, f := cache("wine")

	if v != true || f != nil {
		t.Error("cache should be empty, no errors")
	}
	v, f = cache("wine")
	if v != false || f != nil {
		t.Error("cache should be full, no errors")
	}
	_, f = cache("wines")
	if f == nil {
		t.Error("no cacheable, should get error")
	}
}

func TestFrom(t *testing.T) {

	tpl := `{{from "actor"}}`
	if err := runt(tpl, "Angelina Jolie"); err != nil {
		t.Error(err)
	}
	tpl = `{{from "actors"}}`
	if err := runt(tpl, ""); err != nil {
		t.Error(err)
	}
}

func TestPassword(t *testing.T) {

	tpl := `{{password 5 true "PwD" "!?!"}}`
	if err := runt(tpl, "PwDalYza!?!"); err != nil {
		t.Error(err)
	}
}

func TestIPv6(t *testing.T) {

	tpl := `{{ipv6}}`
	if err := runt(tpl, "face:bf42:e25e:8b14:eafc:81ea:e0d0:f2c"); err != nil {
		t.Error(err)
	}
}

func TestIP(t *testing.T) {

	tpl := `{{ip "10.2.0.0/16"}}`
	if err := runt(tpl, "10.2.238.223"); err != nil {
		t.Error(err)
	}
}

func TestUseragent(t *testing.T) {

	tpl := `{{useragent}}`
	if err := runt(tpl, "Mozilla/5.0 (Windows NT 10.0) AppleWebKit/549.15 (KHTML, like Gecko) Edge/8.3.0.0 Mobile Safari/8.5"); err != nil {
		t.Error(err)
	}
}

func TestCounter(t *testing.T) {

	tpl := `{{counter "A" 0 1}},{{counter "B" 2 2}},{{counter "C" -4 1}},{{counter "D" 0 -1}}`

	if err := runt(tpl, "0,2,-4,0"); err != nil {
		t.Error(err)
	}

	if err := runt(tpl, "1,4,-3,-1"); err != nil {
		t.Error(err)
	}

	if err := runt(tpl, "2,6,-2,-2"); err != nil {
		t.Error(err)
	}
}

func TestArray(t *testing.T) {

	tpl := `{{array 5}}`
	if err := runt(tpl, "[0 0 0 0 0]"); err != nil {
		t.Error(err)
	}

	tpl2 := `{{array 1}}`
	if err := runt(tpl2, "[0]"); err != nil {
		t.Error(err)
	}

	tpl3 := `{{array 0}}`
	if err := runt(tpl3, "[]"); err != nil {
		t.Error(err)
	}
}

func runt(tpl, expect string) error {
	return runtv(tpl, expect, "")
}

func TestRegex(t *testing.T) {

	tpl := `{{regex "Z{2,5}"}}`
	if err := runt(tpl, "ZZ"); err != nil {
		t.Error(err)
	}
	//123[0-2]+.*\w{3}
	//tpl = "{{Regex `123[0-2]+.*\w{3}`}}"
	//if err := runt(tpl, "ZZZ"); err != nil {
	//	t.Error(err)
	//}
}
func runtv(tpl, expect string, vars interface{}) error {

	t := template.Must(template.New("test").Funcs(FunctionsMap()).Parse(tpl))
	var b bytes.Buffer
	err := t.Execute(&b, vars)
	if err != nil {
		return err
	}
	if expect != b.String() {
		return fmt.Errorf("Expected '%s', got '%s'", expect, b.String())
	}
	return nil
}
