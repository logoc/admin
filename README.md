
### 1. 后端代码
#### 安装fresh 热更新-边开发边编译
```
go install github.com/pilu/fresh@latest
```
#### 初始化mod
```
go mod tidy
```
#### 热编译运行
```
bee run 或 fresh
```
#### 打包
```
go build main.go
```
#### 打包（此时会打包成Linux上可运行的二进制文件，不带后缀名的文件）
```
SET GOOS=linux
SET GOARCH=amd64
go build
```
#### widows
```
// 配置环境变量
SET CGO_ENABLED=1
SET GOOS=windows
SET GOARCH=amd64

go build main.go

// 编译命令
```
#### 编译成Linux环境可执行文件
```

// 配置参数
SET CGO_ENABLED=0 
SET GOOS=linux 
SET GOARCH=amd64 

go build main.go

// 编译命令
```
#### 服务器部署
部署是把打包生成的二进制文件(Linux:gofly，windows:gofly.exe)和资源文件resource复制过去即可。
### 2. 前端端代码
#### 初始化依赖
 ```
 npm install 或者 yarn install
 ```
如果第一次使用Arco Design Pro install初始化可以报错，如果保存请运行下面命令（安装项目模版的工具）：
```
npm i -g arco-cli
```
#### 运行
```
npm run serve 或者  yarn serve
```
#### 打包
```
npm run build 或者 yarn build
```

## 八、前端代码安装及源码位置
由于框架是前端后端分离，且在Go本地开发使用fresh热编译，Go目录不能用太多文件影响编译时间，
所以我们开发是建议前端代码放在其他位置。在安装界面填写你前端代码放置位置或者手动在Go项目config/settings.yml配置文件中vueobjroot手动配置前端业务端开发路径：
```
vueobjroot: D:/Project/develop/vue/gofly_base/gofly_business
```
如果你想要手动安装前端代码，源码在代码包的resource/staticfile/template/vuecode目录下文件夹中，自己复制到开发文件夹下即可。

## 服务部署
### 后端
```
/bin
/log
/resource
/start.sh
/stop.sh
```