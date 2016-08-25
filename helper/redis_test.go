package helper

import (
	"github.com/mijia/sweb/log"
	redis "gopkg.in/redis.v4"

	"math/rand"
	"sync"
	"testing"
	"time"
)

const (
	goodsKey = "goods_count"

	usersKey = "user_lists"

	MAX_CLIENT = 100

	COUNTS = 100000

	GOODS = COUNTS / 2

	TIMES = COUNTS / MAX_CLIENT

	RedisAddr = "127.0.0.1:6001"
)

var (
	counts int32 = 0
	lock   sync.Mutex
)

func initTest() {
	client := redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer client.Close()
	client.FlushAll()
	client.Set(goodsKey, GOODS, 0)
}

func thread_test(client *redis.Client, over chan<- struct{}) {
	for i := 0; i < TIMES; i++ {
		id := rand.Intn(COUNTS)
		val, _ := ClientDealWithGoods(client, lua_script_calc, []string{goodsKey, usersKey}, id)
		if val == "-1" {
			break
		} else {
			lock.Lock()
			counts++
			lock.Unlock()
		}
	}
	over <- struct{}{}
}

func startTest(overs chan<- struct{}) {
	client := redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
		PoolSize: MAX_CLIENT,
	})
	for i := 0; i < MAX_CLIENT; i++ {
		go thread_test(client, overs)
	}
}

func verifyTest() bool {
	client := redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if val, err := client.Get(goodsKey).Result(); err == nil {
		l, _ := client.LLen(usersKey).Result()
		return val == "0" && l == GOODS
	}
	return false
}

func Test_ConcurrentTest(t *testing.T) {
	overs := make(chan struct{}, MAX_CLIENT)
	initTest()
	start := time.Now()
	startTest(overs)
	for i := 0; i < MAX_CLIENT; i++ {
		<-overs
	}
	log.Info("verified result:", verifyTest())
	costs := time.Now().Sub(start).Seconds()
	log.Info(costs, " qps:", (COUNTS / costs))
	log.Info("consume goods' counts:", counts)
}
