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
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// Account returns a random account number of given length
func Account(length int) string {
	account := make([]byte, length)
	for i := range account {
		account[i] = digits[Random.Intn(len(digits))]
	}
	return string(account)
}

// Amount returns an amount of money between min and max, and given currency
func Amount(min float32, max float32, currency string) string {
	amount := min + Random.Float32()*(max-min)
	return fmt.Sprintf("%s%.2f", currency, amount)
}

// Swift returns a swift/bic code
func Swift() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	bankCode := make([]byte, 4)
	for i := range bankCode {
		bankCode[i] = letters[rand.Intn(len(letters))]
	}
	country := Word("country")
	location := rand.Intn(100)
	branch := rand.Intn(1000)

	return string(bankCode) + string(country) + fmt.Sprintf("%02d", location) + fmt.Sprintf("%03d", branch)

}

// Bitcoin returns a bitcoin address
func Bitcoin() string {
	bc, _ := Regex("^(bc1|[13])[a-zA-HJ-NP-Z0-9]{25,39}$")
	return bc
}

// Cusip returns a valid 9 characters Cusip code
func Cusip() string {
	cusip, _ := Regex("^[0-9]{3}[0-9A-Z]{5}")
	check := cusipCheckDigit(cusip)
	return cusip + check
}

// Isin returns a valid 12 characters Isin code
func Isin(country string) string {
	c := country + Cusip()
	return c + isinCheckDigit(c)
}

// returns a valid 7 characters sedol code
func Sedol() string {
	sedol, _ := Regex("[0-9a-zA-Z]{6}")
	return sedol
}

// Valor returns a valid 5-9 digits Valor code
func Valor() string {
	valor, _ := Regex("[0-9a-zA-Z]{5,9}")
	return valor
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

// cusipCheckDigit returns a valid cusip check digit
func cusipCheckDigit(code string) string {
	var sum, v int

	if len(code) != 8 {
		return ""
	}

	code = strings.ToUpper(code)

	for i := 0; i < 8; i++ {
		c := code[i]
		if c >= '0' && c <= '9' {
			v = int(c - '0')
		} else if c == '*' {
			v = 36
		} else if c == '@' {
			v = 37
		} else if c == '#' {
			v = 38
		} else if c >= 'A' && c <= 'Z' {
			v = int(c - 'A' + 10)
		}
		if i%2 != 0 {
			v = v * 2
		}
		sum += (v / 10) + v%10
	}
	return strconv.Itoa((10 - sum%10) % 10)

}

// luhnCheckDigit returns a valid luhn check digit
func luhnCheckDigit(code string) string {
	var sum, v int
	l := len(code)
	for i := 0; i < l; i++ {
		c := code[l-i-1]
		if c >= '0' && c <= '9' {
			v = int(c - '0')
		} else if c >= 'A' && c <= 'Z' {
			v = int(c - 'A' + 10)
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

func isinCheckDigit(code string) string {
	return "TBD"
}
