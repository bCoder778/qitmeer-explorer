package api

import (
	"os"
	"os/signal"
)

type Api struct {
	rest *RestApi
}

func NewApi(addr string) *Api {
	return &Api{
		rest: NewRestApi(addr),
	}
}

func (a *Api) Run() {
	a.listenInterrupt()
	a.rest.Start()
}

func (a *Api) listenInterrupt() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)

	go func() {
		<-c
		a.rest.Stop()
	}()
}
