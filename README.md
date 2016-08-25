# highlycon

[![MIT license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://opensource.org/licenses/MIT)

## Architecture
![Architecture](Architecture.png)

## redis qps test
1. 开启一个redis节点 port 6001
2. cd helper && go test

![redis test](redis_test.png)

## mysql transaction qps test

1. cd mysql
2. create table
3. go test

![mysql test](mysql_test.png)


## web qps test(use redis as cache)

1. go run web
2. cd test
3. go run contest.go

![web test](web_test.png)
