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

package functions_test

import (
	"github.com/jrnd-io/jr/pkg/functions"
	"testing"
)

type cardTestcase struct {
	card string
}

var cards = []cardTestcase{
	{`37144963539843`},
	{`3056930902590`},
	{`601111111111111`},
	{`353011133330000`},
	{`555555555555444`},
	{`411111111111111`},
	{`7992739871`},
}

var luhnCheck = []string{"1", "4", "7", "0", "4", "1", "3"}

var cusips = []cardTestcase{
	{`252IXRF6`},
	{`929WRP5E`},
	{`93114210`},
	{`03783310`},
	{`17275R10`},
	{`38259P50`},
	{`59491810`},
	{`68389X10`},
}

var cusipCheck = []string{"2", "9", "3", "0", "2", "8", "4", "5"}

func TestLuhn(t *testing.T) {
	for i, test := range cards {
		d := functions.LuhnCheckDigit(test.card)
		if d != luhnCheck[i] {
			t.Error("Luhn digit for", test.card, " is not right, it's ", d, " should be ", luhnCheck[i])
		}
	}
}

func TestCusip(t *testing.T) {
	for i, test := range cusips {
		d := functions.CusipCheckDigit(test.card)
		if d != cusipCheck[i] {
			t.Error("Cusip digit for", test.card, " is not right, it's ", d, " should be ", cusipCheck[i])
		}
	}
}

func TestSedol(t *testing.T) {
	sedol := "026349"
	check := functions.SedolCheckDigit(sedol)
	if check != "4" {
		t.Error("Cusip digit for", sedol, " is not right, it's ", check, " should be ", "4")
	}
}
