1. 使用 redis benchmark 工具, 测试 10 20 50 100 200 1k 5k 字节 value 大小，redis get set 性能。

- docker 启动redis

```go
docker run -p 6379:6379 --name redis -d redis
```

- benchmark

```go
  redis-benchmark -t set,get -n 100000 -d 10
```
分别使用10 20 50 100 200 1k 5k 字节 value 大小进行get，set测试
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1639289434612-bc781429-cf21-4370-b46b-6af49179b02c.png#clientId=u9bed3ecf-cc72-4&crop=0&crop=0&crop=1&crop=1&from=paste&height=390&id=u02e13d21&margin=%5Bobject%20Object%5D&name=image.png&originHeight=780&originWidth=1102&originalType=binary&ratio=1&rotation=0&showTitle=false&size=519173&status=done&style=none&taskId=uc47bc284-2352-42ea-aad9-1cf4f5acf4f&title=&width=551)


2. 写入一定量的 kv 数据, 根据数据大小 1w-50w 自己评估, 结合写入前后的 info memory 信息 , 分析上述不同 value 大小下，平均每个 key 的占用内存空间。

[代码实现](./main.go)

![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1639295593559-7170ace6-1d4e-4c4e-b15b-dc1e83c553d6.png#clientId=u39a2db41-69c9-4&crop=0&crop=0&crop=1&crop=1&from=paste&height=488&id=Z0Fhs&margin=%5Bobject%20Object%5D&name=image.png&originHeight=976&originWidth=1362&originalType=binary&ratio=1&rotation=0&showTitle=false&size=141848&status=done&style=none&taskId=u0ab4b641-0342-4957-8881-93b5b0ed901&title=&width=681)
