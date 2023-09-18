package main

import (
	"myRPC/gee"
	"net/http"
)

func main() {
	r := gee.New()

	r.GET(
		"/",
		func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>hello Gee</h1>")
		},
	)

	// -> "/hello?name=caiqj"
	r.GET(
		"/hello",
		func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you are at %s\n", c.Query("name"), c.Path)
		},
	)

	// -> "/hello/caiqj"
	r.GET(
		"/hello/:name",
		func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you are at %s\n", c.Param("name"), c.Path)
		},
	)

	r.GET(
		"/assets/*filepath",
		func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
		},
	)

	r.Run(":9999")
}
