package helper

import (
	"github.com/mijia/sweb/log"
	redis "gopkg.in/redis.v4"
)

type lua_script func() string

func Lua_script_calc() string {
	return "if redis.call('GET', KEYS[1]) <= '0' then return '-1' else redis.call('RPUSH', KEYS[2], ARGV[1]) return redis.call('INCRBY',KEYS[1], -1) end"
}

func ClientDealWithGoods(client *redis.Client, script lua_script, args []string, argv int) (val interface{}, err error) {
	if val, err = client.Eval(script(), args, argv).Result(); err != nil {
		log.Error(err, val)
	}
	return
}
