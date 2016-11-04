package configv2

import (
	"fmt"
	"reflect"
	"testing"
)

// TestField 测试reflect struct 中的field
func testField(t *testing.T) {
	type S struct {
		a int
		b float32
		s string
	}

	type SS struct {
		st S
		B  bool
	}

	s := &SS{}
	rv := reflect.ValueOf(s)

	fmt.Printf("st value: %v type: %v %v\n", rv, reflect.TypeOf(s), reflect.TypeOf(rv))
	if rv.Kind() == reflect.Ptr {
		//st = st.
		rv = reflect.Indirect(rv)
	}
	t.Log(rv.FieldByName("B").CanSet(), rv.FieldByName("B").CanAddr())
	rv.FieldByName("B").SetBool(true)
	st := reflect.TypeOf(rv)
	t.Logf("indirect rv type: %v\n", st)
	t.Logf("%v\n", s)

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		t.Logf("field %d name: %s\n", i, field.Name)
		if field.Name == "st" {

		}
	}
}
