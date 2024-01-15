## P.I.G Web Service

---
P.I.G 是一个网页服务程序, 提供了基础的洋葱模型, 最简化的核心模块, 用户可以按照
自己的需求, 通过编写插件或者引入第三方插件的方式, 来扩展 P.I.G 的功能模块.

#### 运行图示

---
![img.png](test/img.png)

#### 功能示例

---
##### Hello world
```go
package main

import (
"github.com/gorilla-go/pig"
)

func main() {
	r := pig.NewRouter()
	r.GET("/", func(context *pig.Context) {
		context.Echo("hello world")
	})

	pig.New().Router(r).Run(8088)
}
```

##### 路由

---
> P.I.G 提供了基础的路由功能, 用户可以通过编写路由规则, 来实现不同的功能.
> 也可以自行实现路由接口, 接入自定义路由或者第三方路由.

---
###### 基础路由
```go
package main
import (
	"github.com/gorilla-go/pig"
)

func main() {
	r := pig.NewRouter()
	r.GET("/user", func(context *pig.Context) {
		context.Echo("hello world")
	})

	pig.New().Router(r).Run(8088)
}
```

###### 泛参数路由
```go
package main
import (
    "github.com/gorilla-go/pig"
)
func main() {
	r := pig.NewRouter()
	r.GET("/user/:id", func(context *pig.Context) {
		context.Echo(context.ParamVar().TrimString("id"))
	})

	pig.New().Router(r).Run(8088)
}
```

###### 正则参数路由
```go
package main
import (
    "github.com/gorilla-go/pig"
)
func main() {
	r := pig.NewRouter()
	r.GET("/user/<id:\\d+>", func(context *pig.Context) {
		context.Echo(context.ParamVar().TrimString("id"))
	})

	pig.New().Router(r).Run(8088)
}
```

###### 缺省路由
```go
package main

import (
	"github.com/gorilla-go/pig"
)

func main() {
	r := pig.NewRouter()
	r.Miss(func(context *pig.Context) {
		context.Echo("404")
	})

	pig.New().Router(r).Run(8088)
}
```

##### 中间件

---
> P.I.G 实现了基础的中间件功能, 用户可以通过编写中间件, 来实现不同的功能.
---

###### 前置中间件
```go
package main

import (
	"fmt"
	"github.com/gorilla-go/pig"
)

type Middleware struct {
}

func (m *Middleware) Handle(context *pig.Context, next func(*pig.Context)) {
	fmt.Println("Middleware")
	next(context)
}

func main() {
	r := pig.NewRouter()
	r.GET("/", func(context *pig.Context) {
		context.Echo("Hello, World")
	})

	pig.New().Use(&Middleware{}).Router(r).Run(8088)
}
```

###### 后置中间件
```go
package main

import (
	"fmt"
	"github.com/gorilla-go/pig"
)

type Middleware struct {
}

func (m *Middleware) Handle(context *pig.Context, next func(*pig.Context)) {
	next(context)
	fmt.Println("Middleware")
}

func main() {
	r := pig.NewRouter()
	r.GET("/", func(context *pig.Context) {
		context.Echo("Hello, World")
	})

	pig.New().Use(&Middleware{}).Router(r).Run(8088)
}
```

###### 路由中间件
```go
package main

import (
	"fmt"
	"github.com/gorilla-go/pig"
)

type Middleware struct {
}

func (m *Middleware) Handle(context *pig.Context, next func(*pig.Context)) {
	fmt.Println("Route Middleware")
	next(context)
}

func main() {
	r := pig.NewRouter()
	r.GET("/", func(context *pig.Context) {
		context.Echo("Hello, World")
	}, &Middleware{})

	pig.New().Router(r).Run(8088)
}
```
> 请注意, 路由中间件只会在指定的路由中生效. 并且定义了路由中间件, 全局中间件将不再生效.
> 如果需要需要全局中间件, 请在路由中间件中包含全局中间件.

##### 请求

---
###### 请求参数
```go
package main

import (
	"github.com/gorilla-go/pig"
)

func main() {
	r := pig.NewRouter()
	
	// Get 请求参数
	r.GET("/:id", func(context *pig.Context) {
		context.Json(map[string]interface{}{
			"id": context.ParamVar().Int("id"),
		})
	})

	// Post 请求参数
	r.POST("/post/:id", func(context *pig.Context) {
		context.Json(map[string]interface{}{
			"id":   context.ParamVar().Int("id"),
			"post": context.PostVar().String("post"),
		})
	})

	pig.New().Router(r).Run(8088)
}
```

