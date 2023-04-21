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
	"github.com/squeeze69/generacodicefiscale"
	"log"
	"strconv"
	"unicode"
)

type Person struct {
	Name    string
	Surname string
	Gender  string
}

type City struct {
	Name  string
	Zip   string
	State string
}

// Name returns a random Name (male/female)
func Name() string {
	s := Random.Intn(2)
	if s == 0 {
		return NameM()
	} else {
		return NameF()
	}
}

// NameM returns a random male Name
func NameM() string {
	return Word("nameM")
}

// NameF returns a random female Name
func NameF() string {
	return Word("nameF")
}

// Surname returns a random Surname
func Surname() string {
	return Word("surname")
}

// Middlename returns a random Middlename
func Middlename() string {
	middles := []string{"M", "J", "K", "P", "T", "S"}
	return middles[Random.Intn(len(middles))]
}

// Address returns a random Address
func Address() string {
	//TODO: implement
	addresses := []string{""}
	return addresses[Random.Intn(len(addresses))]
}

// Country returns a random ISO 3166 Country
func Country() string {
	return Word("country")
}

// State returns a random State
func State() string {
	return Word("state")
}

// StateAt returns State at given index
func StateAt(index int) string {
	return WordAt("state", index)
}

// StateShort returns a random short State
func StateShort() string {
	return Word("state_short")
}

// StateShortAt returns short State at given index
func StateShortAt(index int) string {
	return WordAt("state_short", index)
}

// Capital returns a random Capital
func Capital() string {
	return Word("capital")
}

// CapitalAt returns Capital at given index
func CapitalAt(index int) string {
	return WordAt("capital", index)
}

// Zip returns a random Zip code
func Zip() string {
	return Word("zip")
}

// ZipAt returns Zip code at given index
func ZipAt(index int) string {
	return WordAt("zip", index)
}

// Company returns a random Company Name
func Company() string {
	return Word("company")
}

// EmailProvider returns a random mail provider
func EmailProvider() string {
	return Word("mail_provider")
}

// CodiceFiscale return a valid Italian Codice Fiscale
func CodiceFiscale(name string, surname string, sex string, birth string, city string) string {

	codicecitta, erc := generacodicefiscale.CercaComune(city)
	if erc != nil {
		log.Fatal(erc)
	}
	cf, erg := generacodicefiscale.Genera(surname, name, sex, codicecitta.Codice, birth)
	if erg != nil {
		log.Fatal(erg)
	}
	return cf
}

// IsVowel reports whether the rune is an ASCII and vowel case letter.
func isVowel(s rune) bool {
	switch unicode.ToUpper(s) {
	case 'A', 'E', 'I', 'O', 'U':
		return true
	default:
		return false
	}
}

func isConsonant(s rune) bool {
	return !isVowel(s)
}

// Ssn return a valid Social Security Number id
func Ssn() string {
	first := Random.Intn(899) + 1
	second := Random.Intn(99) + 1
	third := Random.Intn(9999) + 1
	return fmt.Sprintf("%03d-%02d-%04d", first, second, third)
}

// Username returns a random Username using Name, Surname and a length
func Username(firstName string, lastName string, size int) string {

	var name string

	useSurname := (Random.Intn(2)) != 0
	shuffleName := (Random.Intn(2)) != 0
	if useSurname || len(name) < size {
		name = firstName + lastName
	} else {
		name = firstName
	}

	if shuffleName {
		nameRunes := []rune(name)
		Random.Shuffle(len(nameRunes), func(i, j int) {
			nameRunes[i], nameRunes[j] = nameRunes[j], nameRunes[i]
		})
		name = string(nameRunes)
	}

	var username string
	for i := 0; i < size && i < len(name); i++ {
		username += string(name[i])
	}

	username += strconv.Itoa(50 + Random.Intn(49))

	return username
}
