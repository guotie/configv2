package configv2

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
)

// 文件循环引用问题
//
type Config map[string]interface{}

var (
	cfg         = make(map[string]interface{})
	KeyNotExist = errors.New("key not exist")
)

// read config file
//
func ReadCfg(fn string) error {
	if fn == "" {
		fn = "./config.json"
	}
	content, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &cfg)

	return err
}

// read config file, return Config
func ReadConfigFile(fn string) (Config, error) {
	var c Config

	content, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(content, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Get
// if the key do not exist, return nil
func Get(key string) interface{} {
	if v, ok := cfg[key]; ok {
		return v
	}

	return nil
}

// GetMust
// if the key do not exist, panic
func GetMust(key string) interface{} {
	v, ok := cfg[key]
	if !ok {
		panic("Not found " + key + " in config.")
	}

	return v
}

// GetString
// Get string type by key
func GetString(key string) string {
	v := Get(key)
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}

	return ""
}

// GetString
// Get string type by key, if not exist, panic
func GetStringMust(key string) string {
	v := Get(key)
	if v == nil {
		panic("Not found key " + key)
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}

	panic("Cannot convert key " + key + " to string.")
	return ""
}

// GetString
// Get string type by key, if key not exist, return the default
func GetStringDefault(key, def string) string {
	s := GetString(key)
	if s == "" {
		return def
	}
	return s
}

// GetInt
func GetInt(key string) (int64, bool) {
	v := Get(key)
	if v == nil {
		return 0, false
	}
	if i, ok := v.(float64); ok {
		return int64(i), true
	}

	return 0, false
}

func GetInt64Default(key string, dv int64) int64 {
	v := Get(key)
	if v == nil {
		return dv
	}
	if i, ok := v.(float64); ok {
		return int64(i)
	}

	return dv
}

func GetIntDefault(key string, dv int) int {
	v := Get(key)
	if v == nil {
		return dv
	}
	if i, ok := v.(float64); ok {
		return int(i)
	}

	return dv
}

func GetFloat(key string) (float64, bool) {
	v := Get(key)
	if v == nil {
		return 0, false
	}
	if i, ok := v.(float64); ok {
		return i, true
	}

	return 0.0, false
}

func GetBoolean(key string) (bool, bool) {
	v := Get(key)
	if v == nil {
		return false, false
	}

	if i, ok := v.(bool); ok {
		return i, true
	}

	return false, false
}

func GetBooleanDefault(key string, dv bool) bool {
	v := Get(key)
	if v == nil {
		return dv
	}
	if i, ok := v.(bool); ok {
		return i
	}

	return dv
}

func GetBytes(key string) ([]byte, bool) {
	v := Get(key)
	if v == nil {
		return []byte(""), false
	}
	if i, ok := v.([]byte); ok {
		return i, true
	}

	return []byte(""), false
}

func GetBytesDefault(key string, dv []byte) []byte {
	v := Get(key)
	if v == nil {
		return dv
	}

	if i, ok := v.([]byte); ok {
		return i
	}

	return dv
}

// Scan
//
func Scan(key string, dest interface{}) (err error) {
	val, ok := cfg[key]
	if !ok {
		return KeyNotExist
	}

	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("param must be pointer and not nil")
	}

	return scan(val, rv)
}

// ScanConfig
func (c Config) ScanConfig(key string, dest interface{}) error {
	val, ok := c[key]
	if !ok {
		return KeyNotExist
	}

	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("param must be pointer and not nil")
	}

	return scan(val, rv)
}
