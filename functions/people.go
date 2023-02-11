package functions

import "math/rand"

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
	ports := []string{"80", "81", "443", "22", "631"}
	return ports[rand.Intn(len(ports))]
}

func surname() string {
	ports := []string{"80", "81", "443", "22", "631"}
	return ports[rand.Intn(len(ports))]
}

func middlename() string {
	ports := []string{"80", "81", "443", "22", "631"}
	return ports[rand.Intn(len(ports))]
}

func address() string {
	ports := []string{"80", "81", "443", "22", "631"}
	return ports[rand.Intn(len(ports))]
}

func city() string {
	ports := []string{"80", "81", "443", "22", "631"}
	return ports[rand.Intn(len(ports))]
}

func state() string {
	ports := []string{"80", "81", "443", "22", "631"}
	return ports[rand.Intn(len(ports))]
}
