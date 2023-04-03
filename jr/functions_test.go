package jr

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"text/template"
)

func TestMain(m *testing.M) {
	JrContext.Seed = 0
	fmt.Println("HERE")
	code := m.Run()
	os.Exit(code)
}

func TestSubstr(t *testing.T) {
	tpl := `{{"fooo" | substr 0 3 }}`
	if err := runt(tpl, "foo"); err != nil {
		t.Error(err)
	}
}

func TestSplit(t *testing.T) {
	tpl := `{{split "a|b" "|"}}`
	if err := runt(tpl, "[a b]"); err != nil {
		t.Error(err)
	}
}

func TestTitle(t *testing.T) {
	tpl := `{{"foo" | title}}`
	if err := runt(tpl, "Foo"); err != nil {
		t.Error(err)
	}
}

func TestMax(t *testing.T) {
	tpl := `{{max 1 4}}`
	if err := runt(tpl, "4"); err != nil {
		t.Error(err)
	}
}

func TestMin(t *testing.T) {
	tpl := `{{min 1 4}}`
	if err := runt(tpl, "1"); err != nil {
		t.Error(err)
	}
}

func TestUSState(t *testing.T) {

	hawaii := "{{capital_at 10}} {{state_at 10}} {{state_short_at 10}} {{zip_at 10}}"
	massachussets := `{{capital_at 20}} {{state_at 20}} {{state_short_at 20}} {{zip_at 20}}`
	newyork := `{{capital_at 31}} {{state_at 31}} {{state_short_at 31}} {{zip_at 31}}`
	texas := `{{capital_at 42}} {{state_at 42}} {{state_short_at 42}} {{zip_at 42}}`
	virginia := `{{capital_at 45}} {{state_at 45}} {{state_short_at 45}} {{zip_at 45}}`
	wyoming := `{{capital_at 49}} {{state_at 49}} {{state_short_at 49}} {{zip_at 49}}`
	if err := runt(hawaii, "Honolulu Hawaii HI 96813"); err != nil {
		t.Error(err)
	}
	if err := runt(massachussets, "Boston Massachusetts MA 02201"); err != nil {
		t.Error(err)
	}
	if err := runt(newyork, "Albany New York NY 12207"); err != nil {
		t.Error(err)
	}
	if err := runt(texas, "Austin Texas TX 78701"); err != nil {
		t.Error(err)
	}
	if err := runt(virginia, "Richmond Virginia VA 23219"); err != nil {
		t.Error(err)
	}
	if err := runt(wyoming, "Cheyenne Wyoming WY 82001"); err != nil {
		t.Error(err)
	}
}

func TestShuffle(t *testing.T) {
	tpl := `{{from_shuffle "state"}}`
	if err := runt(tpl, "[North Dakota Michigan Colorado Montana California Mississippi South Carolina Indiana "+
		"North Carolina Virginia New York Texas Alaska Wyoming Oregon Florida Maryland Ohio Minnesota Pennsylvania "+
		"Kansas Arkansas Nebraska Arizona Hawaii Louisiana Washington South Dakota Massachusetts Connecticut Vermont "+
		"Kentucky Wisconsin Utah Nevada New Hampshire Alabama New Jersey West Virginia Missouri Idaho Oklahoma "+
		"Rhode Island Illinois Delaware Tennessee New Mexico Georgia Iowa Maine]"); err != nil {
		t.Error(err)
	}
}

func TestShuffleN(t *testing.T) {
	tpl := `{{from_n "state" 3}}`
	if err := runt(tpl, "[Massachusetts Utah Kansas]"); err != nil {
		t.Error(err)
	}
}

func TestCache(t *testing.T) {
	v, f := cache("wine")

	if v != true || f != nil {
		t.Error("cache should be empty, no errors")
	}
	v, f = cache("wine")
	if v != false || f != nil {
		t.Error("cache should be full, no errors")
	}
	v, f = cache("wines")
	if f == nil {
		t.Error("no cacheable, should get error")
	}
}

func TestFrom(t *testing.T) {
	tpl := `{{from "actor"}}`
	if err := runt(tpl, "Kate Winslet"); err != nil {
		t.Error(err)
	}
	tpl = `{{from "actors"}}`
	if err := runt(tpl, ""); err != nil {
		t.Error(err)
	}
}

func TestPassword(t *testing.T) {
	tpl := `{{password 5 true "PwD" "!?!"}}`
	if err := runt(tpl, "PwDUSiqU!?!"); err != nil {
		t.Error(err)
	}
}

func TestIPv6(t *testing.T) {
	tpl := `{{ipv6}}`
	if err := runt(tpl, "9e0f:87e3:f10:20c1:3384:666c:f7a9:ffad"); err != nil {
		t.Error(err)
	}
}

func TestIP(t *testing.T) {
	tpl := `{{ip "10.2.0.0/16"}}`
	if err := runt(tpl, "10.2.64.220"); err != nil {
		t.Error(err)
	}
}

func TestUseragent(t *testing.T) {
	tpl := `{{useragent}}`
	if err := runt(tpl, "Mozilla/5.0 (iOS 14_4_2) AppleWebKit/515.46 (KHTML, like Gecko) Edge Mobile/47.0.2.2 Mobile Safari/7.1"); err != nil {
		t.Error(err)
	}
}

func TestCounter(t *testing.T) {
	tpl := `{{counter "A" 0 1}},{{counter "B" 2 2}},{{counter "C" -4 1}},{{counter "D" 0 -1}}`

	if err := runt(tpl, "0,2,-4,0"); err != nil {
		t.Error(err)
	}

	if err := runt(tpl, "1,4,-3,-1"); err != nil {
		t.Error(err)
	}

	if err := runt(tpl, "2,6,-2,-2"); err != nil {
		t.Error(err)
	}
}

func TestArray(t *testing.T) {
	tpl := `{{array 5}}`
	if err := runt(tpl, "[0 0 0 0 0]"); err != nil {
		t.Error(err)
	}

	tpl2 := `{{array 1}}`
	if err := runt(tpl2, "[0]"); err != nil {
		t.Error(err)
	}

	tpl3 := `{{array 0}}`
	if err := runt(tpl3, "[]"); err != nil {
		t.Error(err)
	}
}

func runt(tpl, expect string) error {
	return runtv(tpl, expect, "")
}

func TestRegex(t *testing.T) {
	tpl := `{{regex "Z{2,5}"}}`
	if err := runt(tpl, "ZZZZ"); err != nil {
		t.Error(err)
	}
	//123[0-2]+.*\w{3}
	//tpl = "{{regex `123[0-2]+.*\w{3}`}}"
	//if err := runt(tpl, "ZZZ"); err != nil {
	//	t.Error(err)
	//}
}
func runtv(tpl, expect string, vars interface{}) error {
	t := template.Must(template.New("test").Funcs(FunctionsMap()).Parse(tpl))
	var b bytes.Buffer
	err := t.Execute(&b, vars)
	if err != nil {
		return err
	}
	if expect != b.String() {
		return fmt.Errorf("Expected '%s', got '%s'", expect, b.String())
	}
	return nil
}
