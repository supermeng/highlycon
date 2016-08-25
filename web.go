package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strconv"

	"github.com/mijia/sweb/log"
	"github.com/supermeng/highlycon/helper"
	"github.com/supermeng/highlycon/netutil"
	// "golang.org/x/net/netutil"
	redis "gopkg.in/redis.v4"
)

const (
	goodsKey = "goods_count"

	usersKey = "user_lists"

	MAX_CLIENT = 100
	MAX_CONN   = 1000

	COUNTS = 100000

	GOODS = COUNTS / 2

	TIMES = COUNTS / MAX_CLIENT

	RedisAddr = "127.0.0.1:6001"
)

type GoodsServer struct {
	client *redis.Client
}

func NewGoodsServer() *GoodsServer {
	client := redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
		PoolSize: MAX_CLIENT,
	})
	return &GoodsServer{client: client}
}

func (gs *GoodsServer) ConsumeGoods() (interface{}, error) {
	id := rand.Intn(COUNTS)
	return helper.ClientDealWithGoods(gs.client, helper.lua_script_calc, []string{goodsKey, usersKey}, id)
}

func (gs GoodsServer) testHandler(w http.ResponseWriter, r *http.Request) {
	if val, err := gs.ConsumeGoods(); err == nil {
		if iv, ok := val.(int64); ok {

			w.Write([]byte(strconv.FormatInt(iv, 10)))
		} else if sv, ok := val.(string); ok {
			w.Write([]byte(sv))
		}
	}
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	b := "hello world!"
	w.Write([]byte(b))
}

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

func main() {
	initTest()

	gs := NewGoodsServer()

	http.HandleFunc("/", gs.testHandler)

	l, err := net.Listen("tcp", "0.0.0.0:8888")
	if err != nil {
		fmt.Printf("Listen: %v\n", err)
	}
	defer l.Close()

	l = netutil.NewLimitListener(l, MAX_CONN)

	if err := http.Serve(l, nil); err != nil {
		log.Error(err)
	}
	// http.ListenAndServe(":8888", nil)
}
