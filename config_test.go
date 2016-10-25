package config

import (
	"fmt"
	//"reflect"
	"testing"
)

func init() {
	ReadCfg("./test.json")
}

func TestString(t *testing.T) {
	if GetString("redisProto") != "tcp" {
		t.Fatal("get redisProto failed!")
	}
}

func TestInt(t *testing.T) {
	maxNews, ok := GetInt("maxNews")
	if (!ok) || (maxNews != 5000) {
		t.Fatal("get maxnews failed:", ok, maxNews)
	}
}

func TestFloat(t *testing.T) {
	ft, ok := GetFloat("floatTest")
	if (!ok) || (ft != 1.67) {
		t.Fatal("get float failed:", ok, ft)
	}
}

func TestScanStruct(t *testing.T) {
	type Is struct {
		Strin string `json:"str"`
		Iin   int64
	}
	type OS struct {
		Inner Is
		Str   string
		I     int
		Slc   []int
	}
	var (
		s1, s2 OS
		s3     *OS
		s4     ***OS
	)

	s1 = OS{}
	err := Scan("ts", &s1)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(s1)

	Scan("ts", &s2)
	fmt.Println(s2)

	Scan("ts", &s3)
	fmt.Println(s3)

	Scan("ts", &s4)
	fmt.Println(s4, *s4, **s4)
}

func TestScanMap(t *testing.T) {
	type IE3 struct {
		Ie1 int
		Ie2 string
		Ie3 float32
	}

	var (
		m11 map[string]map[string]int
		m12 map[string]map[string]*int
		m2  map[string]*map[string]int
	)

	err := Scan("tm3", &m11)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("m11:", m11)

	err = Scan("tm3", &m12)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("m12:", m12)

	err = Scan("tm3", &m2)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("m2:", m2)

	m3 := make(map[string]map[string]int)
	err = Scan("tm3", &m3)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("m3:", m3)

	m4 := make(map[string]IE3)
	err = Scan("tm4", &m4)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("m4:", m4)

	var m5 map[string]*IE3
	err = Scan("tm4", &m5)
	fmt.Println("m5:", m5, m5["e1"], m5["e2"], m5["e3"])

	var m6 map[string]***IE3
	err = Scan("tm4", &m6)
	fmt.Println(m6, m6["e1"], *m6["e1"], **m6["e1"])

	var (
		m71 map[string][]int
		m72 map[string]*[]int
		m73 map[string]*[]*int
		m8  map[string][]string
	)
	err = Scan("m7x", &m71)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("m71:", m71)

	err = Scan("m7x", &m72)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("m72:", m72)

	err = Scan("m7x", &m73)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("m73:", m73, m73["s1"], *(*m73["s1"])[0])

	err = Scan("m8x", &m8)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("m8:", m8, len(m8["ms3"]))
}

func testScan4(t *testing.T) {
	var s1, is1 []int
	var s2, is2 []string

	err := Scan("sl1", &s1)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("s1:", s1)

	is1 = make([]int, 0)
	err = Scan("sl1", &is1)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("is1:", is1)

	err = Scan("sl2", &s2)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("s2:", s2)

	is2 = make([]string, 0)
	err = Scan("sl2", &is2)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("is2:", is2)
}

func TestScanSlice(t *testing.T) {
	type St struct {
		I   int
		Str string
		Ar  []string
	}
	var (
		st1 []St
		st2 []map[string]int
		st3 [][]int
	)

	var s1, is1 []int
	var ps1 *[]int
	var s2, is2 []string

	err := Scan("sl1", &s1)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("s1:", s1)

	is1 = make([]int, 0)
	err = Scan("sl1", &is1)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("is1:", is1)

	err = Scan("sl1", &ps1)
	fmt.Println("ps1:", ps1)

	err = Scan("sl2", &s2)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("s2:", s2)

	is2 = make([]string, 0)
	err = Scan("sl2", &is2)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("is2:", is2)

	Scan("ss1", &st1)
	fmt.Println("st1", st1)
	Scan("ss2", &st2)
	fmt.Println("st2", st2)
	Scan("ss3", &st3)
	fmt.Println("st3", st3)
}

func TestScanInterface(t *testing.T) {
	var i1 []interface{}

	err := Scan("if1", &i1)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("i1:", i1)

	var i2 map[string]interface{}
	err = Scan("if2", &i2)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("i2:", i2)

	var i3 interface{}
	err = Scan("if3", &i3)
	fmt.Println("i3:", i3)
	var i4 []interface{}
	err = Scan("if3", &i4)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestScanArray(t *testing.T) {
	var (
		a1 [3]int
		a2 [6]int
	)

	err := Scan("a1", &a1)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(a1)
	err = Scan("a1", &a2)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(a2)
}

func TestScanError(t *testing.T) {
	var a int

	err := Scan("bl", &a)
	if err != nil {
		fmt.Println(err.Error())
	}
}
