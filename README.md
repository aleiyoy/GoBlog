# GoBlog
一个用go写的博客而已

## 数据块
### 文章表
### 标签表：方便Seo收录
### 关系表：标签-文章 = 一对多
一个标签可以有多个文章，一个文章一个标签


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


## 处理 model 回调
针对我们的公共字段 created_on、modified_on、deleted_on、is_del 进行处理，如果每一个 DB 操作都去设置公共字段的值，那么不仅多了很多重复的代码，在要调整公共字段时工作量也会翻倍。
采用设置 model callback 的方式去实现公共字段的处理，本项目使用的 ORM 库是 GORM，GORM 本身是提供回调支持的，因此我们可以根据自己的需要自定义 GORM 的回调操作，而在 GORM 中我们可以分别进行如下的回调相关行为：
- 注册一个新的回调。
- 删除现有的回调。
- 替换现有的回调。
- 注册回调的先后顺序。

在本项目中使用到的“替换现有的回调”这一行为

## Gorm不判断零值是 有意义值还是空值

```go
func (t Tag) Update(db *gorm.DB, values interface{}) error {
	// 这里不用结构体更新数据库，是因为gorm很难判断结构体中state字段的0，是空值还是真实并有意义的值
	if err := db.Model(t).Where("id = ? AND is_del = ?", t.ID, 0).Updates(values).Error; err != nil {
		return err
	}
	return nil

	//return db.Model(&Tag{}).Where("id = ? AND is_del = ?", t.ID, 0).Update(t).Error
}
```

## validater验证器会认为 state=0  为空值，并非有意义的0
可以将接口参数结构体的state字段改成指针类型
```go
	newState := convert.StrTo(c.Param("state")).MustUInt8()
	param := service.UpdateTagRequest{
		ID:    convert.StrTo(c.Param("id")).MustUInt32(),
		State: &newState,  // 这里不用指针验证器会认为0是空值
	}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
```


## 测试用的cmd
```json
curl -X POST http://127.0.0.1:8000/api/v1/articles -F 'title=Go的Gin框架教程' -F 'tag_id=1' -F 'content=内容' -F 'cover_image_url=https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d08a7d55720d6cf.png'   -F created_by=yl

curl -X POST http://127.0.0.1:8000/api/v1/articles -F 'title=Go的Gorm教程' -F 'tag_id=1' -F 'content=内容2' -F 'cover_image_url=https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d08a7d55720d6cf.png'   -F created_by=yl



curl -X PUT http://127.0.0.1:8000/api/v1/articles/1 -F 'title=Go的Gin框架教程' -F 'tag_id=1' -F 'content=内容2--更新后' -F 'cover_image_url=https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d08a7d55720d6cf.png'   -F desc=简述--新增的    -F modified_by=yl2


curl -X GET http://127.0.0.1:8000/api/v1/articles/1


curl -X GET http://127.0.0.1:8000/api/v1/articles -F tag_id=1
{
    "list":[
        {
            "id":1,
            "title":"Go的Gin框架教程",
            "desc":"",
            "content":"内容",
            "cover_image_url":"https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d08a7d55720d6cf.png",
            "state":0,
            "tag":{
                "id":1,
                "created_by":"",
                "modified_by":"",
                "created_on":0,
                "modified_on":0,
                "deleted_on":0,
                "is_del":0,
                "name":"Go",
                "state":0
            }
        },
        {
            "id":2,
            "title":"Go的Gorm教程",
            "desc":"",
            "content":"内容2",
            "cover_image_url":"https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d08a7d55720d6cf.png",
            "state":0,
            "tag":{
                "id":1,
                "created_by":"",
                "modified_by":"",
                "created_on":0,
                "modified_on":0,
                "deleted_on":0,
                "is_del":0,
                "name":"Go",
                "state":0
            }
        }
    ],
    "pager":{
        "page":1,
        "page_size":10,
        "total_rows":2
    }
}



curl -X DELETE http://127.0.0.1:8000/api/v1/articles/2 

```