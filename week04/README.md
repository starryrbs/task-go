题目：1. 按照自己的构想，写一个项目满足基本的目录结构和工程，
代码需要包含对数据层、业务层、API 注册，以及 main 函数对于服务的注册和启动，
信号处理，使用 Wire 构建依赖。可以使用自己熟悉的框架。

答：

我做了一个名叫"big"（随便取得）的app，主要做了一下用户的curd的操作

- api: 简单的定义数据返回体，没有定义API proto文件

- cmd: 负责服务的启动，使用wire构建依赖，初始化biz，service，data，conf等对象

- configs/config.yaml: 定义项目的配置文件

- internal: 

    - biz: 定义了UserRepo接口（包含增删改查方法）

    - conf: 定义了配置的格式，包含了Database和Server配置

    - data: 业务数据访问，使用ent建模以及建立对sqllite的访问，并实现了biz的UserRepo的接口

    - services: 实现了api定义的访问层，简单了处理了DTO=>DO的对象转换，以及简单业务逻辑校验

## 笔记链接

https://www.yuque.com/docs/share/268542be-9b7a-442c-95e6-3d8d535cd979?# 《Go工程化实践》