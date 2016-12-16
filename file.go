package configv2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"reflect"
	"regexp"
)

const (
	fieldTagName = "field"
)

//
// 主配置文件可以包含子配置文件
//
// 子配置文件
type subConfig struct {
	SubConfig map[string]string `json:"includeConfigs"`
}

type fileConfig struct {
	files []string // 所有待选config 文件
	file  string   // 选中的config文件
	data  []byte   // 文件内容

	val interface{}
	// 配置文件必须是一个json object对象, 可以映射为一个map对象
	m map[string]interface{}
	// 包含的子配置文件
	subconfs subConfig
}

var (
	// 换行
	LF = byte('\n')
	// 注释的行
	commentLine = regexp.MustCompile(`\s*#`)
)

// Typ config类型
func (fc *fileConfig) Typ() int {
	return TypFile
}

// Location config file位置
func (fc *fileConfig) Location() []string {
	return fc.files
}

// 读config文件，填充到制定的结构体中
func (fc *fileConfig) Read(v interface{}) error {
	// 传入参数必须为指针
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("param should be pointer")
	}

	// 读取文件内容
	data, err := fc.readFiles()
	if err != nil {
		return err
	}
	// 映射到结构体中
	err = fc.readConfig(data, rv)
	if err != nil {
		return err
	}
	if err := fc.readSubconfig(data, rv); err != nil {
		return err
	}
	return nil
}

// 写config文件
func (fc *fileConfig) Save(v interface{}) error {
	return nil
}

// 依次读取files，直到有一个成功为止
func (fc *fileConfig) readFiles() ([]byte, error) {
	if len(fc.files) == 0 {
		return nil, fmt.Errorf("No files found")
	}

	errs := []error{}
	for _, fn := range fc.files {
		data, err := fc.readfile(fn)
		if err == nil {
			fc.file = fn
			return data, nil
		}
		errs = append(errs, err)
	}

	return nil, fmt.Errorf("Read all file failed: %v", errs)
}

func (fc *fileConfig) readConfig(data []byte, rv reflect.Value) (err error) {
	//fmt.Println(string(data))
	//printValueFileds(rv, fmt.Sprint(rv.Type()))
	err = json.Unmarshal(data, rv.Interface())
	if err != nil {
		return err
	}
	return nil
}

//
// 读配置文件内容, 移除所有以 # 开始的行
func (fc *fileConfig) readfile(fn string) (data []byte, err error) {
	data, err = ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	return removeComments(data), nil
}

func removeComments(data []byte) []byte {
	var (
		line  []byte
		ndata = []byte{}
		err   error
	)

	rd := bytes.NewBuffer(data)
	for line, err = rd.ReadBytes(LF); err == io.EOF || err == nil; line, err = rd.ReadBytes(LF) {
		if commentLine.Match(line) == false {
			line = bytes.TrimSpace(line)
			if len(line) != 0 {
				// print("    :", string(line))
				ndata = append(ndata, line...)
			}
		}
		if err == io.EOF {
			break
		}
	}

	return ndata
}

// 读子配置
func (fc *fileConfig) readSubconfig(data []byte, rv reflect.Value) error {
	var (
		err      error
		subconfs subConfig
	)

	err = json.Unmarshal(data, &subconfs)
	if err != nil {
		return fmt.Errorf("read subconfig failed: %v", err)
	}

	dir := path.Dir(fc.file)
	rv = indirect(rv)
	//printValueFileds(rv, fmt.Sprint(rv.Type()))
	t := rv.Type()
	for fieldName, confName := range subconfs.SubConfig {
		// 暂不考虑匿名对象的情况
		if _, ok := t.FieldByName(fieldName); ok {
			subv := rv.FieldByName(fieldName)
			subv = indirect(subv)
			//printValueFileds(subv, fieldName)

			if subv.IsValid() == false {
				return fmt.Errorf("reflect.Value has no field %s", fieldName)
			}
			if subv.CanAddr() == false {
				return fmt.Errorf("value field cannot addr")
			}
			subv = subv.Addr()
			data, ierr := fc.readfile(path.Join(dir, confName))
			if ierr != nil {
				fmt.Printf("Read sub config file %v failed: %v\n",
					path.Join(dir, confName), ierr)
				continue
			}
			ierr = fc.readConfig(data, subv)
			if ierr != nil {
				fmt.Printf("get field/subconfig %v/%v failed: %v",
					fieldName, confName, ierr)
			}
		} else {
			// 根据 field name 找不到 field 的情况
			fmt.Printf("Not found field %s", fieldName)
		}
	}

	return nil
}

//
// print value fields
func printValueFileds(v reflect.Value, name string) {
	if v.IsValid() == false {
		fmt.Printf("printValueFileds: v %v is zero value.\n", v)
		return
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		fmt.Printf("printValueFileds: param v should be struct, but %v %v\n", v.Kind(), name)
		return
	}
	fmt.Printf("printValueFileds: %s\n", name)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fmt.Printf("field %d: %v %v\n", i, field.Kind(), field)
	}
}
