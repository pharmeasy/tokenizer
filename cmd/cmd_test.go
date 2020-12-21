package cmd

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestNewRootCmd_Version(t *testing.T) {
	configFile = "tokenizer-dev.yml"
	os.Args[2] = "--config=tokenizer-dev.yml"
	cmd := NewRootCmd()
	buf := bytes.Buffer{}
	cmd.SetOut(&buf)
	os.Args[1] = "version"
	cmd.Execute()

	o, _ := ioutil.ReadAll(&buf)

	data, _ := ioutil.ReadFile("../VERSION")
	fmt.Println("version " + string(data))
	fmt.Println(string(o))
	assert.True(t, strings.Contains(string(o), "version "))
}
