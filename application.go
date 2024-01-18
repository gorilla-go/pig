package pig

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/gorilla-go/pig/foundation"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Application struct {
	address    net.IP
	port       int
	middleware []IMiddleware
	router     IRouter
	version    string
}

func New() *Application {
	return &Application{
		address:    net.ParseIP("0.0.0.0"),
		port:       8080,
		middleware: []IMiddleware{},
		version:    "1.0.0-beta",
	}
}

func (a *Application) Use(m ...IMiddleware) *Application {
	a.middleware = append(a.middleware, m...)
	return a
}

func (a *Application) Router(router IRouter) *Application {
	a.router = router
	return a
}

func (a *Application) Run(port ...int) {
	s := http.NewServeMux()
	s.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		NewKernel(a.router).Through(a.middleware).Handle(w, req)
	})

	a.port = foundation.DefaultParam(port, 8080)
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", a.address.String(), a.port),
		Handler: s,
	}

	go func() {
		a.PrintMeta()
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigs
	err := server.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
}

func (a *Application) PrintMeta() {
	color.Cyan(
		"   ___    ____ _____  _      __    __     ____             _        \n  / _ \\  /  _// ___/ | | /| / /__ / /    / __/__ _____  __(_)______ \n / ___/ _/ /_/ (_ /  | |/ |/ / -_) _ \\  _\\ \\/ -_) __/ |/ / / __/ -_)\n/_/  (_)___(_)___/   |__/|__/\\__/_.__/ /___/\\__/_/  |___/_/\\__/\\__/\n\n",
	)
	color.Green(fmt.Sprintf("version: %s", a.version))
	color.Green(fmt.Sprintf("listen:  %s:%d", a.address.String(), a.port))
}
