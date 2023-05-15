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

// CountryCode returns a random Country Code prefix
func CountryCode() string {
	countryIndex := JrContext.CountryIndex
	if countryIndex == -1 {
		return Word("country_code")
	} else {
		return WordAt("country_code", countryIndex)
	}
}

// CountryCodeAt returns a Country Code prefix at a given index
func CountryCodeAt(index int) string {
	return WordAt("country_code", index)
}

// Imei returns a random imei number of 15 digits
func Imei() string {
	account := make([]byte, 14)
	for i := range account {
		account[i] = digits[Random.Intn(len(digits))]
	}
	first14 := string(account)
	return first14 + luhnCheckDigit(first14)
}

// Phone returns a random land prefix
func Phone() string {
	cityIndex := JrContext.CityIndex
	if cityIndex == -1 {
		l := Word("land_prefix")
		lp, _ := Regex(l)
		return lp
	} else {
		return PhoneAt(cityIndex)
	}

}

// PhoneAt returns a land prefix at a given index
func PhoneAt(index int) string {
	l := WordAt("land_prefix", index)
	lp, _ := Regex(l)
	return lp
}

// MobilePhone returns a random mobile phone
func MobilePhone() string {
	countryIndex := JrContext.CountryIndex
	if countryIndex == -1 {
		m := Word("mobile_phone")
		mp, _ := Regex(m)
		return mp
	} else {
		return MobilePhoneAt(countryIndex)
	}
}

// MobilePhoneAt returns a mobile phone at a given index
func MobilePhoneAt(index int) string {
	m := WordAt("mobile_phone", index)
	mp, _ := Regex(m)
	return mp
}
