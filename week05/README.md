## 作业

1. 参考 Hystrix 实现一个滑动窗口计数器。

### 回答

总结滑动窗口计数法的思路是：

1. 将时间划分为细粒度的区间，每个区间维持一个计数器，每进入一个请求则将计数器加一。
2. 多个区间组成一个时间窗口，每流逝一个区间时间后，则抛弃最老的一个区间，纳入新区间。
3. 若当前窗口的区间计数器总和超过设定的限制数量，则本窗口内的后续请求都被丢弃。


### 使用

`CounterRateLimit`方法实现了固定窗口限流算法，`SlideRateLimit`实现了滑动窗口限流算法

通过传入配置来控制限流参数

```
config:=SlideRateLimitConfig{
    windowSize: 1000,
    splitNum: 100,
    startTime:  time.Now(),
    limit: 1,
}
r := gin.Default()
r.Use(SlideRateLimit(&config))
```

1. 启动gin应用

```shell
go run main.go
```

2. 单独写一个请求gin接口的demo go应用

```go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main()  {
	for{
		response,err:= http.Get("http://localhost:8080/ping")
		if err != nil {
			log.Fatal(err)
		}else {
			body, _ := io.ReadAll(response.Body)
			fmt.Printf("status: %d,message: %s\n",response.StatusCode,body )
		}
		time.Sleep(time.Millisecond*100)
	}
}
```

测试结果如下：

```shell
status: 400,message: {"message":"请求数量超过限制"}
status: 400,message: {"message":"请求数量超过限制"}
status: 400,message: {"message":"请求数量超过限制"}
status: 200,message: {"message":"success"}{"message":"pong"}
status: 200,message: {"message":"success"}{"message":"pong"}
status: 400,message: {"message":"请求数量超过限制"}
status: 400,message: {"message":"请求数量超过限制"}
```

## 笔记
https://www.yuque.com/docs/share/1dad9950-e5e4-4859-8069-3b33086c1808?# 《可用性设计》