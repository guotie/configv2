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
	Location() string
	Get(key string) (interface{}, bool)
}

// Options Config options, 用于生成Config
type Options struct {
	typ      int
	location string
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
func (op *Options) Location(loc string) *Options {
	op.location = loc
	return op
}
