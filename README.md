# GoBlog

一个用go写的博客而已

## Swagger生成接口文档

Swagger 相关的工具集会根据 OpenAPI 规范去生成各式各类的与接口相关联的内容，常见的流程是编写注解 =》调用生成库-》生成标准描述文件 =》生成/导入到对应的 Swagger 工具。

### 安装 Go 对应的开源 Swagger 相关联的库

```cmd
$ go get -u github.com/swaggo/swag/cmd/swag@v1.6.5
$ go get -u github.com/swaggo/gin-swagger@v1.2.0 
$ go get -u github.com/swaggo/files
$ go get -u github.com/alecthomas/template


```

在$GOPATH/pkg中找到对应的包cmd/swag/ 执行`go install`生成bin/swag
将$GOPATH/bin放入环境变量

```cmd
$ swag -v
swag version v1.6.5
```

### 注解
每个业务func上面写这些注释
@Summary	摘要
@Produce	API 可以产生的 MIME 类型的列表，MIME 类型你可以简单的理解为响应类型，例如：json、xml、html 等等
@Param	参数格式，从左到右分别为：参数名、入参类型、数据类型、是否必填、注释
@Success	响应成功，从左到右分别为：状态码、参数类型、数据类型、注释
@Failure	响应失败，从左到右分别为：状态码、参数类型、数据类型、注释
@Router	路由，从左到右分别为：路由地址，HTTP 方法

### 返回不在model中的字段
```go
// 用Tag举例， model/tag.go， 再包一层即可
type TagSwagger struct {
	List  []*Tag
	Pager *app.Pager
}

```


### 生成及注册路由
```cmd
$ swag init
```
```go

import (
	...
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func NewRouter() *gin.Engine {
    r := gin.New()
    r.Use(gin.Logger())
    r.Use(gin.Recovery())
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    ...
    return r
}
```
### 访问
http://127.0.0.1:8000/swagger/index.html
