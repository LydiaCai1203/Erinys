package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func trace(message string) string {
	var str strings.Builder
	str.WriteString(message + "\nTraceback:")

	// uintptr 无符号整数类型
	// 用于存储 Go 程序中的内存地址
	var pcs [32]uintptr
	// 3: 指定从调用堆栈的哪个深度开始收集信息
	// pcs: 存储调用堆栈中函数的程序计数器的地址
	n := runtime.Callers(3, pcs[:])
	for _, pc := range pcs[:n] {
		// 获取与程序计数器 pc 关联的函数信息
		fn := runtime.FuncForPC(pc)
		// 获取函数所在的文件和行号
		file, line := fn.FileLine(pc)
		// 将文件还有行号信息加入字符串构建器
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

// Recovery 中间件
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				// 返回 500 内部服务器错误
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		c.Next()
	}
}
