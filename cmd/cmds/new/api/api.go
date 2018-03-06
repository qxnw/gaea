package api

var mainTmpl = `package main

import (
	"github.com/qxnw/hydra/hydra"
)

func main() {
	hydra := hydra.New(loader())
	defer hydra.Close()
	hydra.Start()
}

`
var TmplMap map[string]string

func init() {
	TmplMap = map[string]string{
		"main.go":                       mainTmpl,
		"loader.go":                     registryTmpl,
		"services/order/order.query.go": orderQueryTmpl,
	}
}
