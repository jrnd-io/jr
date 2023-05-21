//Adapted for JR from https://github.com/lucasjones/reggen

package functions_test

import (
	"github.com/ugol/jr/pkg/functions"
	"regexp"
	"testing"
)

type testCase struct {
	regex string
}

var c = []testCase{
	{`123[0-2]+.*\w{3}`},
	{`^\d{1,2}[/](1[0-2]|[1-9])[/]((19|20)\d{2})$`},
	{`^((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])$`},
	{`^\d+$`},
	{`\D{3}`},
	{`((123)?){3}`},
	{`(ab|bc)def`},
	{`[^abcdef]{5}`},
	{`[^1]{3,5}`},
	{`[[:upper:]]{5}`},
	{`[^0-5a-z\s]{5}`},
	{`Z{2,5}`},
	{`[a-zA-Z]{100}`},
	{`^[a-z]{5,10}@[a-z]{5,10}\.(com|net|org)$`},
}

func TestGenerate(t *testing.T) {
	for _, test := range c {
		r, err := functions.Regex(test.regex)
		if err != nil {
			t.Fatal("Error creating generator: ", err)
		}

		re, err := regexp.Compile(test.regex)
		if err != nil {
			t.Fatal("Invalid test case. Regex: ", test.regex, " failed to compile:", err)
		}
		if !re.MatchString(r) {
			t.Error("Generated data does not match Regex. Regex: ", test.regex, " output: ", r)
		}
	}
}
