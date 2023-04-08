//Adapted for JR from https://github.com/lucasjones/reggen

package functions

import (
	"fmt"
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
		for i := 0; i < 10; i++ {
			r, err := regex(test.regex)
			if err != nil {
				t.Fatal("Error creating generator: ", err)
			}

			if i < 1 {
				fmt.Printf("Regex: %v Result: \"%s\"\n", test.regex, r)
			}
			re, err := regexp.Compile(test.regex)
			if err != nil {
				t.Fatal("Invalid test case. regex: ", test.regex, " failed to compile:", err)
			}
			if !re.MatchString(r) {
				t.Error("Generated data does not match regex. Regex: ", test.regex, " output: ", r)
			}
		}
	}
}
