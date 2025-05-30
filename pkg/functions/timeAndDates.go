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
	"time"

	"github.com/rs/zerolog/log"
)

// UnixTimeStamp returns a random unix timestamp not older than the given number of days (in seconds)
func UnixTimeStamp(days int) int64 {
	return UnixTS(days, false)
}

// UnixTimeStampMS returns a random unix timestamp not older than the given number of days (in milliseconds)
func UnixTimeStampMS(days int) int64 {
	return UnixTS(days, true)
}

func UnixTS(days int, millisecondPrecision bool) int64 {
	unixEpoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	now := time.Now()
	if millisecondPrecision {
		first := now.AddDate(0, 0, -days).Sub(unixEpoch).Milliseconds()
		last := now.Sub(unixEpoch).Milliseconds()
		return Random.Int63n(last-first) + first
	} else {
		first := now.AddDate(0, 0, -days).Sub(unixEpoch).Seconds()
		last := now.Sub(unixEpoch).Seconds()
		return Random.Int63n(int64(last-first)) + int64(first)
	}
}

// DateBetween returns a date between fromDate and toDate
func DateBetween(fromDate string, toDate string) string {
	start, err := time.Parse(time.DateOnly, fromDate)
	if err != nil {
		log.Fatal().Err(err).Msg("Error parsing date")
	}

	end, err := time.Parse(time.DateOnly, toDate)
	if err != nil {
		log.Fatal().Err(err).Msg("Error parsing date")
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

// Justpassed returns a date in the past not before the given milliseconds
func Justpassed(milliseconds int64) string {
	now := time.Now()

	duration := time.Duration(Random.Int63n(milliseconds)) * time.Millisecond
	pastTime := now.Add(-duration)

	return pastTime.Format(time.DateTime)
}

// Now returns the current time as a Unix millisecond timestamp
func Now() int64 {
	return time.Now().UnixMilli()
}

// FormatTimestamp formats a unix millisecond timestamp with the given pattern
func FormatTimestamp(timestamp int64, format string) string {
	t := time.Unix(0, timestamp*int64(time.Millisecond))
	formattedDate := t.Format(format)
	return formattedDate
}

// Nowsub returns a date in the past of given milliseconds
func Nowsub(milliseconds int64) string {
	now := time.Now()

	duration := time.Duration(milliseconds) * time.Millisecond
	pastTime := now.Add(-duration)

	return pastTime.Format(time.DateTime)
}

// Nowadd returns a date in the future of given milliseconds
func Nowadd(milliseconds int64) string {
	now := time.Now()

	duration := time.Duration(milliseconds) * time.Millisecond
	pastTime := now.Add(duration)

	return pastTime.Format(time.DateTime)
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

// Past returns a date in the past not before the given years
func Past(years int) string {
	now := time.Now().UTC()
	start := now.AddDate(-years, 0, 0)
	delta := now.Sub(start).Nanoseconds()
	randNsec := Random.Int63n(delta)
	d := start.Add(time.Duration(randNsec))
	return d.Format(time.DateOnly)
}

// Future returns a date in the future not after the given years
func Future(years int) string {
	now := time.Now().UTC()
	start := now.AddDate(years, 0, 0)
	delta := start.Sub(now).Nanoseconds()
	randNsec := Random.Int63n(delta)
	d := now.Add(time.Duration(randNsec))
	return d.Format(time.DateOnly)
}

// Recent returns a date in the past not before the given days
func Recent(days int) string {
	now := time.Now().UTC()
	start := now.AddDate(0, 0, -days)
	delta := now.Sub(start).Nanoseconds()
	randNsec := Random.Int63n(delta)
	d := start.Add(time.Duration(randNsec))
	return d.Format(time.DateOnly)
}

// Soon returns a date in the future not after the given days
func Soon(days int) string {
	now := time.Now().UTC()
	start := now.AddDate(0, 0, days)
	delta := start.Sub(now).Nanoseconds()
	randNsec := Random.Int63n(delta)
	d := now.Add(time.Duration(randNsec))
	return d.Format(time.DateOnly)
}
