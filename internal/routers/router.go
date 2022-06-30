package routers

import (
	v1 "GoBlog/internal/routers/api/v1"
	"github.com/gin-gonic/gin"
)


// 路由管理
func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	airticle := v1.NewArticle()
	tag := v1.NewTag()

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
	}

	return r
}
