package main

// 固定窗口计数法

import (
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"time"
)

// CounterRateLimit 限流之固定窗口计数
func CounterRateLimit(maxCount int, windowSize int64) gin.HandlerFunc {
	counter := 1
	lastTime := time.Now()
	return func(c *gin.Context) {
		now := time.Now()
		if now.Sub(lastTime).Milliseconds() < windowSize {
			if counter < maxCount {
				counter++
				c.JSON(200, gin.H{"message": "success"})
				return
			} else {
				c.AbortWithStatusJSON(400, map[string]string{"message": "请求数量超过限制"})
				return
			}
		}
		lastTime = now
		counter = 0
		return
	}
}

type SlideRateLimitConfig struct {
	windowSize int64     // 窗口大小，单位为毫秒数
	splitNum   uint      // 窗口分割的个数
	startTime  time.Time // 起始时间
	limit      int       //限制请求个数
}

// SlideRateLimit 限流之滑动窗口计数器
func SlideRateLimit(config *SlideRateLimitConfig) gin.HandlerFunc {
	counters := make([]int, config.splitNum)
	splitNum := config.splitNum
	index := 0

	slideWindow := func(windowsNum float64) {
		if windowsNum <= 0 {
			return
		}
		slideNum := int(math.Min(windowsNum, float64(config.splitNum)))
		for i := 0; i < slideNum; i++ {
			index = (index + 1) % int(config.splitNum)
			counters[index] = 0
		}

		addTime := int64(windowsNum) * (config.windowSize / int64(config.splitNum))
		config.startTime = config.startTime.Add(time.Duration(addTime) * time.Millisecond)
	}

	return func(c *gin.Context) {
		now := time.Now()
		// 比较当前时间和起始时间移动了多少个小窗口
		diffTime := float64(now.Sub(config.startTime).Milliseconds() - config.windowSize)
		windowTime := float64(config.windowSize) / float64(config.splitNum)
		windowsNum := math.Max(diffTime, 0) / (windowTime)
		// 滑动窗口，重新设置起始时间
		slideWindow(windowsNum)
		count := 0
		for i := 0; i < int(splitNum); i++ {
			count += counters[i]
		}

		if count > config.limit {
			c.AbortWithStatusJSON(400, map[string]string{"message": "请求数量超过限制"})
			return
		} else {
			counters[index]++
			c.JSON(200, gin.H{"message": "success"})
			return
		}
	}
}

func main() {
	r := gin.Default()
	//r.Use(CounterRateLimit(10, 1000))
	config := SlideRateLimitConfig{
		windowSize: 1000,
		splitNum:   100,
		startTime:  time.Now(),
		limit:      1,
	}

	r.Use(SlideRateLimit(&config))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	log.Fatal(r.Run()) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
