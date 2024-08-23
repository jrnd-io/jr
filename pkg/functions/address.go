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

package functions

import (
	"fmt"
	"github.com/jrnd-io/jr/pkg/ctx"
	"math"
)

const (
	earthRadius     = 6371000 // in meters
	degreesPerMeter = 1.0 / earthRadius * 180.0 / math.Pi
)

// BuildingNumber generates a random building number of max n digits
func BuildingNumber(n int) string {
	building := make([]byte, Random.Intn(n)+1)
	for i := range building {
		building[i] = digits[Random.Intn(len(digits))]
	}
	return string(building)
}

// Capital returns a random Capital
func Capital() string {
	return Word("capital")
}

// CapitalAt returns Capital at given index
func CapitalAt(index int) string {
	return WordAt("capital", index)
}

// Cardinal return a random cardinal direction, in long or short form
func Cardinal(short bool) string {
	if short {
		directions := []string{"N", "S", "E", "O", "NE", "NO", "SE", "SO"}
		return directions[Random.Intn(len(directions))]
	}

	directions := []string{"North", "South", "East", "Ovest", "North-East", "North-Ovest", "South-East", "South-Ovest"}
	return directions[Random.Intn(len(directions))]
}

// City returns a random City
func City() string {
	c := Word("city")
	ctx.JrContext.Ctx["_city"] = c
	ctx.JrContext.CityIndex = ctx.JrContext.LastIndex
	return c
}

// CityAt returns City at given index
func CityAt(index int) string {
	return WordAt("city", index)
}

// Country returns the ISO 3166 Country selected with locale
func Country() string {
	countryIndex := ctx.JrContext.CountryIndex
	if countryIndex == -1 {
		return Word("country")
	}

	return WordAt("country", countryIndex)
}

// CountryRandom returns a random ISO 3166 Country
func CountryRandom() string {
	return Word("country")
}

// CountryAt returns an ISO 3166 Country at a given index
func CountryAt(index int) string {
	return WordAt("country", index)
}

// Latitude returns a random latitude between -90 and 90
func Latitude() string {
	latitude := -90 + Random.Float64()*(180)
	return fmt.Sprintf("%.4f", latitude)
}

// Longitude returns a random longitude between -180 and 180
func Longitude() string {
	longitude := -180 + Random.Float64()*(360)
	return fmt.Sprintf("%.4f", longitude)
}

// NearbyGPS returns a random latitude longitude within a given radius in meters
func NearbyGPS(latitude float64, longitude float64, radius int) string {
	radiusInMeters := float64(radius)

	// Generate a random angle in radians
	randomAngle := Random.Float64() * 2 * math.Pi

	// Calculate the distance from the center point
	distanceInMeters := Random.Float64() * radiusInMeters

	// Convert the distance to degrees
	distanceInDegrees := distanceInMeters * degreesPerMeter

	// Calculate the new latitude and longitude
	newLatitude := latitude + (distanceInDegrees * math.Cos(randomAngle))
	newLongitude := longitude + (distanceInDegrees * math.Sin(randomAngle))

	return fmt.Sprintf("%.4f %.4f", newLatitude, newLongitude)

}

// State returns a random State
func State() string {
	s := Word("state")
	ctx.JrContext.Ctx["_state"] = s
	ctx.JrContext.CountryIndex = ctx.JrContext.LastIndex
	return s
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

// Street returns a random street
func Street() string {
	return Word("street")
}

// Zip returns a random Zip code
func Zip() string {
	cityIndex := ctx.JrContext.CityIndex

	if cityIndex == -1 {
		z := Word("zip")
		zip, _ := Regex(z)
		return zip
	}

	return ZipAt(cityIndex)
}

// ZipAt returns Zip code at given index
func ZipAt(index int) string {
	z := WordAt("zip", index)
	zip, _ := Regex(z)
	return zip
}
