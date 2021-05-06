# GO语言入门学习

## 基础语法demo

## Go Modelus用法

go mod init `package name`

package name表示包名，默认为文件夹的名字

使用了go mod之后，目录都是以go mod中的modules为基准，不再支持 import <相对路径>,因此，geecache中day3下的cache.go中引用lru的话，需要在第四行加入example

## GeeCache缓存系统

### 简介

> 参考了GeekTuTu的教程，表示感谢！

`GeeCache`模仿了groupcache的实现，支持的特性有：

* 单机缓存和基于HTTP的分布式缓存
* LRU缓存策略
* 使用GO的锁机制防止缓存击穿
* 使用一致性哈希算法选择节点，实现负载均衡
* 使用protobuf优化节点间的通信

### LRU缓存淘汰策略

`GeeCache`是内存缓存，大小受到限制，因此当内存空间已经满了以后，我们需要将内存中的部分数据按照一定的策略淘汰出来，然后再添加新的缓存。

`LRU`是指最近最少使用，算法思想就是维护一个队列，如果某条记录被访问了，则移动到队尾，那么队首则是最近最少访问的数据，淘汰该条记录即可。

```go
这里有个注意的地方：
kv := ele.Value.(*entry)
这个地方的语法是类型断言，用来判断空接口中的值
语法格式： x.(T) 
 	x：表示类型为interface{}的变量
    T：表示断言x可能是的类型。
该语法返回两个参数，第一个参数是x转化为T类型后的变量，第二个值是一个布尔值，
若为true则表示断言成功，为false则表示断言失败。
```

### 单机并发缓存

#### 数组和切片

数组：是同一种数据类型的固定长度的序列。

```
数组定义：var a [len]int，比如：var a [5]int，数组长度必须是常量，且是类型的组成部分。一旦定义，长度不能变。
数组是值类型，赋值和传参会复制整个数组，而不是指针。因此改变副本的值，不会改变本身的值。
```

切片：切片是数组的一个引用，因此切片是引用类型。但自身是结构体，值拷贝传递。

```
切片的定义：var 变量名 []类型，比如 var str []string  var arr []int
切片的长度可以改变，因此，切片是一个可变的数组。
```

对切片的操作会影响到数组

### protobuf

* protoc --go_out=. *.proto可以将目录下的proto文件转换成Go代码，但是可能会报错

```GO
protoc-gen-go: unable to determine Go import path for "student.proto"
```

解决方法：在proto文件头中加入`option go_package ="XXXX";`

`XXXX`代表包名，详情可见：[google docs](https://developers.google.com/protocol-buffers/docs/reference/go-generated#package)

### TODO

* geecache/day5的测试过程存在疑问
* 关于【分布式】存在疑问
* 可以使用gobenchmark来测试一下

### 参考资料

* [GeekTuTu七天从零实现教程](https://geektutu.com/post/geecache.html)
* 

## 小计
* `VS code中markdown文档，用ctrl+k v可以实现双侧显示`