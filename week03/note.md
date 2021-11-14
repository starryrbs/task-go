## Goroutine
### 进程和线程
进程是对运行时程序的封装，是**系统进行资源调度和分配的的基本单位，实现了操作系统的并发**；
线程是进程的子任务，**是CPU调度和分派的基本单位**，**用于保证程序的实时性，实现进程内部的并发；线程是操作系统可识别的最小执行和调度单位**。
### Goroutines and Parallelism
Go 语言层面支持的 go 关键字，可以快速的让一个函数创建为 goroutine，我们可以认为 main 函数就是作为 goroutine 执行的。操作系统调度线程在可用处理器上运行，Go运行时调度 goroutine 在绑定到单个操作系统线程的逻辑处理器中运行（P）。即使使用这个单一的逻辑处理器和操作系统线程，也可以调度数十万 goroutine 以惊人的效率和性能并发运行。
> Concurrency is not Parallelism.

并发不是并行。并行是指两个或多个线程同时在不同的处理器执行代码。如果将运行时配置为使用多个逻辑处理器，则调度程序将在这些逻辑处理器之间分配 goroutine，这将导致 goroutine 在不同的操作系统线程上运行。但是，要获得真正的并行性，您需要在具有多个物理处理器的计算机上运行程序。否则，goroutine 将针对单个物理处理器并发运行，即使 Go 运行时使用多个逻辑处理器。
### Keep yourself busy or do the work yourself
空的 select 语句将永远阻塞。
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1635861413478-758f35f2-32f9-43f2-be6a-1f5532b0811e.png#clientId=ude34c636-e662-4&from=paste&height=720&id=uba9662a8&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1440&originWidth=2213&originalType=binary&ratio=1&size=1170023&status=done&style=none&taskId=u534c7e3c-8a98-4e6a-bbed-5ce8ff1e6ff&width=1106.5)
如果你的 goroutine 在从另一个 goroutine 获得结果之前无法取得进展，那么通常情况下，你自己去做这项工作比委托它( go func() )更简单。（简单点说就是你必须等待这个goroutine结束才能执行后续流程，就不要用goroutine了）
这通常消除了将结果从 goroutine 返回到其启动器所需的大量状态跟踪和 chan 操作。
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1635861600074-ff0e208c-643d-449e-8090-0ed78cfea382.png#clientId=ude34c636-e662-4&from=paste&height=661&id=ua1617f1f&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1322&originWidth=2560&originalType=binary&ratio=1&size=1521195&status=done&style=none&taskId=u7fea4140-3b34-46f4-9a41-ca723f28500&width=1280)
### Leave concurrency to the caller
#### 这两个 API 有什么区别？
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1635861657622-34879fa2-58c7-445b-864d-1fdd64ee950b.png#clientId=ude34c636-e662-4&from=paste&height=393&id=ue387aa19&margin=%5Bobject%20Object%5D&name=image.png&originHeight=785&originWidth=2560&originalType=binary&ratio=1&size=967348&status=done&style=none&taskId=ud7712ccd-8396-46bc-a486-cc45257a933&width=1280)

- 将目录读取到一个 slice 中，然后返回整个切片，或者如果出现错误，则返回错误。这是同步调用的，ListDirectory 的调用方会阻塞，直到读取所有目录条目。根据目录的大小，这可能需要很长时间，并且可能会分配大量内存来构建目录条目名称的 slice。

问题：1. 目录很大会占用很长时间

- ListDirectory 返回一个 chan string，将通过该 chan 传递目录。当通道关闭时，这表示不再有目录。由于在 ListDirectory 返回后发生通道的填充，ListDirectory 可能内部启动 goroutine 来填充通道。
> ListDirectory chan 版本还有两个问题

- 通过使用一个关闭的通道作为不再需要处理的项目的信号，ListDirectory **无法告诉调用者通过通道返回的项目集不完整**，因为中途**遇到了错误**。调用方无法区分空目录与完全从目录读取的错误之间的区别。这两种方法都会导致从 ListDirectory 返回的通道会立即关闭。
- **调用者必须持续从通道读取，直到它关闭**，因为这是调用者知道填充 chan 的 goroutine 已经停止的唯一方法。这对 ListDirectory 的使用是一个严重的限制，调用者必须花时间从通道读取数据，即使它可能已经收到了它想要的答案。对于大中型目录，它可能在内存使用方面更为高效，但这种方法并不比原始的基于 slice 的方法快。
#### 解决办法
通过传入一个func让调用者决定何时取消并发
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1635861934646-ff95d6fd-45f7-4a20-9697-8786cb964827.png#clientId=ude34c636-e662-4&from=paste&height=103&id=u57305c5a&margin=%5Bobject%20Object%5D&name=image.png&originHeight=205&originWidth=2560&originalType=binary&ratio=1&size=274017&status=done&style=none&taskId=u55ad5550-277f-48ec-a2cc-c953bc19800&width=1280)
filepath.WalkDir 也是类似的模型，如果函数启动 goroutine，则必须向调用方提供显式停止该goroutine 的方法。通常，将异步执行函数的决定权交给该函数的调用方通常更容易。
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1635862149535-2d9ee7b7-455e-419f-b96d-e19093f0e694.png#clientId=ude34c636-e662-4&from=paste&height=148&id=u8e885167&margin=%5Bobject%20Object%5D&name=image.png&originHeight=260&originWidth=1308&originalType=binary&ratio=1&size=39697&status=done&style=none&taskId=u0eb5613c-1753-43a6-b296-7e29430e8d8&width=746)
### Never start a goroutine without knowing when it will stop
#### 例子一
在这个例子中，goroutine 可以在 code review 快速识别出来。不幸的是，生产代码中的 goroutine 泄漏通常更难找到。我无法说明 goroutine 泄漏可能发生的所有可能方式，您可能会遇到：
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1635862317162-d409474a-9e31-4632-8722-2d6fbb74c6b0.png#clientId=ude34c636-e662-4&from=paste&height=613&id=u71d08c3e&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1226&originWidth=2560&originalType=binary&ratio=1&size=1062664&status=done&style=none&taskId=u6b2f59c9-6e97-4867-a046-df590a735c6&width=1280)
ch这个chan将不会关闭，会导致内存泄漏
#### 例子二
search 函数是一个模拟实现，用于模拟长时间运行的操作，如数据库查询或 rpc 调用。在本例中，硬编码为200ms。定义了一个名为 process 的函数，接受字符串参数，传递给 search。对于某些应用程序，顺序调用产生的延迟可能是不可接受的。
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1635862404664-52f63f4b-693e-4a32-b52c-411422730b55.png#clientId=ude34c636-e662-4&from=paste&height=191&id=dKpnS&margin=%5Bobject%20Object%5D&name=image.png&originHeight=705&originWidth=2560&originalType=binary&ratio=1&size=926140&status=done&style=none&taskId=u04ece040-43ac-49fa-b388-390ec6080b4&width=695)
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1635862419946-6d0780eb-2c44-4cbb-b235-8cbb866cf2d7.png#clientId=ude34c636-e662-4&from=paste&height=293&id=f4j7P&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1159&originWidth=2560&originalType=binary&ratio=1&size=1325894&status=done&style=none&taskId=ub729bb4d-9742-4d60-85a9-b02f1e58bef&width=647)
同步较慢，这里改成异步
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1635862390078-d3373bd3-0981-41ea-b952-ff3d7a6d0ba7.png#clientId=ude34c636-e662-4&from=paste&height=347&id=ua365d9fc&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1440&originWidth=2002&originalType=binary&ratio=1&size=1231555&status=done&style=none&taskId=u76c8c415-f4de-47c9-a071-8cb9f1218a2&width=482)
上面的代码也会导致内存泄漏，当ctx.Done()先返回，ch这个chan会一直阻塞，这将会导致内存泄漏，解决办法就是将ch设置buffer为1的chan，让goroutine可以结束
#### Any time you start a Goroutine you must ask yourself

1. Goroutine什么时候结束？
1. 有没有办法结束Goroutine？
> 示例

这个简单的应用程序在两个不同的端口上提供 http 流量，端口8080用于应用程序流量，端口8001用于访问 /debug/pprof 端点。
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1635863038555-cad7a124-3bc0-4173-ac4e-39267192bb7a.png#clientId=ude34c636-e662-4&from=paste&height=580&id=uc87e6306&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1159&originWidth=2560&originalType=binary&ratio=1&size=1610239&status=done&style=none&taskId=u78409933-dab5-4551-98fd-e633a48a4ef&width=1280)
通过将 serveApp 和 serveDebug 处理程序分解为各自的函数，我们将它们与 main.main 解耦，我们还遵循了上面的建议，并确保 serveApp 和 serveDebug 将它们的并发性留给调用者。如果 serveApp 返回，则 main.main 将返回导致程序关闭，只能靠类似 supervisor 进程管理来重新启动。
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1635863229391-b57ad6e2-09a6-465e-9fe6-7279c35353a2.png#clientId=ude34c636-e662-4&from=paste&height=628&id=u6caad90f&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1255&originWidth=2560&originalType=binary&ratio=1&size=1449594&status=done&style=none&taskId=uaeacb277-cf9c-49aa-a7c0-ebc94ac0567&width=1280)
然而，serveDebug 是在一个单独的 goroutine 中运行的，如果它返回，那么所在的 goroutine 将退出，而程序的其余部分继续运行。由于 /debug 处理程序很久以前就停止工作了，所以其他同学会很不高兴地发现他们无法在需要时从您的应用程序中获取统计信息。
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1635863287675-f5d584da-f3b4-424f-9467-752f2f6acf70.png#clientId=ude34c636-e662-4&from=paste&height=713&id=uf63bf4b5&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1426&originWidth=2560&originalType=binary&ratio=1&size=1240985&status=done&style=none&taskId=uba1f2a1c-2069-4248-a754-1ee933fe4c7&width=1280)
ListenAndServer 返回 nil error，最终 main.main 无法退出。log.Fatal 调用了 os.Exit，会无条件终止程序；defers 不会被调用到。
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1635863340122-7f9a0026-2d14-4b01-a055-614eb55658dd.png#clientId=ude34c636-e662-4&from=paste&height=713&id=uc953553a&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1426&originWidth=2560&originalType=binary&ratio=1&size=1240985&status=done&style=none&taskId=ub8a070b3-a653-4d24-9c1a-df1e46cfcca&width=1280)
> log.Fatal只应该用在init函数或者main函数中，否则会导致一些未知的错误

