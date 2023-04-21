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
	"log"
	"time"
)

// UnixTimeStamp returns a random unix timestamp not older than the given number of days
func UnixTimeStamp(days int) int64 {
	unixEpoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	now := time.Now()
	first := now.AddDate(0, 0, -days).Sub(unixEpoch).Seconds()
	last := now.Sub(unixEpoch).Seconds()
	return Random.Int63n(int64(last-first)) + int64(first)
}

// DateBetween returns a date between fromDate and toDate
func DateBetween(fromDate string, toDate string) string {
	start, err := time.Parse(time.DateOnly, fromDate)
	if err != nil {
		log.Fatal(err)
	}

	end, err := time.Parse(time.DateOnly, toDate)
	if err != nil {
		log.Fatal(err)
	}

	delta := end.Sub(start).Nanoseconds()
	randNsec := Random.Int63n(delta)

	d := start.Add(time.Duration(randNsec))
	return d.Format(time.DateOnly)
}

// DatesBetween returns an array of num dates between fromDate and toDate
func DatesBetween(fromDate string, toDate string, num int) []string {

	dates := make([]string, num)
	for i := 0; i < len(dates); i++ {
		dates[i] = DateBetween(fromDate, toDate)
	}
	return dates
}

// BirthDate returns a birthdate between minAge and maxAge
func BirthDate(minAge int, maxAge int) string {

	maxBirthYear := time.Now().Year() - minAge
	minBirthYear := maxBirthYear - (maxAge - minAge)

	birthYear := Random.Intn(maxBirthYear-minBirthYear+1) + minBirthYear

	birthMonth := Random.Intn(12) + 1
	lastDayOfMonth := time.Date(birthYear, time.Month(birthMonth+1), 0, 0, 0, 0, 0, time.UTC).Day()
	birthDay := Random.Intn(lastDayOfMonth) + 1

	d := time.Date(birthYear, time.Month(birthMonth), birthDay, 0, 0, 0, 0, time.UTC)
	return d.Format(time.DateOnly)
}
