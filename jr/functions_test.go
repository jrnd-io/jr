package jr

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"
)

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

func TestRandomness(t *testing.T) {
	tpl := `{{integer 0 100}}`
	if err := runt(tpl, "74"); err != nil {
		t.Error(err)
	}
	tpl_seed := `{{seed 100}}{{integer 0 100}}`
	if err := runt(tpl_seed, "83"); err != nil {
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
	tpl := `{{shuffle "state"}}`
	if err := runt(tpl, "[Virginia Connecticut North Carolina Pennsylvania Nevada Maryland Alaska Indiana Idaho "+
		"Delaware California Oklahoma Vermont South Dakota Tennessee Colorado Hawaii Wisconsin New York Arizona Iowa Florida "+
		"Ohio Minnesota West Virginia Rhode Island Louisiana Washington Arkansas Massachusetts Illinois Texas New Jersey Kentucky "+
		"Mississippi New Mexico Maine Oregon Nebraska North Dakota Kansas Montana South Carolina Wyoming Alabama Utah Missouri "+
		"Michigan New Hampshire Georgia]"); err != nil {
		t.Error(err)
	}
}

func TestShuffleN(t *testing.T) {
	tpl := `{{shuffle_n "state" 3}}`
	if err := runt(tpl, "[Utah Idaho New Mexico]"); err != nil {
		t.Error(err)
	}
}

func TestPassword(t *testing.T) {
	tpl := `{{password 5 true "PwD" "!?!"}}`
	if err := runt(tpl, "PwDarOSA!?!"); err != nil {
		t.Error(err)
	}
}

func TestIPv6(t *testing.T) {
	tpl := `{{ipv6}}`
	if err := runt(tpl, "3e5f:6418:9bd0:f7b7:d9d5:7121:ddf9:e87c"); err != nil {
		t.Error(err)
	}
}

func TestIP(t *testing.T) {
	tpl := `{{ip "10.2.0.0/16"}}`
	if err := runt(tpl, "10.2.55.217"); err != nil {
		t.Error(err)
	}
}

func TestUseragent(t *testing.T) {
	tpl := `{{useragent}}`
	if err := runt(tpl, "Mozilla/5.0 (iOS 14_0) AppleWebKit/570.1 (KHTML, like Gecko) Firefox Mobile/6.3 Mobile Safari/7.6"); err != nil {
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