###### 文件上传
```go
package main

import (
	"github.com/gorilla-go/pig"
)

func main() {
	r := pig.NewRouter()
	r.POST("/upload", func(context *pig.Context) {
		filePath := context.FileVar()["file"].FilePath
		context.Echo(filePath)
	})
	pig.New().Router(r).Run(8088)
}
```

###### 原始请求
```go
package main

import (
	"github.com/gorilla-go/pig"
	"github.com/samber/do/v2"
	"net/http"
)

func main() {
	r := pig.NewRouter()
	r.GET("/:id", func(context *pig.Context) {
		req := context.Request()
		context.Json(map[string]interface{}{
			"method": req.Method,
			"uri":    req.RequestURI,
		})
	})

	pig.New().Router(r).Run(8088)
}
```
##### 响应

---
###### 文本
```go
package main

import (
	"github.com/gorilla-go/pig"
)

func main() {
	r := pig.NewRouter()
	r.GET("/", func(context *pig.Context) {
		context.Echo("hello world")
	})

	pig.New().Router(r).Run(8088)
}
```

###### JSON
```go
package main

import (
	"github.com/gorilla-go/pig"
)

func main() {
	r := pig.NewRouter()
	r.GET("/", func(context *pig.Context) {
		context.Json(map[string]interface{}{
			"id":   context.ParamVar().Int("id"),
		})
	})

	pig.New().Router(r).Run(8088)
}
```

###### 文件下载
```go
package main

import (
	"github.com/gorilla-go/pig"
)

func main() {
	r := pig.NewRouter()
	r.GET("/download", func(context *pig.Context) {
		context.Download(
			pig.NewFile("/your/file/path.jpg"),
			"filename.jpg",
		)
	})

	pig.New().Router(r).Run(8088)
}
```

###### 重定向
```go
package main

import (
	"github.com/gorilla-go/pig"
)

func main() {
	r := pig.NewRouter()
	r.GET("/redirect", func(context *pig.Context) {
		context.Redirect("/redirected", 302)
	})

	pig.New().Router(r).Run(8088)
}
```
###### 原生响应
```go
package main

import (
	"github.com/gorilla-go/pig"
)

func main() {
	r := pig.NewRouter()
	r.GET("/", func(context *pig.Context) {
		context.ResponseWriter().Write([]byte("hello world"))
	})

	pig.New().Router(r).Run(8088)
}
```

##### 最佳实践

---
###### 默认参数
```go
package main

import (
	"fmt"
	"github.com/gorilla-go/pig/foundation"
)

func main() {
	DefaultParams("no default")
	DefaultParams("default", 200)
}

func DefaultParams(p string, i ...int) {
	defCode := foundation.DefaultParam(i, 0)
	fmt.Println(defCode)
}
```

###### 日志
```go
package main

import (
	"github.com/gorilla-go/pig"
	"github.com/samber/do/v2"
	"log"
)

type Logger struct{}

func (*Logger) Info(message string, c *pig.Context) {
	log.Println(message)
}

func (*Logger) Debug(message string, c *pig.Context) {
	log.Println(message)
}

func (*Logger) Warning(message string, c *pig.Context) {
	log.Println(message)
}

func (*Logger) Fatal(message string, c *pig.Context) {
	log.Println(message)
}

type Middleware struct {
}

func (*Middleware) Handle(c *pig.Context, next func(*pig.Context)) {
	do.ProvideValue[pig.ILogger](c.Injector(), &Logger{})
	next(c)
}

func main() {
	r := pig.NewRouter()
	r.GET("/", func(c *pig.Context) {
		c.Logger().Info("Hello World!", c)
	})

	pig.New().Use(&Middleware{}).Router(r).Run(8088)
}
```
> 用户自行实现日志处理或加载第三方日志框架.

###### 错误处理
```go
package main

import (
	"fmt"
	"github.com/gorilla-go/pig"
	"github.com/samber/do/v2"
)

type HttpErrorHandler struct {
}

func (h *HttpErrorHandler) Handle(a any, context *pig.Context) {
	fmt.Println("error targeted")
	context.Echo("500", 500)
}

type Middleware struct {
}

func (*Middleware) Handle(c *pig.Context, next func(*pig.Context)) {
	do.ProvideValue[pig.IHttpErrorHandler](c.Injector(), &HttpErrorHandler{})
	next(c)
}

func main() {
	r := pig.NewRouter()
	r.GET("/", func(c *pig.Context) {
		panic("error")
	})

	pig.New().Use(&Middleware{}).Router(r).Run(8088)
}

```


