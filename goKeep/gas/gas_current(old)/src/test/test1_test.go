package main

import (
	"testing"
)

var gc gasComponent

var invDates = []string{"lalalal"}
var vdates = []string{"2013-12-15T09:26:14.635"}

func TestDate(t *testing.T) {
	gc = gasComponent{}
	gc.LoadJsonGas()
	for _, d := range invDates {
		if gc.IcheckDate(d) == "valid" {
			t.Fatalf("invalid date is valid: \n%v", )
		}
	}

	for _, d := range vdates {
		if gc.IcheckDate(d) == "invalid" {
			t.Fatalf("valid date is invalid: \n%v", )
		}
	}
}
