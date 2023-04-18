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
	"strconv"
	"strings"
)

// Cusip returns a valid 9 characters Cusip code
func Cusip() string {
	cusip, _ := Regex("^[0-9]{3}[0-9A-Z]{5}")
	check := cusipCheckDigit(cusip)
	return cusip + check
}

// Valor returns a valid 5-9 digits Valor code
func Valor() string {
	valor, _ := Regex("[0-9a-zA-Z]{5,9}")
	return valor
}

// Isin returns a valid 12 characters Isin code
func Isin() string {
	return "TBD"
}

// returns a valid 7 characters sedol code
func Sedol() string {
	sedol, _ := Regex("[0-9a-zA-Z]{6}")
	return sedol
}

// returns a valid 6 characters wkn code
func Wkn() string {
	wkn, _ := Regex("[0-9a-zA-Z]{6}")
	return wkn
}

// CreditCard returns a valid credit card
func CreditCard(issuer string) string {

	var regex string
	switch issuer {
	case "visa":
		regex = "^4\\d{12}(?:\\d{2})?$"
	case "mastercard":
		regex = "^5[1-5]\\d{13}$"
	case "amex":
		regex = "^3[47]\\d{12,13}$"
	case "discover":
		regex = "^(?:6011\\d{12})|(?:65\\d{13})$"
	default:
		return ""
	}
	card, _ := Regex(regex)
	check := luhnCheckDigit(card)
	return card + check
}

func cusipCheckDigit(code string) string {
	var sum, v int

	if len(code) != 8 {
		return ""
	}

	code = strings.ToUpper(code)

	for i := 0; i < 8; i++ {
		c := code[i]
		if c >= 48 && c <= 57 {
			v = int(c - 48)
		} else if c == 42 {
			v = 36
		} else if c == 64 {
			v = 37
		} else if c == 35 {
			v = 38
		} else if c >= 65 && c <= 90 {
			v = int(c - 55)
		}
		if (7-i)%2 == 0 {
			v = v * 2
		}
		sum += (v / 10) + v%10
	}
	return strconv.Itoa((10 - sum%10) % 10)

}

func luhnCheckDigit(code string) string {
	var sum, v int
	l := len(code)
	for i := 0; i < l; i++ {
		c := code[l-i-1]
		if c >= 48 && c <= 57 {
			v = int(c - 48)
		} else {
			v = int(c - 55)
		}
		if i%2 == 0 {
			v = v * 2
			if v > 9 {
				v = v - 9
			}
		}
		sum += v
	}
	return strconv.Itoa((10 - sum%10) % 10)

}
