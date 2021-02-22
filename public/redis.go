package public

import (
	"github.com/garyburd/redigo/redis"
	"github.com/zhao890420/dc-gateway/common"
)

func RedisConfPipline(pip ...func(c redis.Conn)) error {
	//c, err := lib.RedisConnFactory("default")
	//if err != nil {
	//	return err
	//}
	c := common.DefaultRedisPool.Get()
	defer c.Close()
	for _, f := range pip {
		f(c)
	}
	c.Flush()
	return nil
}

func RedisConfDo(commandName string, args ...interface{}) (interface{}, error) {
	c := common.DefaultRedisPool.Get()
	defer c.Close()
	return c.Do(commandName, args...)
	return nil, nil
}
