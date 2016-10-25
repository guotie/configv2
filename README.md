# Config Readme

------

This package is used to get config params in json format text file.

------

## ReadCfg

```
func ReadCfg(fn string)
```

read and parse config file, if there has any error, it will panic.

## Get

```
func Get(key string) interface{}
```

Get config param by key, if key does not exists, it return nil.

## GetMust

```
func GetMust(key string) interface{}
```

Just like Get, but if the key does not exist, it will panic.

## GetString, GetStringDefault, GetStringMust

```
func GetString(key string) string
func GetStringMust(key string) string
func GetStringDefault(key, def string) string
```

Get config param by key, return a string. if the key does not exist, GetString return "", GetStringMust will panic,
and GetStringDefault will return def.

## GetInt, GetInt64Default, GetIntDefault, GetFloat, 
```
func GetInt(key string) (int64, bool)
func GetInt64Default(key string, dv int64) int64
func GetIntDefault(key string, dv int) int
func GetFloat(key string) (float64, bool)
```
Just like GetString, but return int or float

## GetBoolean, GetBooleanDefault
```
func GetBoolean(key string) (bool, bool)
func GetBooleanDefault(key string, dv bool) bool
```

Just like GetString, but return bool

## Scan
```
func Scan(key string, dest interface{}) (err error) 
```

Scan can convert to any type, dest MUST be a pointer.

Use it like this:
```
var var1 []int

Scan("key1", &var1)
```

more examples is in config_test.go