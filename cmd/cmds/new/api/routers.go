package api

var routerTmpl = `package main

import (
	"{@pImportName}/handlers"

	"github.com/qxnw/goplugin"
)

var reg *goplugin.Registry

func init() {
	reg = goplugin.NewRegistry()
}

func Register(name string, handler goplugin.Handler) {
	reg.Register(name, handler)
}
func GetServices() []string {
	return reg.Services
}
func GetHandlers() map[string]goplugin.Handler {
	return reg.ServiceHandlers
}

func init() {
	Register("/order/query", handlers.NewOrderQuery())
}
`
