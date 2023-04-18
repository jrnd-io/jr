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
)

// Cusip returns a valid 9 characters Cusip code
func Cusip() string {
	return "TBD"
}

// Valor returns a valid 5-9 digits Valor code
func Valor() string {
	return "TBD"
}

// Isin returns a valid 12 characters Isin code
func Isin() string {
	return "TBD"
}

// returns a valid 7 characters sedol code
func Sedol() string {
	return "TBD"
}

// returns a valid 6 characters wkn code
func Wkn() string {
	return "TBD"
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

func luhnCheckDigit(number string) string {
	sum := 0
	l := len(number)
	for i := 0; i < l; i++ {
		digit, _ := strconv.Atoi(string(number[l-i-1]))
		if i%2 == 0 {
			digit = digit * 2
			if digit > 9 {
				digit = digit - 9
			}
		}
		sum += digit
	}
	return strconv.Itoa((10 - sum%10) % 10)

}
