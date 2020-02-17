# LearnUE

[真学虚幻网 - UE4中文社区](https://learnue.com "真学虚幻网 - UE4中文社区") 源码

本项目由UE4中国社区经理[大钊](https://www.zhihu.com/people/fjz13)发起,由[罗传月武](https://www.zhihu.com/people/luochuanyuewu)维护。


## 关于此仓库的说明
本仓库最初fork自[StudyGolang](https://github.com/studygolang/studygolang),并在此基础上继续开发LearnUE,由我开发的内容将在learnue分支持续提交。其他分支以及过往提交内容版权归属于原作者。

在此对StudyGolang的[贡献者们](https://github.com/studygolang/studygolang/graphs/contributors)表示感谢。

同时我会对已有代码根据需求不断改进,目前的开发计划如下(*号表示未完成):

短期目标

    1.适配learnue.com的需求进行代码修改
    2.开发部署流程优化(使用mkcert以支持本地证书签发，并使用caddy配置本地证书实现https反向代理，并使用docker进行开发环境搭建和部署)
    3.配置模块采用viper和.yml配置文件替换旧的基于.ini模式的配置,viper是go社区最火的一个配置库
    4.配置优化
    5.*去掉无用或者不相关的，以及代码写死的模块
   
长期目标

    1.*路由库由echo改为gin(待定),gin是go社区最火的一个配置库。
    2.*ui重构,目前界面ui采用传统mvc方式开发,考虑前端引入vue之类的前端库/或者前端完全工程化，实现前后端分离.

## 参与开发LearnUE

### 环境搭建
LearnUE采用go语言开发。需要你提前安装以下工具。

    1.go语言 1.12+
    
    2.mkcert,本地证书签发工具,目的是生成本地https证书,确保开发环境和生产环境尽可能接近，并给下边的caddy使用。
    
    3.caddy,作为反向代理工具,同时也是部署时用到的web服务器，配合mkcert生成的本地证书，让开发环境尽可能和部署环境一致（参考：learnue/caddy）
    
    4.docker,容器工具,开发中用到的redis,mysql等服务通过配置docker-compose来搭建环境,无需手动下载安装，学会了百利无一害。

部署环境单独分离出一个仓库,感兴趣的可以参考[Lernue_Deploy](https://github.com/ue4-community/learnue_deploy)

### 安装编译

1、下载源码到本地某个非gopath目录

```shell
git clone https://github.com/ue4-community/learnue.git
```

2、编译

进入 learnue 项目目录，执行如下命令：

```shell
// unix/linux系统
export GOPROXY=https://goproxy.cn
go mod vendor && go build -o goapp -mod vendor github.com/ue4-community/learnue/goapp
// windows
set GOPROXY=https://goproxy.cn
go mod vendor && go build -o goapp -mod vendor github.com/ue4-community/learnue/goapp
```

这样便编译好了 learnue

3、在 learnue 源码中的 goapp 目录下应该有了 goapp 可执行文件。

接下来启动 learnue。

```shell
// unix
goapp/goapp
// windows
goapp\goapp.exe
```

一切顺利的话，网站应该就启动了。

4、验证

在浏览器中输入：http://127.0.0.1:8088

应该就能看到了。

接下来你会看到图形化安装界面，一步步照做吧。

* 如果之后有出现页面空白，请查看 error.log 是否有错误

## 参与我们

fork + PR。如果有修改 js 和 css，请执行 gulp （需要先安装 gulp）。注意，Node 版本为：v10.16.2

## 使用该项目搭建的网站

- [真学虚幻网](https://learnue.com)
- [Go语言中文网](https://studygolang.com)
- [Kotlin中国](https://kotlintc.com)
