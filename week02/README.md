## Error vs Exception
### 痛点
大量的if err!=nil，减少该类代码
### error正确使用
使用规范：errors.New("package name : error")
errors源码
```
package errors

// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
func New(text string) error {
	return &errorString{text}
}


// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

```
#### 为什么errors.New要返回指针?
在go中判断两个Error是否相等，一定是要要判断内存地址是否相等，不能是简单的字符串相等，这是不合理的
### 异常演进历史
#### C语言
单参数返回
#### C++
可以抛出异常，但是调用方不知道会抛出什么异常
#### Java
checked exception，可以指定调用者必须捕获异常，但是存在被滥用的情况
### panic
panic就相当于fatal error（挂了）
异常是不可恢复的，err是需要调用者处理的
### go为什么要这样实现Error？

- 简单
- 考虑失败，而不是成功
- 没有隐藏的控制流
- 完全交给你控制Error
- Error are values

​

## Error Type
#### Sentinel Error
预定义的特定错误
例子：
```
EOF = errors.New("EOF")
```
缺点：

1. 使Error成为你的API的一部分，增加API的表面积
1. 不能携带更多信息
#### Error Yypes
实现了Error接口的自定义类型
例子:
```
type MyError struct {
	Msg string
	Name string
}

func (e *MyError) Error() string {
	return e.Msg
}

os.PathError
```
调用者可通过类型断言，来获取更多的上下文信息
虽然比Sentinel Error多了一些上下文信息，但是还是存在一些相同的问题
#### Opaque Errors
不透明的错误，代码和调用者之间的耦合最少
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1635345472889-9782cdd9-2677-4753-98d0-c8b8d3fbccdd.png#clientId=ua4d1cf12-3d46-4&from=paste&height=239&id=u08ead92f&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1440&originWidth=2480&originalType=binary&ratio=1&size=1422795&status=done&style=none&taskId=u81274735-9199-4805-a3de-ce0529eafea&width=411)
在少数情况下，这种二分错误处理方法是不够的。例如，与进程外的世界进行交互（如网络活动），需要调用方调查错误的性质，以确定重试该操作是否合理。在这种情况下，我们可以断言错误实现了特定的行为，而不是断言错误是特定的类型或值。考虑这个例子
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1635345564411-0ff30b43-0eee-4751-a960-2a90351880fa.png#clientId=ua4d1cf12-3d46-4&from=paste&height=233&id=ue8d74e07&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1106&originWidth=2560&originalType=binary&ratio=1&size=1453805&status=done&style=none&taskId=u10ecb30b-6a9a-4ae8-a0aa-6543b56aca1&width=539)
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1635345576170-8725b982-8a51-4758-8173-0116bf9edfb3.png#clientId=ua4d1cf12-3d46-4&from=paste&height=185&id=u175aa724&margin=%5Bobject%20Object%5D&name=image.png&originHeight=868&originWidth=2560&originalType=binary&ratio=1&size=476533&status=done&style=none&taskId=ue4b0f5a2-8181-41f3-bb12-bc4c94be41d&width=547)
![image.png](https://cdn.nlark.com/yuque/0/2021/png/757992/1635345590044-3f9644bf-1e3b-4505-951e-58f0b265f861.png#clientId=ua4d1cf12-3d46-4&from=paste&height=169&id=uf981f38b&margin=%5Bobject%20Object%5D&name=image.png&originHeight=799&originWidth=2560&originalType=binary&ratio=1&size=365307&status=done&style=none&taskId=uce34efba-c583-49ed-ab0c-f917b171018&width=542)
## Handing Error
#### 如何减少 if err!=nil?
例子1,是一个io.Writer的例子
```
_, err = fd.Write(p0[a:b])
if err != nil {
    return err
}
_, err = fd.Write(p1[c:d])
if err != nil {
    return err
}
_, err = fd.Write(p2[e:f])
if err != nil {
    return err
}
```
优化后的代码
```
type errWriter struct {
	w   io.Writer
	err error
}

func (ew *errWriter) Write(buf []byte) {
	if ew.err != nil {
		return
	}
	_, ew.err = ew.w.Write(buf)
}

ew := &errWriter{w: fd}
ew.write(p0[a:b])
ew.write(p1[c:d])
ew.write(p2[e:f])
if ew.err != nil {
    return ew.err
}
```
使用errWriter来封装Error，自定义实现了Write方法，减少了大量的error的判空逻辑
#### wrap error
解决没有堆栈信息的问题
[pkg/errors](https://github.com/pkg/errors)
```
_, err := ioutil.ReadAll(r)
if err != nil {
        return errors.Wrap(err, "read failed")
}
```
errors.Wrap包装error，并附带额外信息，包含堆栈信息
errors.WithMessage,和Wrap类似，但是没有堆栈信息
errors.Cause,返回根error
errors.Errorf打印堆栈信息
## go 1.13error
go1.13 errors 包包含两个用于检查错误的新函数：Is 和 As。
fmt.Errorf可以向错误添加附加信息
## 原则

- 错误只处理一次
- 第三方库尽量不要wrap error，应该在业务代码中wrap error
## 参考

- [https://blog.huoding.com/2019/04/11/728](https://blog.huoding.com/2019/04/11/728)
