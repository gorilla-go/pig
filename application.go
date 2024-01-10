package p_i_g

import "net"

type Application struct {
	port    int
	address net.IP
}

func New() *Application {
	return &Application{}
}

func (a *Application) Start() {

}
