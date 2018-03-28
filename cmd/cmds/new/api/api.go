package api

var mainTmpl = `
package main

import (
	"github.com/qxnw/hydra/engines"
	"github.com/qxnw/hydra/hydra"
)

func main() {
	engines.AddServiceLoader(loader())
	app := hydra.NewApp(
		hydra.WithPlatName("@pShortName"),
		hydra.WithSystemName("test"),
		hydra.WithServerTypes("api"),
		hydra.WithAutoCreateConf(true),
		hydra.WithDebug())
	app.Start()
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
