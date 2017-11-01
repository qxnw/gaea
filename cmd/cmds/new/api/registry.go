package api

var routerTmpl = `package main

import (
	"{@pImportName}/services"

	"github.com/qxnw/hydra/context"
)

var reg *context.Registry

func init() {
	reg = context.NewRegistry()
}

func Register(name string, handler interface{}) {
	reg.Register(name, handler)
}
func GetServices() []string {
	return reg.Services
}
func Handle(name string, mode string, service string, ctx *context.Context) (r context.Response, err error) {
	return reg.Handle(name, mode, service, ctx)
}

func init() {
	Register("/myservice", services.NewMyService())
}`
