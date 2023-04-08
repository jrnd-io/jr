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

func name() string {
	s := Random.Intn(2)
	if s == 0 {
		return nameM()
	} else {
		return nameF()
	}
}

func nameM() string {
	return word("nameM")
}

func nameF() string {
	return word("nameF")
}

func surname() string {
	return word("surname")
}

func middlename() string {
	middles := []string{"M", "J", "K", "P", "T", "S"}
	return middles[Random.Intn(len(middles))]
}

func address() string {
	//TODO: implement
	addresses := []string{""}
	return addresses[Random.Intn(len(addresses))]
}

func state() string {
	return word("state")
}

func stateAt(index int) string {
	return wordAt("state", index)
}

func stateShort() string {
	return word("state_short")
}

func stateShortAt(index int) string {
	return wordAt("state_short", index)
}

func capital() string {
	return word("capital")
}

func capitalAt(index int) string {
	return wordAt("capital", index)
}

func zip() string {
	return word("zip")
}

func zipAt(index int) string {
	return wordAt("zip", index)
}

func company() string {
	return word("company")
}

func emailProvider() string {
	return word("mail_provider")
}

func username(firstName string, lastName string, size int) string {

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
