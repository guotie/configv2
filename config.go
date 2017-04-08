package configv2

//"encoding/json"
//"fmt"
//"io/ioutil"
//"reflect"

const (
	// TypFile json文件类型配置
	TypFile = 0
	// TypEtcd etcd类型配置
	TypEtcd = 1
)

// 文件循环引用问题

// Config Config接口
type Config interface {
	Typ() int
	Location() []string
	Read(v interface{}) error
	Save(v interface{}) error
	Get(key string) (interface{}, bool)
}

// Options Config options, 用于生成Config
type Options struct {
	typ      int
	location []string
	prefix   string
	username string
	password string
}

var _config Config

// Get get key from config
func Get(key string) (interface{}, bool) {
	if _config != nil {
		return _config.Get(key)
	}

	panic("config has not initialized")
}

// NewConfig 根据options生成Config
func NewConfig(o *Options) (c Config) {
	switch o.typ {
	case TypFile:
		c = &fileConfig{
			files: o.location,
		}
		_config = c

	case TypEtcd:
		panic("Typ etcd not implement yet.")
	default:
		panic("invalid option typ")
	}

	return
}

// NewFileConfig 生成TypFile类型的Config
func NewFileConfig(fn string, fns ...string) Config {
	return &fileConfig{
		files: append([]string{fn}, fns...),
	}
}

//
// NewOptions 新建Opitons
//
// Usage: configv2.NewOptions().Typ(int).Location(string)
//
func NewOptions() *Options {
	return &Options{}
}

// Typ 设置options类型
func (op *Options) Typ(t int) *Options {
	op.typ = t
	return op
}

// Location 设置options地址
func (op *Options) Location(loc []string) *Options {
	op.location = loc
	return op
}

// Prefix 设置options prefix
func (op *Options) Prefix(prefix string) *Options {
	op.prefix = prefix
	return op
}

// Location 设置options的username
func (op *Options) Username(un string) *Options {
	op.username = un
	return op
}

// Location 设置options的password
func (op *Options) Password(pwd string) *Options {
	op.password = pwd
	return op
}
