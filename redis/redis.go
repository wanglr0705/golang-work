package redis

import "github.com/gomodule/redigo/redis"

var (
	RDB      redis.Conn
	network  string = "tcp"
	address  string = "127.0.0.1:6379"
	password string = "123456"
)

func Redis() error {
	var err error
	RDB, err = redis.Dial(network, address, redis.DialPassword(password))
	return err
}
