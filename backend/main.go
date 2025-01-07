package main

import (
	"todo-go/app/api"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	_ "github.com/mattn/go-sqlite3"
)

func CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func main() {
	// 设置配置文件路径
	g.Cfg().SetFileName("manifest/config/config.toml")

	s := g.Server()

	// 添加 CORS 中间件
	s.Use(CORS)

	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.GET("/todos", api.Todo.List)
		group.POST("/todos", api.Todo.Create)
		group.PUT("/todos/{id}", api.Todo.Update)
		group.DELETE("/todos/{id}", api.Todo.Delete)

		group.GET("/categories", api.Category.List)
		group.POST("/categories", api.Category.Create)
		group.PUT("/categories", api.Category.Update)
		group.DELETE("/categories/{id}", api.Category.Delete)
	})
	s.Run()
}
