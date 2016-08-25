package helper

import (
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	CONCURRENT = 100
	COUNTS     = 1000
)

var (
	goods int32 = 0
	lock  sync.Mutex
)

/**
* run sql before
* create database test;
* create table test_t(id int PRIMARY KEY AUTO_INCREMENT, value int);
* create table goods_list(id int PRIMARY KEY AUTO_INCREMENT, u_id int);
 */
func initTest(db *sql.DB) {
	if _, err := db.Exec("TRUNCATE table goods_list"); err != nil {
		panic(err)
	}
	if _, err := db.Exec("update test_t set value=50000  where id=1"); err != nil {
		panic(err)
	}
}
func update(db *sql.DB) (rows int64) {
	tx, err := db.Begin()
	defer tx.Commit()
	stmtUpdt, err := tx.Prepare("UPDATE test_t SET value=value-1 where id = ? and value > 0")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtUpdt.Close()
	stmtInt, err := tx.Prepare("INSERT goods_list VALUES(0, ?)")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtInt.Close()
	if res, err := stmtUpdt.Exec(1); err == nil {
		rows, _ = res.RowsAffected()
		if rows != 0 {
			lock.Lock()
			defer lock.Unlock()
			goods++
			if _, e := stmtInt.Exec(goods); e != nil {
				fmt.Println("err:", e)
				tx.Rollback()
			}
		}
	}
	return
}

func thread_test(over chan<- struct{}) {
	db, err := sql.Open("mysql", "root:admin@/test")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()
	for i := 0; i < COUNTS; i++ {
		if update(db) == 0 {
			break
		}
	}
	over <- struct{}{}
}

func Test_UpdataQps(t *testing.T) {
	db, err := sql.Open("mysql", "root:admin@/test")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	stmtIns, err := db.Prepare("INSERT INTO test_t VALUES( ?, ?)") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates
	stmtIns.Exec(2, 2)

	initTest(db)

	overs := make(chan struct{}, 1)
	start := time.Now()

	for i := 0; i < CONCURRENT; i++ {
		go thread_test(overs)
	}

	for i := 0; i < CONCURRENT; i++ {
		<-overs
	}

	duration := time.Now().Sub(start).Seconds()
	t.Log(duration, "goods:", goods)
	fmt.Println("duration:", duration, "goods:", goods)
	fmt.Println("qps:", CONCURRENT*COUNTS/duration)
}
