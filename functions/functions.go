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
	"bufio"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"text/template"
)

func FunctionsMap() template.FuncMap {
	return template.FuncMap(fmap)
}

var Random = rand.New(rand.NewSource(0))
var data = map[string][]string{}
var fmap = map[string]interface{}{

	// text utilities
	"upper":                    strings.ToUpper,
	"lower":                    strings.ToLower,
	"title":                    cases.Title(language.English).String,
	"trim":                     strings.TrimSpace,
	"squeeze":                  func(s string) string { return strings.ReplaceAll(s, " ", "") },
	"substr":                   func(start, length int, s string) string { return s[start:length] },
	"first":                    func(s string) string { return s[:1] },
	"firstword":                func(s string) string { return strings.Split(s, " ")[0] },
	"join":                     strings.Join,
	"repeat":                   strings.Repeat,
	"squeezechars":             func(s, c string) string { return strings.ReplaceAll(s, c, "") },
	"replaceall":               strings.ReplaceAll,
	"trimchars":                strings.Trim,
	"atoi":                     strconv.Atoi,
	"split":                    strings.Split,
	"markov":                   nonsense,
	"lorem":                    lorem,
	"sentence":                 sentence,
	"sentence_prefix":          sentencePrefix,
	"regex":                    regex,
	"random":                   func(s []string) string { return s[Random.Intn(len(s))] },
	"randoms":                  func(s string) string { a := strings.Split(s, "|"); return a[Random.Intn(len(a))] },
	"counter":                  counter,
	"random_string":            randomString,
	"random_string_vocabulary": randomStringVocabulary,
	"from":                     word,
	"from_at":                  wordAt,
	"from_shuffle":             wordShuffle,
	"from_n":                   wordShuffleN,

	// math utilities
	"add":       func(a, b int) int { return a + b },
	"sub":       func(a, b int) int { return a - b },
	"div":       func(a, b int) int { return a / b },
	"mod":       func(a, b int) int { return a % b },
	"mul":       func(a, b int) int { return a * b },
	"max":       math.Max,
	"min":       math.Min,
	"integer":   func(min, max int) int { return min + Random.Intn(max-min) },
	"integer64": func(min, max int64) int64 { return min + Random.Int63n(max-min) },
	"floating":  func(min, max float32) float32 { return min + Random.Float32()*(max-min) },

	// networking and time utilities
	"ip":                ip,
	"ipv6":              ipv6,
	"mac":               mac,
	"http_method":       httpMethod,
	"ip_known_protocol": ipKnownProtocol,
	"ip_known_port":     ipKnownPort,
	"password":          password,
	"useragent":         userAgent,
	"unix_time_stamp":   unixTimeStamp,

	// people related utilities
	"name":           name,
	"name_m":         nameM,
	"name_f":         nameF,
	"middlename":     middlename,
	"surname":        surname,
	"username":       username,
	"address":        address,
	"capital":        capital,
	"capital_at":     capitalAt,
	"state":          state,
	"state_at":       stateAt,
	"state_short":    stateShort,
	"state_short_at": stateShortAt,
	"zip":            zip,
	"zip_at":         zipAt,
	"company":        company,
	"email_provider": emailProvider,

	// generic utilities
	"key":     func(name string, n int) string { return fmt.Sprintf("%s%d", name, Random.Intn(n)) },
	"seed":    func(rndSeed int64) string { Random.Seed(rndSeed); return "" },
	"uuid":    uniqueId,
	"bool":    randomBool,
	"yesorno": yesOrNo,
	"array":   func(count int) []int { return make([]int, count) },
}

func initialize(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Failed to open file: %s", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error in closing file: %s", err)
		}
	}(file)

	var words []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	return words
}

// word returns a random string from a list of strings in a file.
func word(name string) string {
	_, err := cache(name)
	if err != nil {
		return ""
	}
	words := data[name]
	return words[Random.Intn(len(words))]
}

// wordAt returns a string at a given position in a list of strings in a file.
func wordAt(name string, index int) string {
	_, err := cache(name)
	if err != nil {
		return ""
	}
	words := data[name]
	return words[index]
}

// wordShuffle returns a shuffled list of strings in a file.
func wordShuffle(name string) []string {
	_, err := cache(name)
	if err != nil {
		return []string{""}
	}
	words := data[name]
	return wordShuffleN(name, len(words))
}

// wordShuffleN return a subset of n elements in a list of string in a file.
func wordShuffleN(name string, n int) []string {
	_, err := cache(name)
	if err != nil {
		return []string{""}
	}
	words := data[name]
	Random.Shuffle(len(words), func(i, j int) {
		words[i], words[j] = words[j], words[i]
	})
	return words[:n]
}

// cache is used to internally cache data from word files
func cache(name string) (bool, error) {

	v := data[name]
	if v == nil {
		locale := JrContext.Locales[Random.Intn(len(JrContext.Locales))]
		filename := fmt.Sprintf("%s/data/%s/%s", JrContext.TemplateDir, locale, name)
		data[name] = initialize(filename)
		if len(data[name]) == 0 {
			return false, fmt.Errorf("no words found in %s", filename)
		} else {
			return true, nil
		}
	}
	return false, nil
}
