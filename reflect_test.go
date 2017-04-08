package configv2

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/guotie/assert"
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

func TestUnmarshal(t *testing.T) {
	var S = struct {
		S string
	}{"abcd"}

	err := json.Unmarshal([]byte(`{}`), &S)
	assert.Assert(err == nil, "err should be nil")
	assert.Assertf(S.S == "abcd", "S.S should be abcd, but %v", S.S)
}

// 测试从struct中提取tag, 并设置struct field的默认值
func TestStructField(t *testing.T) {
	type Sj struct {
		H4 uint64 `json:"h4" default:"640"`
		H  uint   `json:"h" default:"10"`
		B  string `json:"b" default:"hijk"`
		D  bool   `json:"d" default:"true"`
	}
	type Si struct {
		A   []int
		S   string  `default:"bcde"`
		F   float64 `default:"0.3"`
		Sub Sj
	}
	type Ss struct {
		A1  int8    `json:"a1" default:"8"`
		A2  int16   `json:"a2" default:"16"`
		A3  int32   `json:"a3" default:"32"`
		A4  int64   `json:"a4" default:"64"`
		A   int     `json:"a" default:"1"`
		H1  uint8   `json:"h1" default:"80"`
		H2  uint16  `json:"h2" default:"160"`
		H3  uint32  `json:"h3" default:"320"`
		H4  uint64  `json:"h4" default:"640"`
		H   uint    `json:"h" default:"10"`
		B   string  `json:"b" default:"abcd"`
		D   bool    `json:"d" default:"true"`
		C1  float32 `json:"c1" default:"0.32"`
		C2  float64 `json:"c2" default:"0.64"`
		F   []float64
		Sub Si
	}

	var (
		mS Ss
		//mSi Si
	)

	setDefaultValue(reflect.ValueOf(&mS))
	assert.Assert(mS.A1 == 8, "ms.A1!=8")
	assert.Assert(mS.A2 == 16, "ms.A2!=8")
	assert.Assert(mS.A3 == 32, "ms.A3!=8")
	assert.Assert(mS.A4 == 64, "ms.A4!=8")
	assert.Assert(mS.A == 1, "ms.A1!=8")

	assert.Assert(mS.H1 == 80, "ms.H1!=80")
	assert.Assert(mS.H2 == 160, "ms.H2!=160")
	assert.Assert(mS.H3 == 320, "ms.H3!=320")
	assert.Assert(mS.H4 == 640, "ms.H4!=640")
	assert.Assert(mS.H == 10, "ms.H!=10")

	assert.Assert(mS.B == "abcd", "ms.B!=abcd")
	assert.Assert(mS.D == true, "ms.D!=true")
	assert.Assert(mS.C1 == 0.32, "ms.C1!=0.32")
	assert.Assert(mS.C2 == 0.64, "ms.C2!=0.64")

	assert.Assert(mS.Sub.S == "bcde", "ms.Sub.S!=bcde")
	assert.Assert(mS.Sub.F == 0.3, "ms.Sub.F!=0.3")

	assert.Assert(mS.Sub.Sub.B == "hijk", "ms.Sub.Sub.B!=hijk")
	assert.Assert(mS.Sub.Sub.H4 == 640, "ms.Sub.Sub.H4!=640")
	assert.Assert(mS.Sub.Sub.H == 10, "ms.Sub.Sub.H!=10")
	assert.Assert(mS.Sub.Sub.D == true, "ms.Sub.Sub.D!=true")
}

func TestGetField(t *testing.T) {
	type (
		SA struct {
			FA int
			FB bool
			FC float32
			FD string
		}
		SB struct {
			FA SA
			FB string
		}
		SC struct {
			FA map[string]SA
			FB SB
			FC string
		}
	)

	var (
		a = SA{
			FA: 100,
			FB: true,
			FC: 4.321,
			FD: "Struct SA",
		}
		b = SB{
			FA: SA{
				101,
				false,
				3.456,
				"Struct SA in SB",
			},
			FB: "struct SB",
		}
		c = SC{
			FA: map[string]SA{
				"key1": a,
				"key2": b.FA,
			},
			FB: b,
			FC: "Struct SC",
		}

		cases = []struct {
			obj    interface{}
			field  string
			result interface{}
			exist  bool
			expect interface{}
		}{
			{a, "FA", 100, true, a.FA},
			{a, "FB", true, true, a.FB},
			{a, "FE", nil, false, nil},
			{&a, "FC", 4.321, true, a.FC},
			{&a, "FD", "Struct SA", true, a.FD},
			{&a, "FE", nil, false, nil},

			{c, "FA.key1.FA", 100, true, c.FA["key1"].FA},
			{&c, "FA.key1.FA", 100, true, c.FA["key1"].FA},
			{&c, "FB.FA.FA", 101, true, c.FB.FA.FA},
		}
	)

	for idx, c := range cases {
		result, ok := getFields(c.obj, strings.Split(c.field, "."))
		assert.Assertf(ok == c.exist, "exist NOT equal expect: index=%d, obj: %v field: %v", idx, c.obj, c.field)
		if ok {
			assert.Assertf(reflect.DeepEqual(result, c.expect),
				"result NOT equal expect: index=%d result=%v expect=%v result_type=%v expect_type=%v",
				idx, result, c.expect, reflect.TypeOf(result), reflect.TypeOf(c.expect))
		}
	}
}
