package configv2

import (
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
)

type etcdConfig struct {
	endpoint string
	prefix   string
	username string
	password string

	endpoints      []string
	dialTimeout    time.Duration
	requestTimeout time.Duration
}

// Typ config类型
func (ec *etcdConfig) Typ() int {
	return TypEtcd
}

// Location etcd服务器endpoints
func (ec *etcdConfig) Location() string {
	return ec.endpoint
}

// Read 从etcd中读取配置
func (ec *etcdConfig) Read(v interface{}) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   ec.endpoints,
		DialTimeout: ec.dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	_, err = cli.Put(context.TODO(), "foo", "bar")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), ec.requestTimeout)
	resp, err := cli.Get(ctx, "foo")
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	for _, ev := range resp.Kvs {
		log.Printf("%s : %s\n", ev.Key, ev.Value)
	}
	return nil
}
