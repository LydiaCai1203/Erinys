package gee

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		c.Next()
		log.Printf(
			"[%d] %s in %v",
			c.StatusCode,
			c.Req.RequestURI,
			time.Since(t),
		)
	}
}

func LoggerV2() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		c.Fail(500, "Internal Server Error")
		log.Printf(
			"[%d] %s in %v for group v1",
			c.StatusCode,
			c.Req.RequestURI,
			time.Since(t),
		)
	}
}
