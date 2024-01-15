## P.I.G Web Service

---
P.I.G 是一个网页服务程序, 提供了基础的洋葱模型, 最简化的核心模块, 用户可以按照
自己的需求, 通过编写插件或者引入第三方插件的方式, 来扩展 P.I.G 的功能模块.

#### 运行图示

---
![img.png](img.png)

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

- 缺省路由
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
###### 获取请求参数
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

###### 获取原始请求
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
		req := do.MustInvoke[*http.Request](context.Injector())
		context.Json(map[string]interface{}{
			"method": req.Method,
			"uri":    req.RequestURI,
		})
	})

	pig.New().Router(r).Run(8088)
}
```
