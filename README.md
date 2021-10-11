githubIpTool
=====
github可用IP获取工具 golang

> 有关该项目的说明详见 https://del.pub/githubiptool

![预览图](https://cdn.jsdelivr.net/gh/mopo/githubiptool@master/preview.png)

## 1. 依赖
* 纯真IP库
> 下载地址：https://update.cz88.net/soft/setup.zip

* github官方API
> 地址: https://api.github.com/meta

## 2. 使用
```bash
git clone git@github.com:mopo/githubiptool.git
cd githubiptool
go mod tidy
go run main.go
```

## 3. 其它
> 线程超时时间是10秒，若是全部失败，可以考虑加长时间再试,但是程序执行时间会加长

## 4. license
Licensed under MIT.