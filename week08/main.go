package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"regexp"
	"time"
)

var (
	ctx                     = context.Background()
	ip                      = "localhost"
	port             uint16 = 6379
	redisMemoryRegex        = regexp.MustCompile("used_memory_human:.*?\n")
	batchCount              = 1 * 10000
)

//set 通过pipeline向redis写入指定count的kv数据
func set(redisClient *redis.Client, key string, value string, count int) {
	fmt.Printf("插入key[%s]value:[%s],个数:[%v]....\n", key, value, count)

	pipe := redisClient.Pipeline()
	for i := 0; i < count; i++ {
		newKey := fmt.Sprintf("%s:%v", key, i)
		pipe.Set(ctx, newKey, value, -1)
		if i%batchCount == 0 {
			execPipe(pipe)
		}
	}
	execPipe(pipe)
}

func execPipe(pipe redis.Pipeliner) {
	_, err := pipe.Exec(ctx)
	if err != nil {
		panic(err)
	}
}

//getValue 生成指定字节的value
func getValue(dataSize int) string {
	bytes := make([]byte, dataSize)
	for i, _ := range bytes {
		bytes[i] = 'a'
	}
	return string(bytes)
}

//analysis 分析redis的内存,并将内存信息输出到控制台
func analysis(redisClient *redis.Client, f func()) {
	redisClient.FlushDB(ctx)
	f()
	val, err := redisClient.Info(ctx, "memory").Result()
	fmt.Printf("%q\n", redisMemoryRegex.Find([]byte(val)))
	if err != nil {
		panic(err)
	}
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%v", ip, port),
		WriteTimeout: time.Second * 10,
	})
	analysis(client, func() {
		set(client, "len_10_10000", getValue(10), 10000)
	})

	analysis(client, func() {
		set(client, "len_10_500000", getValue(10), 500000)
	})

	analysis(client, func() {
		set(client, "len_20_10000", getValue(20), 10000)
	})

	analysis(client, func() {
		set(client, "len_20_500000", getValue(20), 500000)
	})

	analysis(client, func() {
		set(client, "len_200_10000", getValue(200), 10000)
	})

	analysis(client, func() {
		set(client, "len_200_500000", getValue(200), 500000)
	})

	analysis(client, func() {
		set(client, "len_1024_10000", getValue(1024), 10000)
	})

	analysis(client, func() {
		set(client, "len_1024_500000", getValue(1024), 500000)
	})

	analysis(client, func() {
		set(client, "len_5120_10000", getValue(5120), 10000)
	})

	analysis(client, func() {
		set(client, "len_5120_50000", getValue(5120), 50000)
	})
}
