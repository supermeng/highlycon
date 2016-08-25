package main

import (
	"github.com/mijia/sweb/log"
	redis "gopkg.in/redis.v4"
	"net"
	"strconv"
	"sync"
)

const (
	Port = 8080

	RedisAddr = "127.0.0.1:6001"

	MAX_CONN = 1000
)

var (
	value      int32 = 0
	goods            = 0
	goods_lock sync.Mutex
)

var count int32 = 0

type LimitListener struct {
	server_listener net.Listener
	cond            chan struct{}
}

func NewLimitListener() *LimitListener {
	p := &LimitListener{cond: make(chan struct{}, MAX_CONN)}
	hostAndPort := "0.0.0.0:" + strconv.Itoa(Port)
	var err error
	serverAddr, err := net.ResolveTCPAddr("tcp", hostAndPort)
	if err != nil {
		return err
	}
	p.server_listener, err = net.ListenTCP("tcp", serverAddr)
	if err != nil {
		return err
	}
	log.Debug("Listening to: ", p.server_listener.Addr().String())
	return p
}

func (l *limitListener) Accept() (net.Conn, error) {
	<-cond
	return l.Accept()
}

type RedisServer struct {
	listener *net.TCPListener
}

func NewRedisServer() *RedisServer {
	rs := &RedisServer{}
	rs.InitServer()
	return rs
}

func (rs *RedisServer) InitServer() error {
	hostAndPort := "0.0.0.0:" + strconv.Itoa(Port)
	var err error
	serverAddr, err := net.ResolveTCPAddr("tcp", hostAndPort)
	if err != nil {
		return err
	}
	rs.listener, err = net.ListenTCP("tcp", serverAddr)
	if err != nil {
		return err
	}
	return nil
}

func (rs *RedisServer) StartListen() {
	for {
		if conn, err := rs.listener.Accept(); err != nil {
			log.Error(err)
			break
		} else {
			go rs.HandlerConn(conn)
		}
	}
}

func (rs *RedisServer) HandlerConn(conn net.Conn) {
	rs.redisPool.Get()
}

func MemDealWithGoods() int {
	goods_lock.Lock()
	defer goods_lock.Unlock()
	if goods > 0 {
		goods--
		return goods
	} else {
		return -1
	}
}
