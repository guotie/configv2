package configv2

import (
	"bytes"
	"fmt"
	"testing"

	"reflect"

	"github.com/guotie/assert"
)

func TestRemoveComments(t *testing.T) {
	var data = []byte("{\r\n    # ssfdiso\n####\r\n\t\t#\r\n\r\n}")
	ndata := removeComments(data)
	if bytes.Compare(ndata, []byte("{}")) != 0 {
		t.Log(string(ndata))
		t.Fail()
	}
}

func TestFileConfig(t *testing.T) {
	type Su struct {
		F float64 `json:"fl"`
		S string  `default:"abcdefg"`
	}
	type ST struct {
		A   int
		F   float64 `json:"f"`
		Str string  `default:"12345"`
		B   bool

		Sub  *Su
		Arrs []string
	}

	var (
		s ST
	)
	fc := fileConfig{files: []string{"./tests/config.json"}}

	err := fc.Read(&s)
	if err != nil {
		fmt.Printf("Read failed: %v\n", err)
		t.Fail()
		return
	}

	assert.Assert(s.F == 3.1415, "s.F!=3.1415")
	assert.Assert(s.Str == "12345", "s.Str!=12345")
	assert.Assert(s.Sub.F == 1.234, "s.Sub.F!=1.234")
	assert.Assert(s.Sub.S == "abcdefg", "s.Sub.S!=abcdefg")
	assert.Assert(reflect.DeepEqual(s.Arrs, []string{"hello", "world", "you can you up", "no can no bb"}), "s.Attr Not equal")

	v, ok := fc.Get("F")
	assert.Assert(ok == true, "F filed should exist")
	assert.Assert(v == 3.1415, "v==3.1415")
	//fmt.Println(s, *s.Sub)
}
