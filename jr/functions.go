package jr

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"text/template"
)

func FunctionsMap() template.FuncMap {
	return template.FuncMap(fmap)
}

var Random = rand.New(rand.NewSource(0))

var fmap = map[string]interface{}{

	// text utilities
	"trim":      strings.TrimSpace,
	"upper":     strings.ToUpper,
	"lower":     strings.ToLower,
	"title":     strings.Title,
	"substr":    func(start, length int, s string) string { return s[start:length] },
	"first":     func(s string) string { return s[:1] },
	"firstword": func(s string) string { return strings.Split(s, " ")[0] },
	"join":      func(sep string, ss []string) string { return strings.Join(ss, sep) },
	"repeat":    func(count int, str string) string { return strings.Repeat(str, count) },
	"trimall":   func(c, s string) string { return strings.Trim(s, c) },
	"atoi":      func(a string) int { i, _ := strconv.Atoi(a); return i },
	"split":     func(sep, s string) []string { return strings.Split(s, sep) },
	"markov":    func(prefixLen, numWords int, baseText string) string { return Nonsense(prefixLen, numWords, baseText) },
	"lorem":     func(size int) string { return Lorem(size) },

	//math utilities
	"add":      func(a, b int) int { return a + b },
	"sub":      func(a, b int) int { return a - b },
	"div":      func(a, b int) int { return a / b },
	"mod":      func(a, b int) int { return a % b },
	"mul":      func(a, b int) int { return a * b },
	"max":      math.Max,
	"min":      math.Min,
	"integer":  func(min, max int) int { return min + Random.Intn(max-min) },
	"floating": func(min, max float32) float32 { return min + Random.Float32()*(max-min) },
	"Random":   func(s []string) string { return s[Random.Intn(len(s))] },
	"randoms":  func(s string) string { a := strings.Split(s, "|"); return a[Random.Intn(len(a))] },

	//networking and time utilities
	"ip":                 func(s string) string { return ip(s) },
	"ip_known_protocols": ipKnownProtocols,
	"ip_known_ports":     ipKnownPorts,
	"unix_time_stamp":    func(days int) int64 { return unixTimeStamp(days) },

	//people
	"name":   name,
	"name_m": nameM,
	"name_f": nameF,

	"surname":        surname,
	"middlename":     middlename,
	"address":        address,
	"capital":        capital,
	"capital_at":     capitalAt,
	"state":          state,
	"state_at":       func(index int) string { return stateAt(index) },
	"state_short":    stateShort,
	"state_short_at": func(index int) string { return stateShortAt(index) },
	"zip":            zip,
	"zip_at":         func(index int) string { return zipAt(index) },
	"company":        company,
	"email_provider": emailProvider,

	//generic
	"seed":    func(rndSeed int64) string { Random.Seed(rndSeed); return "" },
	"uuid":    uniqueId,
	"bool":    randomBool,
	"yesorno": yesOrNo,
}
