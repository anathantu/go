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

对切片的操作会影响到数组，可见`slice\slice.go`

#### ByteView

> 　Go语言没有像其它语言一样有public、protected、private等访问控制修饰符，它是通过字母大小写来控制可见性的，如果定义的常量、变量、类型、接口、结构、函数等的名称是大写字母开头表示能被其它包访问或调用（相当于public），非大写开头就只能在包内使用（相当于private，变量或常量也可以下划线开头）
>

`ByteView`中的变量是`b`，只能在包内访问，所以在没有定义写入方法的情况下，它是只读的。

##### ByteView的只读的特性的作用：

`geecache.go`中的`getLocally`函数如下：

```go
func (g *Group) getLocally(key string) (ByteView, error) {
   //调用Group中的回调函数
    bytes, err := g.getter.Get(key)     //这里的bytes是[]byte数组
   if err != nil {
      return ByteView{}, err
   }
   //保存时赋值一份
   value := ByteView{b: cloneBytes(bytes)}
   g.populateCache(key, value)
   return value, nil
}
```

`bytes`是切片，那么为了防止被外部程序修改，需要将`bytes`转成`ByteView`并且保存到缓存中。

然后在`Get`方法中，有一段代码

```go
if v, ok := g.mainCache.get(key); ok {
   log.Println("[GeeCache] hit")
   return v, nil
}
```

这里的`v`是`ByteView`类型的，可以直接返回给外部程序，因为`ByteView`只读的特性（这里就与`[]byte`类型的处理有差别了)

#### Group

`group`最核心的数据结构，包含了缓存、`Getter`接口等

### GeeCache HTTP服务端

分布式缓存需要实现节点间通信，建立基于 HTTP 的通信机制是比较常见和简单的做法。如果一个节点启动了 `HTTP` 服务，那么这个节点就可以被其他节点访问。

> Go中提供了http的包，想要启动一个http服务，调用http.ListenAndServe即可。
>
> `http.ListenAndServe` 接收 2 个参数，第一个参数是服务启动的地址，第二个参数是 Handler，任何实现了 `ServeHTTP` 方法的对象都可以作为 HTTP 的 Handler。

#### HTTPPool

作为承载节点间 `HTTP` 通信的核心数据结构（包括服务端和客户端）

```go
type HTTPPool struct {
   //self 用来记录自己的ip地址，basePath作为节点间通讯的前缀
   self     string
   basePath string
}
```

实现了`ServeHTTP`方法，那么`HTTPPool`结构体就可以作为`http.ListenAndServe`的第二个参数了。

### 一致性哈希

#### 一致性哈希算法简介

##### 分布式哈希

分布式缓存通过hash将不同的缓存存储到特定的节点上，这样每次查询给定的key都会从同一个节点查询，有两个有优点：

1. 缓存查找的效率很高
2. 每个节点存储的都是映射到该节点的缓存值，节省空间

