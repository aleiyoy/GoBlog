package routers

import (
	_ "GoBlog/docs"
	"GoBlog/global"
	"GoBlog/internal/middleware"
	v1 "GoBlog/internal/routers/api/v1"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)


// 路由管理
func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// 自定义错误信息翻译中间件
	// go-playground/validator 默认的错误信息是英文，但我们的错误信息不一定是用的英文，有可能要简体中文
	r.Use(middleware.Translations())

	// 接口文档的路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.StaticFS("/static", http.Dir(global.AppSetting.UploadSavePath))

	airticle := v1.NewArticle()
	tag := v1.NewTag()
	upload := v1.NewUpload()

	apiv1 := r.Group("/api/v1")
	{
		apiv1.POST("/tags", tag.Create)
		apiv1.DELETE("/tags/:id", tag.Delete)
		apiv1.PUT("/tags/:id", tag.Update)
		apiv1.PATCH("/tags/:id/state", tag.Update)
		apiv1.GET("/tags", tag.List)

		apiv1.POST("/articles", airticle.Create)
		apiv1.DELETE("/articles/:id", airticle.Delete)
		apiv1.PUT("/articles/:id", airticle.Update)
		apiv1.PATCH("/articles/:id/state", airticle.Update)
		apiv1.GET("/articles/:id", airticle.Get)
		apiv1.GET("/articles", airticle.List)

		r.POST("/upload/file", upload.UploadFile)
	}

	return r
}
