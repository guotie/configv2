package configv2

import (
	"bytes"
	"fmt"
	"testing"
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
	}
	type ST struct {
		A   int
		F   float64 `json:"f"`
		Str string
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

	fmt.Println(s, *s.Sub)
}
