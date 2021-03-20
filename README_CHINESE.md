<h1>
  <a href="https://github.com/dstotijn/hetty">
    <img src="https://hetty.xyz/assets/logo.png" width="293">
  </a>
</h1>

[![Latest GitHub release](https://img.shields.io/github/v/release/dstotijn/hetty?color=18BA91&style=flat-square)](https://github.com/dstotijn/hetty/releases/latest)
![GitHub download count](https://img.shields.io/github/downloads/dstotijn/hetty/total?color=18BA91&style=flat-square)
[![GitHub](https://img.shields.io/github/license/dstotijn/hetty?color=18BA91&style=flat-square)](https://github.com/dstotijn/hetty/blob/master/LICENSE)
[![Documentation](https://img.shields.io/badge/hetty-docs-18BA91?style=flat-square)](https://hetty.xyz/)


**Hetty**是一个用来做安全研究的HTTP工具箱，它的目标是以开源替代Burp Suite Pro这款商业软件，以强大的功能来用于信息安全和赏金漏洞社区

<img src="https://hetty.xyz/assets/hetty_v0.2.0_header.png">

## 特征
- Man-in-the-middle (MITM) HTTP/1.1 proxy with logs
- 基于SQLite数据库
- Scope 支持
- GraphQL无界面管理
- 嵌入式web接口(Next.js)

ℹ️ Hetty 目前处于早期开发状态，额外的功能目前计划在v1.0版本发布，具体请看 <a href="https://github.com/dstotijn/hetty/projects/1">backlog</a>

## 文档
📖 [Read the docs.](https://hetty.xyz/)

## 安装
Hetty编译成一个二进制文件，内嵌SQLite数据库和基于web的管理接口

### 安装构建好的版本(推荐)
👉 下载Linux，macOS和windows可以在具体看这里 [releases page](https://github.com/dstotijn/hetty/releases).

### 从源代码构建
#### 前期准备
- [Go](https://golang.org/)
- [Yarn](https://yarnpkg.com/)
- [go.rice](https://github.com/GeertJohan/go.rice)

Hetty依赖于SQLite([mattn/go-sqlite3](https://github.com/mattn/go-sqlite3))

同时需要cgo去编译，除此之外i，管理界面的几个静态资源接口Next.js需要通过[Yarn](https://yarnpkg.com/) 生成，同时通过[go.rice](https://github.com/GeertJohan/go.rice)嵌入到.go文件中

从github下载之后通过build来创建二进制包:

```
$ git clone git@github.com:dstotijn/hetty.git
$ cd hetty
$ make build
```
### Docker
在Docker Hub中有Docker的镜像:[`dstotijn/hetty`](https://hub.docker.com/r/dstotijn/hetty)，如何需要持久化CA证书或者文件数据库，挂在一个容器:

```
$ mkdir -p $HOME/.hetty
$ docker run -v $HOME/.hetty:/root/.hetty -p 8080:8080 dstotijn/hetty
```

## 用法
当Hetty在运行时，默认监听的是8080端口的，可以通过http://localhost:8080来访问，

默认情况下，项目的数据库文件和CA证书存储在.hetty目录中，在用户home目录下(Linux/macOS的$HOME，Windows的%USERPROFILE%)

启动的时候，注意在$PAHT配置hetty环境变量，然后启动
```
hetty
```
我们可以查看一下配置项
```
$ hetty -h
Usage of ./hetty:
  -addr string
        TCP address to listen on, in the form "host:port" (default ":8080")
  -adminPath string
        File path to admin build
  -cert string
        CA certificate filepath. Creates a new CA certificate is file doesn't exist (default "~/.hetty/hetty_cert.pem")
  -key string
        CA private key filepath. Creates a new CA private key if file doesn't exist (default "~/.hetty/hetty_key.pem")
  -projects string
        Projects directory path (default "~/.hetty/projects")
```

你就可以看到
```
2020/11/01 14:47:10 [INFO] Running server on :8080 ...
```

然后，通过访问 [http://localhost:8080](http://localhost:8080) 来启动吧

详细的文档正在开发中会很快推出

## 证书安装
为了让Hetty代理HTTPS请求，一个根CA证书需要Hetty去安装，然后这个证书还需要被你的浏览器所信任，下面的步骤是教你如何生成根证书，并将他们提供给hetty，然后你可以把它们安装到本地

下面是在linux机器上完成的，但是也可以对windows和macOS系统提供一定的参考

### 生成CA证书
你可以通过两种方式生成CA密钥对，第一种方式是和Hetty息息相关，流程非常简单，第二种方式是使用OpenSSL来生成它们，这种方式提供了更多对过期时间和加密算法的控制，而且你需要安装OpenSSL工具，还是第一种方式比较适合初学者

#### 通过hetty来生成CA证书
Hetty如果在~/.hetty/中没有找到证书，会生成默认的key和证书，我们只需要通过hetty命令启动hetty就可以了

启动之后，你可以在`~/.hetty/hetty_key.pem`和`~/.hetty/hetty_cert.pem`看到它们

#### 通过OpenSSL来生成CA证书
你可以生成一个在一个月之后失效的CA证书和key

```
mkdir ~/.hetty
openssl req -newkey rsa:2048 -new -nodes -x509 -days 31 -keyout ~/.hetty/hetty_key.pem -out ~/.hetty/hetty_cert.pem
```

hetty会默认在`~/.hetty/`目录下检查CA证书和key，也就是分别对应`hetty_key.pem` 和 `hetty_cert.pem`
，你也可以指定这两个位置参数
```
hetty -key key.pem -cert cert.pem
```
### 信任CA证书
为了让你的浏览器允许流量来本地Hetty代理，你可能需要吧这些证书安装到你的本地CA store

在Ubuntu中，你可以通过下面的命令更新你本地CA store:
```
sudo cp ~/.hetty/hetty_cert.pem /usr/local/share/ca-certificates/hetty.crt
sudo update-ca-certificates
```
在Windows中，你可以通过证书管理器来添加你的证书，通过下面的命令
```
certmgr.msc
```
在MAC中，你可以通过使用Keychain Access来添加证书，这个程序可以在`Application/Utilities/Keychain Access.app`找到，打开之后，把证书拖入到APP，然后在app中打开证书，进入_Trust_部分，在_When using this certificate_里面选择_Always Trust_.

注意:不同的Linux版本可能需要不同的步骤和命令，可以通过查看Linux发行版的文档来让系统信任你的自签名证书

## 愿景
* 用go语言构建的超快速的引擎，最小的内存占用量
* 好用的管理接口，通过Next.js和Material UI来构建
* 通过GraphQL API无界面管理
* 扩展性是核心，所有的模块是go包的形式，可以被用于其他软件
* 通过可插拔的结构构建一个组件式的系统
* 基于渗透测试这和漏洞社区使用者的反馈开发
* 目标是一个相对小的核心特征但是可以满足大部分安全研究者需求

## 支持
使用[issues](https://github.com/dstotijn/hetty/issues)来报告BUG和添加新特性，以及 [discussions](https://github.com/dstotijn/hetty/discussions)来解决问题

## 社区
💬 [Join the Hetty Discord server](https://discord.gg/3HVsj5pTFP).

## 贡献
希望作出贡献，细节可以看这里[Contribution Guidelines](CONTRIBUTING.md)

## 致谢
- 感谢[Hacker101 community on Discord](https://www.hacker101.com/discord)的鼓励和反馈
  
- 字体感谢 [JetBrains Mono](https://www.jetbrains.com/lp/mono/).


## License
[MIT License](LICENSE)

---

© 2020 David Stotijn — [Twitter](https://twitter.com/dstotijn), [Email](mailto:dstotijn@gmail.com)