![](https://geektutu.com/post/geecache-day4/hash_select.jpg)

##### 考虑节点数量变化的场景

若分布式节点数量发生了变化，比如减少一台服务器或者增加一台服务器，那么`hash(key)%10`就变成了`hash(key)%9`或者`hash(key)%11`。这意味着，无法通过给定`key`找到存储缓存值的节点了，需要重新去数据源获取数据，容易引起缓存雪崩。

> 缓存在同一时刻全部失效，造成瞬时DB请求量大、压力骤增，引起雪崩。常因为缓存服务器宕机，或缓存设置了相同的过期时间引起。

##### 一致性哈希算法原理

**步骤**

一致性哈希算法将`key`映射到`2^32`的空间中，将这个数字首尾相连形成一个环。

* 计算节点/机器(通常使用**节点的名称**、**编号**和 **IP 地址**)的哈希值，放置在环上。
* 计算 `key` 的哈希值，放置在环上，顺时针寻找到的第一个节点，就是应选取的节点/机器。

<img src="https://geektutu.com/post/geecache-day4/add_peer.jpg" style="zoom: 67%;" />

一开始有`peer2`、`peer4`和`peer6`三个节点，`key11`、`key2`和`key27`都映射到`peer2`中，然后插入`peer8`节点后，`key27`映射到`peer8`（`peer6`到`peer8`之间的都映射到`peer8`），其他的节点并没有改变。

也就是说，一致性哈希算法在新增/删除节点后，只会影响到插入节点附近的缓存值，就可以避免缓存雪崩的问题。

**数据倾斜与虚拟节点**

数据倾斜是指节点的位置不均匀，造成大量的`key`映射到一个节点中，造成节点间缓存负载不均衡。比如，大部分的`key`值都会映射到`peer2`节点。

为了解决负载不均衡的问题，引入了虚拟节点，一个真实的节点对应多个虚拟节点。

假设 1 个真实节点对应 3 个虚拟节点，那么 peer1 对应的虚拟节点是 peer1-1、 peer1-2、 peer1-3（通常以添加编号的方式实现），其余节点也以相同的方式操作。

* 第一步，计算虚拟节点的 Hash 值，放置在环上。
* 第二步，计算 key 的 Hash 值，在环上顺时针寻找到应选取的虚拟节点，例如是 peer2-1，那么就对应真实节点 peer2。

采用这种方式，缓存值的映射会更加均衡了，我们需要做的就是记录真实节点与虚拟节点的映射关系。

### 分布式节点

需要利用到一致性哈希来查找远程节点。

从main方法来看，每次启动程序都会创建一个Group

```go
type Group struct {
   name      string
   getter    Getter
   mainCache cache
   peers PeerPicker
}
```

并且goup的peers会记录其他节点，然后查照时，利用一致性哈希找到对应节点，然后取得缓存值。

另外http.go中，定义了httpGetter，并实现了Get方法，那么向其他节点请求缓存值时，发起的就是http请求了。

```go
type httpGetter struct {
   baseURL string
}
```

请求后得到返回值，这个返回值就是我们想要的缓存值了。

### singleflight

singleflight可以用来防止缓存击穿，原理就是发起多个相同的请求，只有一个能够真实地进行请求，当这个请求得到结果以后，其他的请求也共享该结果。

实现上来说，通过一个map来记录请求，每一个key对应一类请求，通过加锁的机制让多个同类请求只有一个能够完成请求，其他都被阻塞。

### protobuf

> protobuf 即 Protocol Buffers，Google 开发的一种数据描述语言，是一种轻便高效的结构化数据存储格式，与语言、平台无关，可扩展可序列化。protobuf 以二进制方式存储，占用空间小。

`protobuf` 广泛地应用于远程过程调用`(RPC)` 的二进制传输，使用 `protobuf` 的目的非常简单，为了获得更高的性能。传输前使用 `protobuf` 编码，接收方再进行解码，可以显著地降低二进制传输的大小。另外一方面，`protobuf` 可非常适合传输结构化数据，便于通信字段的扩展。

使用方式分两步：

* 按照protobuf的语法，在.proto文件中定义数据结构，并用protoc生成Go代码
* 在项目代码中引用生成的Go代码

```GO
protoc --go_out=. *.proto可以将目录下的proto文件转换成Go代码，但是可能会报错
protoc-gen-go: unable to determine Go import path for "student.proto"
解决方法：在proto文件头中加入`option go_package ="XXXX";`
`XXXX`代表包名
```

详情可见：[google docs](https://developers.google.com/protocol-buffers/docs/reference/go-generated#package)

### TODO

* var _ PeerPicker = (*HTTPPool)(nil) 的作用是什么？原理是什么

* geecache/day5的测试过程存在疑问
* 关于【分布式】存在疑问
* 可以使用gobenchmark来测试一下

### 参考资料

* [GeekTuTu七天从零实现教程](https://geektutu.com/post/geecache.html)

## 小计
* `VS code中markdown文档，用ctrl+k v可以实现双侧显示`