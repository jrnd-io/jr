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

import "strconv"

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
	return word("NameM")
}

// NameF returns a random female Name
func NameF() string {
	return word("NameF")
}

// Surname returns a random Surname
func Surname() string {
	return word("Surname")
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

// State returns a random State
func State() string {
	return word("State")
}

// StateAt returns State at given index
func StateAt(index int) string {
	return wordAt("State", index)
}

// StateShort returns a random short State
func StateShort() string {
	return word("state_short")
}

// StateShortAt returns short State at given index
func StateShortAt(index int) string {
	return wordAt("state_short", index)
}

// Capital returns a random Capital
func Capital() string {
	return word("Capital")
}

// CapitalAt returns Capital at given index
func CapitalAt(index int) string {
	return wordAt("Capital", index)
}

// Zip returns a random Zip code
func Zip() string {
	return word("Zip")
}

// ZipAt returns Zip code at given index
func ZipAt(index int) string {
	return wordAt("Zip", index)
}

// Company returns a random Company Name
func Company() string {
	return word("Company")
}

// EmailProvider returns a random mail provider
func EmailProvider() string {
	return word("mail_provider")
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
