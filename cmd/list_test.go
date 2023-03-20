package cmd

import (
	"testing"
)

func TestRightTemplate(t *testing.T) {
	tpl := `{{"fooo" | substr 0 3 }}`
	if !isValidTemplate([]byte(tpl)) {
		t.Error("Template should be valid")
	}
}

func TestRightTemplate2(t *testing.T) {
	tpl := `{"fooo" | substr 0 3 }}`
	if !isValidTemplate([]byte(tpl)) {
		t.Error("Template should be valid")
	}
}

func TestRightTemplate3(t *testing.T) {
	tpl := `{"fooo" | subsTr 0 3 }}`
	if !isValidTemplate([]byte(tpl)) {
		t.Error("Template should be valid")
	}
}

func TestWrongTemplate(t *testing.T) {
	tpl := `{{"fooo" | subsTr 0 3 }}`
	if isValidTemplate([]byte(tpl)) {
		t.Error("Template should be invalid")
	}
}