继续改进
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1635863406390-3aae479e-1e3b-42d6-9be4-81f36acf3b95.png#clientId=ude34c636-e662-4&from=paste&height=501&id=ua985a780&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1001&originWidth=2560&originalType=binary&ratio=1&size=787101&status=done&style=none&taskId=ub05801f5-0e30-4cc3-a14d-d77d2e03e3c&width=1280)
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1635863426668-4c49264c-1757-4891-810c-8b0825e956bc.png#clientId=ude34c636-e662-4&from=paste&height=551&id=u5fc27089&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1624&originWidth=1440&originalType=binary&ratio=1&size=917812&status=done&style=none&taskId=u6cbd0190-c561-46fc-ba2a-5e1be1b6392&width=489)
定义done和stop两个chan，done chan接收到值后，将关闭stop chan，serve函数中goroutine则可以执行shutdown。
详细参考：[https://github.com/da440dil/go-workgroup](https://github.com/da440dil/go-workgroup)
### Application Lifecycle
对于应用的服务的管理，一般会抽象一个 application lifecycle 的管理，方便服务的启动/停止等。我们 go-kratos kit 库也按照类似的思路做了应用的生命周期托管。

- 应用的信息
- 服务的 start/stop
- 信号处理
- 服务注册

kit 的使用者可以非常方便的对整个 application 级别的资源进行托管，kratos 中使用了 errgroup + functional options 的方式进行了设计。
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1635863859837-ef69429e-4ff0-4657-beaa-fccf08dad9ac.png#clientId=ude34c636-e662-4&from=paste&height=651&id=ubdee4344&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1472&originWidth=1440&originalType=binary&ratio=1&size=293537&status=done&style=none&taskId=u4bab4937-7366-43b4-91f6-b30070b69cb&width=637)
#### 案例：埋点事件上报
无法保证创建的 goroutine 生命周期管理，会导致最常见的问题，就是在服务关闭时候，有一些事件丢失。
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1635863983178-2616e15b-0591-4812-8c75-d599032c4491.png#clientId=ude34c636-e662-4&from=paste&height=390&id=ucb9ff9a7&margin=%5Bobject%20Object%5D&name=image.png&originHeight=779&originWidth=2560&originalType=binary&ratio=1&size=817123&status=done&style=none&taskId=u4aa107cf-ab9d-4945-9f03-7c5651f3f33&width=1280)
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1635863991282-771778cf-792e-4b1c-847e-9c91cddf3c15.png#clientId=ude34c636-e662-4&from=paste&height=720&id=uc89fc03a&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1440&originWidth=2254&originalType=binary&ratio=1&size=1118885&status=done&style=none&taskId=u0f159a2f-a5ab-4c58-9ff9-940c687a1d6&width=1127)
使用 sync.WaitGroup来追踪每一个创建的 goroutine。
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1635864050568-78ae6f9d-2f38-44e6-b499-aee8a3151e0e.png#clientId=ude34c636-e662-4&from=paste&height=607&id=u1fc2b212&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1213&originWidth=2560&originalType=binary&ratio=1&size=1568056&status=done&style=none&taskId=u18466532-cd18-40b6-9a27-921dccf01a7&width=1280)
将 wg.Wait() 操作托管到其他 goroutine，owner goroutine 使用 context 处理超时。
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1635864105470-bb46375b-4436-46d4-865c-08c0b4a6bdc7.png#clientId=ude34c636-e662-4&from=paste&height=287&id=u1e61e420&margin=%5Bobject%20Object%5D&name=image.png&originHeight=573&originWidth=2560&originalType=binary&ratio=1&size=802513&status=done&style=none&taskId=u61658066-3644-45a6-8443-542b3746c92&width=1280)
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1635864112529-dc493421-49c1-4578-b47b-3b53143977a5.png#clientId=ude34c636-e662-4&from=paste&height=720&id=ucf3961f8&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1440&originWidth=1799&originalType=binary&ratio=1&size=1327469&status=done&style=none&taskId=u17789c0a-b90c-4e53-9557-0a623682fb4&width=899.5)
这种方式大量创建goroutine，代价高
可以参考下面这种方式：
```
package main

import (
	"context"
	"time"
)

type Tracker struct {
	ch   chan string
	stop chan struct{}
}

func NewTracker() *Tracker {
	return &Tracker{
		ch: make(chan string, 10),
	}
}

func (t *Tracker) Event(ctx context.Context, data string) error {
	select {
	case t.ch <- data:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *Tracker) Run() {
	for data := range t.ch {
		time.Sleep(time.Second)
		println(data)
	}
	t.stop <- struct{}{}
}

func (t *Tracker) Shutdown(ctx context.Context) {
	close(t.ch)
	select {
	case <-t.stop:
	case <-ctx.Done():
	}
}

func main() {
	tr := NewTracker()
	go tr.Run()
	_ = tr.Event(context.Background(), "test1")
	_ = tr.Event(context.Background(), "test2")
	_ = tr.Event(context.Background(), "test3")
	time.Sleep(3 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()
	tr.Shutdown(ctx)
}

```


## Memory Model
[https://golang.org/ref/mem](https://golang.org/ref/mem)
如何保证在一个 goroutine 中看到在另一个 goroutine 修改的变量的值，如果程序中修改数据时有其他 goroutine 同时读取，那么必须将读取串行化。为了串行化访问，请使用 channel 或其他同步原语，例如 sync 和 sync/atomic 来保护数据。
### Happens-Before
在一个 goroutine 中，读和写一定是按照程序中的顺序执行的。即编译器和处理器只有在不会改变这个 goroutine 的行为时才可能修改读和写的执行顺序。由于**指令重排**，不同的 goroutine 可能会看到不同的执行顺序。例如，一个goroutine 执行 a = 1;b = 2;，另一个 goroutine 可能看到 b 在 a 之前更新。
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636028092146-32dafd61-dce6-4568-b290-47d2cface544.png#clientId=u646af123-0f4f-4&from=paste&height=353&id=u3a49d8e9&margin=%5Bobject%20Object%5D&name=image.png&originHeight=706&originWidth=2560&originalType=binary&ratio=1&size=272307&status=done&style=none&taskId=u94cc6a4d-01e3-416a-9a78-dbe0e3aa46d&width=1280)
### Memory Reordering
用户写下的代码，先要编译成汇编代码，也就是各种指令，包括读写内存的指令。CPU 的设计者们，为了榨干 CPU 的性能，无所不用其极，各种手段都用上了，你可能听过不少，像流水线、分支预测等等。其中，为了提高读写内存的效率，会对读写指令进行重新排列，这就是所谓的内存重排，英文为 MemoryReordering。
> 这一部分说的是 CPU 重排，其实还有编译器重排。比如:

![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636028159109-61529ad2-807e-4e1c-bff1-2bdc4173cb22.png#clientId=u646af123-0f4f-4&from=paste&height=224&id=u102da3b8&margin=%5Bobject%20Object%5D&name=image.png&originHeight=448&originWidth=2470&originalType=binary&ratio=1&size=103873&status=done&style=none&taskId=ue3da7fda-fded-4df7-97c7-62e24af8d88&width=1235)
但是，如果这时有另外一个线程同时干了这么一件事：
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636028191391-6678d897-80da-42d9-9bd1-93960293b310.png#clientId=u646af123-0f4f-4&from=paste&height=74&id=u3f8b7790&margin=%5Bobject%20Object%5D&name=image.png&originHeight=635&originWidth=2560&originalType=binary&ratio=1&size=138655&status=done&style=none&taskId=uf15fa8f5-75c0-44da-88bc-60a7e175c14&width=297)
在多核心场景下，没有办法轻易地判断两段程序是“等价”的。
现代 CPU 为了“抚平” 内核、内存、硬盘之间的速度差异，搞出了各种策略，例如三级缓存等。为了让 (2) 不必等待 (1) 的执行“效果”可见之后才能执行，我们可以把 (1) 的效果保存到 store buffer：
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636028260870-7a35af52-1e8b-4c24-8432-fddfdd335eb0.png#clientId=u646af123-0f4f-4&from=paste&height=449&id=ue38e4c94&margin=%5Bobject%20Object%5D&name=image.png&originHeight=897&originWidth=2560&originalType=binary&ratio=1&size=1022559&status=done&style=none&taskId=u0ce8d6ae-15a8-49f6-8d4f-2ab3492a937&width=1280)
> 上面的例子会有四种情况，即 1，1、1，0、0，1、0，0

先执行 (1) 和 (3)，将他们直接写入 store buffer，接着执行 (2) 和 (4)。“奇迹”要发生了：(2) 看了下 store buffer，并没有发现有 B 的值，于是从 Memory 读出了 0，(4) 同样从 Memory 读出了 0。最后，打印出了 00。
​

因此，对于多线程的程序，所有的 CPU 都会提供“锁”支持，称之为 barrier，或者 fence。它要求：barrier 指令要求所有对内存的操作都必须要“扩散”到 memory 之后才能继续执行其他对 memory 的操作。因此，我们可以用高级点的 atomic compare-and-swap，或者直接用更高级的锁，通常是标准库提供。
​

为了说明读和写的必要条件，我们定义了先行发生（Happens Before）。如果事件 e1 发生在 e2 前，我们可以说 e2 发生在 e1 后。如果 e1不发生在 e2 前也不发生在 e2 后，我们就说 e1 和 e2 是并发的。（简单点说就是e1和e2的发生的顺序是不固定的）
在单一的独立的 goroutine 中先行发生的顺序即是程序中表达的顺序。
当下面条件满足时，对变量 v 的读操作 r 是被允许看到对 v 的写操作 w 的：

- r 不先行发生于 w
- 在 w 后 r 前没有对 v 的其他写操作

为了保证对变量 v 的读操作 r 看到对 v 的写操作 w，要确保 w 是 r 允许看到的唯一写操作。即当下面条件满足时，r 被保证看到 w：

- w 先行发生于 r
- 其他对共享变量 v 的写操作要么在 w 前，要么在 r 后。

这一对条件比前面的条件更严格，需要没有其他写操作与 w 或 r 并发发生。
## Package sync
### 共享内存通信
传统的线程之间通信是通过共享内存实现的。
Go 没有显式地使用锁来协调对共享数据的访问，而是鼓励使用 chan 在 goroutine 之间传递对数据的引用。这种方法确保在给定的时间只有一个 goroutine 可以访问数据。
### Detecting Race Conditions With Go
data race 是两个或多个 goroutine 访问同一个资源
检测方式：

- go build -race
- go test -race

示例：
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636280301818-92d1f49a-7ac2-4837-9d9a-cd44ea6d4126.png#clientId=u646af123-0f4f-4&from=paste&height=616&id=u6e0b7e7e&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1232&originWidth=2560&originalType=binary&ratio=1&size=1024598&status=done&style=none&taskId=u9f363b30-8c9e-41c8-b61a-755d5bd6d75&width=1280)
全局计数器变量的值为 2 或者 4。（正确应该是2，但是会偶尔会输出4）
试图通过 i++ 方式来解决原子赋值的问题，但是我们通过查看底层汇编:
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636280530906-cf7d70eb-9039-4c4d-a196-04167e58aa51.png#clientId=u646af123-0f4f-4&from=paste&height=131&id=u2a17b2f0&margin=%5Bobject%20Object%5D&name=image.png&originHeight=262&originWidth=2560&originalType=binary&ratio=1&size=169761&status=done&style=none&taskId=u9245bc27-b8d9-41bc-98d5-b40e54cf2f0&width=1280)
这几行汇编指令在CPU上执行时会有上下文切换，不能保证安全的数据访问。
**什么是Single Machine World？**
比如现在我的操作系统是**Ubuntu20.04 TLS,**是64位的操作系统，这就意味着我的机器字是8Byte,也就是说我单次能处理的最大的数据容量是**8Byte**，我刚刚好可以利用这点来进行原子赋值操作。虽然**Single Machine Word**的特性确实可以提供原子赋值的能力，不过我始终觉得这是一个风险很大的操作，比如说我在64bit的操作系统对一个16byte长的数据进行操作，比如说interface，在存在并发的场景下，无法保证是前一半的8Byte先写入还是后一半的8Byte被写入，所以这就不能保证原子写入了，所以说**Single Machine World**存在风险，还请谨慎使用。
下面是一段阐述**Single Machine World**存在风险这个观点的佐证代码
```
package main

import "fmt"

type IceCreamMaker interface {
	Hello()
}

type Ben struct {
	name string
}

func (b *Ben) Hello() {
	fmt.Printf("Ben says,\"Hello my name is %s\"\n", b.name)
}

type Jerry struct {
	name string
}

func (j *Jerry) Hello() {
	fmt.Printf("Jerry says,\"Hello my name is %s\"\n", j.name)
}

func main() {
	var ben = &Ben{name: "Ben"}
	var jerry = &Jerry{"Jerry"}
	var maker IceCreamMaker = ben
	var loop0, loop1 func()
	loop0 = func() {
		maker = ben
		go loop1()
	}
	loop1 = func() {
		maker = jerry
		go loop0()
	}

	go loop0()

	for {
		maker.Hello()
	}

}
```
上述的代码的逻辑比较好理解，我定义了一个叫IceCreamMaker的Interface，这个Interface只有一个叫Hello的方法，然后还定义了俩Struct，一个叫Ben，另一个叫Jerry,这两个对象的除了名字不一样，内容是一样的，都只有一个叫name的字段，这两个Struct都有一个方法叫Hello,main函数有一个叫loop0和一个叫loop1的函数，通过一个for的死循环在不停的相互调用，互相打印，我们期待得到的结果一直是
```
Jery says,"Hello my name is Jerry"
Ben says,"Hello my name is Ben"
```
可以我在运行的过程中竟然出现了
```
Jerry says,"Hello my name is Bem"
Ben says,"Hello my name is Jerry"
```
出现上述的结果，用专业的术语来描述就是出现了**data race**我们大概都可以猜到**interface**的大小绝对不止**8byte**,至少也是2个**8byte**,我们可以看下interface的底层结构
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636282706391-6be2baaf-43f0-4de7-8620-8b8d1f3b9809.png#clientId=u646af123-0f4f-4&from=paste&height=90&id=u210369a8&margin=%5Bobject%20Object%5D&name=image.png&originHeight=180&originWidth=1172&originalType=binary&ratio=1&size=23235&status=done&style=none&taskId=ub43a68d5-3551-4e82-9e79-9ce7a8cc999&width=586)
**Interface**由两个字段组成，分别是Type和Data，这两个字段都是uintptr,是用于指针运算的整数类型指针，指针是8byte的，interface有两个指针对象，所以就是16byte。
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636282726117-14479aed-80c6-4d3c-b76f-0bbdd5be6ada.png#clientId=u646af123-0f4f-4&from=paste&height=176&id=ue79b6a8c&margin=%5Bobject%20Object%5D&name=image.png&originHeight=352&originWidth=1018&originalType=binary&ratio=1&size=46850&status=done&style=none&taskId=ub4f43b8d-c4d1-487e-9ffd-63840736a93&width=509)
Type 指向实现了接口的 struct，Data 指向了实际的值。Data 作为通过 interface 中任何方法调用的接收方传递。
对于语句 var maker IceCreamMaker=ben，编译器将生成执行以下操作的代码。
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636282745929-b5244547-6092-450f-a79e-d87f8bc05410.png#clientId=u646af123-0f4f-4&from=paste&height=321&id=u6d47039c&margin=%5Bobject%20Object%5D&name=image.png&originHeight=642&originWidth=1126&originalType=binary&ratio=1&size=64047&status=done&style=none&taskId=uc398b1a8-12be-403d-b22a-7fb1f7cb169&width=563)
当 loop1() 执行 maker=jerry 语句时，必须更新接口值的两个字段。
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636282768786-859450b7-fb1a-4443-8e95-dd6caab8ff3d.png#clientId=u646af123-0f4f-4&from=paste&height=313&id=uc5cad2bd&margin=%5Bobject%20Object%5D&name=image.png&originHeight=626&originWidth=1170&originalType=binary&ratio=1&size=79910&status=done&style=none&taskId=u2214b485-ed01-43c8-bcff-d795bbbd45d&width=585)
Go memory model 提到过: 表示写入单个 machine word 将是原子的，但 interface 内部是是两个 machine word 的值。另一个goroutine 可能在更改接口值时观察到它的内容，也就是说go routine看到的是上一个data的值。
在这个例子中，Ben 和 Jerry 内存结构布局是相同的，因此它们在某种意义上是兼容的。想象一下，如果他们有不同的内存布局会发生什么混乱？
```
package main

import "fmt"

type IceCreamMaker interface {
	Hello()
}

type Ben struct {
	Id int // 这里加了一个 字段叫id
	name string
}

func (b *Ben) Hello() {
	fmt.Printf("Ben says,\"Hello my name is %s and id is %d\"\n", b.name,b.Id) //这里也要答应id
}

type Jerry struct {	
	name string
}

func (j *Jerry) Hello() {
	fmt.Printf("Jerry says,\"Hello my name is %s\"\n", j.name)
}

func main() {
	var ben = &Ben{name: "Ben",Id: 1} // id赋值
	var jerry = &Jerry{"Jerry"}
	var maker IceCreamMaker = ben
	var loop0, loop1 func()
	loop0 = func() {
		maker = ben
		go loop1()
	}
	loop1 = func() {
		maker = jerry
		go loop0()
	}

	go loop0()

	for {
		maker.Hello()
	}

}
```
对代码进行了一点小修改，Ben的struct加了一个叫id的字段，这样Ben和Jerry就是两个完全不同的struct，这意味着内存布局就不同了，我们现在看看会出现什么结果
```
Ben says,"Hello my name is Ben and id is 1"
Ben says,"Hello my name is Ben and id is 1"
Jerry says,"Hello my name is Jerry"
Jerry says,"Hello my name is Jerry"
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x2 addr=0x1 pc=0x100d7a860]

goroutine 1 [running]:
fmt.(*buffer).writeString(...)
	/opt/homebrew/Cellar/go/1.16.4/libexec/src/fmt/print.go:82
fmt.(*fmt).padString(0x14000108d40, 0x1, 0x100db3fb0)
	/opt/homebrew/Cellar/go/1.16.4/libexec/src/fmt/format.go:110 +0x78
fmt.(*fmt).fmtS(0x14000108d40, 0x1, 0x100db3fb0)
	/opt/homebrew/Cellar/go/1.16.4/libexec/src/fmt/format.go:359 +0x54
fmt.(*pp).fmtString(0x14000108d00, 0x1, 0x100db3fb0, 0x73)
	/opt/homebrew/Cellar/go/1.16.4/libexec/src/fmt/print.go:446 +0x18c
fmt.(*pp).printArg(0x14000108d00, 0x100ddf0e0, 0x140001366f0, 0x73)
	/opt/homebrew/Cellar/go/1.16.4/libexec/src/fmt/print.go:694 +0x7d8
fmt.(*pp).doPrintf(0x14000108d00, 0x100db9032, 0x21, 0x1400011ff08, 0x1, 0x1)
	/opt/homebrew/Cellar/go/1.16.4/libexec/src/fmt/print.go:1026 +0x12c
fmt.Fprintf(0x100dfa9a8, 0x14000128008, 0x100db9032, 0x21, 0x1400011ff08, 0x1, 0x1, 0x24, 0x0, 0x0)
	/opt/homebrew/Cellar/go/1.16.4/libexec/src/fmt/print.go:204 +0x54
fmt.Printf(...)
	/opt/homebrew/Cellar/go/1.16.4/libexec/src/fmt/print.go:213
main.(*Jerry).Hello(0x14000114018)
	/Users/ttlv/codes/go/src/test/main.go:23 +0x94
main.main()
	/Users/ttlv/codes/go/src/test/main.go:43 +0x190
exit status 2
```
运行**go run main.go**出现了一个空指针的错误，为什么会出现这个错，其实比较好理解，根本原因是两个的内存布局不一样了，Ben比Jerry多了一个id,也就是说Type是Ben，结果Data是Jerry,Jerry的内存数据结构就没有id,那这个时候自然就出现上面的错了。
继续再修改
```
package main

import "fmt"

type IceCreamMaker interface {
	Hello()
}

type Ben struct {
	name string
}

func (b *Ben) Hello() {
	fmt.Printf("Ben says,\"Hello my name is %s\"\n", b.name)
}

type Jerry struct {
	field1 *[5]byte
	field2 int
}

func (j *Jerry) Hello() {
	fmt.Printf("Jerry says,\"Hello my name is %s\"\n", j.field1)
}

func main() {
	var ben = &Ben{name: "Ben"}
	var jerry = &Jerry{field1:&[5]byte{},field2: 5}
	var maker IceCreamMaker = ben
	var loop0, loop1 func()
	loop0 = func() {
		maker = ben
		go loop1()
	}
	loop1 = func() {
		maker = jerry
		go loop0()
	}

	go loop0()

	for {
		maker.Hello()
	}

}
```
这里的Jerry的结构有两个字段，不过要是熟悉string底层结构的同学肯定就看出来了，这个结构就是string的底层数据结构，这段代码跑起来没有问题。底层的内存布局是一样的在某种程度上是兼容的。是不是很神奇！！！
那么如果是一个普通的指针、map、slice 可以安全的更新吗？
如果您非常熟悉这些数据类型的底层数据结构，我觉得您可以尝试用机器字的特性去完成原子赋值，我还是要阐述一个观点，就是
```
没有安全的 data race(safe data race)。您的程序要么没有 data race，要么其操作未定义。如果能保证
- 原子性
- 可见行
这两点那么您的代码一定是安全的，不会出现data race.
```
[并发编程三大特性-原子性、可见性、有序性](https://www.cnblogs.com/yeyang/p/13576636.html)
## sync.atomic
解决data race的方式在go中有几种。这一小结来讨论一下 sync.atomic的用法，请看下面这段代码
```
package main

import (
	"fmt"
	"sync"
)

type Config struct {
	a []int
}

func main() {
	cfg := &Config{}

	go func() {
		i := 0
		for {
			i++
			cfg.a = []int{i, i + 1, i + 2, i + 3, i + 4, i + 5}
		}
	}()
	var wg sync.WaitGroup
	for n := 0; n < 4; n++ {
		wg.Add(1)
		go func() {
			for n := 0; n < 100; n++ {
				fmt.Printf("%v\n", cfg)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
```
有1个goroutine负责写入，有4个goroutine负责读取，很明显这段代码是一段存在**data race**的代码，我们可以用**go build --race**的命令去做测试
```
==================
WARNING: DATA RACE
Read at 0x00c00012a018 by goroutine 8:
  reflect.typedmemmove()
      /opt/homebrew/Cellar/go/1.16.4/libexec/src/runtime/mbarrier.go:177 +0x0
  reflect.packEface()
      /opt/homebrew/Cellar/go/1.16.4/libexec/src/reflect/value.go:121 +0xf0
  reflect.valueInterface()
      /opt/homebrew/Cellar/go/1.16.4/libexec/src/reflect/value.go:1046 +0x160
  reflect.Value.Interface()
      /opt/homebrew/Cellar/go/1.16.4/libexec/src/reflect/value.go:1016 +0x2d08
  fmt.(*pp).printValue()
      /opt/homebrew/Cellar/go/1.16.4/libexec/src/fmt/print.go:722 +0x2d0c
  fmt.(*pp).printValue()
      /opt/homebrew/Cellar/go/1.16.4/libexec/src/fmt/print.go:876 +0x1cf0
  fmt.(*pp).printArg()
      /opt/homebrew/Cellar/go/1.16.4/libexec/src/fmt/print.go:712 +0x1e4
  fmt.(*pp).doPrintf()
      /opt/homebrew/Cellar/go/1.16.4/libexec/src/fmt/print.go:1026 +0x264
  fmt.Fprintf()
      /opt/homebrew/Cellar/go/1.16.4/libexec/src/fmt/print.go:204 +0x54
  fmt.Printf()
      /opt/homebrew/Cellar/go/1.16.4/libexec/src/fmt/print.go:213 +0x98
  main.main.func2()
      /Users/ttlv/codes/go/src/sync_atmoic/atomic.go:27 +0x2c

Previous write at 0x00c00012a018 by goroutine 7:
  main.main.func1()
      /Users/ttlv/codes/go/src/sync_atmoic/atomic.go:19 +0xfc

Goroutine 8 (running) created at:
  main.main()
      /Users/ttlv/codes/go/src/sync_atmoic/atomic.go:25 +0xec

Goroutine 7 (running) created at:
  main.main()
      /Users/ttlv/codes/go/src/sync_atmoic/atomic.go:15 +0x74
==================
{[154280 154281 154282 154283 154284 154285]}
&{[154113 154281 154282 154283 154284 154285]}
.........
.........
Found 7 data race(s)
```
我随便截取了一段打印的数据，发现数据有连续的也有不连续的，连续的是我们期望得到的，不连续的就是产生了**data race**的结果，当然解决这个问题的方式有很多，我相信很多同学自然而然的就会想到加锁🔐，互斥锁Mutex或者是读写锁RWMutex等等，不可否认，加这些锁使能解决问题，但是加锁毕竟会对性能产生影响，比如说上面这段代码是一个goroutine在写4个goutine在读的情况，频繁的加锁开锁对会产生性能瓶颈，所以就引出了本小结的主角，Atomic。
众所周知，Benchmark出真理，我们可以尝试用Benchmark去对比两者的性能区别
```
package test

import (
	"sync"
	`sync/atomic`
	"testing"
)

type Config struct {
	a []int
}
func(c *Config)T(){}

func BenchmarkAtmoic(b *testing.B){
	var v atomic.Value
	v.Store(&Config{})
	go func() {
		i := 0
		for {
			i++
			cfg := &Config{a: []int{i, i + 1, i + 2, i + 3, i + 4, i + 5}}
			v.Store(cfg)
		}
	}()
	var wg sync.WaitGroup
	for n := 0; n < 4; n++ {
		wg.Add(1)
		go func() {
			for n := 0; n < b.N; n++ {
				cfg :=v.Load().(*Config)
				cfg.T()
				//fmt.Printf("%v\n", cfg)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkMutex(b *testing.B) {
	var l sync.RWMutex
	var cfg *Config
	go func() {
		i := 0

		for {
			i++
			l.Lock()
			cfg = &Config{a: []int{i, i + 1, i + 2, i + 3, i + 4, i + 5}}
			l.Unlock()
		}
	}()
	var wg sync.WaitGroup
	for n := 0; n < 4; n++ {
		wg.Add(1)
		go func() {
			for n := 0; n < b.N; n++ {
				l.RLock()
				cfg.T()
				//fmt.Printf("%v\n", cfg)
				l.RUnlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
```
为了可以看清对比的结果我去掉了打印,用一个啥逻辑都没有的函数代替。
```
go test -bench=. benchmark_test.go                                                                 
goos: darwin
goarch: arm64
BenchmarkAtmoic-8   	248052165	         5.121 ns/op
BenchmarkMutex-8    	  907572	      1438 ns/op
PASS
ok  	command-line-arguments	3.092s
```
Atmoic平均读取一次是5.121ns,Mutex平均读取一次是1438ns,量级差距直接肉眼可见。虽然是得出了结论，我们也知道Mutex确实相比之下更重。原因是什么，因为涉及到更多的goroutine之间的上下文切换**pack blocking goroutine**以及唤醒**goroutine**。goroutine虽然是轻量级的线程模型，但是不管怎么说频繁的上下文切换的花销还是很大。如果觉得这数据不够有说服力,那我可以用go tool trace来显示更加详细的CPU性能指标。
先来分析**Atmoic**
```
package main

import (
	`fmt`
	"os"
	"runtime/trace"
	"sync"
	"sync/atomic"
)

type Config struct {
	a []int
}

func main() {
	f, err := os.Create("trace-atmoic.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()
	var v atomic.Value
	v.Store(&Config{})
	go func() {
		i := 0
		for {
			i++
			cfg := &Config{a: []int{i, i + 1, i + 2, i + 3, i + 4, i + 5}}
			v.Store(cfg)
		}
	}()
	var wg sync.WaitGroup
	for n := 0; n < 4; n++ {
		wg.Add(1)
		go func() {
			for n := 0; n < 1000; n++ {
				cfg := v.Load().(*Config)
				fmt.Printf("%v\n", cfg)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
```
执行**go run main.go**之后会在当前目录下生成一个叫**trace-atmoic.out**的文件。执行**go tool trace trace-atmoic.out**
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636284204990-cf95a3ae-9ac8-4503-b37d-dd4b15ad5a39.png#clientId=u646af123-0f4f-4&from=paste&height=900&id=ufdb30271&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1800&originWidth=2880&originalType=binary&ratio=1&size=388466&status=done&style=none&taskId=uc3917fe3-1e01-4279-a718-b566a9927be&width=1440)
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636284221627-638795fc-3386-440d-8c91-4337dab394ab.png#clientId=u646af123-0f4f-4&from=paste&height=900&id=u1ad70147&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1800&originWidth=2880&originalType=binary&ratio=1&size=1305359&status=done&style=none&taskId=ud184d82b-a6a1-4453-b393-8ec7b1fc811&width=1440)
再来分析**Mutex**
```
package main

import (
	`fmt`
	"os"
	"runtime/trace"
	"sync"
)

type Config struct {
	a []int
}

func main() {
	f, err := os.Create("trace-mutex.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()
	var l sync.RWMutex
	var cfg *Config
	go func() {
		i := 0
		for {
			i++
			l.Lock()
			cfg = &Config{a: []int{i, i + 1, i + 2, i + 3, i + 4, i + 5}}
			l.Unlock()
		}
	}()
	var wg sync.WaitGroup
	for n := 0; n < 4; n++ {
		wg.Add(1)
		go func() {
			for n := 0; n < 100; n++ {
				l.RLock()
				fmt.Printf("%v\n", cfg)
				l.RUnlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
```
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636284297292-935cf2aa-e352-4758-937a-832ed3f53cad.png#clientId=u646af123-0f4f-4&from=paste&height=900&id=u95237ce7&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1800&originWidth=2880&originalType=binary&ratio=1&size=333388&status=done&style=none&taskId=ube112268-b3c9-49a1-8bdd-faa8385fc1d&width=1440)
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636284306831-aaae81d7-bd60-4f2b-8d2f-470ae63ed8e7.png#clientId=u646af123-0f4f-4&from=paste&height=900&id=ubcc2beb4&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1800&originWidth=2880&originalType=binary&ratio=1&size=252450&status=done&style=none&taskId=u529a6eb7-59dc-4aca-a5f2-7eb4c07933d&width=1440)
我分别给出了**Atmoic**和**Mutex**的性能数据，明显加锁的**Mutex**在CPU占用的时间上更加的零碎，**Atmoic**占用CPU持续的时间更长而且总的时间也是更长。但是BenchMark的数据来看，**Atmoic**是绝对优于**Mutex**。
这一小节的主角是**Atomic**所以我们要仔细分析下这个是怎么玩的。**Atomic Value**主要用到**Copy-On-Write** 中文翻译就是写时复制。这个思路在微服务降级或者 local cache 场景中经常使用。写时复制指的是，写操作时候复制全量老数据到一个新的对象中，携带上本次新写的数据，之后利用原子替换(atomic.Value)，更新调用者的变量。来完成无锁访问共享数据。**copy-on-write 适合读多写少的场景**这点要切记。熟悉redis的同学是不是直接就联想到了bgsave，没错，bgsave也是**Copy-On-Write**
```
package main
 
import (
	"sync/atomic"
	"time"
)
 
func loadConfig() map[string]string {
	return make(map[string]string)
}
 
func requests() chan int {
	return make(chan int)
}
 
func main() {
	var config atomic.Value // holds current server configuration
	// Create initial config value and store into config.
	config.Store(loadConfig())
	go func() {
		// Reload config every 10 seconds
		// and update config value with the new version.
		for {
			time.Sleep(10 * time.Second)
			config.Store(loadConfig())
		}
	}()
	// Create worker goroutines that handle incoming requests
	// using the latest config value.
	for i := 0; i < 10; i++ {
		go func() {
			for r := range requests() {
				c := config.Load()
				// Handle request r using config c.
				_, _ = r, c
			}
		}()
	}
```
下面的示例展示了如何使用写时复制习惯用法维护可伸缩、频繁读取但不频繁更新的数据结构。大家有没有想过一个问题，对于上面的这段代码，如果我Store了一个v2版本的数据,对于下面读取的goroutine来说，他们不会全部都读到v2版本的数据，就是一部分出现读到了v1,一部分读到了v2,这就是使用**Atomic**出现的弊端，但是如果用的是读写锁，就不会出现这样的问题。
```
package main
 
import (
	"sync"
	"sync/atomic"
)
 
func main() {
	type Map map[string]string
	var m atomic.Value
	m.Store(make(Map))
	var mu sync.Mutex // used only by writers
	// read function can be used to read the data without further synchronization
	read := func(key string) (val string) {
		m1 := m.Load().(Map)
		return m1[key]
	}
	// insert function can be used to update the data without further synchronization
	insert := func(key, val string) {
		mu.Lock() // synchronize with other potential writers
		defer mu.Unlock()
		m1 := m.Load().(Map) // load current value of the data structure
		m2 := make(Map)      // create a new value
		for k, v := range m1 {
			m2[k] = v // copy all data from the current object to the new one
		}
		m2[key] = val // do the update that we need
		m.Store(m2)   // atomically replace the current object with the new one
		// At this point all new readers start working with the new version.
		// The old version will be garbage collected once the existing readers
		// (if any) are done with it.
	}
	_, _ = read, insert
}
```
为什么说是**Atomic**适合读多写少的场景，核心的关键问题就在于拷贝，**Copy-On-Write**每次都要拷贝一份原始数据出来，如果频繁的写就意味着要频繁的拷贝，这样就导致拷贝的成本会非常高，所以说建议在读多写少的场景下使用**Atomic**.
## Mutex
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636284888451-ae0a9e49-604f-4f48-ae4c-bafdf8d4d5e7.png#clientId=u646af123-0f4f-4&from=paste&height=620&id=u90d97767&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1240&originWidth=1102&originalType=binary&ratio=1&size=150348&status=done&style=none&taskId=u085313ad-99a5-4d72-b2ed-69f6a992114&width=551)
这个案例基于两个 goroutine:
goroutine 1 持有100毫秒的锁
goroutine 2 每100ms 持有一次锁
都是100ms 的周期，但是由于 goroutine 1 不断的请求锁，可预期它会更频繁的持续到锁。我们基于 Go 1.8 循环了10次，下面是锁的请求占用分布:
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636284913727-ee171571-bf72-4ddf-b5c8-238b11dd2895.png#clientId=u646af123-0f4f-4&from=paste&height=69&id=u1465a404&margin=%5Bobject%20Object%5D&name=image.png&originHeight=138&originWidth=680&originalType=binary&ratio=1&size=18034&status=done&style=none&taskId=uf1571c2c-b99e-44f9-a01e-b0fcc328ee9&width=340)
g1拿到了**7200216**次锁🔐而g2却拿到了**10**次，g2确实是在一个for循环里，而且就10次，那理论上确实是只能拿到10次锁，我想表达的观点是这两个量级查了这么多，g1是持有100毫秒的锁，g2是每100ms 持有一次锁，按照持锁的结果来看就意味着g2他不怎么能争到锁，这就意味着会出现锁饥饿的现象，这是我们要关注的重点。你的业务代码因为抢不到锁而blocking，业务逻辑就卡主了。所以1.8版本的go的Mutex其实还是有一定的问题的。以上都是我以结果为导向所阐述的观点，下面我将阐述Mutex持锁的原理，先不讲源码，源码这个等到后面我时间了再去做解析，因为源码这个东西每个版本都在变，特别是实现的细节，所以就先不要这么早就关注源码，先说清楚背后的原理。
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636285446275-fa8a01cc-10eb-4fbe-a355-ed5bdb1430ca.png#clientId=u646af123-0f4f-4&from=paste&height=321&id=u2ce7720a&margin=%5Bobject%20Object%5D&name=image.png&originHeight=642&originWidth=1320&originalType=binary&ratio=1&size=108241&status=done&style=none&taskId=u20b48819-7bdf-40be-805c-dc357b3e1f4&width=660)
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636285463108-c3453b1a-c38a-479e-a090-a5e2ac09de4c.png#clientId=u646af123-0f4f-4&from=paste&height=321&id=ua78df233&margin=%5Bobject%20Object%5D&name=image.png&originHeight=642&originWidth=1424&originalType=binary&ratio=1&size=106139&status=done&style=none&taskId=uca63c134-6a15-4bd7-9734-383c472000e&width=712)
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636285473062-3a130339-d6d6-48aa-b5d6-874f93c3ebc6.png#clientId=u646af123-0f4f-4&from=paste&height=356&id=ufcdef5ea&margin=%5Bobject%20Object%5D&name=image.png&originHeight=712&originWidth=1364&originalType=binary&ratio=1&size=127887&status=done&style=none&taskId=u1b67fa54-54f5-4e11-8819-94209423689&width=682)
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636285486353-4d56135a-b8da-44a6-a676-d109ada2fd69.png#clientId=u646af123-0f4f-4&from=paste&height=340&id=uab4f649d&margin=%5Bobject%20Object%5D&name=image.png&originHeight=680&originWidth=1338&originalType=binary&ratio=1&size=131298&status=done&style=none&taskId=u1297f665-90ba-44f2-a367-627440af6c8&width=669)
首先，goroutine1 将获得锁并休眠100ms。当goroutine2 试图获取锁时，它将被添加到锁的队列中- FIFO 顺序，goroutine 将进入等待状态。
然后，当 goroutine1 完成它的工作时，它将释放锁。此版本将通知队列唤醒 goroutine2。goroutine2 将被标记为可运行的，并且正在等待 Go 调度程序在线程上运行.然而，当 goroutine2 等待运行时，goroutine1将再次请求锁。goroutine2 尝试去获取锁，结果悲剧的发现锁又被人持有了，它自己继续进入到等待模式。就这样不停的重复上演悲剧，所以就出现了开头的g1拿到了七百多万次锁。
我们看看几种 Mutex 锁的实现:

- Barging_(插入)_. 这种模式是为了提高吞吐量，当锁被释放时，它会唤醒第一个等待者，然后把锁给第一个等待者或者给第一个请求锁的人。这就是Go 1.8的设计和反映我们之前看到的结果.

![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636285633420-ca117add-e807-47c1-b16c-ee96d4e82765.png#clientId=u646af123-0f4f-4&from=paste&height=225&id=uc30e986b&margin=%5Bobject%20Object%5D&name=image.png&originHeight=450&originWidth=1502&originalType=binary&ratio=1&size=111007&status=done&style=none&taskId=ubed1804c-e196-4a1f-bd46-8be72b375e2&width=751)

- Handsoff_(切换)_. 当锁释放时候，锁会一直持有直到第一个等待者准备好获取锁。它降低了吞吐量，因为锁被持有，即使另一个 goroutine 准备获取它。

![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636285705230-a2b0e241-92f2-439a-83b6-8284260c2738.png#clientId=u646af123-0f4f-4&from=paste&height=325&id=u2438195d&margin=%5Bobject%20Object%5D&name=image.png&originHeight=650&originWidth=1546&originalType=binary&ratio=1&size=169318&status=done&style=none&taskId=u70118bed-078d-4b60-bdf7-f5194ec5ced&width=773)
我们可以在Linux内核的mutex中找到这个逻辑： Mutex Starvation是可能的,因为mutex_lock（）允许锁窃取,其中运行（或乐观spinning）任务优先于唤醒等待者而获取锁. 锁窃取是一项重要的性能优化,因为等待等待者唤醒并获得运行时间可能需要很长时间,在此期间每个人都会在锁定时停止. […]这重新引入了一些等待时间,因为一旦我们进行了切换,我们必须等待等待者再次醒来.
一个互斥锁的 handsoff 会完美地平衡两个goroutine 之间的锁分配，但是会降低性能，因为它会迫使第一个 goroutine 等待锁。

- Spinning_(自旋锁)_. 自旋在等待队列为空或者应用程序重度使用锁时效果不错。当服务员的队列为空或应用程序大量使用互斥锁时,Spinning可能很有用. Parking 和 unparking goroutines是有成本的,可能比等待下一次锁定获取Spinning慢。

当一个线程尝试去获取某一把锁的时候，如果这个锁此时已经被别人获取(占用)，那么此线程就无法获取到这把锁，该线程将会等待，间隔一段时间后会再次尝试获取。这种采用循环加锁 -> 等待的机制被称为自旋锁(spinlock)
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636289156620-7e2dcf9c-23dd-416e-bf48-a38cf1b48a9c.png#clientId=u646af123-0f4f-4&from=paste&height=232&id=ue1d9897a&margin=%5Bobject%20Object%5D&name=image.png&originHeight=463&originWidth=800&originalType=binary&ratio=1&size=98465&status=done&style=none&taskId=u01dc783d-0fbc-4c3b-ab7c-5931a50e8c3&width=400)
Go 1.8也使用了这种策略.当尝试获取已经保持的锁时,如果本地队列为空并且处理器的数量大于1,则goroutine将spinning几次,使用一个处理器spinning将仅阻止该程序. spinning后,goroutine将停放.在程序密集使用锁的情况下,它充当快速路径. 有关如何设计锁定的更多信息 -barging, handoff, spinlock,Filip Pizlo撰写了一篇必读文章[“Locking in WebKit”](https://webkit.org/blog/6161/locking-in-webkit/).
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636285948321-b527125f-93b9-4af5-8d38-0e4d1fa043bd.png#clientId=u646af123-0f4f-4&from=paste&height=331&id=u472cadc5&margin=%5Bobject%20Object%5D&name=image.png&originHeight=662&originWidth=1594&originalType=binary&ratio=1&size=142448&status=done&style=none&taskId=ucc53f2f0-e25b-4e30-a5ca-746c72d0ed0&width=797)
总结：Go 1.8 使用了 Barging 和 Spining 的结合实现。当试图获取已经被持有的锁时，如果本地队列为空并且 P [理解成goroutine的队列]的数量大于1，goroutine 将自旋几次(用一个 P 旋转会阻塞程序)。自旋后，goroutine park。在程序高频使用锁的情况下，它充当了一个快速路径(fast path)
Go 1.9 通过添加一个新的饥饿模式来解决先前解释的问题，该模式将会在释放时候触发 handsoff。所有等待锁超过一毫秒的 goroutine(也称为有界等待)将被诊断为饥饿。当被标记为饥饿状态时，unlock 方法会 handsoff 把锁直接扔给第一个等待者。
在饥饿模式下，自旋也被停用，因为传入的goroutines 将没有机会获取为下一个等待者保留的锁。我们来看下最后Go1.9跑本小节刚刚开始的代码的结果
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636289233272-d09c05a9-e41d-4583-8932-c5e5e873c431.png#clientId=u646af123-0f4f-4&from=paste&height=118&id=u88622098&margin=%5Bobject%20Object%5D&name=image.png&originHeight=236&originWidth=982&originalType=binary&ratio=1&size=24952&status=done&style=none&taskId=uaf426e85-c66a-413e-8867-3208fb62692&width=491)
g1值拿到了57次锁,极大的降低了锁饥饿的情况的发生，非常有效。👍🏻
## errorgroup
在最后推荐一个包，就是**errorgroup**在并发请求处理的场景下特别好用。
我们把一个复杂的任务，尤其是依赖多个微服务 rpc 需要聚合数据的任务，分解为依赖和并行，依赖的意思为: 需要上游 a 的数据才能访问下游 b 的数据进行组合。但是并行的意思为: 分解为多个小任务并行执行，最终等全部执行完毕。
```
package main

func main() {
	var a, b int
	var err1, err2 error
	var ch chan result
	// call rpc1
	go func() {
		a, err1 := rpc_servive1()
		ch <- result{a,err1}
	}()
	// call rpc2
	go func() {
		b, err2 := rpc_servive2()
			ch <- result{b,err2}
	}()
	<- ch
}
```
上面的代码是一段伪代码，不能编译，就想说明的是用goroutine并发的去请求比如rpc服务，如果不用**errorgroup**,那只能按照上面遮掩写，而且要处理每一个rpc调用的返回值和error还会特别的麻烦，就像我上面的要定一个**chan**去接收。**errorgroup**的做法就是把这些操作打包处理好。 现在我要讲**errgroup**这个包怎么使用了。
```
package main

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	var a []int
	// 调用服务a 正常没有error
	g.Go(func() error {
		a = []int{0}
		return nil
	})

	// 调用服务b异常 有error
	g.Go(func() error {
		a = []int{1}
		return errors.New("error b")
	})

	err := g.Wait()
	fmt.Println(a)
	fmt.Println(err)
	fmt.Println(ctx.Err())
}
```
这一看代码就和弄清楚了，对于并发调用的错误处理和服务降级，**errgroup**就完全能胜任了。这个包的核心工作原理就是利用 sync.Waitgroup 管理并行执行的 goroutine。

- 并行工作流
- 误处理 或者 优雅降级
- 利用局部变量+闭包

![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636289591004-87e61fce-f939-4e56-9e70-18391ac3372c.png#clientId=u646af123-0f4f-4&from=paste&height=362&id=uc0b6bd4f&margin=%5Bobject%20Object%5D&name=image.png&originHeight=724&originWidth=1620&originalType=binary&ratio=1&size=79540&status=done&style=none&taskId=uae3144db-e7a2-49bd-985f-5b2e9423ed5&width=810)
这个包虽然很好用但是并不是完美的，也有缺陷

- 比如说我在goroutine里面加了一个error的handler是直接panic，这直接一整个进程就退出了，很可怕。
- 第二个就是调用者创建了大量的goroutine,完全超出了**GOMAXPROCS**的值。
- 返回的ctx,可能会和别人串用导致大量报错，比如说用了一个已经close的ctx。

这几点其实在**Kratos**中已经有对应的解决方案了，**Kratos**其实也是基于原生的**errgroup**包做了一些优化，解决了上面的三个问题。
## sync.Pool
在 golang 中有一个池，它特别神奇，你只要和它有个约定，你要什么它就给什么，你用完了还可以还回去，但是下次拿的时候呢，却不一定是你上次存的那个，这个池就是 sync.Pool。
首先我们来看看这个 sync.Pool 是如何使用的，其实非常的简单。
它一共只有三个方法我们需要知道的：New、Put、Get
```

package main

import (
    "fmt"
    "sync"
)

var strPool = sync.Pool{
    New: func() interface{} {
        return "test str"
    },
}

func main() {
    str := strPool.Get()
    fmt.Println(str)
    strPool.Put(str)
}
```

- 通过**New**去定义你这个池子里面放的究竟是什么东西，在这个池子里面你只能放一种类型的东西。比如在上面的例子中我就在池子里面放了字符串。
- 我们随时可以通过**Get**方法从池子里面获取我们之前在New里面定义类型的数据。
- 当我们用完了之后可以通过**Put**方法放回去，或者放别的同类型的数据进去。

那么这个池子的目的是什么呢？其实一句话就可以说明白，就是为了复用已经使用过的对象，来达到优化内存使用和回收的目的。说白了，一开始这个池子会初始化一些对象供你使用，如果不够了呢，自己会通过new产生一些，当你放回去了之后这些对象会被别人进行复用，当对象特别大并且使用非常频繁的时候可以大大的减少对象的创建和回收的时间。所以sync.Pool 的场景是用来保存和复用临时对象，以减少内存分配，降低 GC 压力(Request-Driven 特别合适)。sync.Pool 的场景是用来保存和复用临时对象，以减少内存分配，降低 GC 压力(Request-Driven 特别合适)。
## chan
channels 是一种类型安全的消息队列，充当两个 goroutine 之间的管道，将通过它同步的进行任意资源的交换。chan 控制 goroutines 交互的能力从而创建了 Go 同步机制。**当创建的 chan 没有容量时，称为无缓冲通道**。反过来，**使用容量创建的 chan 称为缓冲通道**。要了解通过 chan 交互的 goroutine 的同步行为是什么，我们需要知道通道的类型和状态。根据我们使用的是无缓冲通道还是缓冲通道，场景会有所不同，所以让我们单独讨论每个场景。
### 无缓冲通道
```
ch := make(chan struct{})
```
无缓冲信道的本质是保证同步。

- Receive 先于 Send 发生。
- 好处: 100% 保证能收到。
- 代价: 延迟时间未知。
### 有缓冲通道

1. 通道中有空间，发送可以立即进行
1. 缓冲通道为空时，该通道将锁住 goroutine 并使其等待资源被发送。
- Send 先于 Receive 发生。
- 好处: 延迟更小。
- 代价: 不保证数据到达，越大的 buffer，越小的保障到达。buffer = 1 时，给你延迟一个消息的保障。
### Latencies due to under-sized buffer（缓冲区不足导致的延迟）
缓冲区不是越大性能就越好，空间换时间也是有瓶颈的
### Go Concurrency Patterns

- Timing out
- Moving on
- Pipeline
- Fan-out, Fan-in
- CancellationClose
    - 先于 Receive 发生（类似 Buffered）。
    - 不需要传递数据，或者传递 nil。
    - 非常适合取消和超时控制。
- Context
#### 参考

- [https://blog.golang.org/concurrency-timeouts](https://go.dev/blog/concurrency-timeouts)
- [https://blog.golang.org/pipelines](https://blog.golang.org/pipelines)
- [https://talks.golang.org/2013/advconc.slide#1](https://talks.golang.org/2013/advconc.slide#1)
- [https://github.com/go-kratos/kratos/tree/master/pkg/sync](https://github.com/go-kratos/kratos/tree/master/pkg/sync)
### Design Philosophy

- If any given Send on a channel CAN cause the sending goroutine to block（如果通道上任何send会导致goroutine阻塞）:
    - Not allowed to use a Buffered channel larger than 1.（buffed channel size不要大于1）
        - Buffers larger than 1 must have reason/measurements.（大于1必须有原因）
    - Must know what happens when the sending goroutine blocks.（必须知道goroutine阻塞时会发生什么）
- If any given Send on a channel WON’T cause the sending goroutine to block（如果通道上任何send不会导致goroutine阻塞）:
    - You have the exact number of buffers for each send（你有每次发送的确切缓冲size）.
        - Fan Out pattern
    - You have the buffer measured for max capacity（您测量了最大容量的缓冲区）.
        - Drop pattern
- Less is more with buffers（buffer少即是多）
    - Don’t think about performance when thinking about buffers.（考虑缓冲区时不要考虑性能）
    - Buffers can help to reduce blocking latency between signaling（缓冲区可以帮助减少信号之间的阻塞延迟）.
        - Reducing blocking latency towards zero does not necessarily mean better throughput（将阻塞延迟降低到零并不一定意味着更高的吞吐量。）.
        - If a buffer of one is giving you good enough throughput then keep it（如果一个缓冲区为您提供足够好的吞吐量，请保留它。）.
        - Question buffers that are larger than one and measure for size（大于1 并测量大小的问题缓冲区）.
        - Find the smallest buffer possible that provides good enough throughput（找到可能提供足够高吞吐量的最小缓冲区。）.
## Package context
在 Go 服务中，每个传入的请求都在其自己的goroutine 中处理。请求处理程序通常启动额外的 goroutine 来访问其他后端，如数据库和 RPC 服务。处理请求的 goroutine 通常需要访问特定于请求（request-specific context）的值，例如最终用户的身份、授权令牌和请求的截止日期（deadline）。
当一个请求被取消或超时时，处理该请求的所有 goroutine 都应该快速退出（fail fast），这样系统就可以回收它们正在使用的任何资源。
Go 1.7 引入一个 context 包，它使得跨 API 边界的请求范围元数据、取消信号和截止日期很容易传递给处理请求所涉及的所有 goroutine（显示传递）。
### 将 context 对象集成到 API 中

1. The first parameter of a function call

![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636292189365-cb1aed80-0859-429e-b12c-d1a173e698c5.png#clientId=u646af123-0f4f-4&from=paste&height=77&id=uc2728179&margin=%5Bobject%20Object%5D&name=image.png&originHeight=290&originWidth=2560&originalType=binary&ratio=1&size=111867&status=done&style=none&taskId=u23add16a-0b32-4b57-9972-c24a67ad04f&width=680)

2. Optional config on a request structure

例如net/http库
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636292216542-69a7ba2c-5ead-4abd-85f4-69a2e0eaa137.png#clientId=u646af123-0f4f-4&from=paste&height=64&id=u27962018&margin=%5Bobject%20Object%5D&name=image.png&originHeight=230&originWidth=2560&originalType=binary&ratio=1&size=92545&status=done&style=none&taskId=u964fba07-d7ed-4f64-bf80-a1c8709b91d&width=715)
### context.WithValue
context.WithValue 内部基于 valueCtx 实现
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636292310290-b847f678-2716-453a-a67c-bc7c984611f5.png#clientId=u646af123-0f4f-4&from=paste&height=116&id=u0f168f42&margin=%5Bobject%20Object%5D&name=image.png&originHeight=510&originWidth=2560&originalType=binary&ratio=1&size=559864&status=done&style=none&taskId=u8e0d05d8-455c-464c-b7f9-3751ba82167&width=580)
为了实现不断的 WithValue，构建新的 context，**内部在查找 key 时候，使用递归方式不断从当前，从父节点寻找匹配的 key，直到 root context**（Backgrond 和 TODO Value 函数会返回）
​

![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636292343746-50e9e0be-2a81-4d61-98fc-4d26f9a8534c.png#clientId=u646af123-0f4f-4&from=paste&height=151&id=u7c8fb1a5&margin=%5Bobject%20Object%5D&name=image.png&originHeight=689&originWidth=2560&originalType=binary&ratio=1&size=640197&status=done&style=none&taskId=u8512f03c-2640-498b-9170-6d690f375dd&width=562)
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636292358001-e5a91f60-7db0-4f76-9f90-5aa76420c643.png#clientId=u646af123-0f4f-4&from=paste&height=410&id=uecb29ce8&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1440&originWidth=1849&originalType=binary&ratio=1&size=246178&status=done&style=none&taskId=u895b3ed3-9fbf-4f1b-892c-5f06dc076b3&width=527)
### When a Context is canceled, all Contexts derived from it are also canceled
当一个 context 被取消时，从它派生的所有 context 也将被取消。WithCancel(ctx) 参数 ctx 认为是 parent ctx，在内部会进行一个传播关系链的关联。Done() 返回 一个 chan，当我们取消某个parent context, 实际上上会递归层层 cancel 掉自己的 child context 的 done chan 从而让整个调用链中所有监听 cancel 的 goroutine 退出。
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636292464048-4e9050b1-b2c3-4869-b258-11e89f55990e.png#clientId=u646af123-0f4f-4&from=paste&height=401&id=u6953ee6a&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1440&originWidth=1743&originalType=binary&ratio=1&size=171447&status=done&style=none&taskId=u1a43a875-a577-4a66-9e92-235a5cd1226&width=485)
### All blocking/long operations should be cancelable
如果要实现一个超时控制，通过上面的 context 的 parent/child 机制，其实我们只需要启动一个定时器，然后在超时的时候，直接将当前的 context 给 cancel 掉，就可以实现监听在当前和下层的额 context.Done() 的 goroutine 的退出。
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1636292540658-7a390154-9ee7-42d9-aacb-0f1e2426ae62.png#clientId=u646af123-0f4f-4&from=paste&height=584&id=u9a997bd1&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1168&originWidth=2560&originalType=binary&ratio=1&size=391699&status=done&style=none&taskId=u456f765b-414c-4517-96d0-1b0d7d6d555&width=1280)
## 参考
[https://www.ardanlabs.com/blog/2018/11/goroutine-leaks-the-forgotten-sender.html](https://www.ardanlabs.com/blog/2018/11/goroutine-leaks-the-forgotten-sender.html)
[https://www.ardanlabs.com/blog/2019/04/concurrency-trap-2-incomplete-work.html](https://www.ardanlabs.com/blog/2019/04/concurrency-trap-2-incomplete-work.html)
[https://www.ardanlabs.com/blog/2014/01/concurrency-goroutines-and-gomaxprocs.html](https://www.ardanlabs.com/blog/2014/01/concurrency-goroutines-and-gomaxprocs.html)
[https://dave.cheney.net/practical-go/presentations/qcon-china.html#_concurrency](https://dave.cheney.net/practical-go/presentations/qcon-china.html#_concurrency)
[https://golang.org/ref/mem](https://golang.org/ref/mem)
[https://blog.csdn.net/caoshangpa/article/details/78853919](https://blog.csdn.net/caoshangpa/article/details/78853919)
[https://blog.csdn.net/qcrao/article/details/92759907](https://blog.csdn.net/qcrao/article/details/92759907)
[https://cch123.github.io/ooo/](https://cch123.github.io/ooo/)
[https://blog.golang.org/codelab-share](https://blog.golang.org/codelab-share)
[https://dave.cheney.net/2018/01/06/if-aligned-memory-writes-are-atomic-why-do-we-need-the-sync-atomic-package](https://dave.cheney.net/2018/01/06/if-aligned-memory-writes-are-atomic-why-do-we-need-the-sync-atomic-package)
[http://blog.golang.org/race-detector](http://blog.golang.org/race-detector)
[https://dave.cheney.net/2014/06/27/ice-cream-makers-and-data-races](https://dave.cheney.net/2014/06/27/ice-cream-makers-and-data-races)
[https://www.ardanlabs.com/blog/2014/06/ice-cream-makers-and-data-races-part-ii.html](https://www.ardanlabs.com/blog/2014/06/ice-cream-makers-and-data-races-part-ii.html)
[https://medium.com/a-journey-with-go/go-how-to-reduce-lock-contention-with-the-atomic-package-ba3b2664b549](https://medium.com/a-journey-with-go/go-how-to-reduce-lock-contention-with-the-atomic-package-ba3b2664b549)
[https://medium.com/a-journey-with-go/go-discovery-of-the-trace-package-e5a821743c3c](https://medium.com/a-journey-with-go/go-discovery-of-the-trace-package-e5a821743c3c)
[https://medium.com/a-journey-with-go/go-mutex-and-starvation-3f4f4e75ad50](https://medium.com/a-journey-with-go/go-mutex-and-starvation-3f4f4e75ad50)
[https://www.ardanlabs.com/blog/2017/10/the-behavior-of-channels.html](https://www.ardanlabs.com/blog/2017/10/the-behavior-of-channels.html)
[https://medium.com/a-journey-with-go/go-buffered-and-unbuffered-channels-29a107c00268](https://medium.com/a-journey-with-go/go-buffered-and-unbuffered-channels-29a107c00268)
[https://medium.com/a-journey-with-go/go-ordering-in-select-statements-fd0ff80fd8d6](https://medium.com/a-journey-with-go/go-ordering-in-select-statements-fd0ff80fd8d6)
[https://www.ardanlabs.com/blog/2017/10/the-behavior-of-channels.html](https://www.ardanlabs.com/blog/2017/10/the-behavior-of-channels.html)
[https://www.ardanlabs.com/blog/2014/02/the-nature-of-channels-in-go.html](https://www.ardanlabs.com/blog/2014/02/the-nature-of-channels-in-go.html)
[https://www.ardanlabs.com/blog/2013/10/my-channel-select-bug.html](https://www.ardanlabs.com/blog/2013/10/my-channel-select-bug.html)
[https://blog.golang.org/io2013-talk-concurrency](https://blog.golang.org/io2013-talk-concurrency)
[https://blog.golang.org/waza-talk](https://blog.golang.org/waza-talk)
[https://blog.golang.org/io2012-videos](https://blog.golang.org/io2012-videos)
[https://blog.golang.org/concurrency-timeouts](https://blog.golang.org/concurrency-timeouts)
[https://blog.golang.org/pipelines](https://blog.golang.org/pipelines)
[https://www.ardanlabs.com/blog/2014/02/running-queries-concurrently-against.html](https://www.ardanlabs.com/blog/2014/02/running-queries-concurrently-against.html)
[https://blogtitle.github.io/go-advanced-concurrency-patterns-part-3-channels/](https://blogtitle.github.io/go-advanced-concurrency-patterns-part-3-channels/)
[https://www.ardanlabs.com/blog/2013/05/thread-pooling-in-go-programming.html](https://www.ardanlabs.com/blog/2013/05/thread-pooling-in-go-programming.html)
[https://www.ardanlabs.com/blog/2013/09/pool-go-routines-to-process-task.html](https://www.ardanlabs.com/blog/2013/09/pool-go-routines-to-process-task.html)
[https://blogtitle.github.io/categories/concurrency/](https://blogtitle.github.io/categories/concurrency/)
[https://medium.com/a-journey-with-go/go-context-and-cancellation-by-propagation-7a808bbc889c](https://medium.com/a-journey-with-go/go-context-and-cancellation-by-propagation-7a808bbc889c)
[https://blog.golang.org/context](https://blog.golang.org/context)
[https://www.ardanlabs.com/blog/2019/09/context-package-semantics-in-go.html](https://www.ardanlabs.com/blog/2019/09/context-package-semantics-in-go.html)
[https://golang.org/ref/spec#Channel_types](https://golang.org/ref/spec#Channel_types)
[https://drive.google.com/file/d/1nPdvhB0PutEJzdCq5ms6UI58dp50fcAN/view](https://drive.google.com/file/d/1nPdvhB0PutEJzdCq5ms6UI58dp50fcAN/view)
[https://medium.com/a-journey-with-go/go-context-and-cancellation-by-propagation-7a808bbc889c](https://medium.com/a-journey-with-go/go-context-and-cancellation-by-propagation-7a808bbc889c)
[https://blog.golang.org/context](https://blog.golang.org/context)
[https://www.ardanlabs.com/blog/2019/09/context-package-semantics-in-go.html](https://www.ardanlabs.com/blog/2019/09/context-package-semantics-in-go.html)
[https://golang.org/doc/effective_go.html#concurrency](https://golang.org/doc/effective_go.html#concurrency)
[https://zhuanlan.zhihu.com/p/34417106?hmsr=toutiao.io](https://zhuanlan.zhihu.com/p/34417106?hmsr=toutiao.io)
[https://talks.golang.org/2014/gotham-context.slide#1](https://talks.golang.org/2014/gotham-context.slide#1)
[https://medium.com/@cep21/how-to-correctly-use-context-context-in-go-1-7-8f2c0fafdf39](https://medium.com/@cep21/how-to-correctly-use-context-context-in-go-1-7-8f2c0fafdf39)
[

](https://www.ardanlabs.com/blog/2018/11/goroutine-leaks-the-forgotten-sender.html)
 
