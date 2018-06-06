package raz_utils_test

import (
	"bytes"
	"errors"
	"log"
	rz "raz_utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogError(t *testing.T) {
	rz.LogError(errors.New("This is err1"))
	assert.Equal(t, "[error] in raz_utils_test.TestLogError[/home/raz/IdeaProjects1/projects1/gas_03_2017/src/raz_utils/raz_utils_test.go:14:0] This is err1\n", buf.String(), "This should be equal")
}

var buf bytes.Buffer

func init() {

	log.SetOutput(&buf)
}
