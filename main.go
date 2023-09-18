package main

import (
	"myRPC/gee"
	"net/http"
)

func main() {
	r := gee.New()
	// r *Engine: 继承了 RouterGroup 的所有方法，可以直接调用 GET
	r.GET(
		"/index",
		func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>Index Page</h1>")
		},
	)

	v1 := r.Group("/v1")
	{
		v1.GET(
			"/",
			func(c *gee.Context) {
				c.HTML(http.StatusOK, "<h1>hello Gee</h1>")
			},
		)

		// -> "/hello?name=caiqj"
		v1.GET(
			"/hello",
			func(c *gee.Context) {
				c.String(http.StatusOK, "hello %s, you are at %s\n", c.Query("name"), c.Path)
			},
		)
	}
	v2 := r.Group("/v2")
	{
		// -> "/hello/caiqj"
		v2.GET(
			"/hello/:name",
			func(c *gee.Context) {
				c.String(http.StatusOK, "hello %s, you are at %s\n", c.Param("name"), c.Path)
			},
		)

		v2.POST(
			"/login",
			func(c *gee.Context) {
				c.JSON(
					http.StatusOK,
					gee.H{
						"username": c.PostForm("username"),
						"password": c.PostForm("password"),
					},
				)
			},
		)
	}
	r.Run(":9999")
}
