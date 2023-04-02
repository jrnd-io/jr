package jr

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
	return word("state-short")
}

func stateShortAt(index int) string {
	return wordAt("state-short", index)
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
