#### 简介
- 只是一款protocol buffer编译的插件，作用类似protobuf的grpc插件，用于生成dubbogo的go stub

#### 准备
- 首先要有protoc编译工具, 具体安装方法google一下就可以了

```shell script
go get -u github.com/skyitachi/protoc-gen-dubbogo
```

#### 使用指南

```shell script
cd example
# 如果protoc-gen-dubbgo在PATH里面的话
protoc --dubbogo_out=plugins=dubbogo:. user/user.proto

# 指定plugins path
protoc --plugin={plugin_path} --dubbogo_out=plugins=dubbogo:. user/user.proto
```

ps: 除了plugin的名称不一样以外，其他使用方式和grpc的plugin是一样的