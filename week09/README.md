# go进阶训练营第九周-go语言实践-网络编程

## 作业

1. 总结几种 socket 粘包的解包方式：fix length/delimiter based/length field based frame decoder。尝试举例其应用。

   - 什么是粘包？
    
      若干个数据包在接收方接收时，粘成了一个包

   - 怎么解决粘包问题？
    
     - fix length based 每次发送固定长度的数据，并且不超过缓冲区，接收方每次按照固定长度去读取数据
     - delimiter based 发送在数据包加入特殊分隔符，标识数据包开始或结束
     - length field based 在消息数据包头加入长度字段

2. 实现一个从 socket connection 中解码出 goim 协议的解码器。

    实现了goim协议的`decode`和`encode`,运行protocol_test查看

    ![goim 协议包](https://cdn.nlark.com/yuque/0/2021/png/757992/1639878307667-1b25ab75-becb-4817-916f-c2f1df2db127.png?x-oss-process=image%2Fresize%2Cw_1274%2Climit_0) 