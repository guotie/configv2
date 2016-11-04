package configv2

import (
	"fmt"
	"testing"
)

func TestFileConfig(t *testing.T) {
	type Su struct {
		F float64 `json:"fl"`
	}
	type ST struct {
		A   int
		F   float64 `json:"f"`
		Str string
		B   bool

		Sub Su
	}

	var (
		s ST
	)
	fc := fileConfig{filename: "./tests/config.json"}

	err := fc.Read(&s)
	if err != nil {
		fmt.Printf("Read failed: %v\n", err)
		t.Fail()
		return
	}

	fmt.Println(s)
}
